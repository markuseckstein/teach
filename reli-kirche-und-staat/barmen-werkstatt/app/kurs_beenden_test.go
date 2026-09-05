package main

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestKursBeendenLoeschtAlleZeilen(t *testing.T) {
	app := newTestApp(t)
	kursID, _, gruppenIDs := app.KursAnlegen(t, 2)

	for _, gruppeID := range gruppenIDs {
		app.GeraetBeitreten(t, gruppeID)
		if _, err := app.db.Exec(
			`INSERT INTO these (gruppe_id, bibelstelle) VALUES (?, ?)`, gruppeID, "Röm 13",
		); err != nil {
			t.Fatalf("these anlegen: %v", err)
		}
	}
	// Für Zeitkapsel-Einträge gibt es noch keinen HTTP-Endpunkt (kommt erst mit
	// Vorgang 0019/0020) — direkt per SQL angelegt.
	if _, err := app.db.Exec(
		`INSERT INTO kapsel (kurs_id, art, beschriftung) VALUES (?, ?, ?)`, kursID, "foto", "Wäscheleine DS 4",
	); err != nil {
		t.Fatalf("kapsel anlegen: %v", err)
	}

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/beenden", nil)
	if err != nil {
		t.Fatalf("POST beenden: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	for tabelle, bedingung := range map[string]string{
		"kurs":   "id = " + strconv.FormatInt(kursID, 10),
		"gruppe": "kurs_id = " + strconv.FormatInt(kursID, 10),
		"geraet": "gruppe_id IN (SELECT id FROM gruppe WHERE kurs_id = " + strconv.FormatInt(kursID, 10) + ")",
		"these":  "gruppe_id IN (SELECT id FROM gruppe WHERE kurs_id = " + strconv.FormatInt(kursID, 10) + ")",
		"kapsel": "kurs_id = " + strconv.FormatInt(kursID, 10),
	} {
		var anzahl int
		if err := app.db.QueryRow(`SELECT count(*) FROM ` + tabelle + ` WHERE ` + bedingung).Scan(&anzahl); err != nil {
			t.Fatalf("%s zählen: %v", tabelle, err)
		}
		if anzahl != 0 {
			t.Errorf("erwarte 0 Zeilen in %s nach Kurs beenden, habe %d", tabelle, anzahl)
		}
	}
}

func TestKursBeendenBenenntWasVerlorenGeht(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/beenden")
	if err != nil {
		t.Fatalf("GET beenden: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}

	var body strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		body.Write(buf[:n])
		if err != nil {
			break
		}
	}
	text := body.String()
	for _, erwartet := range []string{"Kurs", "Gruppe", "Geräte", "These", "Zeitkapsel"} {
		if !strings.Contains(text, erwartet) {
			t.Errorf("erwarte %q auf der Bestätigungsseite", erwartet)
		}
	}
}

func TestUnbekannterKursBeimBeendenLiefert404UndLoeschtNichts(t *testing.T) {
	app := newTestApp(t)
	app.KursAnlegen(t, 1) // ein realer Kurs, der unberührt bleiben soll

	resp, err := http.PostForm(app.server.URL+"/kurse/999999/beenden", nil)
	if err != nil {
		t.Fatalf("POST beenden: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("erwarte 404, habe %d", resp.StatusCode)
	}

	var anzahlKurse int
	if err := app.db.QueryRow(`SELECT count(*) FROM kurs`).Scan(&anzahlKurse); err != nil {
		t.Fatalf("kurs zählen: %v", err)
	}
	if anzahlKurse != 1 {
		t.Errorf("erwarte weiterhin 1 Kurs nach 404, habe %d", anzahlKurse)
	}

	respGet, err := http.Get(app.server.URL + "/kurse/999999/beenden")
	if err != nil {
		t.Fatalf("GET beenden: %v", err)
	}
	defer respGet.Body.Close()
	if respGet.StatusCode != http.StatusNotFound {
		t.Fatalf("erwarte 404 für GET auf unbekannten Kurs, habe %d", respGet.StatusCode)
	}
}
