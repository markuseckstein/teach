package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

// themenfelder sind die sechs zulässigen Themenfelder für Gruppen, aus
// SPEZIFIKATION.md, Abschnitt "Inhalte". Vorgang 0004 braucht nur diese
// Namen, um Zuweisungen zu prüfen; die zugehörigen Bibelstellen-Vorräte
// liefert Vorgang 0007 als eigene Daten.
var themenfelder = []string{
	"Umgang mit Geflüchteten",
	"Umgang mit Minderheiten",
	"Umgang mit der Wahrheit in der Öffentlichkeit",
	"Schöpfung und Klima",
	"Nation und Volkszugehörigkeit",
	"Freiheit der Kirche gegenüber jeder Regierung",
}

func istGueltigesThemenfeld(name string) bool {
	for _, t := range themenfelder {
		if t == name {
			return true
		}
	}
	return false
}

// beitrittscodeAlphabet lässt Zeichen aus, die laut vorgelesen oder auf
// einem Tabletbildschirm leicht verwechselt werden: 0/O, 1/I/l.
const beitrittscodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func neuerBeitrittscode() string {
	zeichen := make([]byte, 6)
	for i := range zeichen {
		zeichen[i] = beitrittscodeAlphabet[rand.IntN(len(beitrittscodeAlphabet))]
	}
	return string(zeichen)
}

func istEindeutigkeitsfehler(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// gruppeIDVon löst die Gruppe eines Kurses über ihre Nummer auf.
func gruppeIDVon(database *sql.DB, kursID int64, nummer int) (int64, error) {
	var gruppeID int64
	err := database.QueryRow(`SELECT id FROM gruppe WHERE kurs_id = ? AND nummer = ?`, kursID, nummer).Scan(&gruppeID)
	return gruppeID, err
}

type gruppenAnsicht struct {
	Nummer     int
	Themenfeld string
}

type kursAnsicht struct {
	ID            int64
	Name          string
	Beitrittscode string
	Gruppen       []gruppenAnsicht
	Themenfelder  []string
}

var kursTmpl = template.Must(template.ParseFS(templatesFS, "templates/kurs.html"))

var maxGruppenzahl = len(themenfelder)

// handleKursAnlegen nimmt Name und Gruppenzahl entgegen, legt einen Kurs mit
// eindeutigem Beitrittscode und die zugehörigen Gruppen an.
func handleKursAnlegen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		if name == "" {
			name = "Kurs"
		}

		gruppenzahl, err := strconv.Atoi(r.FormValue("gruppenzahl"))
		if err != nil || gruppenzahl < 1 || gruppenzahl > maxGruppenzahl {
			http.Error(w, fmt.Sprintf("Gruppenzahl muss zwischen 1 und %d liegen", maxGruppenzahl), http.StatusBadRequest)
			return
		}

		tx, err := database.Begin()
		if err != nil {
			http.Error(w, "Kurs konnte nicht angelegt werden", http.StatusInternalServerError)
			log.Printf("kurs anlegen: transaktion: %v", err)
			return
		}
		defer tx.Rollback()

		var kursID int64
		var beitrittscode string
		const maxVersuche = 5
		for versuch := 0; versuch < maxVersuche; versuch++ {
			beitrittscode = neuerBeitrittscode()
			res, err := tx.Exec(`INSERT INTO kurs (name, beitrittscode) VALUES (?, ?)`, name, beitrittscode)
			if err == nil {
				kursID, _ = res.LastInsertId()
				break
			}
			if !istEindeutigkeitsfehler(err) {
				http.Error(w, "Kurs konnte nicht angelegt werden", http.StatusInternalServerError)
				log.Printf("kurs anlegen: %v", err)
				return
			}
		}
		if kursID == 0 {
			http.Error(w, "Kurs konnte nicht angelegt werden: kein eindeutiger Beitrittscode gefunden", http.StatusInternalServerError)
			return
		}

		for nummer := 1; nummer <= gruppenzahl; nummer++ {
			if _, err := tx.Exec(`INSERT INTO gruppe (kurs_id, nummer) VALUES (?, ?)`, kursID, nummer); err != nil {
				http.Error(w, "Gruppen konnten nicht angelegt werden", http.StatusInternalServerError)
				log.Printf("gruppe anlegen: %v", err)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Kurs konnte nicht angelegt werden", http.StatusInternalServerError)
			log.Printf("kurs anlegen: commit: %v", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/kurse/%d", kursID), http.StatusSeeOther)
	}
}

// handleKursAnzeigen zeigt Beitrittscode und Gruppen eines Kurses samt Formular
// zur Themenfeld-Zuweisung.
func handleKursAnzeigen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var name, beitrittscode string
		err = database.QueryRow(`SELECT name, beitrittscode FROM kurs WHERE id = ?`, kursID).Scan(&name, &beitrittscode)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Kurs konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("kurs laden: %v", err)
			return
		}

		rows, err := database.Query(`SELECT nummer, themenfeld FROM gruppe WHERE kurs_id = ? ORDER BY nummer`, kursID)
		if err != nil {
			http.Error(w, "Gruppen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppen laden: %v", err)
			return
		}
		defer rows.Close()

		var gruppen []gruppenAnsicht
		for rows.Next() {
			var g gruppenAnsicht
			if err := rows.Scan(&g.Nummer, &g.Themenfeld); err != nil {
				http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("gruppe lesen: %v", err)
				return
			}
			gruppen = append(gruppen, g)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("gruppen lesen: %v", err)
			return
		}

		ansicht := kursAnsicht{
			ID:            kursID,
			Name:          name,
			Beitrittscode: beitrittscode,
			Gruppen:       gruppen,
			Themenfelder:  themenfelder,
		}

		renderTemplate(w, kursTmpl, ansicht)
	}
}

// handleThemenfeldZuweisen weist einer Gruppe ein Themenfeld zu. Ist es im
// selben Kurs bereits einer anderen Gruppe zugewiesen, wird abgelehnt.
//
// Die Eindeutigkeit prüft allein der UNIQUE-Index auf gruppe(kurs_id,
// themenfeld) aus dem Datenbankschema — ein vorgelagertes SELECT wäre ein
// Prüfen-dann-Handeln ohne Schutz vor zwei gleichzeitigen Zuweisungen aus
// verschiedenen Gruppen.
func handleThemenfeldZuweisen(database *sql.DB) http.HandlerFunc {
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
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formular ungültig", http.StatusBadRequest)
			return
		}

		themenfeld := r.FormValue("themenfeld")
		if !istGueltigesThemenfeld(themenfeld) {
			http.Error(w, "Unbekanntes Themenfeld", http.StatusBadRequest)
			return
		}

		res, err := database.Exec(
			`UPDATE gruppe SET themenfeld = ? WHERE kurs_id = ? AND nummer = ?`,
			themenfeld, kursID, nummer,
		)
		if istEindeutigkeitsfehler(err) {
			http.Error(w, "Themenfeld ist in diesem Kurs bereits vergeben", http.StatusConflict)
			return
		}
		if err != nil {
			http.Error(w, "Themenfeld konnte nicht zugewiesen werden", http.StatusInternalServerError)
			log.Printf("themenfeld zuweisen: %v", err)
			return
		}
		betroffen, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Themenfeld konnte nicht zugewiesen werden", http.StatusInternalServerError)
			log.Printf("themenfeld zuweisen: %v", err)
			return
		}
		if betroffen == 0 {
			http.NotFound(w, r)
			return
		}

		if gruppeID, err := gruppeIDVon(database, kursID, nummer); err != nil {
			log.Printf("gruppe-id für sse-benachrichtigung laden: %v", err)
		} else {
			werkstattHub.benachrichtigeGruppe(gruppeID)
		}

		http.Redirect(w, r, fmt.Sprintf("/kurse/%d", kursID), http.StatusSeeOther)
	}
}
