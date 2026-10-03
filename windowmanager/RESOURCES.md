# Windows Multi-Monitor Desktop Management — Resources

## Knowledge

- [PowerToys Workspaces Utility for Windows Desktop Management — Microsoft Learn](https://learn.microsoft.com/en-us/windows/powertoys/workspaces)
  Offizielle Doku zum Workspaces-Modul: Setup erfassen, mit CLI-Argumenten starten, Desktop-Verknüpfung anlegen. Use for: das tägliche Neustart-Problem lösen (Lektion 1).
- [FancyZones Window Manager for Windows — Microsoft Learn](https://learn.microsoft.com/en-us/windows/powertoys/fancyzones)
  Vollständige Referenz: Zonen-Editor (Grid/Canvas), alle Tastenkombinationen, Multi-Monitor-Settings, FancyZones-CLI, bekannte Inkompatibilitäten (z. B. Terminal braucht Admin-Modus zum Snappen). Use for: Zonen-Layouts pro Monitor bauen, Shortcuts, Troubleshooting.
- [Microsoft PowerToys — Übersicht aller Utilities — Microsoft Learn](https://learn.microsoft.com/en-us/windows/powertoys/)
  Einstiegspunkt für alle PowerToys-Module (Workspaces, FancyZones, Awake, PowerToys Run, …). Use for: prüfen, ob es für ein Problem bereits ein passendes Modul gibt, bevor extern gesucht wird.
- [PowerToys Run Quick Launcher for Windows — Microsoft Learn](https://learn.microsoft.com/en-us/windows/powertoys/run)
  Vollständige Plugin-Referenz (Window Walker, Program, Shell, Calculator, Value Generator, …), Direct-Activation-Präfixe, Shortcuts. Use for: tastaturbasiertes App-/Fenster-Switching, Lektion 5. Hinweis im Dokument selbst: Nachfolger "Command Palette" ist in Entwicklung — bewusst noch nicht Teil der Mission (Out of scope, bis relevant).
- [Snap Layouts — Windows Learning Center — Microsoft](https://www.microsoft.com/en-us/windows/learning-center/organize-screen-with-snap-layouts)
  Offizielle Erklärung der Bordmittel-Snap-Funktion (Win+Z) als Windows-natives Pendant/Ergänzung zu FancyZones.
- Windows Virtuelle Desktops & Task View — Microsoft-Tastenkombinationen (Win+Ctrl+D neuer Desktop, Win+Ctrl+←/→ wechseln, Win+Ctrl+F4 schließen, Win+Tab Task View)
  Use for: Kontexte (Kommunikation/Frontend/Backend/Recherche) sauber trennen, ohne Dauer-Apps zu verlieren.

- [PowerToys GitHub Issue #38221 — Workspaces & multiple virtual desktops](https://github.com/microsoft/PowerToys/issues/38221) und [#35407](https://github.com/microsoft/PowerToys/issues/35407)
  Bestätigt: Workspaces unterstützt (Stand der Recherche) keine native Zuordnung von Layouts zu einzelnen virtuellen Desktops — offener Feature-Request, kein Bedienfehler. Use for: Lektion 4 (Workspaces + virtuelle Desktops kombinieren), bei PowerToys-Updates erneut prüfen.

## Wisdom (Communities)

- [r/PowerToys](https://reddit.com/r/PowerToys)
  Aktive Community rund um PowerToys-Module inkl. FancyZones/Workspaces-Tipps und Layout-Sharing (`custom-layouts.json` Austausch). Use for: fertige Layouts für ungewöhnliche Monitor-Kombinationen (z. B. Portrait+Landscape-Mischsetups) finden.
- [PowerToys GitHub Issues](https://github.com/microsoft/PowerToys/issues)
  Offizieller Issue-Tracker — dort stehen auch die dokumentierten App-Inkompatibilitäten (Terminal, Citrix, etc.) und Workarounds. Use for: konkrete Bugs/Grenzfälle mit einzelnen Apps im eigenen Setup.

## Gaps

- Noch keine hochwertige Quelle zu GitHub Copilot CLI + Terminal-Multiplexing-Workflows im Zusammenspiel mit FancyZones gefunden — bei Bedarf in einer späteren Session gezielt recherchieren.
- Unklar (nicht dokumentiert): ob "Show windows from this app on all desktops" einen Neustart übersteht — durch Ausprobieren im Alltag zu klären.

## Community-Präferenz

Noch nicht erfragt, ob Markus überhaupt an Community-Austausch (r/PowerToys etc.) interessiert ist — bei Gelegenheit nachfragen statt ungefragt vorauszusetzen.
