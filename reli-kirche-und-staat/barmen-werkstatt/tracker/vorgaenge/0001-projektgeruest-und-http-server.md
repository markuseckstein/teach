---
id: 0001
titel: Projektgerüst und HTTP-Server
status: erledigt
label: ready-for-agent
haengt-an: []
---

# 0001 — Projektgerüst und HTTP-Server

## Worum es geht

Es gibt noch nichts. Bevor irgendein Verhalten entstehen kann, braucht es ein
lauffähiges Go-Binary, das auf einem Port lauscht und eine Seite ausliefert.

## Fertig, wenn

- `go run .` startet einen HTTP-Server auf einem konfigurierbaren Port.
- Statische Dateien (CSS, wenige Bilder) sind über `embed.FS` ins Binary eingebettet;
  ein gebautes Binary läuft ohne Begleitdateien.
- Vorlagen werden mit `html/template` gerendert; eine Startseite zeigt, dass es lebt.
- `go build` erzeugt ein einzelnes Binary; Cross-Compile nach linux/amd64 funktioniert
  ohne cgo.
- `go vet` ist sauber.

## Gehört nicht dazu

Datenbank, Fachlogik, Gestaltung. Diese Aufgabe endet, sobald ein leeres, aber echtes
Gerüst steht.
