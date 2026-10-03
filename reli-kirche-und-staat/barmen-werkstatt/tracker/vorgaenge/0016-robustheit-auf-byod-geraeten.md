---
id: 0016
titel: Robustheit auf BYOD-Geräten
status: verworfen
label: blockiert
haengt-an: ["0010"]
---

# 0016 — Robustheit auf BYOD-Geräten

## Worum es geht

Dreißig Jugendliche bringen dreißig verschiedene Android-Geräte mit, teils alt. Das
größte Risiko des Projekts ist nicht die Fachlogik, sondern die eine Gruppe, bei der es
nicht läuft.

## Fertig, wenn

- Moderne CSS- und HTML5-Mittel (Container Queries, `:has()`, `dialog`, View
  Transitions) werden verwendet, aber als progressive enhancement: Ohne sie bleibt die
  Werkstatt vollständig bedienbar.
- Kein Ablauf hängt an JavaScript allein; jeder Schritt ist auch als normales Formular
  absendbar.
- Bedienelemente sind mit dem Finger auf kleinen Bildschirmen erreichbar.
- Getestet auf mindestens einem bewusst alten Browserprofil.

## Hinweise

Der erste echte Klassendurchlauf wird trotzdem etwas zutage fördern. Diese Aufgabe
senkt die Wahrscheinlichkeit, sie beseitigt sie nicht.

## Verworfen (2026-09-06)

Won't fix — die tatsächliche Geräteflotte der Klasse besteht durchweg aus solider
Samsung-Mittelklasse, kein Gerät älter als drei Jahre. Das im "Worum es geht"
unterstellte Risiko (dreißig verschiedene, teils alte Geräte) trifft hier nicht zu; ein
gezielter Alt-Geräte-Test würde ein Problem absichern, das nicht besteht. Die generelle
Progressive-Enhancement-Haltung aus SPEZIFIKATION.md (kein Ablauf hängt an JavaScript
allein) ist ohnehin durchgehend im Code umgesetzt — das war kein separater Aufwand,
sondern der Standardweg jedes Formulars in dieser Anwendung.
