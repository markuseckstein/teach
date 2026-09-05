package main

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestKursWirdMitBeitrittscodeUndGruppenAngelegt(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.PostForm(app.server.URL+"/kurse", url.Values{
		"name":        {"9c"},
		"gruppenzahl": {"3"},
	})
	if err != nil {
		t.Fatalf("POST /kurse: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}

	var anzahlKurse, anzahlGruppen int
	if err := app.db.QueryRow(`SELECT count(*) FROM kurs`).Scan(&anzahlKurse); err != nil {
		t.Fatalf("kurs zählen: %v", err)
	}
	if anzahlKurse != 1 {
		t.Fatalf("erwarte 1 Kurs, habe %d", anzahlKurse)
	}

	if err := app.db.QueryRow(`SELECT count(*) FROM gruppe`).Scan(&anzahlGruppen); err != nil {
		t.Fatalf("gruppe zählen: %v", err)
	}
	if anzahlGruppen != 3 {
		t.Fatalf("erwarte 3 Gruppen, habe %d", anzahlGruppen)
	}

	var beitrittscode string
	if err := app.db.QueryRow(`SELECT beitrittscode FROM kurs`).Scan(&beitrittscode); err != nil {
		t.Fatalf("beitrittscode lesen: %v", err)
	}
	if len(beitrittscode) != 6 {
		t.Errorf("erwarte 6-stelligen Beitrittscode, habe %q", beitrittscode)
	}
	if !strings.Contains(string(body), beitrittscode) {
		t.Errorf("erwarte Beitrittscode %q in der Antwortseite", beitrittscode)
	}
}

func TestGruppenzahlAusserhalbDesGueltigenBereichsWirdAbgelehnt(t *testing.T) {
	app := newTestApp(t)

	for _, gruppenzahl := range []string{"0", "7", "abc"} {
		resp, err := http.PostForm(app.server.URL+"/kurse", url.Values{
			"name":        {"9c"},
			"gruppenzahl": {gruppenzahl},
		})
		if err != nil {
			t.Fatalf("POST /kurse (gruppenzahl=%s): %v", gruppenzahl, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("gruppenzahl=%s: erwarte 400, habe %d", gruppenzahl, resp.StatusCode)
		}
	}

	var anzahlKurse int
	if err := app.db.QueryRow(`SELECT count(*) FROM kurs`).Scan(&anzahlKurse); err != nil {
		t.Fatalf("kurs zählen: %v", err)
	}
	if anzahlKurse != 0 {
		t.Errorf("erwarte 0 Kurse nach abgelehnten Anfragen, habe %d", anzahlKurse)
	}
}

func TestThemenfeldWirdEinerGruppeZugewiesen(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 2)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/themenfeld", url.Values{
		"themenfeld": {"Schöpfung und Klima"},
	})
	if err != nil {
		t.Fatalf("POST themenfeld: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var themenfeld string
	if err := app.db.QueryRow(`SELECT themenfeld FROM gruppe WHERE kurs_id = ? AND nummer = 1`, kursID).Scan(&themenfeld); err != nil {
		t.Fatalf("themenfeld lesen: %v", err)
	}
	if themenfeld != "Schöpfung und Klima" {
		t.Errorf("erwarte zugewiesenes Themenfeld, habe %q", themenfeld)
	}
}

func TestThemenfeldKannImSelbenKursNichtZweimalVergebenWerden(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 2)

	zuweisen := func(gruppenNummer int) *http.Response {
		resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/"+strconv.Itoa(gruppenNummer)+"/themenfeld", url.Values{
			"themenfeld": {"Schöpfung und Klima"},
		})
		if err != nil {
			t.Fatalf("POST themenfeld (gruppe %d): %v", gruppenNummer, err)
		}
		return resp
	}

	erste := zuweisen(1)
	erste.Body.Close()
	if erste.StatusCode != http.StatusOK {
		t.Fatalf("erste Zuweisung: erwarte 200, habe %d", erste.StatusCode)
	}

	zweite := zuweisen(2)
	defer zweite.Body.Close()
	if zweite.StatusCode != http.StatusConflict {
		t.Fatalf("doppelte Zuweisung: erwarte 409, habe %d", zweite.StatusCode)
	}

	var themenfeld string
	if err := app.db.QueryRow(`SELECT themenfeld FROM gruppe WHERE kurs_id = ? AND nummer = 2`, kursID).Scan(&themenfeld); err != nil {
		t.Fatalf("themenfeld lesen: %v", err)
	}
	if themenfeld != "" {
		t.Errorf("erwarte kein zugewiesenes Themenfeld für Gruppe 2, habe %q", themenfeld)
	}
}

func TestThemenfeldWirdBeiGleichzeitigerZuweisungNurEinmalVergeben(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 2)

	zuweisen := func(gruppenNummer int) (*http.Response, error) {
		return http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/"+strconv.Itoa(gruppenNummer)+"/themenfeld", url.Values{
			"themenfeld": {"Schöpfung und Klima"},
		})
	}

	var wg sync.WaitGroup
	antworten := make([]*http.Response, 2)
	fehler := make([]error, 2)
	for i, gruppenNummer := range []int{1, 2} {
		wg.Add(1)
		go func(i, gruppenNummer int) {
			defer wg.Done()
			antworten[i], fehler[i] = zuweisen(gruppenNummer)
		}(i, gruppenNummer)
	}
	wg.Wait()

	erfolge := 0
	for i, resp := range antworten {
		if fehler[i] != nil {
			t.Fatalf("POST themenfeld: %v", fehler[i])
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			erfolge++
		} else if resp.StatusCode != http.StatusConflict {
			t.Errorf("erwarte 200 oder 409, habe %d", resp.StatusCode)
		}
	}
	if erfolge != 1 {
		t.Errorf("erwarte genau eine erfolgreiche Zuweisung bei gleichzeitiger Anfrage, habe %d", erfolge)
	}

	var anzahlZugewiesen int
	if err := app.db.QueryRow(
		`SELECT count(*) FROM gruppe WHERE kurs_id = ? AND themenfeld = 'Schöpfung und Klima'`, kursID,
	).Scan(&anzahlZugewiesen); err != nil {
		t.Fatalf("zuweisungen zählen: %v", err)
	}
	if anzahlZugewiesen != 1 {
		t.Errorf("erwarte genau eine Gruppe mit diesem Themenfeld, habe %d", anzahlZugewiesen)
	}
}

func TestUnbekanntesThemenfeldWirdAbgelehnt(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/themenfeld", url.Values{
		"themenfeld": {"Erfundenes Themenfeld"},
	})
	if err != nil {
		t.Fatalf("POST themenfeld: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("erwarte 400, habe %d", resp.StatusCode)
	}
}
