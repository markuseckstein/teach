package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// schreibrechtWarnungAb ist die Schwelle, ab der das Regiepult meldet, dass
// in einer Gruppe seit Langem dasselbe Gerät schreibt — ein Hinweis auf eine
// stille Übernahme durch eine einzelne Person. Die Werkstatt-Phase dauert
// insgesamt 22 Minuten; die Hälfte davon ohne Wechsel ist ein plausibler
// Anlass, kurz hinzuschauen, ohne bei jedem kurzen Halten des Tablets
// Alarm zu schlagen.
const schreibrechtWarnungAb = 11 * time.Minute

type gruppenFortschritt struct {
	Nummer              int
	Themenfeld          string
	Schritt             string
	Ueberarbeitungen    int
	SchreibrechtWarnung bool
}

type regiepultAnsicht struct {
	KursID              int64
	Name                string
	Beitrittscode       string
	Phase               string
	WerkstattGestartet  bool
	WerkstattStart      string // HH:MM, Serverzeit
	WerkstattEnde       string // HH:MM, Serverzeit
	WerkstattAbgelaufen bool
	Gruppen             []gruppenFortschritt
}

var regiepultTmpl = template.Must(template.ParseFS(templatesFS, "templates/regiepult.html"))

// handleRegiepultAnzeigen zeigt die Steuerzentrale der Lehrkraft: Phase,
// Timer und den Fortschritt jeder Gruppe auf einen Blick — einschließlich
// der Zahl der Überarbeitungen, die bewusst nirgendwo sonst in der
// Anwendung sichtbar ist (weder auf den Geräten noch später am Beamer).
func handleRegiepultAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var name, beitrittscode, phase string
		var werkstattGestartetUm sql.NullString
		err = database.QueryRow(
			`SELECT name, beitrittscode, aktive_phase, werkstatt_gestartet_um FROM kurs WHERE id = ?`,
			kursID,
		).Scan(&name, &beitrittscode, &phase, &werkstattGestartetUm)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Kurs konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("kurs laden: %v", err)
			return
		}

		ansicht := regiepultAnsicht{
			KursID:        kursID,
			Name:          name,
			Beitrittscode: beitrittscode,
			Phase:         phase,
		}
		if werkstattGestartetUm.Valid {
			start, err := time.Parse(time.RFC3339, werkstattGestartetUm.String)
			if err != nil {
				http.Error(w, "Timer konnte nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("timer parsen: %v", err)
				return
			}
			ende := start.Add(werkstattDauer)
			ansicht.WerkstattGestartet = true
			ansicht.WerkstattStart = start.Local().Format("15:04")
			ansicht.WerkstattEnde = ende.Local().Format("15:04")
			ansicht.WerkstattAbgelaufen = time.Now().After(ende)
		}

		rows, err := database.Query(
			`SELECT g.nummer, g.themenfeld, COALESCE(g.schreibrecht_seit, ''),
			        COALESCE(t.schritt, ?), COALESCE(t.ueberarbeitungen, 0)
			 FROM gruppe g LEFT JOIN these t ON t.gruppe_id = g.id
			 WHERE g.kurs_id = ?
			 ORDER BY g.nummer`,
			schrittBibelstelle, kursID,
		)
		if err != nil {
			http.Error(w, "Gruppen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppen laden: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var g gruppenFortschritt
			var schreibrechtSeit string
			if err := rows.Scan(&g.Nummer, &g.Themenfeld, &schreibrechtSeit, &g.Schritt, &g.Ueberarbeitungen); err != nil {
				http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("gruppe lesen: %v", err)
				return
			}
			if schreibrechtSeit != "" {
				if seit, err := time.Parse(time.RFC3339, schreibrechtSeit); err == nil {
					g.SchreibrechtWarnung = time.Since(seit) > schreibrechtWarnungAb
				}
			}
			ansicht.Gruppen = append(ansicht.Gruppen, g)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("gruppen lesen: %v", err)
			return
		}

		renderTemplate(w, regiepultTmpl, ansicht)
	}
}

// handlePhaseSetzen schaltet die Phase des Kurses um. "Alle Geräte folgen
// sofort" heißt hier: die SSE-Benachrichtigung an den ganzen Kurs, die schon
// aus Vorgang 0010 dafür vorbereitet ist.
func handlePhaseSetzen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		phase := r.FormValue("phase")
		switch phase {
		case phaseVorbereitung, phaseWerkstatt, phaseAbschluss:
		default:
			http.Error(w, "Unbekannte Phase", http.StatusBadRequest)
			return
		}

		res, err := database.Exec(`UPDATE kurs SET aktive_phase = ? WHERE id = ?`, phase, kursID)
		if err != nil {
			http.Error(w, "Phase konnte nicht gesetzt werden", http.StatusInternalServerError)
			log.Printf("phase setzen: %v", err)
			return
		}
		if betroffen, _ := res.RowsAffected(); betroffen == 0 {
			http.NotFound(w, r)
			return
		}

		if err := werkstattHub.benachrichtigeKurs(database, kursID); err != nil {
			log.Printf("kurs benachrichtigen: %v", err)
		}
		http.Redirect(w, r, fmt.Sprintf("/kurse/%d/regiepult", kursID), http.StatusSeeOther)
	}
}

// handleTimerStarten setzt den Startzeitpunkt der 22-minütigen
// Werkstatt-Phase neu — unabhängig vom Phasenwechsel selbst, damit die
// Lehrkraft den Timer erst nach den letzten Instruktionen an die Klasse
// auslösen kann.
func handleTimerStarten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		res, err := database.Exec(`UPDATE kurs SET werkstatt_gestartet_um = ? WHERE id = ?`, jetzt(), kursID)
		if err != nil {
			http.Error(w, "Timer konnte nicht gestartet werden", http.StatusInternalServerError)
			log.Printf("timer starten: %v", err)
			return
		}
		if betroffen, _ := res.RowsAffected(); betroffen == 0 {
			http.NotFound(w, r)
			return
		}

		if err := werkstattHub.benachrichtigeKurs(database, kursID); err != nil {
			log.Printf("kurs benachrichtigen: %v", err)
		}
		http.Redirect(w, r, fmt.Sprintf("/kurse/%d/regiepult", kursID), http.StatusSeeOther)
	}
}

// handleGruppeManuellWeiterschalten ist die Notfall-Weiterschaltung: Sie
// prüft nichts inhaltlich und zählt nicht als Überarbeitung — sie
// überspringt lediglich einen technisch hängenden Schritt in der normalen
// Vorwärtsreihenfolge des Automaten.
func handleGruppeManuellWeiterschalten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		nummer, err := strconv.Atoi(r.PathValue("nummer"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		gruppeID, err := gruppeIDVon(database, kursID, nummer)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Gruppe konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppe laden: %v", err)
			return
		}

		these, err := ladeThese(database, gruppeID)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "These konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("these laden: %v", err)
			return
		}

		naechster, gibtsWeiter := naechsterSchrittNach(these.Schritt)
		if !gibtsWeiter {
			http.Error(w, "Diese Gruppe ist bereits freigegeben", http.StatusConflict)
			return
		}

		gewechselt, err := wirdSchrittGewechselt(database,
			`UPDATE these SET schritt = ? WHERE gruppe_id = ? AND schritt = ?`,
			naechster, gruppeID, these.Schritt,
		)
		if err != nil {
			http.Error(w, "Konnte nicht weitergeschaltet werden", http.StatusInternalServerError)
			log.Printf("manuell weiterschalten: %v", err)
			return
		}
		if !gewechselt {
			http.Error(w, "Falscher Schritt", http.StatusConflict)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, fmt.Sprintf("/kurse/%d/regiepult", kursID), http.StatusSeeOther)
	}
}
