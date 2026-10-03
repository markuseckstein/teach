# Notizen

## Präferenzen des Nutzers

- **Sprache: Deutsch.** Lektionen, Referenzen, Glossar — alles auf Deutsch. Englische Fachbegriffe in Klammern dazu.
- **Bilder und Diagramme sind ausdrücklich erwünscht.** Wo ein Sachverhalt ein Bild verträgt, gehört ein Diagramm hin (inline SVG, damit es offline und im Druck funktioniert).
- **Videos verlinken**, wo sie wirklich helfen — nicht pflichtschuldig.
- **Inhaltsverzeichnis + Navigation:** Es gibt eine zentrale `index.html`. Die Nav-Links in jeder Lektion (zurück / weiter / Sprungliste) müssen bei jeder Änderung mitgezogen werden.

## Architekturentscheidung: eine Quelle für die Navigation

Alle Lektionen, Referenzen und deren Reihenfolge stehen **ausschließlich** in `assets/course.js` in der Konstante `COURSE`.
`assets/lesson.js` rendert daraus zur Laufzeit die Kopf- und Fußnavigation jeder Lektion, `index.html` das Inhaltsverzeichnis.

**Konsequenz für künftige Sessions:** Wer eine Lektion hinzufügt, ändert genau zwei Dinge — die neue HTML-Datei anlegen und einen Eintrag in `COURSE.lessons` ergänzen. Nav-Links in bestehenden Lektionen werden nie von Hand angefasst. Sie stimmen dann automatisch.

Technische Randbedingung: Die Seiten laufen per Doppelklick über `file://`. Deshalb **keine** ES-Module und **kein** `fetch()` — beides blockiert der Browser bei lokalen Dateien. Alles läuft über klassische `<script src>`-Tags mit globalen Variablen.

## Zielgruppen-Beobachtungen

- Zielgruppe sind Lehrkräfte. Beispiele deshalb konsequent aus dem Schulalltag: Elternbriefe, Klassenarbeiten, Differenzierung, Förderpläne, Konferenzprotokolle, Vertretungsplan.
- Heikler Punkt bei Lehrkräften: personenbezogene Schülerdaten. Das muss früh und unmissverständlich kommen, sonst ist alles andere fahrlässig.
- Erwartbarer Widerstand: „Das dauert ja länger, als es selbst zu machen." Das ist ein *berechtigter* Einwand bei reiner Chatbot-Nutzung — er wird nicht weggeredet, sondern in Lektion 9 zum Angelpunkt gemacht.

## Zu tun / offene Fäden

- Nach dem Workshop: Rückmeldung einsammeln, welche Analogie gezündet hat und welche nicht. Das ist die wertvollste Information für Version 2.
- Prüfen, ob eine der Lehrkräfte an einem Dienstrechner sitzt, auf dem sich lokal nichts installieren lässt — das ändert die Empfehlung in Lektion 8.
