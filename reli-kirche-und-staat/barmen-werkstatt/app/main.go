// Kommando barmen-werkstatt startet die begleitende Webanwendung zur
// Unterrichtseinheit "Kirche und Staat".
package main

import (
	"bytes"
	"database/sql"
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"

	"barmen-werkstatt/internal/db"
)

//go:embed static
var staticFS embed.FS

//go:embed templates
var templatesFS embed.FS

var indexTmpl = template.Must(template.ParseFS(templatesFS, "templates/index.html"))

// newMux baut den Router der Anwendung. database ist derzeit ungenutzt und
// dient künftigen Vorgängen als Anschlussstelle für Fachlogik.
func newMux(database *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		var body bytes.Buffer
		if err := indexTmpl.Execute(&body, nil); err != nil {
			http.Error(w, "Vorlage konnte nicht gerendert werden", http.StatusInternalServerError)
			log.Printf("Vorlagenfehler: %v", err)
			return
		}
		w.Write(body.Bytes())
	})
	return mux
}

func main() {
	port := flag.String("port", "8080", "Port, auf dem der Server lauscht")
	dbPath := flag.String("db", "barmen.db", "Pfad zur SQLite-Datenbankdatei")
	flag.Parse()

	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("Datenbank konnte nicht geöffnet werden: %v", err)
	}
	defer database.Close()

	addr := ":" + *port
	log.Printf("Barmen-Werkstatt lauscht auf %s (Datenbank: %s)", addr, *dbPath)
	if err := http.ListenAndServe(addr, newMux(database)); err != nil {
		log.Fatalf("Server beendet: %v", err)
	}
}
