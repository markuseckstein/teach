package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"barmen-werkstatt/internal/db"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	pfad := filepath.Join(t.TempDir(), "werkstatt.db")
	database, err := db.Open(pfad)
	if err != nil {
		t.Fatalf("Datenbank öffnen: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	server := httptest.NewServer(newMux(database))
	t.Cleanup(server.Close)
	return server
}

func TestStartseiteZeigtDassDieAnwendungLebt(t *testing.T) {
	server := testServer(t)

	resp, err := http.Get(server.URL + "/")
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
	server := testServer(t)

	resp, err := http.Get(server.URL + "/static/style.css")
	if err != nil {
		t.Fatalf("GET /static/style.css: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}
}
