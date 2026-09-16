package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestUnterschriftWirdHinzugefuegtUndErscheintInDerVorschau(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna"}})
	if err != nil {
		t.Fatalf("unterschrift hinzufügen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	body := koerperText(t, resp)
	if !strings.Contains(body, "Anna") {
		t.Errorf("erwarte 'Anna' in der Vorschau, habe: %s", body)
	}
}

func TestUnterschriftIstNichtErzwungenUndDieTheseBleibtFreigebbar(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	// Absichtlich KEINE Unterschrift eingetragen.
	bringeBisFreigegeben(t, app, client,
		"Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	)

	resp, err := client.Get(app.server.URL + "/gruppe/freigegeben")
	if err != nil {
		t.Fatalf("GET /gruppe/freigegeben: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("erwarte, dass eine These ohne Unterschrift trotzdem freigebbar ist, habe %d", resp.StatusCode)
	}
}

func TestLeererVornameWirdAbgelehnt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"   "}})
	if err != nil {
		t.Fatalf("unterschrift hinzufügen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("erwarte 400 für einen leeren Vornamen, habe %d", resp.StatusCode)
	}
}

// TestVornameMitTrennzeichenKannKeineFremdeUnterschriftVortaeuschen deckt die
// Eingabebereinigung ab: Ein Vorname, der zufällig das interne Trennzeichen
// oder einen Zeilenumbruch enthält, darf die gespeicherte Liste nicht in
// mehrere Einträge zerlegen.
func TestVornameMitTrennzeichenKannKeineFremdeUnterschriftVortaeuschen(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna; Ben\nChiara"}})
	if err != nil {
		t.Fatalf("unterschrift hinzufügen: %v", err)
	}
	resp.Body.Close()

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	these := aktuellerSchritt(t, app, gruppeID)
	liste := unterschriftenTeilen(these.Unterschriften)
	if len(liste) != 1 {
		t.Fatalf("erwarte genau einen Listeneintrag, habe %v", liste)
	}
	if liste[0] != "Anna Ben Chiara" {
		t.Errorf("erwarte bereinigten Namen 'Anna Ben Chiara', habe %q", liste[0])
	}
}

func TestUnterschriftKannEinzelnWiederEntferntWerden(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")

	for _, name := range []string{"Anna", "Ben", "Chiara"} {
		resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {name}})
		if err != nil {
			t.Fatalf("unterschrift %q hinzufügen: %v", name, err)
		}
		resp.Body.Close()
	}

	// Ben steht an Index 1 — nur er soll verschwinden, Anna und Chiara bleiben.
	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/entfernen", url.Values{"index": {"1"}})
	if err != nil {
		t.Fatalf("unterschrift entfernen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	these := aktuellerSchritt(t, app, gruppeID)
	liste := unterschriftenTeilen(these.Unterschriften)
	if len(liste) != 2 || liste[0] != "Anna" || liste[1] != "Chiara" {
		t.Errorf("erwarte ['Anna', 'Chiara'] nach dem Entfernen, habe %v", liste)
	}
}

func TestUnterschriftHinzufuegenNachFreigabeLiefertKonflikt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisFreigegeben(t, app, client,
		"Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	)

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna"}})
	if err != nil {
		t.Fatalf("unterschrift hinzufügen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("erwarte 409 nach Freigabe, habe %d", resp.StatusCode)
	}
}

func TestUnterschriftEntfernenNachFreigabeLiefertKonflikt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna"}})
	resp.Body.Close()

	resp, _ = client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"ja"}})
	resp.Body.Close()

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/entfernen", url.Values{"index": {"0"}})
	if err != nil {
		t.Fatalf("unterschrift entfernen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("erwarte 409 nach Freigabe, habe %d", resp.StatusCode)
	}
}

func TestUngueltigerIndexBeimEntfernenWirdAbgelehnt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna"}})
	resp.Body.Close()

	resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/entfernen", url.Values{"index": {"5"}})
	if err != nil {
		t.Fatalf("unterschrift entfernen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("erwarte 400 für einen Index außerhalb der Liste, habe %d", resp.StatusCode)
	}
}

// TestOptimistischesSchlossVerwirftEineVeralteteEntfernung ruft
// unterschriftEntfernenWennUnveraendert direkt auf — das Wettlauffenster
// zwischen Lesen und Schreiben, gegen das es schützt, lässt sich über echte
// HTTP-Anfragen nicht deterministisch erzwingen (kein künstlicher Timing-Hebel
// in den Handlern, und keiner soll dafür eingebaut werden). Der Test simuliert
// die dazwischenkommende Änderung stattdessen direkt in der Datenbank.
func TestOptimistischesSchlossVerwirftEineVeralteteEntfernung(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")
	for _, name := range []string{"Anna", "Ben"} {
		resp, err := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {name}})
		if err != nil {
			t.Fatalf("unterschrift %q hinzufügen: %v", name, err)
		}
		resp.Body.Close()
	}

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	altwertVorDerRace := aktuellerSchritt(t, app, gruppeID).Unterschriften // "Anna; Ben"

	// Ein zweites Gerät derselben Gruppe entfernt "Ben", während unser
	// (simulierter) erster Aufruf noch den alten Stand "Anna; Ben" im Kopf hat.
	if _, err := app.db.Exec(`UPDATE these SET unterschriften = ? WHERE gruppe_id = ?`, "Anna", gruppeID); err != nil {
		t.Fatalf("konkurrierende änderung simulieren: %v", err)
	}

	// Unser Aufruf versucht nun, ausgehend vom veralteten Stand, "Anna" zu
	// entfernen und durch "Ben" zu ersetzen — das darf die zwischenzeitliche
	// Änderung nicht überschreiben.
	geaendert, err := unterschriftEntfernenWennUnveraendert(app.db, gruppeID, altwertVorDerRace, "Ben")
	if err != nil {
		t.Fatalf("unterschriftEntfernenWennUnveraendert: %v", err)
	}
	if geaendert {
		t.Errorf("erwarte, dass eine veraltete Entfernung verworfen wird")
	}

	these := aktuellerSchritt(t, app, gruppeID)
	if these.Unterschriften != "Anna" {
		t.Errorf("erwarte, dass die konkurrierende Änderung erhalten bleibt, habe %q", these.Unterschriften)
	}
}

func TestUnterschriftenErscheinenAufDemDruckblattAberNichtAmBeamer(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/unterschriften/hinzufuegen", url.Values{"vorname": {"Anna"}})
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"ja"}})
	resp.Body.Close()

	druckResp, err := client.Get(app.server.URL + "/gruppe/freigegeben")
	if err != nil {
		t.Fatalf("GET /gruppe/freigegeben: %v", err)
	}
	defer druckResp.Body.Close()
	if !strings.Contains(koerperText(t, druckResp), "Anna") {
		t.Errorf("erwarte 'Anna' auf dem Druckblatt")
	}

	setzePhase(t, app, kursID, phaseAbschluss)
	beamerResp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beamer")
	if err != nil {
		t.Fatalf("beamer laden: %v", err)
	}
	defer beamerResp.Body.Close()
	if strings.Contains(koerperText(t, beamerResp), "Anna") {
		t.Errorf("erwarte KEINEN Namen auf der Beamer-Ansicht")
	}
}
