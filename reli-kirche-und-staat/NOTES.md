# Arbeitsnotizen

## Präferenzen der Nutzerin
- Deutsch. Alles unterrichtsfertig ausformuliert, keine Stichpunktskizzen.
- Überleitungen **wörtlich** ausformulieren (Lehrersprache, zum Vorlesen/Anlehnen).
- Doppelstundenformat (90 min).
- Aktivierung vor Vollständigkeit: lieber weniger Stoff, dafür Beteiligung.

- **Duzen.** Die SuS werden geduzt — beim Ansprechen der ganzen Klasse "ihr",
  bei Aufgaben an Einzelne "du". Auch der Regie- und Hinweistext an die Lehrkraft
  wird geduzt (ausdrueckliche Erlaubnis). Ausnahme: Im Rollenspiel von DS 3 siezen
  die SuS die Rollen (Barth, Althaus, Meiser) — das bleibt so, es ist historisch
  und dramaturgisch richtig. Umgestellt wird mit `python3 duzen.py`.

## Technische Konventionen dieser Workspace
- `assets/style.css` ist die geteilte Gestaltungsschicht (Tokens, Typo, Komponenten).
- Lektionen in `lessons/` binden diese Tokens **inline** ein, weil sie als Artifact
  veröffentlicht werden und dort keine relativen Dateipfade auflösbar sind.
  `assets/style.css` ist damit die Quelle, aus der neue Lektionen kopieren.

## Pflicht bei jeder neuen Lektion
- **`<meta charset="utf-8">` als allererste Zeile.** Die Artifact-Plattform ergaenzt
  die Angabe beim Veroeffentlichen selbst, die lokale Datei aber nicht — ohne die
  Zeile raet der Browser beim Doppelklick auf Windows-1252 und alle Umlaute brechen.
- **Bilder als `data:`-URI einbetten, nie extern verlinken.** Die Artifact-CSP
  blockt Bilder von fremden Hosts ohne sichtbaren Fehler; eingebettete Bilder
  funktionieren zusaetzlich offline am Beamer.
- **Querlinks zwischen Lektionen als volle Artifact-URL**, nicht relativ —
  relative Pfade loesen im veroeffentlichten Artifact nicht auf.

## Offene Punkte / vor dem Unterricht prüfen
- Rechtlicher Status der AfD-Einstufung durch den Verfassungsschutz ändert sich laufend.
  In DS 6 ist ein Prüfhinweis eingebaut — Stand vor der Stunde aktualisieren.
- Beutelsbacher Konsens (Überwältigungsverbot, Kontroversitätsgebot) ist in DS 6
  didaktisch eingearbeitet; bei Elternnachfragen darauf verweisen.
