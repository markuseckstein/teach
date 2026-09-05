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

//go:embed data
var dataFS embed.FS

var indexTmpl = template.Must(template.ParseFS(templatesFS, "templates/index.html"))

// lessonsDirFlag zeigt auf das lessons-Verzeichnis mit den fertigen
// Kopiervorlagen (Vorgang 0017). Es liegt zwei Ebenen über diesem Go-Modul
// (barmen-werkstatt/app/../../lessons — go:embed erlaubt kein ".." im
// Pattern), deshalb wird es zur Laufzeit über einen Dateisystem-Pfad
// ausgeliefert statt eingebettet. Als Paket-Variable deklariert (statt lokal
// in main()), damit Tests den Wert für die Dauer eines Tests umbiegen können,
// ohne dass newMux einen zusätzlichen Parameter braucht.
var lessonsDirFlag = flag.String("lessons-dir", "../../lessons", "Pfad zum lessons-Verzeichnis mit den Kopiervorlagen")

// newMux baut den Router der Anwendung.
func newMux(database *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.Handle("/lessons/", http.StripPrefix("/lessons/", http.FileServer(http.Dir(*lessonsDirFlag))))
	mux.HandleFunc("GET /inhalte", handleInhalteAnzeigen())
	mux.HandleFunc("POST /kurse", handleKursAnlegen(database))
	mux.HandleFunc("GET /kurse/{id}", handleKursAnzeigen(database))
	mux.HandleFunc("POST /kurse/{id}/gruppen/{nummer}/themenfeld", handleThemenfeldZuweisen(database))
	mux.HandleFunc("GET /kurse/{id}/qr.png", handleKursQR(database))
	mux.HandleFunc("GET /kurse/{id}/beenden", handleKursBeendenAnzeigen(database))
	mux.HandleFunc("POST /kurse/{id}/beenden", handleKursBeenden(database))
	mux.HandleFunc("GET /kurse/{id}/gruppe-waehlen", handleGruppeWaehlen(database))
	mux.HandleFunc("GET /kurse/{id}/regiepult", handleRegiepultAnzeigen(database))
	mux.HandleFunc("POST /kurse/{id}/phase", handlePhaseSetzen(database))
	mux.HandleFunc("POST /kurse/{id}/timer/starten", handleTimerStarten(database))
	mux.HandleFunc("POST /kurse/{id}/gruppen/{nummer}/weiterschalten", handleGruppeManuellWeiterschalten(database))
	mux.HandleFunc("POST /kurse/{id}/gruppen/{nummer}/beitreten", handleGruppeBeitreten(database))
	mux.HandleFunc("GET /beitreten", handleBeitreten(database))
	mux.HandleFunc("GET /gruppe", handleGruppeStatus(database))
	mux.HandleFunc("POST /gruppe/schreibrecht", handleSchreibrechtUebernehmen(database))
	mux.HandleFunc("GET /gruppe/bibelstelle", handleBibelstelleAnzeigen(database))
	mux.HandleFunc("POST /gruppe/bibelstelle", handleBibelstelleWaehlen(database))
	mux.HandleFunc("GET /gruppe/positiv", handlePositivAnzeigen(database))
	mux.HandleFunc("POST /gruppe/positiv", handlePositivSpeichern(database))
	mux.HandleFunc("GET /gruppe/verwerfung", handleVerwerfungAnzeigen(database))
	mux.HandleFunc("POST /gruppe/verwerfung", handleVerwerfungSpeichern(database))
	mux.HandleFunc("GET /gruppe/vorschau", handleVorschauAnzeigen(database))
	mux.HandleFunc("POST /gruppe/vorschau/einreichen", handleVorschauEinreichen(database))
	mux.HandleFunc("GET /gruppe/pruefung1", handlePruefung1Anzeigen(database))
	mux.HandleFunc("POST /gruppe/pruefung1", handlePruefung1Beantworten(database))
	mux.HandleFunc("GET /gruppe/pruefung2", handlePruefung2Anzeigen(database))
	mux.HandleFunc("POST /gruppe/pruefung2", handlePruefung2Beantworten(database))
	mux.HandleFunc("GET /gruppe/freigegeben", handleFreigegebenAnzeigen(database))
	mux.HandleFunc("GET /gruppe/live", handleGruppeLive(database))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderTemplate(w, indexTmpl, nil)
	})
	return mux
}

// renderTemplate rendert tmpl zuerst in einen Puffer, bevor es geschrieben
// wird — ein Renderfehler in der Mitte der Vorlage darf nicht mit einem
// bereits gesendeten 200 kollidieren.
func renderTemplate(w http.ResponseWriter, tmpl *template.Template, daten any) {
	var body bytes.Buffer
	if err := tmpl.Execute(&body, daten); err != nil {
		http.Error(w, "Vorlage konnte nicht gerendert werden", http.StatusInternalServerError)
		log.Printf("Vorlagenfehler: %v", err)
		return
	}
	w.Write(body.Bytes())
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
