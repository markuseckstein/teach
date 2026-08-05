# Der Windows-Rechner ist reine Anzeige — kein PAWinS, kein Decoder

Markus hat am alten Rechner nachgesehen: **weder PAWinS noch FMS32, POC32 oder
monitord** laufen dort. Der Rechner macht ausschließlich Anzeige. Damit kommt der
Alarm bereits fertig bei poweralarm.de an (Leitstellen-Schnittstelle wie katsys /
EDI / eMID, oder per Mail/Fax/SMS) und wird von dort nur abgeholt.

**Warum das wichtig ist:** Der PA Web-Monitor ist damit ein **vollständiger**
Ersatz für die Rolle dieses Rechners, nicht nur ein Teilersatz. Es gibt keinen
zweiten Job, der an Windows hängt. Der Umstieg ist eine reine Anzeigefrage —
Windows kann restlos verschwinden.

**Implications:**

- Der Pfad „PAWinS ohne Windows-Arbeitsplatz" ist gestrichen. Nicht mehr planen.
- Die einzige verbleibende Unbekannte am Web-Monitor selbst ist der
  **Alarmton** (und der genaue Funktionsumfang) — offen bei FITT.
- Der Kiosk-Aufbau ist **unabhängig** von der FITT-Antwort und kann sofort
  gelernt und geübt werden: das Gerät kann mit einer beliebigen Test-URL
  aufgesetzt und in Betrieb genommen werden, lange bevor die Lizenz da ist.
  Genau das ist Lektion 2 — sie entkoppelt Markus vom Support-Ticket.
- Reihenfolge damit: erst Kiosk-Hardware lauffähig (übbar, ohne Risiko), dann
  Web-Monitor-URL einsetzen, dann Parallelbetrieb, dann alten Rechner abschalten.
