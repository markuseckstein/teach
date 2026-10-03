package main

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setzePhase(t *testing.T, app *testApp, kursID int64, phase string) {
	t.Helper()
	if _, err := app.db.Exec(`UPDATE kurs SET aktive_phase = ? WHERE id = ?`, phase, kursID); err != nil {
		t.Fatalf("phase setzen: %v", err)
	}
}

func TestWerkstattIstVorFreischaltungGesperrt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	setzePhase(t, app, kursID, phaseVorbereitung)

	getResp, err := client.Get(app.server.URL + "/gruppe/bibelstelle")
	if err != nil {
		t.Fatalf("GET /gruppe/bibelstelle: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.Request.URL.Path != "/gruppe" {
		t.Errorf("erwarte Weiterleitung zu /gruppe vor Freischaltung, gelandet auf %q", getResp.Request.URL.Path)
	}

	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	postResp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("POST /gruppe/bibelstelle: %v", err)
	}
	defer postResp.Body.Close()
	if postResp.StatusCode != http.StatusConflict {
		t.Errorf("erwarte 409 vor Freischaltung, habe %d", postResp.StatusCode)
	}
}

func TestPhaseFreischaltenErlaubtDieWerkstatt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	setzePhase(t, app, kursID, phaseVorbereitung)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/phase", url.Values{"phase": {"werkstatt"}})
	if err != nil {
		t.Fatalf("phase setzen: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	postResp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("POST /gruppe/bibelstelle: %v", err)
	}
	defer postResp.Body.Close()
	if postResp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Freischaltung, habe %d", postResp.StatusCode)
	}
}

func TestPhaseWechselWirdUeberSSEAnDieGruppeGemeldet(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)
	setzePhase(t, app, kursID, phaseVorbereitung)

	v := sseVerbinden(t, app, geraetCookieAus(t, app, client))
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte das anfängliche Sync-Signal")
	}

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/phase", url.Values{"phase": {"werkstatt"}})
	if err != nil {
		t.Fatalf("phase setzen: %v", err)
	}
	resp.Body.Close()

	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte ein Signal nach dem Phasenwechsel (kursweite Benachrichtigung)")
	}
}

func TestTimerStartenErscheintImRegiepult(t *testing.T) {
	app, _ := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/timer/starten", nil)
	if err != nil {
		t.Fatalf("timer starten: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	body := koerperText(t, resp)
	if !strings.Contains(body, "Gestartet um") {
		t.Errorf("erwarte Timer-Anzeige im Regiepult, habe: %s", body)
	}
}

func TestAbgelaufenerTimerWirdAlsSolcherAngezeigt(t *testing.T) {
	app, _ := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)

	laengstVorbei := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	if _, err := app.db.Exec(`UPDATE kurs SET werkstatt_gestartet_um = ? WHERE id = ?`, laengstVorbei, kursID); err != nil {
		t.Fatalf("timer setzen: %v", err)
	}

	resp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/regiepult")
	if err != nil {
		t.Fatalf("regiepult laden: %v", err)
	}
	defer resp.Body.Close()
	body := koerperText(t, resp)
	if !strings.Contains(body, "Zeit ist um") {
		t.Errorf("erwarte Hinweis auf abgelaufene Zeit, habe: %s", body)
	}
}

func TestManuelleWeiterschaltungSpringtEinenSchrittOhneUeberarbeitungZuZaehlen(t *testing.T) {
	app, _ := werkstattApp(t, "Schöpfung und Klima")
	var kursID, gruppeID int64
	app.db.QueryRow(`SELECT kurs_id, id FROM gruppe`).Scan(&kursID, &gruppeID)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/weiterschalten", nil)
	if err != nil {
		t.Fatalf("manuell weiterschalten: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittPositiv {
		t.Errorf("erwarte POSITIV nach Notfall-Weiterschaltung, habe %s", these.Schritt)
	}
	if these.Ueberarbeitungen != 0 {
		t.Errorf("erwarte, dass die Notfall-Weiterschaltung nicht als Überarbeitung zählt, habe %d", these.Ueberarbeitungen)
	}
}

func TestManuelleWeiterschaltungNachFreigabeLiefertKonflikt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID, gruppeID int64
	app.db.QueryRow(`SELECT kurs_id, id FROM gruppe`).Scan(&kursID, &gruppeID)
	bringeBisFreigegeben(t, app, client,
		"Schöpfung und Klima", "Weil-Text", "Gilt-Text", "Verwerfungs-Text",
	)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/weiterschalten", nil)
	if err != nil {
		t.Fatalf("manuell weiterschalten: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("erwarte 409 nach Freigabe, habe %d", resp.StatusCode)
	}
}

// TestManuelleWeiterschaltungOhneBeigetretenesGeraetLiefert404 deckt eine
// Gruppe ab, der noch nie ein Gerät beigetreten ist — es gibt dann noch
// keine These-Zeile, und der Notfall-Weiterschalter darf nicht mit einem
// internen Serverfehler abbrechen.
func TestManuelleWeiterschaltungOhneBeigetretenesGeraetLiefert404(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/weiterschalten", nil)
	if err != nil {
		t.Fatalf("manuell weiterschalten: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("erwarte 404 ohne beigetretenes Gerät, habe %d", resp.StatusCode)
	}
}

// TestUeberarbeitungenNurImRegiepultSichtbar hält die in Vorgang 0011 wörtlich
// geforderte Sichtbarkeitsgrenze als Regressionstest fest: "Die Zahl der
// Überarbeitungen ist hier sichtbar — und nur hier."
func TestUeberarbeitungenNurImRegiepultSichtbar(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID int64
	app.db.QueryRow(`SELECT kurs_id FROM gruppe`).Scan(&kursID)

	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "W", "G", "V")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"nein"}})
	resp.Body.Close()

	regiepultResp, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/regiepult")
	if err != nil {
		t.Fatalf("regiepult laden: %v", err)
	}
	defer regiepultResp.Body.Close()
	regiepultBody := koerperText(t, regiepultResp)
	if !strings.Contains(regiepultBody, "Überarbeitungen") {
		t.Errorf("erwarte 'Überarbeitungen' im Regiepult")
	}

	statusResp, err := client.Get(app.server.URL + "/gruppe")
	if err != nil {
		t.Fatalf("gerätestatus laden: %v", err)
	}
	defer statusResp.Body.Close()
	statusBody := koerperText(t, statusResp)
	if strings.Contains(statusBody, "Überarbeitung") {
		t.Errorf("erwarte KEIN Vorkommen von 'Überarbeitung' auf der Geräteseite, habe: %s", statusBody)
	}
}

func TestSchreibrechtWarnungErscheintErstNachLangerZeit(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	var kursID, gruppeID int64
	app.db.QueryRow(`SELECT kurs_id, id FROM gruppe`).Scan(&kursID, &gruppeID)
	_ = client

	frisch, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/regiepult")
	if err != nil {
		t.Fatalf("regiepult laden: %v", err)
	}
	defer frisch.Body.Close()
	if strings.Contains(koerperText(t, frisch), "schreibt seit Langem") {
		t.Errorf("erwarte keine Warnung direkt nach dem Beitritt")
	}

	laengstHer := time.Now().Add(-30 * time.Minute).UTC().Format(time.RFC3339)
	if _, err := app.db.Exec(`UPDATE gruppe SET schreibrecht_seit = ? WHERE id = ?`, laengstHer, gruppeID); err != nil {
		t.Fatalf("schreibrecht_seit setzen: %v", err)
	}

	spaeter, err := http.Get(app.server.URL + "/kurse/" + strconv.FormatInt(kursID, 10) + "/regiepult")
	if err != nil {
		t.Fatalf("regiepult laden: %v", err)
	}
	defer spaeter.Body.Close()
	if !strings.Contains(koerperText(t, spaeter), "schreibt seit Langem") {
		t.Errorf("erwarte eine Warnung, wenn dasselbe Gerät seit sehr langer Zeit schreibt")
	}
}

func TestUnbekannterKursLiefert404FuerRegiepultRouten(t *testing.T) {
	app := newTestApp(t)

	for _, pfad := range []string{
		"/kurse/999/regiepult",
	} {
		resp, err := http.Get(app.server.URL + pfad)
		if err != nil {
			t.Fatalf("GET %s: %v", pfad, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: erwarte 404, habe %d", pfad, resp.StatusCode)
		}
	}

	for _, pfad := range []string{
		"/kurse/999/phase",
		"/kurse/999/timer/starten",
		"/kurse/999/gruppen/1/weiterschalten",
	} {
		resp, err := http.PostForm(app.server.URL+pfad, url.Values{"phase": {"werkstatt"}})
		if err != nil {
			t.Fatalf("POST %s: %v", pfad, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("POST %s: erwarte 404, habe %d", pfad, resp.StatusCode)
		}
	}
}

func koerperText(t *testing.T, resp *http.Response) string {
	t.Helper()
	daten, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("antwort lesen: %v", err)
	}
	return string(daten)
}
