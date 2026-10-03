# Mission: Effizientes Multi-Monitor Desktop Management unter Windows

## Why
Markus ist Software-Entwickler und arbeitet an 4 Monitoren (Portrait 4K, Landscape 4K, Laptop-Display, Landscape 4K) mit vielen parallelen Kontexten: Dauer-Apps (Outlook, Teams), Web-Recherche, 1-3 Terminals mit GitHub Copilot CLI, mehrere VS Code-Instanzen (Frontend/Backend). Er bootet täglich neu und muss jeden Morgen manuell alle Fenster wieder an ihren Platz ziehen — das kostet Zeit und reißt ihn aus dem Fokus, bevor der Tag überhaupt angefangen hat.

## Success looks like
- Fenster lassen sich ohne Maus per Tastenkombination zwischen Zonen/Monitoren verschieben und in der Größe anpassen
- Feste FancyZones-Layouts sind für alle 4 Monitore passend zur jeweiligen Ausrichtung (Portrait/Landscape) eingerichtet und per Shortcut abrufbar
- Nach jedem Neustart stellt ein Klick/Shortcut das komplette Arbeits-Setup wieder her (Apps offen, an der richtigen Position) — kein manuelles Neuanordnen mehr
- Virtuelle Desktops trennen Kontexte sinnvoll (z. B. Kommunikation, Frontend, Backend, Recherche), ohne dass Dauer-Apps wie Outlook/Teams dabei verloren gehen
- Der komplette Workflow (Zonen wählen, zwischen Desktops wechseln, Setup wiederherstellen) sitzt als Muskelgedächtnis

## Constraints
- Windows 11, PowerToys inkl. FancyZones ist bereits installiert
- Rechner wird jeden Tag neu gestartet — jede Lösung muss diesen Reset überleben bzw. ihn automatisieren, nicht nur einmalig funktionieren
- Wenig Zeit für Setup-Bastelei; will pragmatische, sofort wirksame Verbesserungen, keine Tiefenrecherche in PowerToys-Interna
- Deutsch als Unterrichtssprache

## Out of scope
- Drittanbieter-Tools jenseits von PowerToys/Windows-Bordmitteln (außer sie sind klar die beste Lösung für ein konkretes Problem)
- Linux/macOS Window-Manager-Vergleiche
- Tiefes Customizing von PowerToys-Quellcode/Plugins
