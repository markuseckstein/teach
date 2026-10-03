# Notizen

## Nutzer
- Kommandant einer kleinen FF in Bayern. Sprache: Deutsch.
- Ziel ist praktisch: eine Übung vorbereiten und leiten, nicht Theorie um ihrer selbst willen.
- Material: TSF mit Standardbeladung (Stromerzeuger, Teleskopscheinwerfer, Rosenbauer-Akkuleuchte), 6–8 Kräfte, Privat-Pkw möglich.
- Wunsch (28.09.2026): Die Übung soll Handfunkgeräte **inkl. Kanal-/Gruppenwechsel** einschließen.

## Recherche-Hinweise
- Die Feuerwehr-Lernbar ist eine Vue-SPA; WebFetch sieht nur „Feuerwehr Lernbar“. Die JSON-API funktioniert:
  `https://www.feuerwehr-lernbar.bayern/api/articles?search=<Begriff>` und `/api/articles/<lead>`; PDFs über `/api/media/<media-id>/raw?version=1`.
- Alte `fileadmin/...`-PDF-Links liefern die 404-Seite der SPA.

## Planung
- L1: Regeln Nachtlandeplatz + Ausleuchtung mit eigener Beladung + Funkgrundlagen zum RTH.
- L2 (nächste): Übungsdrehbuch für 6–8 Kräfte mit Funkstrang (DMO ↔ TMO-Übungsgruppe), Rollen, Zeitplan, Sicherheitsbelehrung.
- Später: Erkundung eigener Landeplätze (Checkliste für Tag-Erkundung), Nachbesprechung/Auswertung.
- Antworten (28.09.): Übung Do 1.10.2026, 20:00, Pattenhofen bei Burgthann · Sepura-HRT, Standardgruppe 324 (DMO? unbestätigt) · Sondergruppe bei ILS Nürnberg noch nicht angefragt · kein echter Hubschrauber.
- L2 erstellt + Drehbuch `reference/uebung-2026-10-01-drehbuch.html` (Funkplan A/B/C, Einspielungen E1–E6, Beobachtungsbogen).
- Nächste Sitzung: nach der Übung Erfahrungen abfragen → Learning Record schreiben (Evidenz!), Lektion 3: Auswertung + Einsatzvorplanung eigener Landeplätze (Tag-Erkundung, Karte für die ILS/Alarmordner).
- Sonnenuntergang 1.10. Pattenhofen ≈ 18:54 (selbst berechnet, NOAA-Algorithmus).
- TSF-Rufname (vom Nutzer bestätigt): Florian Pattenhofen 136/44/1.
- 324 = DMO-Gruppe (vom Nutzer bestätigt). Kursübersicht: index.html; Navigation zentral in assets/nav.js (COURSE-Liste).
- Briefing: reference/briefing-leitfaden.html (Übungsleiter, 15 min mit Fragen zuerst) + reference/handout-mannschaft.html (1 A4-Seite, 8× drucken).
- Funkplan (30.09., vom Nutzer): **keine Sondergruppe**; nur DMO. 324 = Einsatzstelle (alle), 320 = simulierte ILS (steht für TMO FW_LAU, nur zum Start: Alarm, Auftrag, Wechselaufforderung), 322 = Pilot („Christoph Übung“), dort auch Lagemeldung. Keine weiteren Zwischenrufe der „Leitstelle“. Unterlagen angepasst; Lektionstexte L1/L2 beschreiben die Sondergruppe noch generisch.
- Funkrollen-Ablauf: `reference/uebungsleiter-funk-cheatsheet.html`. Hubschrauber-Sound: `audio/` (make_heli.py) + `reference/heli-sound.html`.
