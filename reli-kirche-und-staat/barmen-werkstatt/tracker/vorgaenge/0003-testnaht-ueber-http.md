---
id: 0003
titel: Testnaht über HTTP
status: offen
label: ready-for-agent
haengt-an: ["0001", "0002"]
---

# 0003 — Testnaht über HTTP

## Worum es geht

Die Spezifikation legt genau einen Testzugang fest: HTTP. Damit die folgenden Vorgänge
testgetrieben entstehen können, muss dieser Zugang zuerst existieren.

## Fertig, wenn

- Ein Testhelfer startet die vollständige Anwendung mit echtem Router und echter
  SQLite-Datei in einem temporären Verzeichnis.
- Es gibt Helfer für die wiederkehrenden Schritte: Kurs anlegen, Gerät beitreten
  lassen, als bestimmtes Gerät eine Anfrage stellen (Cookie mitführen).
- Ein erster Test prüft nur, dass die Anwendung antwortet — er ist das Muster für alle
  weiteren.
- `go test ./...` läuft ohne Netzwerkzugriff und ohne Aufräumarbeit von Hand.

## Gehört nicht dazu

Testbibliotheken. Standardbibliothek und `net/http/httptest`, sonst nichts. Keine
Mocks: Es wird die echte Anwendung getestet, nur ihre Datei liegt woanders.

## Hinweise

Testnamen als deutsche Verhaltenssätze, passend zum Vokabular des Projekts —
`TestGruppeKommtOhneBibelstelleNichtWeiter`, nicht `TestHandler_Post_400`.
