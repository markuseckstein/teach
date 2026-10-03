---
id: 0009
titel: Neutralitätstest
status: erledigt
label: blockiert
haengt-an: ["0008"]
---

# 0009 — Neutralitätstest

## Worum es geht

Die Zusage „die Anwendung wertet nicht" muss überprüfbar sein, nicht nur behauptet.
Dieser Test ist die ausführbare Fassung der didaktischen Absicherung und die Antwort
auf eine mögliche Elternnachfrage.

## Fertig, wenn

- Ein Test führt zwei inhaltlich völlig verschiedene Thesen durch den Automaten, mit
  identischem Antwortverhalten in den Prüffragen, und stellt fest: gleiches Ergebnis,
  gleicher Zustand, gleiche Zahl an Überarbeitungen.
- Eine der beiden Thesen enthält absichtlich Formulierungen, die ein naiver Filter
  beanstanden würde.
- Im Test steht als Kommentar, warum er existiert und dass er nicht gelöscht werden
  darf.

## Gehört nicht dazu

Nichts. Dieser Vorgang ist bewusst klein und darf nicht mit anderem vermischt werden.
