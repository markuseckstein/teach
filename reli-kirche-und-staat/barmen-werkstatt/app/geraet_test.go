package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"testing"
)

// clientMitCookieJar liefert einen http.Client, der Cookies wie ein
// Browser über mehrere Anfragen hinweg mitführt — für Tests des
// Beitritts- und Wiedereinstiegsflusses.
func clientMitCookieJar(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar: %v", err)
	}
	return &http.Client{Jar: jar}
}

func TestBeitrittMitGueltigemCodeFuehrtZurGruppenwahl(t *testing.T) {
	app := newTestApp(t)
	_, beitrittscode, _ := app.KursAnlegen(t, 3)

	resp, err := http.Get(app.server.URL + "/beitreten?code=" + beitrittscode)
	if err != nil {
		t.Fatalf("GET /beitreten: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	for _, gruppe := range []string{"Gruppe 1", "Gruppe 2", "Gruppe 3"} {
		if !strings.Contains(string(body), gruppe) {
			t.Errorf("erwarte %q in der Gruppenwahl", gruppe)
		}
	}
}

func TestBeitrittMitUnbekanntemCodeZeigtFehler(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/beitreten?code=XXXXXX")
	if err != nil {
		t.Fatalf("GET /beitreten: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 (Formular mit Fehlermeldung), habe %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "gibt es nicht") {
		t.Errorf("erwarte Fehlermeldung im Formular, habe: %s", body)
	}

	var anzahlGeraete int
	if err := app.db.QueryRow(`SELECT count(*) FROM geraet`).Scan(&anzahlGeraete); err != nil {
		t.Fatalf("geraet zählen: %v", err)
	}
	if anzahlGeraete != 0 {
		t.Errorf("erwarte kein angelegtes Gerät, habe %d", anzahlGeraete)
	}
}

func TestGeraetBeitrittSetztCookieUndLegtTheseAn(t *testing.T) {
	app := newTestApp(t)
	kursID, _, gruppenIDs := app.KursAnlegen(t, 2)

	client := clientMitCookieJar(t)
	resp, err := client.PostForm(
		app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/beitreten",
		nil,
	)
	if err != nil {
		t.Fatalf("POST beitreten: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var anzahlGeraete int
	if err := app.db.QueryRow(`SELECT count(*) FROM geraet WHERE gruppe_id = ?`, gruppenIDs[0]).Scan(&anzahlGeraete); err != nil {
		t.Fatalf("geraet zählen: %v", err)
	}
	if anzahlGeraete != 1 {
		t.Fatalf("erwarte 1 Gerät in Gruppe 1, habe %d", anzahlGeraete)
	}

	var anzahlThesen int
	if err := app.db.QueryRow(`SELECT count(*) FROM these WHERE gruppe_id = ?`, gruppenIDs[0]).Scan(&anzahlThesen); err != nil {
		t.Fatalf("these zählen: %v", err)
	}
	if anzahlThesen != 1 {
		t.Errorf("erwarte 1 angelegte These für Gruppe 1, habe %d", anzahlThesen)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Gruppe 1") {
		t.Errorf("erwarte Statusseite mit Gruppe 1, habe: %s", body)
	}
}

func TestGeraetLandetNachWiedereinstiegOhneErneutesBeitretenBeiSeinerGruppe(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	client := clientMitCookieJar(t)
	beitrittResp, err := client.PostForm(
		app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/beitreten",
		nil,
	)
	if err != nil {
		t.Fatalf("POST beitreten: %v", err)
	}
	beitrittResp.Body.Close()

	// Neue Anfrage mit demselben Cookie, ohne den Beitrittsweg erneut zu
	// gehen — simuliert das Tablet nach einem Verbindungsabriss.
	statusResp, err := client.Get(app.server.URL + "/gruppe")
	if err != nil {
		t.Fatalf("GET /gruppe: %v", err)
	}
	defer statusResp.Body.Close()
	if statusResp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", statusResp.StatusCode)
	}
	body, _ := io.ReadAll(statusResp.Body)
	if !strings.Contains(string(body), "Gruppe 1") || !strings.Contains(string(body), "BIBELSTELLE") {
		t.Errorf("erwarte Gruppe 1 im Schritt BIBELSTELLE, habe: %s", body)
	}
}

func TestGeraetStatusOhneCookieLeitetZumBeitrittWeiter(t *testing.T) {
	app := newTestApp(t)

	resp, err := http.Get(app.server.URL + "/gruppe")
	if err != nil {
		t.Fatalf("GET /gruppe: %v", err)
	}
	defer resp.Body.Close()

	if len(resp.Request.URL.Path) == 0 || resp.Request.URL.Path != "/beitreten" {
		t.Errorf("erwarte Weiterleitung nach /beitreten, gelandet auf %q", resp.Request.URL.Path)
	}
}

func TestUnbekanntesGeraeteTokenLeitetZumBeitrittWeiter(t *testing.T) {
	app := newTestApp(t)

	client := clientMitCookieJar(t)
	req, err := http.NewRequest(http.MethodGet, app.server.URL+"/gruppe", nil)
	if err != nil {
		t.Fatalf("Anfrage bauen: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: geraetCookieName, Value: "nie-vergebenes-token"})

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET /gruppe: %v", err)
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path != "/beitreten" {
		t.Errorf("erwarte Weiterleitung nach /beitreten, gelandet auf %q", resp.Request.URL.Path)
	}
}

func TestKursQRLiefertEinPNGBild(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/qr.png")
	if err != nil {
		t.Fatalf("GET qr.png: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200, habe %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Errorf("erwarte image/png, habe %q", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) < 8 || string(body[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Errorf("erwarte gültige PNG-Signatur")
	}
}
