package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// beitreten lässt einen frischen Client (eigenes Cookie-Jar) einer Gruppe
// beitreten und gibt ihn zurück.
func beitreten(t *testing.T, app *testApp, kursID int64, gruppenNummer int) *http.Client {
	t.Helper()
	client := clientMitCookieJar(t)
	resp, err := client.PostForm(
		app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/"+strconv.Itoa(gruppenNummer)+"/beitreten",
		nil,
	)
	if err != nil {
		t.Fatalf("beitreten: %v", err)
	}
	resp.Body.Close()
	return client
}

func gruppenStatusText(t *testing.T, app *testApp, client *http.Client) string {
	t.Helper()
	resp, err := client.Get(app.server.URL + "/gruppe")
	if err != nil {
		t.Fatalf("GET /gruppe: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}
	return string(body)
}

func TestErstesGeraetHaeltAutomatischDasSchreibrecht(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	erstesGeraet := beitreten(t, app, kursID, 1)

	if !strings.Contains(gruppenStatusText(t, app, erstesGeraet), "Du hast das Schreibrecht") {
		t.Errorf("erwarte, dass das erste Gerät automatisch das Schreibrecht hat")
	}
}

func TestZweitesGeraetHatZunaechstKeinSchreibrecht(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	beitreten(t, app, kursID, 1)
	zweitesGeraet := beitreten(t, app, kursID, 1)

	if strings.Contains(gruppenStatusText(t, app, zweitesGeraet), "Du hast das Schreibrecht") {
		t.Errorf("erwarte, dass das zweite Gerät noch kein Schreibrecht hat")
	}
}

func TestGeraetKannSchreibrechtUebernehmen(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	erstesGeraet := beitreten(t, app, kursID, 1)
	zweitesGeraet := beitreten(t, app, kursID, 1)

	resp, err := zweitesGeraet.PostForm(app.server.URL+"/gruppe/schreibrecht", nil)
	if err != nil {
		t.Fatalf("POST /gruppe/schreibrecht: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	if !strings.Contains(gruppenStatusText(t, app, zweitesGeraet), "Du hast das Schreibrecht") {
		t.Errorf("erwarte, dass das zweite Gerät nach der Übernahme das Schreibrecht hat")
	}
	if strings.Contains(gruppenStatusText(t, app, erstesGeraet), "Du hast das Schreibrecht") {
		t.Errorf("erwarte, dass das erste Gerät das Schreibrecht verloren hat")
	}
}

// TestSchreibrechtUebernahmeAusserhalbDerWerkstattPhaseIstGesperrt deckt
// dieselbe Phasensperre ab, die seit Vorgang 0011 für jede andere
// schreibende Aktion gilt (siehe geraetMitSchreibrecht in werkstatt.go) —
// die Übernahme selbst darf davon keine Ausnahme sein.
func TestSchreibrechtUebernahmeAusserhalbDerWerkstattPhaseIstGesperrt(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)
	erstesGeraet := beitreten(t, app, kursID, 1)
	zweitesGeraet := beitreten(t, app, kursID, 1)

	setzePhase(t, app, kursID, phaseVorbereitung)

	resp, err := zweitesGeraet.PostForm(app.server.URL+"/gruppe/schreibrecht", nil)
	if err != nil {
		t.Fatalf("POST /gruppe/schreibrecht: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("erwarte 409 außerhalb der Werkstatt-Phase, habe %d", resp.StatusCode)
	}

	setzePhase(t, app, kursID, phaseWerkstatt)
	if !strings.Contains(gruppenStatusText(t, app, erstesGeraet), "Du hast das Schreibrecht") {
		t.Errorf("erwarte, dass das erste Gerät das Schreibrecht durch die abgelehnte Übernahme behalten hat")
	}
}

func TestGeraetOhneSchreibrechtKannEsNichtVortaeuschen(t *testing.T) {
	// Kein HTTP-Endpunkt speichert heute schon Eingaben (Vorgang 0008
	// bringt den Zustandsautomat) — geprüft wird deshalb die Funktion, die
	// jeder künftige Schreib-Handler vor dem Speichern aufrufen muss, gegen
	// eine echte, temporäre SQLite-Datei.
	app := newTestApp(t)
	kursID, _, gruppenIDs := app.KursAnlegen(t, 1)
	beitreten(t, app, kursID, 1)
	_, ohneRecht := app.GeraetBeitreten(t, gruppenIDs[0])

	hat, err := geraetHatSchreibrecht(app.db, gruppenIDs[0], mustGeraetID(t, app, ohneRecht))
	if err != nil {
		t.Fatalf("geraetHatSchreibrecht: %v", err)
	}
	if hat {
		t.Errorf("erwarte, dass ein zweites, unbeteiligtes Gerät kein Schreibrecht hat")
	}
}

func mustGeraetID(t *testing.T, app *testApp, token string) int64 {
	t.Helper()
	var id int64
	if err := app.db.QueryRow(`SELECT id FROM geraet WHERE token = ?`, token).Scan(&id); err != nil {
		t.Fatalf("geraet-id lesen: %v", err)
	}
	return id
}

func TestSchreibrechtWirdBeiGleichzeitigemBeitrittNurEinmalVergeben(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	const anzahl = 5
	var wg sync.WaitGroup
	clients := make([]*http.Client, anzahl)
	for i := range clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			clients[i] = beitreten(t, app, kursID, 1)
		}(i)
	}
	wg.Wait()

	traegerZahl := 0
	for _, client := range clients {
		if strings.Contains(gruppenStatusText(t, app, client), "Du hast das Schreibrecht") {
			traegerZahl++
		}
	}
	if traegerZahl != 1 {
		t.Errorf("erwarte genau ein Gerät mit Schreibrecht bei gleichzeitigem Beitritt, habe %d", traegerZahl)
	}
}
