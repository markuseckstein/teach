package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// werkstattApp legt einen Kurs mit einer Gruppe an, weist ihr ein
// Themenfeld zu und lässt ein Gerät beitreten — der gemeinsame
// Ausgangspunkt für alle Tests des Zustandsautomaten.
func werkstattApp(t *testing.T, themenfeld string) (*testApp, *http.Client) {
	t.Helper()
	app := newTestApp(t)
	kursID, _, _ := app.KursAnlegen(t, 1)

	resp, err := http.PostForm(
		app.server.URL+"/kurse/"+strconv.FormatInt(kursID, 10)+"/gruppen/1/themenfeld",
		url.Values{"themenfeld": {themenfeld}},
	)
	if err != nil {
		t.Fatalf("themenfeld zuweisen: %v", err)
	}
	resp.Body.Close()

	client := beitreten(t, app, kursID, 1)
	return app, client
}

func ersteBibelstelle(t *testing.T, themenfeld string) string {
	t.Helper()
	vorrat, ok := bibelstellenVorrat(themenfeld)
	if !ok || len(vorrat) == 0 {
		t.Fatalf("kein Bibelstellen-Vorrat für Themenfeld %q", themenfeld)
	}
	return vorrat[0].Angabe
}

// bringeBisVorschau führt eine Gruppe von BIBELSTELLE bis VORSCHAU (nicht
// eingereicht), mit den übergebenen Texten.
func bringeBisVorschau(t *testing.T, app *testApp, client *http.Client, themenfeld, weil, gilt, verwerfung string) {
	t.Helper()

	angabe := ersteBibelstelle(t, themenfeld)
	resp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("bibelstelle wählen: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bibelstelle wählen: erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	resp, err = client.PostForm(app.server.URL+"/gruppe/positiv", url.Values{"weil": {weil}, "gilt": {gilt}})
	if err != nil {
		t.Fatalf("positiv speichern: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("positiv speichern: erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	resp, err = client.PostForm(app.server.URL+"/gruppe/verwerfung", url.Values{"verwerfung": {verwerfung}})
	if err != nil {
		t.Fatalf("verwerfung speichern: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("verwerfung speichern: erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}
}

func aktuellerSchritt(t *testing.T, app *testApp, gruppeID int64) theseZeile {
	t.Helper()
	these, err := ladeThese(app.db, gruppeID)
	if err != nil {
		t.Fatalf("these laden: %v", err)
	}
	return these
}

func TestOhneGewaehlteBibelstelleKeinWeiterkommenAuchNichtDirekt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")

	resp, err := client.PostForm(app.server.URL+"/gruppe/positiv", url.Values{"weil": {"X"}, "gilt": {"Y"}})
	if err != nil {
		t.Fatalf("direkter POST /gruppe/positiv: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("erwarte 409, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	if err := app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID); err != nil {
		t.Fatalf("gruppe-id lesen: %v", err)
	}
	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittBibelstelle {
		t.Errorf("erwarte weiterhin BIBELSTELLE, habe %s", these.Schritt)
	}
	if these.Weil != "" || these.Gilt != "" {
		t.Errorf("erwarte keine gespeicherten Felder, habe weil=%q gilt=%q", these.Weil, these.Gilt)
	}
}

func TestBibelstelleWaehlenSchaltetZuPositivWeiter(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	angabe := ersteBibelstelle(t, "Schöpfung und Klima")

	resp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if err != nil {
		t.Fatalf("bibelstelle wählen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	if err := app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID); err != nil {
		t.Fatalf("gruppe-id lesen: %v", err)
	}
	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittPositiv {
		t.Errorf("erwarte POSITIV, habe %s", these.Schritt)
	}
	if these.Bibelstelle != angabe {
		t.Errorf("erwarte gespeicherte Bibelstelle %q, habe %q", angabe, these.Bibelstelle)
	}
}

func TestErfundeneBibelstelleWirdAbgelehnt(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")

	resp, err := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {"Erfundene Stelle 1,1"}})
	if err != nil {
		t.Fatalf("bibelstelle wählen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("erwarte 400, habe %d", resp.StatusCode)
	}
}

func TestPositivBenoetigtBeideFelder(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	resp, _ := client.PostForm(app.server.URL+"/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	resp.Body.Close()

	resp, err := client.PostForm(app.server.URL+"/gruppe/positiv", url.Values{"weil": {"nur weil, kein gilt"}, "gilt": {""}})
	if err != nil {
		t.Fatalf("positiv mit leerem gilt: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("erwarte 400 bei leerem Feld, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittPositiv {
		t.Errorf("erwarte weiterhin POSITIV, habe %s", these.Schritt)
	}
}

func TestNeinBeiPruefung1FuehrtZurueckNachVerwerfungMitUnveraendertemTextUndZaehler(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	verwerfungstext := "die Bewahrung der Schöpfung wichtiger sei als alles andere"
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", verwerfungstext)

	resp, err := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	if err != nil {
		t.Fatalf("einreichen: %v", err)
	}
	resp.Body.Close()

	resp, err = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"nein"}})
	if err != nil {
		t.Fatalf("pruefung1 nein: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittVerwerfung {
		t.Errorf("erwarte VERWERFUNG (nie an den Anfang), habe %s", these.Schritt)
	}
	if these.Verwerfung != verwerfungstext {
		t.Errorf("erwarte unveränderten Text %q, habe %q", verwerfungstext, these.Verwerfung)
	}
	if these.Ueberarbeitungen != 1 {
		t.Errorf("erwarte ueberarbeitungen=1, habe %d", these.Ueberarbeitungen)
	}
}

func TestNeinBeiPruefung2FuehrtZurueckNachVerwerfungMitUnveraendertemTextUndZaehler(t *testing.T) {
	app, client := werkstattApp(t, "Nation und Volkszugehörigkeit")
	verwerfungstext := "Herkunft über den Wert eines Menschen entscheide"
	bringeBisVorschau(t, app, client, "Nation und Volkszugehörigkeit", "wir vor Gott alle Geschwister sind", "keine Herkunft mehr wert ist als eine andere", verwerfungstext)

	resp, _ := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
	resp.Body.Close()

	resp, err := client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"nein"}})
	if err != nil {
		t.Fatalf("pruefung2 nein: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittVerwerfung {
		t.Errorf("erwarte VERWERFUNG, habe %s", these.Schritt)
	}
	if these.Verwerfung != verwerfungstext {
		t.Errorf("erwarte unveränderten Text %q, habe %q", verwerfungstext, these.Verwerfung)
	}
	if these.Ueberarbeitungen != 1 {
		t.Errorf("erwarte ueberarbeitungen=1, habe %d", these.Ueberarbeitungen)
	}
	if these.Schritt == schrittFreigegeben {
		t.Errorf("erwarte, dass ein nein bei Prüfung 2 nicht freigibt")
	}
}

func TestTheseWirdErstNachZweiJaFreigegebenUndDamitDruckbar(t *testing.T) {
	app, client := werkstattApp(t, "Schöpfung und Klima")
	bringeBisVorschau(t, app, client, "Schöpfung und Klima", "wir Gottes Schöpfung anvertraut sind", "wir verantwortlich mit ihr umgehen", "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung")

	resp, _ := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
	resp.Body.Close()

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)

	// Nach der ersten Prüffrage allein ist die These noch nicht freigegeben.
	resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
	resp.Body.Close()
	if these := aktuellerSchritt(t, app, gruppeID); these.Schritt == schrittFreigegeben {
		t.Fatalf("erwarte keine Freigabe nach nur einem Ja")
	}

	resp, err := client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"ja"}})
	if err != nil {
		t.Fatalf("pruefung2 ja: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwarte 200 nach Redirect, habe %d", resp.StatusCode)
	}

	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittFreigegeben {
		t.Fatalf("erwarte FREIGEGEBEN nach zwei Ja, habe %s", these.Schritt)
	}

	// "Druckbar" heißt hier: die Freigabeseite ist erreichbar.
	freigabeResp, err := client.Get(app.server.URL + "/gruppe/freigegeben")
	if err != nil {
		t.Fatalf("GET /gruppe/freigegeben: %v", err)
	}
	defer freigabeResp.Body.Close()
	if freigabeResp.StatusCode != http.StatusOK {
		t.Errorf("erwarte 200 für die freigegebene These, habe %d", freigabeResp.StatusCode)
	}
}

func TestGeraetOhneSchreibrechtKannKeinenSchrittAendern(t *testing.T) {
	app, _ := werkstattApp(t, "Schöpfung und Klima")

	var gruppeID int64
	app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
	_, tokenOhneRecht := app.GeraetBeitreten(t, gruppeID)
	ohneRecht := app.AlsGeraet(&http.Cookie{Name: geraetCookieName, Value: tokenOhneRecht})

	angabe := ersteBibelstelle(t, "Schöpfung und Klima")
	resp := ohneRecht.PostForm(t, "/gruppe/bibelstelle", url.Values{"bibelstelle": {angabe}})
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("erwarte 403 ohne Schreibrecht, habe %d", resp.StatusCode)
	}

	these := aktuellerSchritt(t, app, gruppeID)
	if these.Schritt != schrittBibelstelle {
		t.Errorf("erwarte, dass der Schritt unverändert bleibt, habe %s", these.Schritt)
	}
}

// TestAutomatWertetNichtInhaltlich ist die ausführbare Fassung der
// didaktischen Zusage aus SPEZIFIKATION.md ("Testentscheidungen") und darf
// nie gelöscht werden: Zwei Thesen mit unterschiedlichem Text, aber
// gleichem Antwortverhalten, kommen zum gleichen Ergebnis. Die Anwendung
// bewertet nichts, erkennt keine Parteinamen, filtert keine Wörter.
func TestAutomatWertetNichtInhaltlich(t *testing.T) {
	texte := []struct {
		weil, gilt, verwerfung string
	}{
		{
			weil:       "wir Gottes Schöpfung anvertraut sind",
			gilt:       "wir verantwortlich mit ihr umgehen",
			verwerfung: "wirtschaftliches Wachstum wichtiger sei als die Bewahrung der Schöpfung",
		},
		{
			weil:       "PARTEI X FORDERT SOFORTIGEN AUSSTIEG",
			gilt:       "wir uns diesem völlig unlogischen Unsinn nie anschließen dürfen",
			verwerfung: "irgendetwas völlig anderes und beliebig Falsches",
		},
	}

	for i, txt := range texte {
		app, client := werkstattApp(t, "Schöpfung und Klima")
		bringeBisVorschau(t, app, client, "Schöpfung und Klima", txt.weil, txt.gilt, txt.verwerfung)

		resp, _ := client.PostForm(app.server.URL+"/gruppe/vorschau/einreichen", nil)
		resp.Body.Close()
		resp, _ = client.PostForm(app.server.URL+"/gruppe/pruefung1", url.Values{"antwort": {"ja"}})
		resp.Body.Close()
		resp, err := client.PostForm(app.server.URL+"/gruppe/pruefung2", url.Values{"antwort": {"ja"}})
		if err != nil {
			t.Fatalf("text %d: pruefung2: %v", i, err)
		}
		resp.Body.Close()

		var gruppeID int64
		app.db.QueryRow(`SELECT id FROM gruppe`).Scan(&gruppeID)
		these := aktuellerSchritt(t, app, gruppeID)
		if these.Schritt != schrittFreigegeben {
			t.Errorf("text %d: erwarte FREIGEGEBEN unabhängig vom Wortlaut, habe %s", i, these.Schritt)
		}
		if !strings.Contains(schrittPfad(these.Schritt), "freigegeben") {
			t.Errorf("text %d: erwarte Pfad zu freigegeben", i)
		}
	}
}
