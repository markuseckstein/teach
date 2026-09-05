---
id: 0002
titel: Datenbankschema
status: erledigt
label: ready-for-agent
haengt-an: []
---

# 0002 — Datenbankschema

## Worum es geht

Das Datenmodell aus der Spezifikation muss als SQLite-Schema existieren, damit alle
folgenden Vorgänge darauf aufsetzen können.

## Fertig, wenn

- SQLite über einen reinen Go-Treiber angebunden ist (kein cgo).
- Die Tabellen `kurs`, `gruppe`, `geraet`, `these`, `kapsel` existieren, wie in
  `SPEZIFIKATION.md` beschrieben.
- Das Schema wird beim Start angelegt, falls die Datei fehlt; ein zweiter Start
  verändert nichts.
- Der Pfad zur Datenbankdatei ist konfigurierbar, damit Tests ein temporäres
  Verzeichnis benutzen können.
- Fremdschlüssel sind aktiviert, und das Löschen eines Kurses löscht alles daran
  Hängende (nötig für Vorgang 0015).

## Gehört nicht dazu

Ein Migrationswerkzeug. Solange das Projekt keine Produktivdaten hat, ist "Schema beim
Start anlegen" ausreichend; Migrationen kommen erst, wenn ein Kurs eine echte Einheit
überlebt hat.
