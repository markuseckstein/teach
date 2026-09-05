---
id: 0016
titel: Robustheit auf BYOD-Geräten
status: offen
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
