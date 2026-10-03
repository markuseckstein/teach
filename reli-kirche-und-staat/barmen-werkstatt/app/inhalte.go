package main

import (
	"encoding/json"
	"log"
)

// Bibelstelle ist ein Eintrag im Bibelstellen-Vorrat eines Themenfelds —
// Angabe und Volltext zusammen, damit eine Gruppe nie nach einer Stelle
// formuliert, die sie nicht gelesen hat.
type Bibelstelle struct {
	Angabe string `json:"angabe"`
	Text   string `json:"text"`
}

type themenfeldEintrag struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Bibelstellen []Bibelstelle `json:"bibelstellen"`
}

type themenfelderDatei struct {
	Uebersetzung string              `json:"uebersetzung"`
	Themenfelder []themenfeldEintrag `json:"themenfelder"`
}

var themenfelderInhalte = ladeThemenfelderInhalte()

func ladeThemenfelderInhalte() themenfelderDatei {
	rohdaten, err := dataFS.ReadFile("data/themenfelder.json")
	if err != nil {
		log.Fatalf("data/themenfelder.json konnte nicht gelesen werden: %v", err)
	}
	var daten themenfelderDatei
	if err := json.Unmarshal(rohdaten, &daten); err != nil {
		log.Fatalf("data/themenfelder.json ist kein gültiges JSON: %v", err)
	}
	return daten
}

// bibelstellenVorrat gibt den Bibelstellen-Vorrat des benannten Themenfelds
// zurück. ok ist false, wenn der Name kein bekanntes Themenfeld ist (etwa
// weil die Lehrkraft noch keins zugewiesen hat).
func bibelstellenVorrat(themenfeldName string) ([]Bibelstelle, bool) {
	for _, tf := range themenfelderInhalte.Themenfelder {
		if tf.Name == themenfeldName {
			return tf.Bibelstellen, true
		}
	}
	return nil, false
}
