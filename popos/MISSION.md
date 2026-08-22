# Mission: PopOS Fenstermanagement für Entwickler

## Why
Markus arbeitet beruflich auf Windows, hat privat aber einen alten Laptop mit Pop!_OS (COSMIC Desktop) hergerichtet. Er möchte dort genauso flüssig Fenster anordnen, wechseln und über virtuelle Desktops verteilen können wie unter Windows (Snap-Layouts, virtuelle Desktops, Win+Tab) — und hat sich bewusst dafür entschieden, statt einer Windows-Kopie COSMICs natives Tiling-Fenstermanagement zu lernen, weil das für einen Entwickler potenziell effizienter ist.

## Success looks like
- Wechselt per Tastatur zwischen mehreren virtuellen Desktops (Workspaces) und verschiebt Fenster gezielt dorthin
- Tilt Fenster per Tastatur, ändert die Ausrichtung, wechselt zwischen Tiling und Floating
- Stapelt Fenster (Stacks) und wechselt zwischen gestapelten Fenstern
- Braucht für alltägliches Fenster-Arrangement kaum noch die Maus
- Weiß, wo in den COSMIC-Settings sich Shortcuts bei Bedarf anpassen lassen

## Constraints
- Lernt in intensiven Deep-Dive-Sessions, nicht in kleinen Häppchen über Wochen verteilt
- Kommt aus der Windows-Welt (Snap-Layouts, Win+Tab, Win+Ctrl+Pfeil für virtuelle Desktops) — Analogien zu Windows helfen beim Einordnen neuer Konzepte
- System: Pop!_OS 24.04 LTS mit COSMIC Desktop (Wayland). Kein GNOME, kein Pop Shell (altes System)

## Out of scope
- Klassisches GNOME-basiertes Pop!_OS / die alte Pop-Shell-Extension — nicht relevant, Markus läuft auf COSMIC
- Ein separater Tiling-WM wie i3 oder sway — COSMIC bietet Tiling nativ, kein Grund zu wechseln
