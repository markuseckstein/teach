# Notes

## Preferences
- Deutsch, duzen.
- Ziel ist „die Welt verstehen“ (Neugier), nicht akute Sorge: sachlich bleiben, keine Panik-Rhetorik.
- Mag: Quiz, historische Fälle, Rechnen/Simulieren, Primärquellen. Etwas längere Lektionen sind ok, wenn sinnvoll.
- 1–2 h/Woche: etwa 2 Lektionen pro Woche.

## Interessen im Hintergrund (in geplante Lektionen einweben, keine eigene Lektion)
- **Psychologie der Inflation** (2026-10-02, Anlass: hohe Spritpreise, Tankrabatt 2026): gefühlte Inflation, Ohnmacht, „die müssen was tun“, Festhalten am Gewohnten. Lernender hält den Tankrabatt für falsch (Pflaster statt Anpassung) – Gegenargumente fair zeigen, nicht bestätigen oder belehren.
  - 0005 (1923): Vertrauensverlust ins Geld, Erwartungen kippen, Flucht in Sachwerte, Suche nach Schuldigen. Kontrast: 2022/2026 blieben 5-Jahres-Erwartungen ≈ 2,5 % verankert.
  - 0006 (1948): Vertrauen durch neue Währung + harte Zumutung (Lastenausgleich) als Gegenmodell zum „Pflaster“.
  - 0008 (Synthese): Ohnmacht vs. Handlungsspielraum – persönliche Inflation durch Anpassung senken; Preissignal vs. Subvention (Tankrabatt) vs. gezielter Transfer.
  - Aufwärm-/Transferfragen dazu in späteren Lektionen.

## Grober Lernpfad (Arbeitsplan, wird angepasst)
1. Woher kommt das Geld? Guthaben = Forderung, Kredit schafft Geld, Tilgung vernichtet es  ← 0001 (6/6 richtig)
2. Zentralbankgeld: Überweisung zwischen Banken, Reserven, Bargeld. Warum es eine Zentralbank braucht  ← 0002 (alles richtig)
3. Leitzins und Geldpolitik: Wie die EZB Kreditvergabe steuert  ← 0003 (9/10)
4. Inflation: Messung (HVPI), Ursachen (Nachfrage, Angebot, Geldmenge), Erwartungen. Wer gewinnt und wer verliert  ← 0004
5. Fall 1923: Kriegsfinanzierung, Reparationen, Notenpresse, Hyperinflation. Sachwerte vs. Geldwerte, Hauszinssteuer
6. Fall 1948: Währungsreform (Umstellung 100:6,5), Lastenausgleich (50 % Vermögensabgabe auf Immobilien)
7. Bankpleite und Bank Run: Einlagensicherung, Zypern 2013
8. Synthese: Konto vs. Depot vs. Immobilie in vier Krisentypen

## Komponenten
- assets/style.css: gemeinsames Tufte-artiges Layout (Randnotizen, Callouts, Druck)
- assets/quiz.js: Multiple-Choice-Quiz mit sofortigem Feedback, mischt die Optionen
- assets/bilanz.js: T-Konten-Simulator (Akteure, Aktiva/Passiva, mehrere Zähler über `counters`)
- assets/linechart.js: Liniendiagramm (Hover, Legende, Endlabels, Tabellenansicht); Palette validiert (dataviz)
- assets/rechner.js: Schieberegler-Rechner mit Presets (0003 Kredit, 0004 persönliche Inflation + Kaufkraft)
- assets/course.js: Kursverzeichnis (einzige Stelle, um Lektionen zu registrieren) → erzeugt index.html-Inhalt, Abschnitts-Inhaltsverzeichnis, Zurück/Weiter

## Beobachtungen
- 0001: 6/6 im ersten Versuch → Tempo kann steigen; Lektionen dürfen dichter werden.
- Möchte ein Inhaltsverzeichnis mit Navigation (umgesetzt in index.html + course.js).
- Jede Lektion beginnt mit 1–2 Aufwärmfragen zu älteren Lektionen (Spacing).
- 0002: alles richtig → ab 0003 Rechenaufgaben und Transferfragen an echten Daten.
- Ideen für später: Quantitative Easing / Überschussliquidität, Negativzinsen, Renditen von Bundesanleihen (Ask-Box 0003).
- 0003: 9/10, falsch nur Frage 9 (Bauzinsen bleiben trotz Leitzinssenkung hoch) → Spacing in 0004 (Aufwärmen) und 0005. Glossar: Begriffe aus 0003 übernommen; die aus 0004 (Inflationsrate, Kerninflation, Realzins, Zweitrundeneffekt) erst nach Rückmeldung.
- 0004: Brücken gelegt zu 0005 (Geldmenge langfristig, Staat als Schuldner) und 0006 (Sachwerte, Lastenausgleich). In 0005 an Realzins/Umverteilung anknüpfen.
- Ideen für später: Deflation (Japan, 1930er), Qualitätsbereinigung im Warenkorb (Ask-Box 0004).
