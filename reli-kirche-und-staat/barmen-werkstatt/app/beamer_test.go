package main

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

// sseVerbindenBeamer verbindet sich wie sseVerbinden (sse_test.go) per SSE,
// aber ohne Geräte-Cookie — der Beamer hat keine Gerätidentität, sondern
// hört auf alle Gruppen eines Kurses.
func sseVerbindenBeamer(t *testing.T, app *testApp, kursID int64) *sseVerbindung {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/beamer/live", nil)
	if err != nil {
		t.Fatalf("beamer-sse-anfrage bauen: %v", err)
	}

	resp, err := app.server.Client().Do(req)
	if err != nil {
		t.Fatalf("beamer-sse verbinden: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		cancel()
		t.Fatalf("beamer-sse verbinden: erwarte 200, habe %d", resp.StatusCode)
	}

	zeilen := make(chan string, 32)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			zeile := scanner.Text()
			if strings.HasPrefix(zeile, "data:") {
				zeilen <- zeile
			}
		}
	}()

	v := &sseVerbindung{resp: resp, zeilen: zeilen}
	t.Cleanup(func() {
		cancel()
		resp.Body.Close()
	})
	return v
}

func TestBeamerZeigtBeitrittscodeVorFreischaltung(t *testing.T) {
	app := newTestApp(t)
	kursID, code, _ := app.KursAnlegen(t, 1)
	setzePhase(t, app, kursID, phaseVorbereitung)

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beamer")
	if err != nil {
		t.Fatalf("beamer laden: %v", err)
	}
	defer resp.Body.Close()
	body := koerperText(t, resp)

	// Der Beitrittscode selbst steht nur im <img>-Link zum QR-Code, nicht als
	// Klartext — geprüft wird deshalb, dass die Wartefläche erscheint und der
	// QR-Endpunkt des richtigen Kurses eingebunden ist.
	if !strings.Contains(body, "Gleich geht's los") {
		t.Errorf("erwarte eine Wartefläche vor der Freischaltung, habe: %s", body)
	}
	if !strings.Contains(body, "/kurse/"+strconv.FormatInt(kursID, 10)+"/qr.png") {
		t.Errorf("erwarte den QR-Code des Kurses auf der Wartefläche")
	}
	_ = code
}

// TestBeamerZeigtWederUeberarbeitungenNochSchrittnamenWaehrendDerWerkstatt
// hält die im Vorgang 0012 wörtlich geforderte Grenze fest: Überarbeitungen
// erscheinen hier nie, und der rohe Automaten-Schrittname (der wie ein
// öffentliches Zurückfallen aussähe) auch nicht — nur ein Fortschrittspunkt.
func TestBeamerZeigtWederUeberarbeitungenNochSchrittnamenWaehrendDerWerkstatt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID, gruppeID int64
	app.db.QueryRow(`SELECT kurs_id, id FROM gruppe`).Scan(&kursID, &gruppeID)

	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"nein"}})
	resp.Body.Close()

	beamerResp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beamer")
	if err != nil {
		t.Fatalf("beamer laden: %v", err)
	}
	defer beamerResp.Body.Close()
	body := koerperText(t, beamerResp)

	if strings.Contains(body, "Überarbeitung") {
		t.Errorf("erwarte KEIN Vorkommen von 'Überarbeitung' am Beamer, habe: %s", body)
	}
	if strings.Contains(body, schrittVerwerfung) {
		t.Errorf("erwarte keinen rohen Automaten-Schrittnamen am Beamer, habe: %s", body)
	}
	if !strings.Contains(body, "Gruppe "+strconv.Itoa(1)) {
		t.Errorf("erwarte die Gruppe im Fortschritt, habe: %s", body)
	}
}

func TestBeamerZeigtKeineGruppenOhneFreigegebeneTheseInDerSchlussansicht(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	_ = client

	setzePhase(t, app, kursID, phaseAbschluss)

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beamer")
	if err != nil {
		t.Fatalf("beamer laden: %v", err)
	}
	defer resp.Body.Close()
	body := koerperText(t, resp)

	if !strings.Contains(body, "Noch keine These freigegeben") {
		t.Errorf("erwarte einen Hinweis ohne freigegebene Thesen, habe: %s", body)
	}
}

func TestBeamerZeigtAlleFreigegebenenThesenNebeneinanderInDerSchlussansicht(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	bringeBisFreigegeben(t, app, client,
		"Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	)
	setzePhase(t, app, kursID, phaseAbschluss)

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beamer")
	if err != nil {
		t.Fatalf("beamer laden: %v", err)
	}
	defer resp.Body.Close()
	body := koerperText(t, resp)

	for _, teil := range []string{
		"Schöpfung und Klima",
		"wir Gottes Schöpfung anvertraut sind",
		"wir verantwortlich mit ihr umgehen",
		"wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	} {
		if !strings.Contains(body, teil) {
			t.Errorf("erwarte %q in der Schlussansicht", teil)
		}
	}
	if strings.Contains(body, "Überarbeitung") {
		t.Errorf("erwarte KEIN Vorkommen von 'Überarbeitung' in der Schlussansicht")
	}
}

// TestBeamerLiveMeldetAenderungAnIrgendeinerGruppeDesKurses deckt genau die
// Eigenheit von handleBeamerLive ab, die es von handleGruppeLive
// unterscheidet: eine einzelne SSE-Verbindung ohne Gruppen-Identität hört
// gleichzeitig auf ALLE Gruppen des Kurses.
func TestBeamerLiveMeldetAenderungAnIrgendeinerGruppeDesKurses(t *testing.T) {
	app := newTestApp(t)
	kursID, _, gruppenIDs := app.KursAnlegen(t, 2)

	v := sseVerbindenBeamer(t, app, kursID)
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte das anfängliche Sync-Signal")
	}

	// Eine Änderung an der ZWEITEN Gruppe muss ankommen, obwohl die
	// Verbindung keiner bestimmten Gruppe zugeordnet ist.
	werkstattHub.benachrichtigeGruppe(gruppenIDs[1])
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte ein Signal nach der Änderung an Gruppe 2")
	}
}

func TestBeamerLiveOhneKursLiefert404(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/kurse/999/beamer/live")
	if err != nil {
		t.Fatalf("GET /kurse/999/beamer/live: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("erwarte 404 für einen unbekannten Kurs, habe %d", resp.StatusCode)
	}
}

func TestUnbekannterKursLiefert404FuerBeamerAnzeigen(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/kurse/999/beamer")
	if err != nil {
		t.Fatalf("GET /kurse/999/beamer: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("erwarte 404 für einen unbekannten Kurs, habe %d", resp.StatusCode)
	}
}
