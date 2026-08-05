# Notes

## Präferenzen
- Sprache: **Deutsch** (Lektionen, Referenzen, Gespräch).
- Fragt sehr konkret und lösungsorientiert ("was ist nötig", "wie komme ich dran")
  — Lektionen also immer mit einem umsetzbaren Schritt beenden, nicht mit Theorie.

## Ausgangslage (Stand 2026-07-27)
- Kleines Feuerwehrhaus, PowerAlarm ist gesetzt.
- PA-Monitor (Windows-Version) läuft aktuell auf altem Rechner mit Windows 10,
  kein Upgrade-Pfad auf Windows 11.
- Ziel: leichtgewichtiger Ersatz, Browser-Kiosk auf Linux / Raspberry Pi.

## Offene Fragen an den Nutzer
1. ~~Läuft auf dem Windows-Rechner auch PAWinS?~~ **Beantwortet 27.07.2026: nein.**
   Nichts von PAWinS / FMS32 / POC32 / monitord läuft dort — der Rechner ist reine
   Anzeige. Windows kann restlos weg. Siehe [[0002-windows-rechner-ist-reine-anzeige]].
2. Braucht der Monitor im Feuerwehrhaus einen **Alarmton**? (Beim Windows-Monitor
   vorhanden; beim Web-Monitor nicht dokumentiert → ggf. Show-Stopper.)
3. Welche Funktionen des Windows-Monitors sind aktuell wirklich aktiv? (Karte,
   Alarmfax, Ausdruck, Bildschirm-Standby, HTTP-Weiterleitung/Türsteuerung)
   → Bestimmt, wie schwer die Lücken der Web-Version wirklich wiegen.
4. Ist die Anzeige im Feuerwehrhaus ein **PC-Monitor oder ein Fernseher**?
   Entscheidet zwischen `wlopm` und HDMI-CEC.
5. Übt er auf dem alten PC mit Linux oder gleich auf einem Pi? (Beides gelehrt,
   Konfiguration identisch.)

## Gelehrt bisher
- **L1** — PowerAlarm-Bausteine (Portal / PAWinS / PA-Monitor / PA Web-Monitor),
  Lizenz-Tor über API-Key-Option, Weg zur Web-Monitor-URL, Funktionslücken der
  Web-Version, Mailvorlage an FITT.
- **L2** — Kiosk-Betrieb: Wayland/labwc statt X11 (die zentrale Falle), Autologin,
  Blanking, Autostart, die vier einsatzkritischen Chromium-Schalter,
  `wlopm` vs. `cec-client`, Abnahme per Stecker-Test.

## Didaktisches Muster, das hier trägt
Immer erst die **Lizenz-/Herstellerfrage** von der **selbst lösbaren Arbeit**
trennen. Markus soll nie auf ein Support-Ticket warten müssen, um weiterzulernen —
L2 wurde bewusst so gebaut, dass sie ohne FITT-Antwort komplett durchführbar ist
(Test-URL statt Web-Monitor-URL).

## Offene Fragen an FITT
- Ist die Lizenz/Option „PA Web-Monitor“ für unseren Account freigeschaltet, und
  was kostet sie?
- Welche Anzeigefelder und Optionen hat der Web-Monitor konkret? Insbesondere:
  Alarmton, Alarmfax, Karte, Uhrzeit?
- Gibt es eine empfohlene/unterstützte Browser-Konfiguration für Kiosk-Betrieb?
- Bleibt die Sitzung dauerhaft angemeldet, oder läuft das Login aus?
