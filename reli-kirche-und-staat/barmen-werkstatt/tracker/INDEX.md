# Übersicht

Erzeugt mit `python3 tracker/index.py` — nicht von Hand ändern.

| | Nr. | Vorgang | Status | Label | hängt an |
|---|---|---|---|---|---|
| ● | 0001 | [Projektgerüst und HTTP-Server](vorgaenge/0001-projektgeruest-und-http-server.md) | erledigt | `ready-for-agent` | — |
| ● | 0002 | [Datenbankschema](vorgaenge/0002-datenbankschema.md) | erledigt | `ready-for-agent` | — |
| ● | 0003 | [Testnaht über HTTP](vorgaenge/0003-testnaht-ueber-http.md) | erledigt | `ready-for-agent` | 0001, 0002 |
| ● | 0004 | [Kurs, Gruppen und Themenfelder anlegen](vorgaenge/0004-kurs-gruppen-und-themenfelder-anlegen.md) | erledigt | `ready-for-agent` | 0002, 0003 |
| ● | 0005 | [Beitritt per Code, Gerätetoken, Wiedereinstieg](vorgaenge/0005-beitritt-per-code-geraetetoken-wiederein.md) | erledigt | `ready-for-agent` | 0004 |
| ● | 0006 | [Schreibrecht innerhalb der Gruppe](vorgaenge/0006-schreibrecht-innerhalb-der-gruppe.md) | erledigt | `blockiert` | 0005 |
| ● | 0007 | [Themenfelder und Bibelstellen als Daten](vorgaenge/0007-themenfelder-und-bibelstellen-als-daten.md) | erledigt | `inhalt` | 0002 |
| ● | 0008 | [Zustandsautomat der Werkstatt](vorgaenge/0008-zustandsautomat-der-werkstatt.md) | erledigt | `blockiert` | 0006, 0007 |
| ○ | 0009 | [Neutralitätstest](vorgaenge/0009-neutralitaetstest.md) | offen | `blockiert` | 0008 |
| ○ | 0010 | [Live-Kanal über SSE](vorgaenge/0010-live-kanal-ueber-sse.md) | offen | `blockiert` | 0006 |
| ○ | 0011 | [Regiepult](vorgaenge/0011-regiepult.md) | offen | `blockiert` | 0008, 0010 |
| ○ | 0012 | [Beamer-Ansicht](vorgaenge/0012-beamer-ansicht.md) | offen | `blockiert` | 0008, 0010 |
| ○ | 0013 | [Druckansicht der freigegebenen These](vorgaenge/0013-druckansicht-der-freigegebenen-these.md) | offen | `blockiert` | 0008 |
| ○ | 0014 | [Freiwillige Unterschriften](vorgaenge/0014-freiwillige-unterschriften.md) | offen | `blockiert` | 0008 |
| ○ | 0015 | [Kurs beenden und Daten löschen](vorgaenge/0015-kurs-beenden-und-daten-loeschen.md) | offen | `blockiert` | 0004 |
| ○ | 0016 | [Robustheit auf BYOD-Geräten](vorgaenge/0016-robustheit-auf-byod-geraeten.md) | offen | `blockiert` | 0010 |
| ○ | 0017 | [Statische Inhalte mit ausliefern](vorgaenge/0017-statische-inhalte-mit-ausliefern.md) | offen | `blockiert` | 0001 |
| ○ | 0018 | [Betrieb auf dem VPS](vorgaenge/0018-betrieb-auf-dem-vps.md) | offen | `betrieb` | 0001 |
| ○ | 0019 | [Zeitkapsel: Fotos aus DS 4 und DS 5](vorgaenge/0019-zeitkapsel-fotos-aus-ds-4-und-ds-5.md) | offen | `spaeter` | 0004 |
| ○ | 0020 | [Nacharbeit-Seiten digital](vorgaenge/0020-nacharbeit-seiten-digital.md) | offen | `spaeter` | 0017 |

12 offen · 0 in Arbeit · 8 erledigt
