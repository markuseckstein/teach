# Notizen

## Präferenzen des Nutzers

- **Anrede: du.** Gilt für alle Lektionen und Referenzen.
- **Sprache: Deutsch.**
- **Haltung: klar positioniert** (in der Eingangsfrage so gewählt). Das heißt konkret:
  Faktenlage zuerst und sauber belegt, danach ein sichtbar abgesetzter Absatz
  „**Meine Einordnung**". Nie beides vermischen. Wo die Forschung uneins ist, wird das gesagt —
  „klar positioniert" ist keine Lizenz, Unsicherheit wegzubügeln.
- **Zweck:** eigene Gespräche (Stammtisch, Familie, Kollegium) + eigene Urteilsbildung.
  Ausdrücklich *nicht* Unterrichtsmaterial. Beispiele also aus dem Erwachsenenalltag,
  nicht aus dem Klassenzimmer.
- **Einstieg gewählt:** Gesprächsführung zuerst, Folgenanalyse danach.
- Bilder und Diagramme sind erwünscht (inline SVG, offline- und drucktauglich).

## Architekturentscheidung: eine Quelle für die Navigation

Übernommen aus dem Workspace `../wasistki` (dieselbe Person, bewährt):
Reihenfolge, Titel und Teaser aller Lektionen stehen **ausschließlich** in
`assets/course.js` (`COURSE`). `assets/lesson.js` rendert daraus Kopf- und Fußnavigation,
`index.html` das Inhaltsverzeichnis.

**Konsequenz:** Neue Lektion = HTML-Datei anlegen + einen Eintrag in `COURSE.lessons`.
Nav-Links werden nie von Hand angefasst.

Technisch: Seiten laufen per `file://`. Deshalb **keine** ES-Module, **kein** `fetch()`.
Klassische `<script src>`-Tags mit globalen Variablen.

## Inhaltliche Leitplanken

- **Kein Strohmann.** Wo etwas im Programm steht, wird zitiert und verlinkt (Primärquelle:
  das Regierungsprogramm selbst). Wo etwas *nicht* drinsteht, aber unterstellt wird, wird das
  markiert. Ein Argument, das an einer Falschbehauptung hängt, ist am Stammtisch wertlos —
  es kippt beim ersten Nachhaken.
- **Der Backfire-Effekt ist überschätzt.** Wood/Porter (2019, N > 10.000) finden ihn praktisch
  nicht. Fakten wirken also durchaus — nur langsam und nicht bei identitätsrelevanten Aussagen.
  Diese Nuance muss durch alle Lektionen tragen, sonst entsteht die falsche Lehre
  „Fakten bringen nichts".
- **Vorsicht mit der Protest-These.** Sie ist bequem und empirisch schwach (FoDEx 2024: ~8 %
  inkonsistente Wahlabsichten). Wer sein Gegenüber als „Protestwähler" einsortiert, redet oft
  an der Person vorbei.

## Zu tun / offene Fäden

- Folgenanalyse: pro Politikfeld eine eigene, kurze Lektion (Haushalt, Bildung, Energie,
  EU/Reisen/Zahlungsverkehr, NATO, Rechtsstaat). Nicht alles in eine Lektion pressen.
- Frage „Bürgerkrieg?" gehört ans Ende der Folgenanalyse, nicht an den Anfang —
  sonst dominiert die Dramatik die Sachfragen.
- Prüfen, ob die GEW-Umsetzbarkeitsanalyse (Rechtsgutachten) als Primärquelle taugt oder
  als Interessenvertretung gekennzeichnet werden muss.
- ~~Referenz „Glossar" anlegen.~~ **Erledigt (10.09.2026):** `reference/glossar.html`,
  25 Einträge in sechs Gruppen, jeweils mit Fundstelle im Kurs und — wo es zählt — einem
  „Nicht verwechseln mit …". Regel für die Aufnahme: ein Begriff kommt hinein, wenn er in
  einer Lektion mit `<span class="term">` ausgezeichnet ist **oder** im Gespräch regelmäßig
  unscharf gebraucht wird. Beim nächsten Ausbau mitpflegen, sonst veraltet er still.
- ~~Der Kurs hat keine Wiederholung.~~ **Erledigt (10.09.2026):**
  `reference/wiederholungsplan.html` — fünf Abrufe an Tag 1, 3, 7, 21, 60, Fragen aus dem
  Gedächtnis, Antworten in `<details class="reveal">`, Datumsfelder zum Eintragen. Ab
  Abruf 4 sind es Anwendungsfragen, keine Wissensfragen mehr. Wenn Lektion 8 kommt, gehört
  ihre Kernfrage in Abruf 3 und ihr Merksatz in die Liste in Abruf 5.
- **Nachgereichte Frage (10.09.2026):** Warum unterstützen Russland und ein Teil der
  US-Tech-Milliardäre die AfD? Eigene Lektion in Teil II. Heikel, weil hier die Grenze
  zwischen Beleg und Verschwörungserzählung schmal ist — daher strikt dreiteilen:
  (1) gerichtsfest/behördlich belegt, (2) gut recherchiert, aber bestritten, (3) plausibel,
  aber unbelegt. Quellen: Verfassungsschutzberichte, EU DisinfoLab/EUvsDisinfo,
  Rechercheverbünde (Correctiv, ZDF/SZ/NDR/WDR), öffentliche Äußerungen der Beteiligten selbst.
- Noch keine Lektion zu „Publikum am Stammtisch". Steht im Roadmap-Abschnitt der index.html.
- **Lektion 7 (10.09.2026): das Muster war nicht inhaltlich, sondern formal.** Die vier
  Kapitel V, XI, XIII und XVII haben thematisch nichts gemeinsam — verbunden sind sie durch
  die Bauform: überwiegend kostenwirksame Forderungen, streckenweise ausgesprochen
  sympathische, und auf gut 40 Seiten genau **eine** Euro-Zahl (180 Mio. €, S. 237) und eine
  Prozentzahl (+20 % Medizinstudienplätze, S. 243). Dagegen steht die Selbstbindung in
  Kapitel XVI. Diese Zählung ist im Volltext nachprüfbar und ist das stärkste Argument der
  Lektion, weil sie ohne jede Wertung auskommt.
- **Vorsicht mit den eigenen Notizen in RESOURCES.md.** Beim Nachprüfen der
  BA-Pressemitteilung für Lektion 7 stand dort etwas anderes, als hier notiert war (keine
  Ärzteschaftszahlen, andere Beschäftigtenzahlen). Die Notiz ist korrigiert. Lehre: Auch die
  eigene Ressourcenliste ist Sekundärliteratur — vor der Verwendung einer Zahl die Quelle
  selbst öffnen, nicht die Zusammenfassung.
- **Offen aus Lektion 7:** Bologna-Ausstieg praktisch (keine Umsetzungsanalyse gefunden),
  Investitionsstau in Zahlen (KfW-Kommunalpanel 2025 ungeprüft), Kostenschätzung speziell
  für diese vier Kapitel (existiert nicht).
- **Erledigt aus Lektion 6 (10.09.2026):** Richterwahlausschuss und bilaterale
  Rückführungsverträge sind nachrecherchiert; beide Stellen in Lektion 6 sind präzisiert,
  die verbliebene Rest-Unsicherheit steht dort ausdrücklich. Details in RESOURCES.md
  unter `## Gaps`.
- **Der Programm-Volltext liegt jetzt im Workspace** (`quellen/`). Vor jeder weiteren
  Programmaussage dort nachschlagen, nie aus zweiter Hand zitieren. Zitierweise, die sich
  eingebürgert hat: Kapitelnummer, Forderungsnummer, Seite — z. B. „XIV.7, S. 213".
  Gedruckte Seite = PDF-Seite.
- **Persönliche Betroffenheit beachten.** Der Nutzer ist evangelische Religionslehrkraft
  am Gymnasium (in Bayern). Das Programm greift den Beutelsbacher Konsens namentlich an
  (S. 77), nennt Bayern ausdrücklich als Vorbild (S. 78), streicht die Evangelische
  Akademie namentlich (S. 67) und bindet Staatsleistungen auf „geistliches Personal und
  Erhalt der Gebäude" (S. 66 f.). Lektion 5 macht das zum Thema — mit dem klaren Hinweis
  vorweg, dass für Bayern kein Wort davon gilt. **Nicht in Unterrichtsmaterial abgleiten**:
  Die Mission sagt eigene Urteilsbildung, nicht Schule. Wenn er Material fürs Kollegium
  will, ist das ein eigener Workspace mit eigenen didaktischen Regeln.

## Layout-Konventionen (nach dem Fehler vom 10.09.2026)

Drei Regeln, die aus einem kaputten Rendering entstanden sind:

1. **Niemals HTML-Tags in `<svg><text>`.** Ein `<em>` dort beendet den SVG-Parse, der Rest
   des Diagramms landet als Fliesstext im Dokument. Kursiv im Diagramm geht nur über
   `<tspan font-style="italic">`. Vor dem Commit prüfen: alle `<svg>`-Blöcke einzeln mit
   `xml.etree.ElementTree` parsen — das findet genau diese Klasse von Fehlern.
2. **`a.source` ist ein Inline-Beleg**, keine Box. `assets/lesson.js` nummeriert diese Links,
   ersetzt den Text durch eine hochgestellte Ziffer und baut daraus am Seitenfuss die
   Quellenliste. Ohne JS bleibt der Klartext lesbar stehen.
3. **`.sidenote` ist eine schmale Marginalie im Gutter** (aus `../wasistki` geerbt) und
   eignet sich nicht für mehrzeilige Quellenangaben unter Tabellen. Dafür gibt es
   `.srcnote`: volle Lesebreite, kleine Schrift, linke Randlinie.

**Vor jedem Abschluss laufen lassen:** Tag-Verschachtelung aller Seiten (`html.parser`),
SVG-Wohlgeformtheit, tote interne Links. Drei kurze Python-Schnipsel, fangen zusammen die
Fehler ab, die man im Terminal sonst nicht sieht.
