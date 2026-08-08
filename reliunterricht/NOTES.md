# Notizen

## Nutzerpräferenzen
- Alles auf Deutsch — auch Dateiinhalte, Kommentare, Materialien.
- Kontroversität ausdrücklich erwünscht („darf ruhig kontrovers sein").
  Rahmen: Beutelsbacher Konsens — kontrovers ja, überwältigend nein.
- Interaktive Elemente erwünscht. Verfügbar: Beamer/Smartboard + Schüler-Smartphones.
  Keine Schüler-PCs → interaktive HTML läuft frontal; Smartphone-Einsatz über
  QR-Codes / externe Umfragetools (z.B. Oncoo, Mentimeter) einplanen.
- Zielgruppe: 7. Klasse (12–13 Jahre). Altersgerecht: keine grausamen Details,
  fiktive bzw. entschärfte Fallbeispiele.
- **Stehende Anweisung (18.07.2026): NavLinks immer aktualisieren.** Jede neue
  Lektion bekommt die `.navleiste` oben; bei jeder neuen Datei werden alle
  bestehenden Dateien (Lektionen + Referenzdokumente) nachgezogen: Navleiste,
  Fußzeilen-Links, „✓"-Links im Einheitsplan.

## Arbeitsstand (18.07.2026)
- Einheit auf 5 Stunden angelegt, siehe `reference/einheitsplan.html`.
- Stunde 1 fertig: `lessons/0001-stunde-1-einstieg.html`.
- Stunde 2 fertig: `lessons/0002-stunde-2-weltlage.html` — QR-Codes für die
  4 Stationen sind als data-URIs eingebettet (offline-fähig, via api.qrserver.com
  erzeugt). Bei Änderung einer Stations-URL muss der QR-Code neu erzeugt werden!
  **Arbeitsregel seit 18.07.2026:** Eingebettete Base64-Daten nie von Hand
  übertragen, sondern skriptgesteuert einsetzen und per `base64 -d | cmp` gegen
  die Originaldatei verifizieren — in Station D hatte ein einziges verdrehtes
  Zeichen den QR-Code unbrauchbar gemacht (von Markus entdeckt, behoben).
- Fallbeispiel Justizirrtum: Ray Krone (USA, 2002 durch DNA entlastet,
  100. Todestrakt-Entlasteter) — bewusst ohne Tatdetails erzählt.
- Stunde 3 fertig: `lessons/0003-stunde-3-religionen.html` — Expertengruppen
  mit Original-Quellenkarten (Judentum/kath./ev./Islam), jede Karte enthält
  bewusst eine *innere* Spannung der Tradition. Gruppen werden GELOST, nie nach
  Herkunft/Religion besetzt. EKD-Quellenlücke geschlossen (Faktenblatt aktualisiert).
- Didaktik-Muster etabliert: Jede Stunde beginnt künftig mit einer kurzen
  Blitz-Wiederholung (Abruf-Fragen zu früheren Stunden, Aufdeck-Format).
- Stunde 4 fertig: `lessons/0004-stunde-4-debatte.html` + Referenzdokument
  `reference/rollenkarten-debatte.html` (8 Rollenkarten P1–P4/C1–C4, Moderation,
  Beobachtungsbogen; druckbar). Szenario: fiktives Land „Silvanien", damit keine
  realen Staaten/Personen getroffen werden. Rolle P3 (Opferangehörige) vor dem
  Losen auf reale Verlusterfahrungen in der Klasse prüfen — ggf. weglassen.
  Neue wiederverwendbare Komponente: Beamer-Timer in `assets/interaktiv.js`
  (auch für Schreibzeit in Stunde 5 nutzbar).
- Stunde 5 fertig: `lessons/0005-stunde-5-urteil.html` — **Einheit vollständig.**
  Vorher/Nachher-Vergleich liest die Tallys `positionierung-vorher` (Std. 1) und
  `positionierung-nachher` aus localStorage → Std. 1 und 5 müssen im selben
  Browser laufen; Fallback (mündlicher Vergleich) steht in der Lektion.
  Bewertungsraster fürs Lernprodukt in Lektion 5, Kurzfassung im Einheitsplan.
  Alle internen Links geprüft (Skript-Check 18.07.2026, alle OK).
- Nächste Schritte nach der Durchführung: Praxis-Feedback von Markus einholen
  und Einheit überarbeiten; offene Idee: Elterninfo-Zettel.
- Positionierungs-Ergebnis aus Stunde 1 wird im Browser (localStorage)
  gespeichert, damit Stunde 5 den Vorher/Nachher-Vergleich zeigen kann —
  dafür muss Stunde 1 und 5 im selben Browser geöffnet werden.

## Arbeitsstand (08.08.2026)
- Neue Startseite `index.html` (Workspace-Root): schönes Inhaltsverzeichnis mit
  den 5 Stunden und dem Referenzmaterial (Einheitsplan, Faktenblatt, Rollenkarten),
  je mit Kurzbeschreibung. Neue wiederverwendbare Komponente `.verzeichnis` in
  `assets/unterricht.css`. Alle bestehenden Dateien (5 Lektionen + 3 Referenzdokumente)
  haben in der Navleiste jetzt vorne einen „Start"-Link auf `index.html` bekommen
  (Stehende Anweisung 18.07.2026 befolgt).

## Ideen für später
- ~~Stunde 4 (Debatte): Rollenkarten als druckbares Referenzdokument.~~ → erledigt 18.07.2026.
- Evtl. Elterninfo-Zettel, falls das Thema zuhause Fragen auslöst.
