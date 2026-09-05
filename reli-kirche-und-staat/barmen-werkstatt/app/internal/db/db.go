// Package db kapselt den Zugriff auf die SQLite-Datei der Barmen-Werkstatt.
package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// schema enthält das vollständige Datenmodell aus SPEZIFIKATION.md. Alle
// Anweisungen sind "CREATE TABLE IF NOT EXISTS", damit ein zweiter Start
// nichts verändert.
const schema = `
CREATE TABLE IF NOT EXISTS kurs (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	name           TEXT NOT NULL,
	beitrittscode  TEXT NOT NULL UNIQUE,
	aktive_phase   TEXT NOT NULL DEFAULT '',
	beendet_am     TEXT
);

CREATE TABLE IF NOT EXISTS gruppe (
	id                    INTEGER PRIMARY KEY AUTOINCREMENT,
	kurs_id               INTEGER NOT NULL REFERENCES kurs(id) ON DELETE CASCADE,
	nummer                INTEGER NOT NULL,
	themenfeld             TEXT NOT NULL DEFAULT '',
	schreibrecht_geraet   INTEGER REFERENCES geraet(id) ON DELETE SET NULL,
	schreibrecht_seit     TEXT,
	UNIQUE (kurs_id, nummer)
);

CREATE TABLE IF NOT EXISTS geraet (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	gruppe_id        INTEGER NOT NULL REFERENCES gruppe(id) ON DELETE CASCADE,
	token            TEXT NOT NULL UNIQUE,
	zuletzt_gesehen  TEXT
);

CREATE TABLE IF NOT EXISTS these (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	gruppe_id         INTEGER NOT NULL REFERENCES gruppe(id) ON DELETE CASCADE,
	bibelstelle       TEXT NOT NULL DEFAULT '',
	weil              TEXT NOT NULL DEFAULT '',
	gilt              TEXT NOT NULL DEFAULT '',
	verwerfung        TEXT NOT NULL DEFAULT '',
	schritt           TEXT NOT NULL DEFAULT 'BIBELSTELLE',
	ueberarbeitungen  INTEGER NOT NULL DEFAULT 0,
	unterschriften    TEXT NOT NULL DEFAULT '',
	freigegeben_am    TEXT
);

CREATE TABLE IF NOT EXISTS kapsel (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	kurs_id       INTEGER NOT NULL REFERENCES kurs(id) ON DELETE CASCADE,
	art           TEXT NOT NULL,
	beschriftung  TEXT NOT NULL DEFAULT '',
	bild          BLOB,
	erstellt_am   TEXT
);
`

// Open öffnet die SQLite-Datei unter path und legt das Schema an, falls die
// Datei fehlt oder Tabellen fehlen. Fremdschlüssel sind für jede Verbindung
// aktiviert, damit das Löschen eines Kurses alles daran Hängende mitnimmt.
func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", path)
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("db öffnen: %w", err)
	}

	// SQLite ist einschreiber-beschränkt; eine einzige Verbindung vermeidet
	// "database is locked" unter dem Formular-POST-Muster der Anwendung.
	database.SetMaxOpenConns(1)

	if _, err := database.Exec(schema); err != nil {
		database.Close()
		return nil, fmt.Errorf("schema anlegen: %w", err)
	}

	return database, nil
}
