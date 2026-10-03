package main

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// sseVerbindung hält eine offene SSE-Verbindung für einen Test und liest
// im Hintergrund jede "data:"-Zeile in einen Kanal, damit der Test darauf
// warten (oder das Ausbleiben feststellen) kann.
type sseVerbindung struct {
	resp   *http.Response
	zeilen chan string
}

func sseVerbinden(t *testing.T, app *testApp, cookie *http.Cookie) *sseVerbindung {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, app.server.URL+"/gruppe/live", nil)
	if err != nil {
		t.Fatalf("sse-anfrage bauen: %v", err)
	}
	req.AddCookie(cookie)

	resp, err := app.server.Client().Do(req)
	if err != nil {
		t.Fatalf("sse verbinden: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		cancel()
		t.Fatalf("sse verbinden: erwarte 200, habe %d", resp.StatusCode)
	}

	zeilen := make(chan string, 32)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			zeile := scanner.Text()
			if strings.HasPrefix(zeile, "data:") {
				zeilen <- zeile
			}
		}
	}()

	v := &sseVerbindung{resp: resp, zeilen: zeilen}
	t.Cleanup(func() {
		cancel()
		resp.Body.Close()
	})
	return v
}

// warteAufEreignis meldet, ob innerhalb von timeout ein weiteres Ereignis
// ankam.
func (v *sseVerbindung) warteAufEreignis(timeout time.Duration) bool {
	select {
	case <-v.zeilen:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestSSEVerbindungBekommtSofortEinSyncSignal(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	cookie := geraetCookieAus(t, app, client)

	v := sseVerbinden(t, app, cookie)
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte sofort ein Sync-Signal nach Verbindungsaufbau")
	}
}

func TestSSESendetSignalBeiSchrittwechsel(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	cookie := geraetCookieAus(t, app, client)

	v := sseVerbinden(t, app, cookie)
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte das anfängliche Sync-Signal")
	}

	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	resp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("bibelstelle wählen: %v", err)
	}
	resp.Body.Close()

	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte ein Signal nach dem Schrittwechsel")
	}
}

func TestSSESendetSignalBeiSchreibrechtUebernahme(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	cookie := geraetCookieAus(t, app, client)

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	_, zweitesToken := app.GeraetBeitreten(t, gruppeID)
	zweitesGeraet := app.AlsGeraet(&http.Cookie{Name: geraetCookieName, Value: zweitesToken})

	v := sseVerbinden(t, app, cookie)
	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte das anfängliche Sync-Signal")
	}

	resp := zweitesGeraet.PostForm(t, "/gruppe/schreibrecht", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	if !v.warteAufEreignis(2 * time.Second) {
		t.Fatalf("erwarte ein Signal nach der Schreibrecht-Übernahme")
	}
}

// TestSSEIsoliertGruppenVoneinander stellt sicher, dass eine Änderung in
// einer Gruppe keine Geräte einer anderen Gruppe weckt — sonst würde der
// Live-Kanal die in SPEZIFIKATION.md geforderte Isolation zwischen Gruppen
// unterlaufen ("Eine Gruppe sieht die Thesen anderer Gruppen erst in der
// Schlussphase").
func TestSSEIsoliertGruppenVoneinander(t *testing.T) {
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 2)
	for _, gruppenNummer := range []int{1, 2} {
		resp, err := http.PostForm(
			app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/"+strconv.Itoa(gruppenNummer)+"/themenfeld",
			url.Values{"themenfeld": {[]string{"Schöpfung und Klima", "Umgang mit Minderheiten"}[gruppenNummer-1]}},
		)
		if err != nil {
			t.Fatalf("themenfeld zuweisen: %v", err)
		}
		resp.Body.Close()
	}

	clientGruppe1 := beitreten(t, app, kursID, 1)
	clientGruppe2 := beitreten(t, app, kursID, 2)

	vGruppe1 := sseVerbinden(t, app, geraetCookieAus(t, app, clientGruppe1))
	vGruppe2 := sseVerbinden(t, app, geraetCookieAus(t, app, clientGruppe2))
	if !vGruppe1.warteAufEreignis(2 * time.Second) {
		t.Fatalf("gruppe 1: erwarte das anfängliche Sync-Signal")
	}
	if !vGruppe2.warteAufEreignis(2 * time.Second) {
		t.Fatalf("gruppe 2: erwarte das anfängliche Sync-Signal")
	}

	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	resp, err := clientGruppe1.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("bibelstelle wählen (gruppe 1): %v", err)
	}
	resp.Body.Close()

	if !vGruppe1.warteAufEreignis(2 * time.Second) {
		t.Fatalf("gruppe 1: erwarte ein Signal nach eigenem Schrittwechsel")
	}
	if vGruppe2.warteAufEreignis(500 * time.Millisecond) {
		t.Errorf("gruppe 2: erwarte KEIN Signal durch einen Schrittwechsel in Gruppe 1")
	}
}

// TestDreissigGleichzeitigeSSEVerbindungen prüft, dass der Kanal unter der
// in Vorgang 0010 verlangten Last aus 30 gleichzeitigen Verbindungen
// funktioniert: fünf Geräte in jeder der sechs Gruppen eines Kurses, ein
// Schrittwechsel in einer Gruppe weckt nur deren fünf Geräte.
func TestDreissigGleichzeitigeSSEVerbindungen(t *testing.T) {
	app := newTestApp(t)
	kursID, _, gruppenIDs := app.KursAnlegen(t, 6)

	resp, err := http.PostForm(app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/themenfeld",
		url.Values{"themenfeld": {"Schöpfung und Klima"}})
	if err != nil {
		t.Fatalf("themenfeld zuweisen: %v", err)
	}
	resp.Body.Close()

	const geraeteJeGruppe = 5
	var verbindungen []*sseVerbindung
	var zielGruppenVerbindungen []*sseVerbindung
	var ersterClientGruppe1 *http.Client

	for gi, gruppeID := range gruppenIDs {
		// Das erste Gerät jeder Gruppe tritt über den echten HTTP-Weg bei —
		// das legt die zugehörige These an (sonst gäbe es später beim
		// Laden des Gruppenstatus keine Zeile) und hält automatisch das
		// Schreibrecht (Vorgang 0006).
		ersterClient := beitreten(t, app, kursID, gi+1)
		if gi == 0 {
			ersterClientGruppe1 = ersterClient
		}
		v := sseVerbinden(t, app, geraetCookieAus(t, app, ersterClient))
		verbindungen = append(verbindungen, v)
		if gi == 0 {
			zielGruppenVerbindungen = append(zielGruppenVerbindungen, v)
		}

		for gerät := 1; gerät < geraeteJeGruppe; gerät++ {
			_, token := app.GeraetBeitreten(t, gruppeID)
			cookie := &http.Cookie{Name: geraetCookieName, Value: token}
			v := sseVerbinden(t, app, cookie)
			verbindungen = append(verbindungen, v)
			if gi == 0 {
				zielGruppenVerbindungen = append(zielGruppenVerbindungen, v)
			}
		}
	}
	if len(verbindungen) != 30 {
		t.Fatalf("erwarte 30 Verbindungen, habe %d", len(verbindungen))
	}

	var wg sync.WaitGroup
	fehlgeschlagen := make([]bool, len(verbindungen))
	for i, v := range verbindungen {
		wg.Add(1)
		go func(i int, v *sseVerbindung) {
			defer wg.Done()
			if !v.warteAufEreignis(3 * time.Second) {
				fehlgeschlagen[i] = true
			}
		}(i, v)
	}
	wg.Wait()
	for i, f := range fehlgeschlagen {
		if f {
			t.Errorf("verbindung %d: erwarte anfängliches Sync-Signal", i)
		}
	}

	// Schrittwechsel in Gruppe 1: Das schreibberechtigte erste Gerät wählt
	// die Bibelstelle.
	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	auswahl, err := ersterClientGruppe1.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("bibelstelle wählen: %v", err)
	}
	auswahl.Body.Close()
	if auswahl.StatusCode != http.StatusOK {
		t.Fatalf("bibelstelle wählen: erwarte 200, habe %d", auswahl.StatusCode)
	}

	var empfangen sync.WaitGroup
	getroffen := make([]bool, len(zielGruppenVerbindungen))
	for i, v := range zielGruppenVerbindungen {
		empfangen.Add(1)
		go func(i int, v *sseVerbindung) {
			defer empfangen.Done()
			// Die Schreibrecht-Übernahme UND die Bibelstellenwahl lösen je
			// ein Signal aus — mindestens eines muss ankommen.
			getroffen[i] = v.warteAufEreignis(3 * time.Second)
		}(i, v)
	}
	empfangen.Wait()
	for i, ok := range getroffen {
		if !ok {
			t.Errorf("gruppe-1-verbindung %d: erwarte ein Signal nach dem Schrittwechsel", i)
		}
	}
}

// geraetCookieAus liest das Geräte-Cookie aus dem Jar eines Clients, der
// bereits über beitreten() einer Gruppe beigetreten ist.
func geraetCookieAus(t *testing.T, app *testApp, client *http.Client) *http.Cookie {
	t.Helper()
	u, err := url.Parse(app.server.URL)
	if err != nil {
		t.Fatalf("server-url parsen: %v", err)
	}
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == geraetCookieName {
			return c
		}
	}
	t.Fatalf("kein Geräte-Cookie im Jar gefunden")
	return nil
}
