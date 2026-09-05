package main

// Testnaht der Anwendung (Vorgang 0003): Jeder Test startet die echte
// Anwendung mit echtem Router gegen eine echte, temporäre SQLite-Datei und
// bedient sie über HTTP — wie ein Tablet oder das Regiepult es täten. Interne
// Funktionen sind nicht Gegenstand der Tests.
//
// "Kurs anlegen" und "Gerät beitreten lassen" schreiben hier direkt in die
// Datenbank, weil die zugehörigen HTTP-Endpunkte noch nicht existieren
// (Vorgänge 0004 und 0005). Sobald sie stehen, ist das die Stelle, an der
// diese Helfer auf echte HTTP-Anfragen umgestellt werden; ihre Signatur nach
// außen kann dabei gleich bleiben.

import (
	"database/sql"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"barmen-werkstatt/internal/db"
)

type testApp struct {
	server *httptest.Server
	db     *sql.DB
}

// newTestApp startet die vollständige Anwendung gegen eine frische SQLite-
// Datei in einem temporären Verzeichnis. Server und Datenbank werden am Ende
// des Tests automatisch geschlossen.
func newTestApp(t *testing.T) *testApp {
	t.Helper()

	pfad := filepath.Join(t.TempDir(), "werkstatt.db")
	database, err := db.Open(pfad)
	if err != nil {
		t.Fatalf("Datenbank öffnen: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	server := httptest.NewServer(newMux(database))
	t.Cleanup(server.Close)

	return &testApp{server: server, db: database}
}

// KursAnlegen legt einen Kurs mit zufälligem Beitrittscode und der
// angegebenen Zahl an Gruppen an und gibt Kurs-ID, Beitrittscode und die IDs
// der Gruppen (in Reihenfolge ihrer Nummer) zurück.
func (a *testApp) KursAnlegen(t *testing.T, gruppenzahl int) (kursID int64, beitrittscode string, gruppenIDs []int64) {
	t.Helper()

	beitrittscode = zufallscode(6)
	res, err := a.db.Exec(`INSERT INTO kurs (name, beitrittscode) VALUES (?, ?)`, "Testkurs", beitrittscode)
	if err != nil {
		t.Fatalf("Kurs anlegen: %v", err)
	}
	kursID, err = res.LastInsertId()
	if err != nil {
		t.Fatalf("Kurs-ID lesen: %v", err)
	}

	for nummer := 1; nummer <= gruppenzahl; nummer++ {
		res, err := a.db.Exec(`INSERT INTO gruppe (kurs_id, nummer) VALUES (?, ?)`, kursID, nummer)
		if err != nil {
			t.Fatalf("Gruppe %d anlegen: %v", nummer, err)
		}
		gruppeID, err := res.LastInsertId()
		if err != nil {
			t.Fatalf("Gruppe-ID lesen: %v", err)
		}
		gruppenIDs = append(gruppenIDs, gruppeID)
	}

	return kursID, beitrittscode, gruppenIDs
}

// GeraetBeitreten registriert ein neues Gerät bei der angegebenen Gruppe und
// gibt dessen Geräte-ID und Token zurück.
func (a *testApp) GeraetBeitreten(t *testing.T, gruppeID int64) (geraetID int64, token string) {
	t.Helper()

	token = zufallscode(16)
	res, err := a.db.Exec(`INSERT INTO geraet (gruppe_id, token) VALUES (?, ?)`, gruppeID, token)
	if err != nil {
		t.Fatalf("Gerät anlegen: %v", err)
	}
	geraetID, err = res.LastInsertId()
	if err != nil {
		t.Fatalf("Geräte-ID lesen: %v", err)
	}
	return geraetID, token
}

// geraetClient stellt Anfragen als ein bestimmtes Gerät, indem es bei jeder
// Anfrage den übergebenen Cookie mitführt.
type geraetClient struct {
	app    *testApp
	cookie *http.Cookie
}

// AlsGeraet liefert einen Client, der jede Anfrage mit cookie versieht — so,
// wie ein wiederkehrendes Tablet über sein Cookie erkannt wird.
func (a *testApp) AlsGeraet(cookie *http.Cookie) *geraetClient {
	return &geraetClient{app: a, cookie: cookie}
}

func (g *geraetClient) Get(t *testing.T, pfad string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, g.app.server.URL+pfad, nil)
	if err != nil {
		t.Fatalf("Anfrage GET %s bauen: %v", pfad, err)
	}
	return g.do(t, req)
}

func (g *geraetClient) PostForm(t *testing.T, pfad string, daten url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, g.app.server.URL+pfad, strings.NewReader(daten.Encode()))
	if err != nil {
		t.Fatalf("Anfrage POST %s bauen: %v", pfad, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return g.do(t, req)
}

func (g *geraetClient) do(t *testing.T, req *http.Request) *http.Response {
	t.Helper()
	req.AddCookie(g.cookie)
	resp, err := g.app.server.Client().Do(req)
	if err != nil {
		t.Fatalf("Anfrage an %s: %v", req.URL.Path, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

const zufallsalphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func zufallscode(laenge int) string {
	zeichen := make([]byte, laenge)
	for i := range zeichen {
		zeichen[i] = zufallsalphabet[rand.IntN(len(zufallsalphabet))]
	}
	return string(zeichen)
}
