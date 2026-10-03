---
id: 0018
titel: Betrieb auf dem VPS
status: offen
label: betrieb
haengt-an: ["0001"]
---

# 0018 — Betrieb auf dem VPS

## Worum es geht

Die Anwendung muss laufen, ohne dass jemand vor der Stunde etwas startet.

## Fertig, wenn

- systemd-Unit, die das Binary betreibt und nach einem Neustart selbst hochkommt.
- Caddy davor für TLS, feste Adresse.
- Sicherung der SQLite-Datei ist dokumentiert und besteht aus einem `cp`.
- Eine knappe Betriebsnotiz beschreibt: neues Binary einspielen, Dienst neu starten,
  Datei sichern, im Notfall zurückrollen.

## Gehört nicht dazu

Container, Orchestrierung, Überwachung. Ein Prozess, eine Datei, ein Dienst.
