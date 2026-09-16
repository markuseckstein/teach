package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// beamerGruppe zeigt von einer Gruppe nur, was am Beamer niemanden
// vorführt: keine Überarbeitungszahl (Vorgang 0011 hält die bewusst nur im
// Regiepult sichtbar), kein Name — nur Themenfeld und ein Fortschrittsbalken
// aus Punkten, damit niemand liest "hängt bei PRUEFUNG_1 fest".
type beamerGruppe struct {
	Nummer      int
	Themenfeld  string
	Punkte      []bool
	Freigegeben bool
}

// beamerThese ist eine einzelne freigegebene These für die Schlussansicht —
// dieselben Felder wie auf dem Ausdruck einer Gruppe (druckansicht.go),
// hier aber alle sechs nebeneinander zum Vergleich.
type beamerThese struct {
	Nummer      int
	Themenfeld  string
	Bibelstelle string
	Weil        string
	Gilt        string
	Verwerfung  string
}

type beamerAnsicht struct {
	KursID             int64
	Phase              string
	Timer              timerStand
	Gruppen            []beamerGruppe
	FreigegebeneThesen []beamerThese
}

var beamerTmpl = template.Must(template.ParseFS(templatesFS, "templates/beamer.html"))

// punkteAus bildet den Schritt einer Gruppe auf einen Fortschrittsbalken ab:
// ein Punkt pro Zustand des Automaten (schrittReihenfolge in werkstatt.go),
// gefüllt bis einschließlich des aktuellen. Ein unbekannter Schritt (kann bei
// intaktem Automaten nicht vorkommen) ergibt einen leeren Balken statt eines
// Ausschlags — sicherer Fehlschlag statt falscher Anzeige am Beamer.
func punkteAus(schritt string) []bool {
	idx := -1
	for i, s := range schrittReihenfolge {
		if s == schritt {
			idx = i
			break
		}
	}
	punkte := make([]bool, len(schrittReihenfolge))
	for i := range punkte {
		punkte[i] = i <= idx
	}
	return punkte
}

// handleBeamerAnzeigen zeigt die einzige Ansicht, die die ganze Klasse sieht:
// Fortschritt aller Gruppen ohne Ueberarbeitungen und ohne Namen, den Timer,
// und — sobald Thesen freigegeben sind — sie nebeneinander zum Vergleich.
func handleBeamerAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var phase string
		var werkstattGestartetUm sql.NullString
		err = database.QueryRow(
			`SELECT aktive_phase, werkstatt_gestartet_um FROM kurs WHERE id = ?`,
			kursID,
		).Scan(&phase, &werkstattGestartetUm)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Kurs konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("kurs laden (beamer): %v", err)
			return
		}

		stand, err := timerStandAus(werkstattGestartetUm)
		if err != nil {
			http.Error(w, "Timer konnte nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("timer parsen (beamer): %v", err)
			return
		}

		ansicht := beamerAnsicht{KursID: kursID, Phase: phase, Timer: stand}

		rows, err := database.Query(
			`SELECT g.nummer, g.themenfeld, COALESCE(t.schritt, ?)
			 FROM gruppe g LEFT JOIN these t ON t.gruppe_id = g.id
			 WHERE g.kurs_id = ?
			 ORDER BY g.nummer`,
			schrittBibelstelle, kursID,
		)
		if err != nil {
			http.Error(w, "Gruppen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppen laden (beamer): %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var g beamerGruppe
			var schritt string
			if err := rows.Scan(&g.Nummer, &g.Themenfeld, &schritt); err != nil {
				http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("gruppe lesen (beamer): %v", err)
				return
			}
			g.Punkte = punkteAus(schritt)
			g.Freigegeben = schritt == schrittFreigegeben
			ansicht.Gruppen = append(ansicht.Gruppen, g)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("gruppen lesen (beamer): %v", err)
			return
		}

		theseRows, err := database.Query(
			`SELECT g.nummer, g.themenfeld, t.bibelstelle, t.weil, t.gilt, t.verwerfung
			 FROM these t JOIN gruppe g ON g.id = t.gruppe_id
			 WHERE g.kurs_id = ? AND t.schritt = ?
			 ORDER BY g.nummer`,
			kursID, schrittFreigegeben,
		)
		if err != nil {
			http.Error(w, "Freigegebene Thesen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("freigegebene thesen laden (beamer): %v", err)
			return
		}
		defer theseRows.Close()

		for theseRows.Next() {
			var t beamerThese
			if err := theseRows.Scan(&t.Nummer, &t.Themenfeld, &t.Bibelstelle, &t.Weil, &t.Gilt, &t.Verwerfung); err != nil {
				http.Error(w, "Freigegebene Thesen konnten nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("freigegebene these lesen (beamer): %v", err)
				return
			}
			ansicht.FreigegebeneThesen = append(ansicht.FreigegebeneThesen, t)
		}
		if err := theseRows.Err(); err != nil {
			http.Error(w, "Freigegebene Thesen konnten nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("freigegebene thesen lesen (beamer): %v", err)
			return
		}

		renderTemplate(w, beamerTmpl, ansicht)
	}
}

// handleBeamerLive hält eine SSE-Verbindung für den Beamer offen. Anders als
// handleGruppeLive (sse.go), das genau einer Gruppe zugeordnet ist, kennt der
// Beamer keine Gruppen-Identität — er meldet sich deshalb bei den Kanälen
// ALLER Gruppen des Kurses zugleich an und fasst sie zu einem Signal
// zusammen. Das braucht keine eigene Datenstruktur im Hub: jede
// schrittändernde Anfrage benachrichtigt ohnehin schon ihre Gruppe
// (werkstattHub.benachrichtigeGruppe), und ein Phasen- oder Timer-Wechsel im
// Regiepult benachrichtigt bereits alle Gruppen eines Kurses
// (benachrichtigeKurs) — der Beamer hört einfach überall mit.
func handleBeamerLive(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		gruppenIDs, err := gruppenIDsVonKurs(database, kursID)
		if err != nil {
			http.Error(w, "Gruppen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppen laden (beamer-live): %v", err)
			return
		}
		if len(gruppenIDs) == 0 {
			http.NotFound(w, r)
			return
		}

		flusher, flushtOk := w.(http.Flusher)
		if !flushtOk {
			http.Error(w, "Streaming wird nicht unterstützt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// sammel bündelt die Signale aller angemeldeten Gruppen-Kanäle zu
		// einem einzigen Strom; fertig beendet die Weiterleitungs-Goroutinen
		// beim Verbindungsende. Non-blocking wie beim Hub selbst: es wartet
		// ohnehin nie mehr als ein ungelesenes "etwas hat sich geändert".
		sammel := make(chan struct{}, 1)
		fertig := make(chan struct{})
		defer close(fertig)

		for _, gruppeID := range gruppenIDs {
			kanal := werkstattHub.abonnieren(gruppeID)
			defer werkstattHub.abbestellen(gruppeID, kanal)

			go func(kanal chan struct{}) {
				for {
					select {
					case <-fertig:
						return
					case <-kanal:
						select {
						case sammel <- struct{}{}:
						default:
						}
					}
				}
			}(kanal)
		}

		if _, err := fmt.Fprint(w, "data: sync\n\n"); err != nil {
			return
		}
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-sammel:
				if _, err := fmt.Fprint(w, "data: sync\n\n"); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}
