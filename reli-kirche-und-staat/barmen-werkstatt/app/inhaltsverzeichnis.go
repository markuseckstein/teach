package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// dsDateiMuster erkennt "ds" gefolgt von der Doppelstundennummer im
// Dateinamen (z. B. "0007-ds6-kopiervorlagen.html" → 6). Case-insensitiv, weil
// die Groß-/Kleinschreibung im Dateinamen keine verlässliche Konvention ist.
var dsDateiMuster = regexp.MustCompile(`(?i)ds(\d+)`)

// doppelstundeEintrag ist eine Zeile im Inhaltsverzeichnis: die
// Kopiervorlagen-Datei einer Doppelstunde und, wo vorhanden, der Einstieg in
// die Anwendung — die im Vorgang geforderte Naht "Inhaltsseiten und Anwendung
// nebeneinanderstellen".
type doppelstundeEintrag struct {
	Nummer    int
	Datei     string
	Anwendung string // Pfad in die Anwendung, leer wenn es keinen gibt
}

type inhalteAnsicht struct {
	Gesamtueberblick string // Dateiname des Gesamtüberblicks, leer wenn keiner gefunden wurde
	Doppelstunden    []doppelstundeEintrag
}

var inhalteTmpl = template.Must(template.ParseFS(templatesFS, "templates/inhalte.html"))

// ladeInhaltsverzeichnis liest lessonsDir und ordnet jede Datei, deren Name
// ein "dsN" enthält, der Doppelstunde N zu. Das geschieht über den Dateinamen
// statt über eine feste Liste, damit später hinzukommende Doppelstunden ohne
// Codeänderung im Verzeichnis auftauchen. Eine Datei ohne "dsN" (der
// Gesamtüberblick) wird separat vermerkt.
func ladeInhaltsverzeichnis(lessonsDir string) (inhalteAnsicht, error) {
	eintraege, err := os.ReadDir(lessonsDir)
	if err != nil {
		return inhalteAnsicht{}, err
	}

	var ansicht inhalteAnsicht
	for _, e := range eintraege {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		match := dsDateiMuster.FindStringSubmatch(name)
		if match == nil {
			if ansicht.Gesamtueberblick == "" {
				ansicht.Gesamtueberblick = name
			}
			continue
		}
		nummer, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}

		eintrag := doppelstundeEintrag{Nummer: nummer, Datei: name}
		if nummer == 6 {
			eintrag.Anwendung = "/beitreten"
		}
		ansicht.Doppelstunden = append(ansicht.Doppelstunden, eintrag)
	}

	sort.Slice(ansicht.Doppelstunden, func(i, j int) bool {
		return ansicht.Doppelstunden[i].Nummer < ansicht.Doppelstunden[j].Nummer
	})

	return ansicht, nil
}

// handleInhalteAnzeigen zeigt das Inhaltsverzeichnis: den Gesamtüberblick und
// je Doppelstunde die Kopiervorlagen-Datei, für Doppelstunde 6 zusätzlich den
// Einstieg in die Anwendung. Der Handler braucht keine Datenbank — er liest
// nur das lessons-Verzeichnis.
func handleInhalteAnzeigen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ansicht, err := ladeInhaltsverzeichnis(*lessonsDirFlag)
		if err != nil {
			http.Error(w, "Inhaltsverzeichnis konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("inhaltsverzeichnis laden: %v", err)
			return
		}
		renderTemplate(w, inhalteTmpl, ansicht)
	}
}
