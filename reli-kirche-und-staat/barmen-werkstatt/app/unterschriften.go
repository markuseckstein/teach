package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// unterschriftenTrenner trennt die Vornamen innerhalb der einen
// unterschriften-Spalte (internal/db/db.go). Kein eigenes Tabellenblatt nötig
// — die Liste ist kurz (höchstens ein paar Namen je Gruppe) und lebt nur
// innerhalb einer These.
const unterschriftenTrenner = "; "

func unterschriftenTeilen(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, unterschriftenTrenner)
}

func unterschriftenZusammenfuegen(liste []string) string {
	return strings.Join(liste, unterschriftenTrenner)
}

// bereinigterVorname entfernt den Trenner und Zeilenumbrüche aus einer
// Nutzereingabe, statt eine erfundene Eingabe abzulehnen — ein Vorname mit
// einem Semikolon oder Zeilenumbruch würde sonst die Liste zerlegen und eine
// fremde "Unterschrift" vortäuschen, die niemand eingegeben hat.
func bereinigterVorname(roh string) string {
	ohneTrenner := strings.ReplaceAll(roh, ";", " ")
	return strings.Join(strings.Fields(ohneTrenner), " ")
}

// handleUnterschriftHinzufuegen trägt einen Vornamen ein. Die einzige Stelle,
// die das Formular dafür zeigt, ist die Vorschau — deshalb ist genau dieser
// Schritt auch die Bedingung in der WHERE-Klausel, nicht bloß "vor der
// Freigabe": Der Redirect danach führt immer echt dorthin zurück, statt sich
// auf einen zufälligen Nebeneffekt von geradeSchritts Weiterleitung zu
// verlassen. Die Prüfung selbst ist atomar, ohne vorgelagertes SELECT, nach
// demselben Muster wie wirdSchrittGewechselt (werkstatt.go).
func handleUnterschriftHinzufuegen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		vorname := bereinigterVorname(r.FormValue("vorname"))
		if vorname == "" {
			http.Error(w, "Vorname darf nicht leer sein", http.StatusBadRequest)
			return
		}

		res, err := database.Exec(
			`UPDATE these
			 SET unterschriften = CASE WHEN unterschriften = '' THEN ? ELSE unterschriften || ? || ? END
			 WHERE gruppe_id = ? AND schritt = ?`,
			vorname, unterschriftenTrenner, vorname, gruppeID, schrittVorschau,
		)
		if err != nil {
			http.Error(w, "Unterschrift konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("unterschrift hinzufügen: %v", err)
			return
		}
		betroffen, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Unterschrift konnte nicht gespeichert werden", http.StatusInternalServerError)
			log.Printf("unterschrift hinzufügen: %v", err)
			return
		}
		if betroffen == 0 {
			http.Error(w, "Unterschriften lassen sich nur in der Vorschau bearbeiten", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittVorschau), http.StatusSeeOther)
	}
}

// unterschriftEntfernenWennUnveraendert schreibt die neue Liste nur, wenn
// altwert noch der aktuelle Datenbankinhalt ist — ein optimistisches Schloss
// gegen den Fall, dass sich die Liste zwischen Laden und Schreiben schon
// geändert hat (zweites Gerät derselben Gruppe, oder die These hat inzwischen
// die Vorschau verlassen). Eine eigene, direkt testbare Funktion, weil dieses
// Wettlauffenster über HTTP allein nicht deterministisch auslösbar ist.
func unterschriftEntfernenWennUnveraendert(database *sql.DB, gruppeID int64, altwert, neuwert string) (geaendert bool, err error) {
	res, err := database.Exec(
		`UPDATE these SET unterschriften = ? WHERE gruppe_id = ? AND unterschriften = ? AND schritt = ?`,
		neuwert, gruppeID, altwert, schrittVorschau,
	)
	if err != nil {
		return false, err
	}
	betroffen, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return betroffen > 0, nil
}

// handleUnterschriftEntfernen entfernt eine einzelne Unterschrift über ihren
// Listenindex (nicht über den Namen, da zwei Mitglieder denselben Vornamen
// tragen können).
func handleUnterschriftEntfernen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gruppeID, ok := geraetMitSchreibrecht(w, r, database)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		index, err := strconv.Atoi(r.FormValue("index"))
		if err != nil || index < 0 {
			http.Error(w, "Ungültiger Index", http.StatusBadRequest)
			return
		}

		these, err := ladeThese(database, gruppeID)
		if err != nil {
			http.Error(w, "These konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("these laden (unterschrift entfernen): %v", err)
			return
		}
		if these.Schritt != schrittVorschau {
			http.Error(w, "Unterschriften lassen sich nur in der Vorschau bearbeiten", http.StatusConflict)
			return
		}

		liste := unterschriftenTeilen(these.Unterschriften)
		if index >= len(liste) {
			http.Error(w, "Ungültiger Index", http.StatusBadRequest)
			return
		}
		neueListe := append(append([]string{}, liste[:index]...), liste[index+1:]...)

		geaendert, err := unterschriftEntfernenWennUnveraendert(database, gruppeID, these.Unterschriften, unterschriftenZusammenfuegen(neueListe))
		if err != nil {
			http.Error(w, "Unterschrift konnte nicht entfernt werden", http.StatusInternalServerError)
			log.Printf("unterschrift entfernen: %v", err)
			return
		}
		if !geaendert {
			http.Error(w, "Die Liste hat sich inzwischen geändert", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, schrittPfad(schrittVorschau), http.StatusSeeOther)
	}
}
