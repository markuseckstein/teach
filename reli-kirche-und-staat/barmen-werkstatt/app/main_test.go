package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestAnwendungAntwortet(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("erwarte text/html, habe %q", ct)
	}
}

func TestStatischeDateienWerdenAusgeliefert(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/static/style.css")
	if err != nil {
		t.Fatalf("GET /static/style.css: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}
}

func TestKursUndGeraetHelferLegenGueltigeZeilenAn(t *testing.T) {
	app := newTestApp(t)

	kursID, beitrittscode, gruppenIDs := app.KursAnlegen(t, 3)
	if kursID == 0 {
		t.Fatalf("erwarte gültige Kurs-ID, habe 0")
	}
	if len(beitrittscode) == 0 {
		t.Fatalf("erwarte nicht-leeren Beitrittscode")
	}
	if len(gruppenIDs) != 3 {
		t.Fatalf("erwarte 3 Gruppen, habe %d", len(gruppenIDs))
	}

	_, token := app.GeraetBeitreten(t, gruppenIDs[0])
	if len(token) == 0 {
		t.Fatalf("erwarte nicht-leeres Geräte-Token")
	}

	client := app.AlsGeraet(&http.Cookie{Name: "geraet", Value: token})
	resp := client.Get(t, "/")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 für Anfrage mit Geräte-Cookie, habe %d", resp.StatusCode)
	}
}
