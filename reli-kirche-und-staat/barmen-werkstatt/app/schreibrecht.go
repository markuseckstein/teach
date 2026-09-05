package main

import (
	"database/sql"
	"log"
	"net/http"
)

// geraetHatSchreibrecht prüft, ob geraetID gerade das Schreibrecht der
// Gruppe gruppeID hält. Künftige Vorgänge, die tatsächlich etwas speichern
// (Vorgang 0008, Zustandsautomat der Werkstatt), rufen dies vor jeder
// schreibenden Aktion auf — "auch nicht durch eine direkt abgesetzte
// Anfrage" heißt: die Prüfung sitzt im Handler, nicht nur im Formular.
func geraetHatSchreibrecht(database *sql.DB, gruppeID, geraetID int64) (bool, error) {
	var schreibrechtGeraet sql.NullInt64
	err := database.QueryRow(`SELECT schreibrecht_geraet FROM gruppe WHERE id = ?`, gruppeID).Scan(&schreibrechtGeraet)
	if err != nil {
		return false, err
	}
	return schreibrechtGeraet.Valid && schreibrechtGeraet.Int64 == geraetID, nil
}

// handleSchreibrechtUebernehmen zieht das Schreibrecht der eigenen Gruppe an
// das anfragende Gerät. Die Übernahme wirkt sofort und ist für jedes Gerät
// der Gruppe über GET /gruppe sichtbar — live über den SSE-Kanal
// (Vorgang 0010), sonst spätestens beim nächsten Neuladen.
func handleSchreibrechtUebernehmen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geraetID, gruppeID, ok := aktuellesGeraet(r, database)
		if !ok {
			http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
			return
		}

		if _, err := database.Exec(
			`UPDATE gruppe SET schreibrecht_geraet = ?, schreibrecht_seit = ? WHERE id = ?`,
			geraetID, jetzt(), gruppeID,
		); err != nil {
			http.Error(w, "Schreibrecht konnte nicht übernommen werden", http.StatusInternalServerError)
			log.Printf("schreibrecht übernehmen: %v", err)
			return
		}

		werkstattHub.benachrichtigeGruppe(gruppeID)
		http.Redirect(w, r, "/gruppe", http.StatusSeeOther)
	}
}
