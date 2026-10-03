# Notes

## Setup
- 4 Monitore, von links nach rechts: Portrait 4K, Landscape 4K, Laptop-Display, Landscape 4K
- Windows 11, PowerToys + FancyZones bereits installiert
- Bootet **jeden Tag neu** → jede Lösung muss den Reset überleben, nicht nur einmalig funktionieren

## Typischer App-Mix
- Dauer-Apps (müssen immer offen/erreichbar sein): Outlook, Teams
- 1-3 Terminals mit GitHub Copilot CLI Sessions
- Mehrere VS Code Instanzen (getrennt: Frontend, Backend)
- Viel Web-Recherche im Browser

## Präferenzen
- Unterrichtssprache: Deutsch
- Will pragmatische, sofort wirksame Verbesserungen — keine Tiefenrecherche in PowerToys-Interna
- Noch nicht erfragt: Interesse an Community-Austausch (r/PowerToys etc.)

## Format-Präferenzen
- Jede Lektion und jedes Referenzdokument bekommt ein Inhaltsverzeichnis (nav.toc, siehe assets/style.css) direkt unter der Lede
- Footer-Navlinks (.lesson-nav) müssen bei jedem neuen Dokument passend verlinkt/aktualisiert werden (nicht nur beim jeweils neuesten)
- Workspace braucht von Anfang an ein `index.html` als Übersichtsseite (Liste aller Lektionen + Referenz + Mission/Resources) — nicht erst am Ende nachreichen. Bei neuem Teach-Workspace gleich mit anlegen, sobald die erste Lektion existiert.

## Reihenfolge-Entscheidung Session 1
Größter genannter Schmerzpunkt: täglicher manueller Neuaufbau des Fenster-Layouts nach dem Reboot.
→ Lektion 1 startet daher nicht bei FancyZones-Grundlagen, sondern direkt bei PowerToys Workspaces
  (löst das konkrete "jeden Morgen neu anordnen"-Problem als sofortiger Quick-Win).
FancyZones-Zonenlayout-Bau und virtuelle Desktops folgen als Lektion 2+, da Workspaces darauf aufbaut
(Workspaces nutzt intern die FancyZones-Engine zur Positionierung).
