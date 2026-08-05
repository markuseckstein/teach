# Mission: PowerAlarm im Feuerwehrhaus

## Why
Der Alarmmonitor im Feuerwehrhaus läuft auf einem alten Windows-10-Rechner, der
nicht mehr auf Windows 11 aktualisiert werden kann und damit keine
Sicherheitsupdates mehr bekommt. Die Alarmanzeige muss auf möglichst
leichtgewichtige, wartungsarme Hardware umziehen — idealerweise Raspberry Pi oder
ein sparsamer Linux-Rechner mit Browser im Kiosk-Modus — ohne dass die
Einsatzmannschaft an Information verliert.

## Success looks like
- Der Alarmmonitor im Feuerwehrhaus läuft ohne Windows-Rechner.
- Markus kennt die PowerAlarm-Bausteine gut genug, um mit FITT gezielt über
  Lizenzen, Freischaltungen und Funktionsumfang zu sprechen.
- Die Web-Monitor-URL und die Zugangsdaten sind eingerichtet und dokumentiert.
- Ein Gerät startet von selbst in die Alarmanzeige — Vollbild, ohne Bedienung,
  überlebt Stromausfall und Neustart.
- Klar dokumentiert ist, welche Funktionen des Windows-Monitors ersetzt sind und
  welche bewusst entfallen (oder anders gelöst wurden).

## Constraints
- Kleines Feuerwehrhaus: kleines Budget, keine IT-Abteilung, Wartung nebenbei.
- Einsatzkritisch — Ausfälle sind nicht akzeptabel. Umstellung muss testbar sein,
  bevor der alte Rechner abgeschaltet wird.
- Manche Antworten liegen nur bei FITT GmbH (Lizenzumfang, Funktionen des
  Web-Monitors). Lernen heißt hier auch: die richtigen Fragen stellen können.
- Sprache: Deutsch.

## Out of scope
- Wechsel zu einer anderen Alarmierungsplattform (Divera, Alamos, FE2).
- Eigenentwicklung eines Alarmmonitors gegen die PowerAlarm-API.
- Umbau der Alarmierungskette selbst (Leitstelle, Sirenen, Melder).
