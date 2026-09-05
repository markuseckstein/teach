package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	Bibelstelle      string
	Weil             string
	Gilt             string
	Verwerfung       string
	Schritt          string
	Ueberarbeitungen int
}

func ladeThese(database *sql.DB, gruppeID int64) (theseZeile, error) {
	var t theseZeile
	err := database.QueryRow(
		`SELECT bibelstelle, weil, gilt, verwerfung, schritt, ueberarbeitungen FROM these WHERE gruppe_id = ?`,
		gruppeID,
	).Scan(&t.Bibelstelle, &t.Weil, &t.Gilt, &t.Verwerfung, &t.Schritt, &t.Ueberarbeitungen)
	return t, err
}

// geradeSchritt löst das Geräte-Cookie auf und lädt die These der Gruppe.
// Steht die Gruppe nicht (mehr) im erwarteten Schritt, leitet es zum
// tatsächlichen Schritt weiter, statt eine veraltete oder übersprungene
// Ansicht zu zeigen — weiter ist dann false, der Aufrufer ist fertig.
func geradeSchritt(w http.ResponseWriter, r *http.Request, database *sql.DB, erwarteterSchritt string) (gruppeID int64, these theseZeile, weiter bool) {
	_, gruppeID, ok := aktuellesGeraet(r, database)
	if !ok {
		http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
		return 0, theseZeile{}, false
	}

	these, err := ladeThese(database, gruppeID)
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

// geraetMitSchreibrecht löst das Geräte-Cookie auf und prüft das
// Schreibrecht. Ohne gültiges Cookie geht es zum Beitritt, ohne
// Schreibrecht gibt es 403 — in beiden Fällen ist der Aufrufer fertig.
func geraetMitSchreibrecht(w http.ResponseWriter, r *http.Request, database *sql.DB) (gruppeID int64, ok bool) {
	geraetID, gruppeID, ok := aktuellesGeraet(r, database)
	if !ok {
		http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
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
	Bibelstelle string
	Weil        string
	Gilt        string
	Verwerfung  string
}

func theseAnsichtAus(these theseZeile) theseAnsicht {
	return theseAnsicht{Bibelstelle: these.Bibelstelle, Weil: these.Weil, Gilt: these.Gilt, Verwerfung: these.Verwerfung}
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
