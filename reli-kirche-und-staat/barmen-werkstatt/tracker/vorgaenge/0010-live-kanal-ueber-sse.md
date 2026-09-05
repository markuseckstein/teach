---
id: 0010
titel: Live-Kanal über SSE
status: offen
label: blockiert
haengt-an: ["0006"]
---

# 0010 — Live-Kanal über SSE

## Worum es geht

Geräte einer Gruppe sollen sehen, was das Schreibgerät tippt; alle Geräte sollen
Phasenwechsel und Timer mitbekommen. Die Gegenrichtung sind normale Formular-POSTs.

## Fertig, wenn

- Ein SSE-Kanal überträgt Änderungen an die Geräte einer Gruppe sowie kursweite
  Ereignisse (Phase, Timer).
- Der Kanal überlebt einen kurzen Verbindungsabriss und stellt den aktuellen Stand
  wieder her — kein halb aktueller Bildschirm.
- Dreißig gleichzeitige Verbindungen sind getestet.
- Ohne SSE bleibt die Anwendung bedienbar: Ein Neuladen zeigt immer den korrekten
  Stand.

## Gehört nicht dazu

WebSocket. Es gibt keinen Kanal, der bidirektional sein müsste.
