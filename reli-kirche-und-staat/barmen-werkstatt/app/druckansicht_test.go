package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func bringeBisFreigegeben(t *testing.T, app *testApp, client *http.Client, themenfeld, weil, gilt, verwerfung string) {
	t.Helper()
	bringeBisVorschau(t, app, client, themenfeld, weil, gilt, verwerfung)

	resp, err := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	if err != nil {
		t.Fatalf("einreichen: %v", err)
	}
	resp.Body.Close()

	resp, err = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
	if err != nil {
		t.Fatalf("pruefung1: %v", err)
	}
	resp.Body.Close()

	resp, err = client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"ja"}})
	if err != nil {
		t.Fatalf("pruefung2: %v", err)
	}
	resp.Body.Close()
}

func TestFreigegebeneTheseZeigtAlleGefordertenAngabenFuerDenAusdruck(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisFreigegeben(t, app, client,
		"Schöpfung und Klima",
		"wir Gottes Schöpfung anvertraut sind",
		"wir verantwortlich mit ihr umgehen",
		"wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	)

	resp, err := client.Get(app.server.URL + "/gruppe/freigegeben")
	if err != nil {
		t.Fatalf("GET /gruppe/freigegeben: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("antwort lesen: %v", err)
	}
	seite := string(body)

	// Fertig-wenn (Vorgang 0013): Themenfeld, Bibelstelle, beide Teile der
	// These stehen auf dem Blatt.
	erwartet := []string{
		"Schöpfung und Klima",
		"wir Gottes Schöpfung anvertraut sind",
		"wir verantwortlich mit ihr umgehen",
		"wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
	}
	for _, teil := range erwartet {
		if !strings.Contains(seite, teil) {
			t.Errorf("erwarte %q auf der Druckseite", teil)
		}
	}

	// Print-CSS statt PDF-Erzeuger: kein Server-seitiges PDF, sondern
	// @media print im HTML.
	if !strings.Contains(seite, "@media print") {
		t.Errorf("erwarte Print-CSS auf der Seite")
	}
}

func TestNichtFreigegebeneTheseIstNichtDruckbar(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	// Absichtlich nicht bis FREIGEGEBEN gebracht — die Gruppe steht noch in
	// BIBELSTELLE.

	resp, err := client.Get(app.server.URL + "/gruppe/freigegeben")
	if err != nil {
		t.Fatalf("GET /gruppe/freigegeben: %v", err)
	}
	defer resp.Body.Close()

	// http.Client folgt der Weiterleitung; sie darf nicht auf der
	// Druckseite landen.
	if resp.Request.URL.Path == "/gruppe/freigegeben" {
		t.Errorf("erwarte, dass eine nicht freigegebene These nicht auf der Druckseite landet")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("antwort lesen: %v", err)
	}
	if strings.Contains(string(body), "@media print") {
		t.Errorf("erwarte keine Druckseite ohne Freigabe")
	}
}
