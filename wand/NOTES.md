# Notes

- Nutzer kommuniziert auf Deutsch — Lektionen, Referenzdokumente und Ansprache auf Deutsch verfassen.
- Baseline: solide allgemeine Heimwerker-Erfahrung, aber keine bisherige Erfahrung mit Verputzen.
- Motivation ist Kostenersparnis (Eigenleistung statt Fachbetrieb) — Lektionen dürfen pragmatisch/effizient bleiben, keine Edeloptik-Verputztechniken.
- Objekt: Nebengebäude (Garage/Schuppen), kein Denkmalschutz, ~5–20 m² Fläche.
- Zeitdruck: soll diese Saison fertig werden (2026), Wetterfenster beachten.
- Vor dem eigentlichen Verputzen entfernt der Nutzer einen großen, kompletten Abschnitt der als nicht-tragend eingeschätzten Wand (Durchbruch) — hat einen großen Bosch SDS-max-Bohrhammer. Deshalb Mission um diesen Vorbereitungsschritt erweitert (Lektion 2).
- Lektions-Template steht: jede Lektion bekommt ein Inhaltsverzeichnis (`nav.toc`) direkt nach dem Header sowie Footer-NavLinks zur vorherigen/nächsten Lektion bzw. Referenzdokument — an diesem Muster für alle künftigen Lektionen festhalten.
- `index.html` (Workspace-Root) ist die Gesamtübersicht/Lernpfad-Seite mit Inhaltsverzeichnis über alle Lektionen + Referenzdokumente. Der Kicker jeder Lektion/Referenz verlinkt "Wand verputzen" zurück auf `../index.html` (CSS: `.index-list`, `.index-item`, `.index-refs`, `.kicker a` in `assets/style.css`). Neue Lektionen/Referenzen müssen in `index.html` ergänzt werden.
- Nutzer möchte zügig durch die Lektionen iterieren ("mach dann auch gleich weiter mit der nächsten Lektion") — nicht nach jeder Lektion auf Bestätigung warten, sondern direkt mit dem nächsten sinnvollen Schritt in der Projekt-Reihenfolge weitermachen.
