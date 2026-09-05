package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// Die sieben Schritte des Zustandsautomaten aus SPEZIFIKATION.md. Rückwärts
// geht es nur aus PRUEFUNG_1/PRUEFUNG_2 nach VERWERFUNG, nie an den Anfang.
const (
	schrittBibelstelle = "BIBELSTELLE"
	schrittPositiv     = "POSITIV"
	schrittVerwerfung  = "VERWERFUNG"
	schrittVorschau    = "VORSCHAU"
	schrittPruefung1   = "PRUEFUNG_1"
	schrittPruefung2   = "PRUEFUNG_2"
	schrittFreigegeben = "FREIGEGEBEN"
)

const verwerfungVorbelegt = "Wir verwerfen die falsche Lehre, als ob …"

// schrittReihenfolge ist die Vorwärtsreihenfolge des Automaten — Grundlage
// für die Notfall-Weiterschaltung im Regiepult (Vorgang 0011). Rückwärts
// (aus den Prüffragen nach VERWERFUNG) bleibt Sache der Handler selbst, hier
// geht es nur um "was kommt normalerweise als Nächstes".
var schrittReihenfolge = []string{
	schrittBibelstelle, schrittPositiv, schrittVerwerfung,
	schrittVorschau, schrittPruefung1, schrittPruefung2, schrittFreigegeben,
}

// naechsterSchrittNach liefert den auf schritt folgenden Schritt. weiter ist
// false bei FREIGEGEBEN (Endzustand) oder einem unbekannten Schritt.
func naechsterSchrittNach(schritt string) (naechster string, weiter bool) {
	for i, s := range schrittReihenfolge {
		if s == schritt && i+1 < len(schrittReihenfolge) {
			return schrittReihenfolge[i+1], true
		}
	}
	return "", false
}

// Die drei Phasen der Doppelstunde 6, die für die Anwendung etwas bedeuten
// (siehe lessons/0007-ds6-kopiervorlagen.html, Ablaufplan): Vorbereitung,
// bevor "Erarbeitung II" beginnt (Gruppen dürfen noch nicht in die
// Werkstatt hineinsehen — sonst würden sie vorauslesen, was die
// Dramaturgie zerstören würde); die 22-minütige Werkstatt-Phase selbst; und
// der Abschluss, in dem nichts mehr geschrieben wird, aber jede Gruppe ihre
// eigene, bereits erreichte Seite noch sehen darf.
const (
	phaseVorbereitung = ""
	phaseWerkstatt    = "werkstatt"
	phaseAbschluss    = "abschluss"
)

// werkstattDauer ist die in SPEZIFIKATION.md fest vorgegebene Dauer der
// Werkstatt-Phase ("die 22 Minuten der Werkstatt") — kein einstellbarer
// Wert, weil die Doppelstunde selbst diese Zahl vorgibt.
const werkstattDauer = 22 * time.Minute

// timerStand fasst den aktuellen Stand des Werkstatt-Timers zusammen —
// gemeinsam genutzt von Regiepult (Vorgang 0011) und Beamer-Ansicht
// (Vorgang 0012), die beide denselben Timer nur unterschiedlich groß zeigen.
// StartISO/EndeISO sind RFC3339 in UTC, damit die Beamer-Ansicht clientseitig
// live herunterzählen kann, ohne bei jeder Sekunde neu vom Server zu laden;
// Start/Ende sind die für Menschen lesbaren Uhrzeiten in Serverzeit.
type timerStand struct {
	Gestartet  bool
	StartISO   string
	EndeISO    string
	Start      string
	Ende       string
	Abgelaufen bool
}

func timerStandAus(werkstattGestartetUm sql.NullString) (timerStand, error) {
	if !werkstattGestartetUm.Valid {
		return timerStand{}, nil
	}
	start, err := time.Parse(time.RFC3339, werkstattGestartetUm.String)
	if err != nil {
		return timerStand{}, err
	}
	ende := start.Add(werkstattDauer)
	return timerStand{
		Gestartet:  true,
		StartISO:   start.UTC().Format(time.RFC3339),
		EndeISO:    ende.UTC().Format(time.RFC3339),
		Start:      start.Local().Format("15:04"),
		Ende:       ende.Local().Format("15:04"),
		Abgelaufen: time.Now().After(ende),
	}, nil
}

// aktivePhase liest die aktuelle Phase des Kurses, zu dem gruppeID gehört.
func aktivePhase(database *sql.DB, gruppeID int64) (string, error) {
	var phase string
	err := database.QueryRow(
		`SELECT k.aktive_phase FROM kurs k JOIN gruppe g ON g.kurs_id = k.id WHERE g.id = ?`,
		gruppeID,
	).Scan(&phase)
	return phase, err
}

func schrittPfad(schritt string) string {
	switch schritt {
	case schrittBibelstelle:
		return "/gruppe/bibelstelle"
	case schrittPositiv:
		return "/gruppe/positiv"
	case schrittVerwerfung:
		return "/gruppe/verwerfung"
	case schrittVorschau:
		return "/gruppe/vorschau"
	case schrittPruefung1:
		return "/gruppe/pruefung1"
	case schrittPruefung2:
		return "/gruppe/pruefung2"
	case schrittFreigegeben:
		return "/gruppe/freigegeben"
	default:
		return "/gruppe"
	}
}

type theseZeile struct {
	Themenfeld       string
	Bibelstelle      string
	Weil             string
	Gilt             string
	Verwerfung       string
	Schritt          string
	Ueberarbeitungen int
	Unterschriften   string
}

func ladeThese(database *sql.DB, gruppeID int64) (theseZeile, error) {
	var t theseZeile
	err := database.QueryRow(
		`SELECT g.themenfeld, t.bibelstelle, t.weil, t.gilt, t.verwerfung, t.schritt, t.ueberarbeitungen, t.unterschriften
		 FROM these t JOIN gruppe g ON g.id = t.gruppe_id
		 WHERE t.gruppe_id = ?`,
		gruppeID,
	).Scan(&t.Themenfeld, &t.Bibelstelle, &t.Weil, &t.Gilt, &t.Verwerfung, &t.Schritt, &t.Ueberarbeitungen, &t.Unterschriften)
	return t, err
}

// geradeSchritt löst das Geräte-Cookie auf und lädt die These der Gruppe.
// Ist die Werkstatt-Phase noch nicht freigeschaltet (Vorgang 0011), geht es
// zur Statusseite statt zu einer Schritt-Ansicht — sonst könnte eine Gruppe
// vorauslesen, was die Dramaturgie der Stunde zerstören würde. Steht die
// Gruppe nicht (mehr) im erwarteten Schritt, leitet es zum tatsächlichen
// Schritt weiter, statt eine veraltete oder übersprungene Ansicht zu
// zeigen — weiter ist dann false, der Aufrufer ist fertig.
func geradeSchritt(w http.ResponseWriter, r *http.Request, database *sql.DB, erwarteterSchritt string) (gruppeID int64, these theseZeile, weiter bool) {
	_, gruppeID, ok := aktuellesGeraet(r, database)
	if !ok {
		http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
		return 0, theseZeile{}, false
	}

	phase, err := aktivePhase(database, gruppeID)
	if err != nil {
		http.Error(w, "Phase konnte nicht geladen werden", http.StatusInternalServerError)
		log.Printf("phase laden: %v", err)
		return 0, theseZeile{}, false
	}
	if phase == phaseVorbereitung {
		http.Redirect(w, r, "/gruppe", http.StatusSeeOther)
		return 0, theseZeile{}, false
	}

	these, err = ladeThese(database, gruppeID)
	if err != nil {
		http.Error(w, "These konnte nicht geladen werden", http.StatusInternalServerError)
		log.Printf("these laden: %v", err)
		return 0, theseZeile{}, false
	}

	if these.Schritt != erwarteterSchritt {
		http.Redirect(w, r, schrittPfad(these.Schritt), http.StatusSeeOther)
		return 0, theseZeile{}, false
	}

	return gruppeID, these, true
}

// geraetMitSchreibrecht löst das Geräte-Cookie auf und prüft Werkstatt-Phase
// und Schreibrecht. Ohne gültiges Cookie geht es zum Beitritt, außerhalb der
// Werkstatt-Phase oder ohne Schreibrecht gibt es 409 bzw. 403 — in allen
// drei Fällen ist der Aufrufer fertig.
func geraetMitSchreibrecht(w http.ResponseWriter, r *http.Request, database *sql.DB) (gruppeID int64, ok bool) {
	geraetID, gruppeID, ok := aktuellesGeraet(r, database)
	if !ok {
		http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
		return 0, false
	}

	phase, err := aktivePhase(database, gruppeID)
	if err != nil {
		http.Error(w, "Phase konnte nicht geladen werden", http.StatusInternalServerError)
		log.Printf("phase laden: %v", err)
		return 0, false
	}
	if phase != phaseWerkstatt {
		http.Error(w, "Die Werkstatt ist gerade nicht freigeschaltet", http.StatusConflict)
		return 0, false
	}

	hat, err := geraetHatSchreibrecht(database, gruppeID, geraetID)
	if err != nil {
		http.Error(w, "Schreibrecht konnte nicht geprüft werden", http.StatusInternalServerError)
		log.Printf("schreibrecht prüfen: %v", err)
		return 0, false
	}
	if !hat {
		http.Error(w, "Kein Schreibrecht", http.StatusForbidden)
		return 0, false
	}

	return gruppeID, true
}

// wirdSchrittGewechselt führt eine schrittändernde Anweisung aus. Sie muss
// selbst "AND schritt = ?" mit dem erwarteten Schritt in der WHERE-Klausel
// enthalten — das macht Prüfen und Ändern in einer Anweisung atomar, ohne
// ein vorgelagertes SELECT, das mit einer zweiten Anfrage wettlaufen könnte
// (siehe die Themenfeld-Zuweisung in kurs.go für dasselbe Muster). gewechselt
// ist false, wenn die Gruppe nicht mehr im erwarteten Schritt stand.
func wirdSchrittGewechselt(database *sql.DB, anweisung string, args ...any) (gewechselt bool, err error) {
	res, err := database.Exec(anweisung, args...)
	if err != nil {
		return false, err
	}
	betroffen, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return betroffen > 0, nil
}

func istGueltigeBibelstelle(vorrat []Bibelstelle, angabe string) bool {
	for _, b := range vorrat {
		if b.Angabe == angabe {
			return true
		}
	}
	return false
}

type bibelstelleAnsicht struct {
	Themenfeld   string
	Bibelstellen []Bibelstelle
}

var bibelstelleTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-bibelstelle.html"))

// handleBibelstelleAnzeigen zeigt den Bibelstellen-Vorrat des zugewiesenen
// Themenfelds. Ohne gewählte Bibelstelle kommt eine Gruppe hier nicht
// vorbei — geradeSchritt leitet jede Anfrage in einem anderen Schritt weiter.
func handleBibelstelleAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, _, weiter := geradeSchritt(w, r, database, schrittBibelstelle)
		if !weiter {
			return
		}

		var themenfeld string
		if err := database.QueryRow(`SELECT themenfeld FROM gruppe WHERE id = ?`, gruppeID).Scan(&themenfeld); err != nil {
			http.Error(w, "Gruppe konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppe laden: %v", err)
			return
		}

		vorrat, ok := bibelstellenVorrat(themenfeld)
		if !ok {
			http.Error(w, "Dieser Gruppe ist noch kein Themenfeld zugewiesen.", http.StatusConflict)
			return
		}

		renderTemplate(w, bibelstelleTmpl, bibelstelleAnsicht{Themenfeld: themenfeld, Bibelstellen: vorrat})
	}
}

// handleBibelstelleWaehlen nimmt die gewählte Bibelstelle entgegen. Eine
// direkt abgesetzte Anfrage mit einer erfundenen Angabe oder in einem
// anderen Schritt als BIBELSTELLE wird abgelehnt.
func handleBibelstelleWaehlen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		angabe := strings.TrimSpace(r.FormValue("bibelstelle"))

		var themenfeld string
		if err := database.QueryRow(`SELECT themenfeld FROM gruppe WHERE id = ?`, gruppeID).Scan(&themenfeld); err != nil {
			http.Error(w, "Gruppe konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppe laden: %v", err)
			return
		}
		vorrat, themenfeldOK := bibelstellenVorrat(themenfeld)
		if !themenfeldOK || !istGueltigeBibelstelle(vorrat, angabe) {
			http.Error(w, "Unbekannte Bibelstelle", http.StatusBadRequest)
			return
		}

		gewechselt, err := wirdSchrittGewechselt(database,
			`UPDATE these SET bibelstelle = ?, schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
			angabe, schrittPositiv, gruppeID, schrittBibelstelle,
		)
		if err != nil {
			http.Error(w, "Bibelstelle konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("bibelstelle speichern: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittPositiv), http.StatusSeeOther)
	}
}

var positivTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-positiv.html"))

func handlePositivAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittPositiv)
		if !weiter {
			return
		}
		renderTemplate(w, positivTmpl, theseAnsichtAus(these))
	}
}

// handlePositivSpeichern verlangt beide Felder gefüllt, bevor es nach
// VERWERFUNG weiterschaltet — sonst zerfällt die Form der These in eine
// bloße Meinung ohne biblische Bindung.
func handlePositivSpeichern(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		weil := strings.TrimSpace(r.FormValue("weil"))
		gilt := strings.TrimSpace(r.FormValue("gilt"))
		if weil == "" || gilt == "" {
			http.Error(w, "Beide Felder müssen ausgefüllt sein", http.StatusBadRequest)
			return
		}

		gewechselt, err := wirdSchrittGewechselt(database,
			`UPDATE these SET weil = ?, gilt = ?, schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
			weil, gilt, schrittVerwerfung, gruppeID, schrittPositiv,
		)
		if err != nil {
			http.Error(w, "Konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("positiv speichern: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittVerwerfung), http.StatusSeeOther)
	}
}

var verwerfungTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-verwerfung.html"))

func handleVerwerfungAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittVerwerfung)
		if !weiter {
			return
		}
		if these.Verwerfung == "" {
			these.Verwerfung = verwerfungVorbelegt
		}
		renderTemplate(w, verwerfungTmpl, theseAnsichtAus(these))
	}
}

func handleVerwerfungSpeichern(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		verwerfung := strings.TrimSpace(r.FormValue("verwerfung"))
		if verwerfung == "" {
			http.Error(w, "Die Verwerfung darf nicht leer sein", http.StatusBadRequest)
			return
		}

		gewechselt, err := wirdSchrittGewechselt(database,
			`UPDATE these SET verwerfung = ?, schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
			verwerfung, schrittVorschau, gruppeID, schrittVerwerfung,
		)
		if err != nil {
			http.Error(w, "Konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("verwerfung speichern: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittVorschau), http.StatusSeeOther)
	}
}

// theseAnsicht zeigt die zusammengesetzte These — für Vorschau, beide
// Prüffragen und die Freigabe identisch, deshalb ein gemeinsamer Typ.
type theseAnsicht struct {
	Themenfeld     string
	Bibelstelle    string
	Weil           string
	Gilt           string
	Verwerfung     string
	Unterschriften string
}

func theseAnsichtAus(these theseZeile) theseAnsicht {
	return theseAnsicht{
		Themenfeld:     these.Themenfeld,
		Bibelstelle:    these.Bibelstelle,
		Weil:           these.Weil,
		Gilt:           these.Gilt,
		Verwerfung:     these.Verwerfung,
		Unterschriften: these.Unterschriften,
	}
}

var vorschauTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-vorschau.html"))

func handleVorschauAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittVorschau)
		if !weiter {
			return
		}
		renderTemplate(w, vorschauTmpl, theseAnsichtAus(these))
	}
}

func handleVorschauEinreichen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}

		gewechselt, err := wirdSchrittGewechselt(database,
			`UPDATE these SET schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
			schrittPruefung1, gruppeID, schrittVorschau,
		)
		if err != nil {
			http.Error(w, "Konnte nicht eingereicht werden", http.StatusInternalServerError)
			log.Printf("vorschau einreichen: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittPruefung1), http.StatusSeeOther)
	}
}

var pruefung1Tmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-pruefung1.html"))

func handlePruefung1Anzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittPruefung1)
		if !weiter {
			return
		}
		renderTemplate(w, pruefung1Tmpl, theseAnsichtAus(these))
	}
}

// handlePruefung1Beantworten: "nein" führt zurück nach VERWERFUNG — nie an
// den Anfang —, der Text bleibt unverändert stehen (nur schritt und
// ueberarbeitungen ändern sich), "ja" schaltet zur zweiten Prüffrage weiter.
func handlePruefung1Beantworten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		var gewechselt bool
		var err error
		var naechsterSchritt string
		switch r.FormValue("antwort") {
		case "ja":
			naechsterSchritt = schrittPruefung2
			gewechselt, err = wirdSchrittGewechselt(database,
				`UPDATE these SET schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
				naechsterSchritt, gruppeID, schrittPruefung1,
			)
		case "nein":
			naechsterSchritt = schrittVerwerfung
			gewechselt, err = wirdSchrittGewechselt(database,
				`UPDATE these SET schritt = ?, ueberarbeitungen = ueberarbeitungen + 1 WHERE gruppe_id = ? AND schritt = ?`,
				naechsterSchritt, gruppeID, schrittPruefung1,
			)
		default:
			http.Error(w, "Antwort muss ja oder nein sein", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("pruefung1 beantworten: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(naechsterSchritt), http.StatusSeeOther)
	}
}

var pruefung2Tmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-pruefung2.html"))

func handlePruefung2Anzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittPruefung2)
		if !weiter {
			return
		}
		renderTemplate(w, pruefung2Tmpl, theseAnsichtAus(these))
	}
}

// handlePruefung2Beantworten: "ja" gibt frei und setzt freigegeben_am,
// "nein" führt wie bei Prüffrage 1 zurück nach VERWERFUNG.
func handlePruefung2Beantworten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		var gewechselt bool
		var err error
		var naechsterSchritt string
		switch r.FormValue("antwort") {
		case "ja":
			naechsterSchritt = schrittFreigegeben
			gewechselt, err = wirdSchrittGewechselt(database,
				`UPDATE these SET schritt = ?, freigegeben_am = ? WHERE gruppe_id = ? AND schritt = ?`,
				naechsterSchritt, jetzt(), gruppeID, schrittPruefung2,
			)
		case "nein":
			naechsterSchritt = schrittVerwerfung
			gewechselt, err = wirdSchrittGewechselt(database,
				`UPDATE these SET schritt = ?, ueberarbeitungen = ueberarbeitungen + 1 WHERE gruppe_id = ? AND schritt = ?`,
				naechsterSchritt, gruppeID, schrittPruefung2,
			)
		default:
			http.Error(w, "Antwort muss ja oder nein sein", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("pruefung2 beantworten: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(naechsterSchritt), http.StatusSeeOther)
	}
}

var freigegebenTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-freigegeben.html"))

func handleFreigegebenAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, these, weiter := geradeSchritt(w, r, database, schrittFreigegeben)
		if !weiter {
			return
		}
		renderTemplate(w, freigegebenTmpl, theseAnsichtAus(these))
	}
}
