package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mitTestLessonsDir legt ein paar Kopiervorlagen-Attrappen in einem
// t.TempDir() an und biegt *lessonsDirFlag für die Dauer des Tests dorthin
// um. Der reale Default "../../lessons" zeigt zwar (verifiziert) auch aus dem
// app-Verzeichnis auf das echte lessons/-Verzeichnis, ein eigenes,
// kontrolliertes Verzeichnis macht den Test aber unabhängig vom jeweiligen
// Stand des echten Unterrichtsmaterials.
func mitTestLessonsDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	dateien := map[string]string{
		"0001-kirche-und-staat.html":   "<title>In Verantwortung vor Gott</title>",
		"0002-ds1-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 1</title>",
		"0003-ds2-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 2</title>",
		"0004-ds3-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 3</title>",
		"0005-ds4-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 4</title>",
		"0006-ds5-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 5</title>",
		"0007-ds6-kopiervorlagen.html": "<title>Kopiervorlagen Doppelstunde 6</title>",
	}
	for name, inhalt := range dateien {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(inhalt), 0o644); err != nil {
			t.Fatalf("Testdatei %s anlegen: %v", name, err)
		}
	}

	alterWert := *lessonsDirFlag
	*lessonsDirFlag = dir
	t.Cleanup(func() { *lessonsDirFlag = alterWert })

	return dir
}

func TestLessonsDateiWirdUnterStabilerAdresseAusgeliefert(t *testing.T) {
	mitTestLessonsDir(t)
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/lessons/0007-ds6-kopiervorlagen.html")
	if err != nil {
		t.Fatalf("GET /lessons/...: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}
	if !strings.Contains(string(body), "Kopiervorlagen Doppelstunde 6") {
		t.Errorf("erwarte Titel der Doppelstunde 6 in der Antwort, habe %q", string(body))
	}
}

func TestInhaltsverzeichnisListetAlleDoppelstundenUndVerlinktBeitreten(t *testing.T) {
	mitTestLessonsDir(t)
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/inhalte")
	if err != nil {
		t.Fatalf("GET /inhalte: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}
	text := string(body)

	for _, datei := range []string{
		"0002-ds1-kopiervorlagen.html",
		"0003-ds2-kopiervorlagen.html",
		"0004-ds3-kopiervorlagen.html",
		"0005-ds4-kopiervorlagen.html",
		"0006-ds5-kopiervorlagen.html",
		"0007-ds6-kopiervorlagen.html",
	} {
		if !strings.Contains(text, "/lessons/"+datei) {
			t.Errorf("erwarte Link auf %q im Inhaltsverzeichnis", datei)
		}
	}

	if !strings.Contains(text, "/beitreten") {
		t.Errorf("erwarte Link auf /beitreten (Einstieg in die Anwendung für DS 6) im Inhaltsverzeichnis")
	}
	if !strings.Contains(text, "0001-kirche-und-staat.html") {
		t.Errorf("erwarte Link auf den Gesamtüberblick im Inhaltsverzeichnis")
	}
}

func TestLadeInhaltsverzeichnisOrdnetDoppelstundenzahlUeberDateinamenZu(t *testing.T) {
	dir := mitTestLessonsDir(t)

	ansicht, err := ladeInhaltsverzeichnis(dir)
	if err != nil {
		t.Fatalf("ladeInhaltsverzeichnis: %v", err)
	}
	if len(ansicht.Doppelstunden) != 6 {
		t.Fatalf("erwarte 6 Doppelstunden, habe %d", len(ansicht.Doppelstunden))
	}
	if ansicht.Gesamtueberblick != "0001-kirche-und-staat.html" {
		t.Errorf("erwarte Gesamtüberblick 0001-kirche-und-staat.html, habe %q", ansicht.Gesamtueberblick)
	}

	var ds6 *doppelstundeEintrag
	for i := range ansicht.Doppelstunden {
		if ansicht.Doppelstunden[i].Nummer == 6 {
			ds6 = &ansicht.Doppelstunden[i]
		}
	}
	if ds6 == nil {
		t.Fatal("erwarte einen Eintrag für Doppelstunde 6")
	}
	if ds6.Anwendung != "/beitreten" {
		t.Errorf("erwarte Anwendungs-Einstieg /beitreten für DS 6, habe %q", ds6.Anwendung)
	}

	for _, e := range ansicht.Doppelstunden {
		if e.Nummer != 6 && e.Anwendung != "" {
			t.Errorf("erwarte keinen Anwendungs-Einstieg für DS %d, habe %q", e.Nummer, e.Anwendung)
		}
	}
}
