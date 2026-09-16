package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestSchemaWirdBeimErstenStartAngelegt(t *testing.T) {
	pfad := filepath.Join(t.TempDir(), "werkstatt.db")

	database, err := Open(pfad)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer database.Close()

	erwartet := []string{"kurs", "gruppe", "geraet", "these", "kapsel"}
	for _, tabelle := range erwartet {
		var name string
		row := database.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tabelle)
		if err := row.Scan(&name); err != nil {
			t.Errorf("Tabelle %q fehlt: %v", tabelle, err)
		}
	}
}

func TestZweiterStartVeraendertNichts(t *testing.T) {
	pfad := filepath.Join(t.TempDir(), "werkstatt.db")

	erste, err := Open(pfad)
	if err != nil {
		t.Fatalf("erster Open: %v", err)
	}
	if _, err := erste.Exec(`INSERT INTO kurs (name, beitrittscode) VALUES ('Testkurs', 'ABC123')`); err != nil {
		t.Fatalf("Testzeile einfügen: %v", err)
	}
	erste.Close()

	zweite, err := Open(pfad)
	if err != nil {
		t.Fatalf("zweiter Open: %v", err)
	}
	defer zweite.Close()

	var anzahl int
	if err := zweite.QueryRow(`SELECT count(*) FROM kurs`).Scan(&anzahl); err != nil {
		t.Fatalf("kurs zählen: %v", err)
	}
	if anzahl != 1 {
		t.Errorf("erwarte 1 Kurs nach zweitem Start, habe %d", anzahl)
	}
}

// TestFehlendeSpalteWirdBeimStartNachtraeglichErgaenzt bildet eine
// Datenbankdatei nach, die mit einer älteren Version der Anwendung entstand
// (kurs-Tabelle noch ohne werkstatt_gestartet_um). Open muss die fehlende
// Spalte ergänzen, statt bei jeder folgenden Abfrage mit "no such column"
// zu scheitern — und darf die bestehende Zeile dabei nicht anfassen.
func TestFehlendeSpalteWirdBeimStartNachtraeglichErgaenzt(t *testing.T) {
	pfad := filepath.Join(t.TempDir(), "werkstatt.db")

	altesSchema, err := sql.Open("sqlite", pfad)
	if err != nil {
		t.Fatalf("altes schema öffnen: %v", err)
	}
	if _, err := altesSchema.Exec(`CREATE TABLE kurs (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		name           TEXT NOT NULL,
		beitrittscode  TEXT NOT NULL UNIQUE,
		aktive_phase   TEXT NOT NULL DEFAULT '',
		beendet_am     TEXT
	)`); err != nil {
		t.Fatalf("altes schema anlegen: %v", err)
	}
	if _, err := altesSchema.Exec(
		`INSERT INTO kurs (name, beitrittscode) VALUES ('Bestandskurs', 'OLD123')`,
	); err != nil {
		t.Fatalf("bestandszeile einfügen: %v", err)
	}
	if err := altesSchema.Close(); err != nil {
		t.Fatalf("altes schema schließen: %v", err)
	}

	erste, err := Open(pfad)
	if err != nil {
		t.Fatalf("Open über eine Datei ohne die neue Spalte: %v", err)
	}

	var beitrittscode string
	if err := erste.QueryRow(
		`SELECT beitrittscode FROM kurs WHERE werkstatt_gestartet_um IS NULL`,
	).Scan(&beitrittscode); err != nil {
		t.Fatalf("spalte nach dem Ergänzen abfragen: %v", err)
	}
	if beitrittscode != "OLD123" {
		t.Errorf("erwarte die bestehende Zeile unverändert, habe beitrittscode=%q", beitrittscode)
	}
	erste.Close()

	// Zweiter Open (die Spalte existiert jetzt bereits) darf nicht erneut
	// versuchen, sie anzulegen, und nicht scheitern.
	zweite, err := Open(pfad)
	if err != nil {
		t.Fatalf("zweiter Open nach dem Ergänzen: %v", err)
	}
	defer zweite.Close()
}

func TestLoeschenEinesKursesLoeschtAllesDaranHaengende(t *testing.T) {
	pfad := filepath.Join(t.TempDir(), "werkstatt.db")

	database, err := Open(pfad)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer database.Close()

	res, err := database.Exec(`INSERT INTO kurs (name, beitrittscode) VALUES ('Testkurs', 'XYZ789')`)
	if err != nil {
		t.Fatalf("kurs einfügen: %v", err)
	}
	kursID, _ := res.LastInsertId()

	res, err = database.Exec(`INSERT INTO gruppe (kurs_id, nummer) VALUES (?, 1)`, kursID)
	if err != nil {
		t.Fatalf("gruppe einfügen: %v", err)
	}
	gruppeID, _ := res.LastInsertId()

	if _, err := database.Exec(`INSERT INTO geraet (gruppe_id, token) VALUES (?, 'tok-1')`, gruppeID); err != nil {
		t.Fatalf("geraet einfügen: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO these (gruppe_id) VALUES (?)`, gruppeID); err != nil {
		t.Fatalf("these einfügen: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO kapsel (kurs_id, art) VALUES (?, 'placemat')`, kursID); err != nil {
		t.Fatalf("kapsel einfügen: %v", err)
	}

	if _, err := database.Exec(`DELETE FROM kurs WHERE id = ?`, kursID); err != nil {
		t.Fatalf("kurs löschen: %v", err)
	}

	for tabelle, spalte := range map[string]string{
		"gruppe": "kurs_id",
		"kapsel": "kurs_id",
	} {
		var anzahl int
		if err := database.QueryRow(`SELECT count(*) FROM `+tabelle+` WHERE `+spalte+` = ?`, kursID).Scan(&anzahl); err != nil {
			t.Fatalf("%s zählen: %v", tabelle, err)
		}
		if anzahl != 0 {
			t.Errorf("erwarte 0 Zeilen in %s nach Kurslöschung, habe %d", tabelle, anzahl)
		}
	}

	var anzahlGeraet, anzahlThese int
	if err := database.QueryRow(`SELECT count(*) FROM geraet WHERE gruppe_id = ?`, gruppeID).Scan(&anzahlGeraet); err != nil {
		t.Fatalf("geraet zählen: %v", err)
	}
	if err := database.QueryRow(`SELECT count(*) FROM these WHERE gruppe_id = ?`, gruppeID).Scan(&anzahlThese); err != nil {
		t.Fatalf("these zählen: %v", err)
	}
	if anzahlGeraet != 0 || anzahlThese != 0 {
		t.Errorf("erwarte 0 Zeilen in geraet/these nach Kaskade, habe geraet=%d these=%d", anzahlGeraet, anzahlThese)
	}
}
