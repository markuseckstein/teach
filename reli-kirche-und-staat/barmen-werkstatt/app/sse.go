package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
)

// sseHub verteilt "etwas hat sich geändert"-Signale an die verbundenen
// Geräte einer Gruppe. Der Kanal transportiert keine Nutzdaten — jedes
// Signal heißt für den Client schlicht "lade neu", und die neu geladene
// Seite liest wie jede andere Anfrage den aktuellen Stand aus der
// Datenbank. Damit bleibt die Anwendung ohne SSE (oder bei abgeschaltetem
// JavaScript) genauso korrekt bedienbar, nur ohne das Live-Update — ein
// Neuladen zeigt immer den richtigen Stand.
type sseHub struct {
	mu      sync.Mutex
	gruppen map[int64]map[chan struct{}]struct{}
}

func newSSEHub() *sseHub {
	return &sseHub{gruppen: make(map[int64]map[chan struct{}]struct{})}
}

var werkstattHub = newSSEHub()

func (h *sseHub) abonnieren(gruppeID int64) chan struct{} {
	kanal := make(chan struct{}, 1)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.gruppen[gruppeID] == nil {
		h.gruppen[gruppeID] = make(map[chan struct{}]struct{})
	}
	h.gruppen[gruppeID][kanal] = struct{}{}
	return kanal
}

func (h *sseHub) abbestellen(gruppeID int64, kanal chan struct{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.gruppen[gruppeID], kanal)
	if len(h.gruppen[gruppeID]) == 0 {
		delete(h.gruppen, gruppeID)
	}
}

// benachrichtigeGruppe weckt alle Geräte einer Gruppe — für Zustandswechsel
// der Werkstatt (Vorgang 0008) und Schreibrecht-Übernahmen (Vorgang 0006).
func (h *sseHub) benachrichtigeGruppe(gruppeID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for kanal := range h.gruppen[gruppeID] {
		select {
		case kanal <- struct{}{}:
		default:
			// Es wartet schon ein ungelesenes Signal — ein weiteres bringt
			// nichts, der Client lädt beim nächsten ohnehin alles neu.
		}
	}
}

// benachrichtigeKurs weckt alle Gruppen eines Kurses — für kursweite
// Ereignisse wie Phase und Timer, die das Regiepult (Vorgang 0011) auslöst.
func (h *sseHub) benachrichtigeKurs(database *sql.DB, kursID int64) error {
	gruppenIDs, err := gruppenIDsVonKurs(database, kursID)
	if err != nil {
		return err
	}

	for _, id := range gruppenIDs {
		h.benachrichtigeGruppe(id)
	}
	return nil
}

// handleGruppeLive hält eine SSE-Verbindung für das Gerät offen. Direkt
// nach dem Verbindungsaufbau (auch nach einem Wiederverbinden durch den
// Browser selbst) kommt sofort ein Signal — kein halb aktueller Bildschirm,
// weil der Client danach ohnehin neu lädt.
func handleGruppeLive(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, gruppeID, ok := aktuellesGeraet(r, database)
		if !ok {
			http.Error(w, "Kein gültiges Gerät", http.StatusUnauthorized)
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

		kanal := werkstattHub.abonnieren(gruppeID)
		defer werkstattHub.abbestellen(gruppeID, kanal)

		if _, err := fmt.Fprint(w, "data: sync\n\n"); err != nil {
			return
		}
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-kanal:
				if _, err := fmt.Fprint(w, "data: sync\n\n"); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}
