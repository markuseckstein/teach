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
- **Kursstruktur seit 22.09.2026:** Teil I · Das Gespräch (1–2), Teil II · Die Folgen
  (3–10), Teil III · Die lange Sicht (11), Teil IV · Die Übung (12 ff. — reale Fälle des
  Nutzers, siehe „Regel für Fallübungen“).
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

## Regel für Fallübungen (seit 22.09.2026)

Wenn der Nutzer eine reale Situation mitbringt, gehört sie in **Teil IV**, nicht in Teil I
— und die Lektion bekommt die nächste freie Nummer. Der Grund ist pragmatisch: Einschieben
heißt Umnummerieren, und das ist die fehleranfälligste Operation in diesem Workspace
(siehe unten). Eine Fallübung ist didaktisch ohnehin etwas anderes als eine Lektion aus
Teil I.

**Inhaltliche Regel, die sich dabei ergeben hat und die ich für übertragbar halte:** Bei
einem Einzelfall ist die Motivfrage („was will die Person damit bezwecken?“) *nicht*
beantwortbar — Forschung beschreibt Muster über viele Fälle. Statt zu diagnostizieren also
mehrere plausible Lesarten nebeneinanderstellen und **eine Antwort suchen, die bei allen
funktioniert.** Das ist die Umkehrung des Vorgehens aus Lektion 1 und der Kern von
Lektion 12.

## Zu tun / offene Fäden

- ~~Folgenanalyse: pro Politikfeld eine eigene, kurze Lektion.~~ **Erledigt (16.09.2026):**
  Alle siebzehn Programmkapitel sind seitengenau ausgewertet. Was auf Bundesebene fehlt
  (Euro, NATO, Reisen, Zahlungsverkehr), braucht das **Bundesprogramm** als eigene
  Primärquelle. **Erledigt (16.09.2026):** Der Nutzer hat
  <https://www.afd.de/grundsatzprogramm/> geliefert; Grundsatzprogramm 2016 und
  Wahlprogramm 2025 liegen maschinenlesbar in `quellen/`, Lektion 5 zitiert seitengenau.
- ~~Frage „Bürgerkrieg?" gehört ans Ende der Folgenanalyse.~~ **Erledigt (16.09.2026):
  Lektion 11**, und die Platzierung am Ende war richtig. Befund: Es war die falsche Frage.
  Es gibt keine Konfliktforschung, die ein Bürgerkriegsmodell auf Deutschland anwendet —
  die einschlägige Literatur ist die zum *democratic backsliding*. Die Lektion bremst
  deshalb **beide** Seiten: die Verharmlosung und die Dramatisierung.
- Prüfen, ob die GEW-Umsetzbarkeitsanalyse (Rechtsgutachten) als Primärquelle taugt oder
  als Interessenvertretung gekennzeichnet werden muss.
- ~~Referenz „Glossar" anlegen.~~ **Erledigt (10.09.2026):** `reference/glossar.html`,
  inzwischen 42 Einträge in acht Gruppen (16.09.2026: Gruppe „Einflussnahme und
  Desinformation“; 22.09.2026: Gruppe „Konflikt im Privaten“ mit Reaktanz, Aufrechnung,
  DARVO, Eskalationsstufen, Opferschaftsanspruch), jeweils mit Fundstelle im Kurs und — wo es zählt — einem
  „Nicht verwechseln mit …". Regel für die Aufnahme: ein Begriff kommt hinein, wenn er in
  einer Lektion mit `<span class="term">` ausgezeichnet ist **oder** im Gespräch regelmäßig
  unscharf gebraucht wird. Beim nächsten Ausbau mitpflegen, sonst veraltet er still.
- ~~Der Kurs hat keine Wiederholung.~~ **Erledigt (10.09.2026):**
  `reference/wiederholungsplan.html` — fünf Abrufe an Tag 1, 3, 7, 21, 60, Fragen aus dem
  Gedächtnis, Antworten in `<details class="reveal">`, Datumsfelder zum Eintragen. Ab
  Abruf 4 sind es Anwendungsfragen, keine Wissensfragen mehr. **Für Lektion 10 erledigt
  (16.09.2026):** Kernfrage als Frage 5 in Abruf 3, Merksatz in der Liste in Abruf 5.
  Für Lektion 2, 9 und 11 am 16.09.2026 ebenfalls erledigt.
- ~~**Nachgereichte Frage (10.09.2026):** Warum unterstützen Russland und ein Teil der
  US-Tech-Milliardäre die AfD?~~ **Erledigt (16.09.2026): Lektion 10.** Die geplante
  Dreiteilung hat sich in der Recherche als überflüssig erwiesen — die vorhandenen drei
  Sicherheitsstufen ①②③ tragen das Thema, ergänzt um eine ausdrückliche Regel:
  *Ermittlungsverfahren ≠ Schuldspruch.* Der Kern der Lektion ist nicht der Befund,
  sondern die Grenze: Der Satz „Die AfD wird aus Moskau gesteuert" kippt beim ersten
  Nachhaken. Details im Learning Record 0006.
- **Bewusst nicht in Lektion 10: Peter Thiel und Palantir.** Es gibt reichlich Belege für
  Thiels Rolle im Umfeld von JD Vance, aber keinen für eine Verbindung zur AfD. Wenn der
  Nutzer danach fragt (die Frage lag als „US-Milliardäre" im Plural vor), ist die ehrliche
  Antwort: Musk ja, Vance und Rubio als Regierung ja, Thiel nicht belegt.
- **Offen aus Lektion 9:** Versammlungsgesetz Sachsen-Anhalt im Wortlaut ungeprüft;
  ob eine Pflicht zur deutschen Sprache auf Versammlungen vor dem BVerfG Bestand hätte,
  steht nur als meine Einschätzung (keine einschlägige Entscheidung gefunden).
- **Offen aus Lektion 11:** Ob sich Landesregierungen sinnvoll mit Nationalstaaten
  vergleichen lassen, ist eine offene Forschungsfrage — die Übertragungstabelle ist
  Stufe ③. Und: in ein paar Monaten nachfassen, was die neue ungarische Regierung mit
  ihrer eigenen Zweidrittelmehrheit macht.
- ~~**Nächste Sitzung vermutlich keine Lektion, sondern eine Übung.**~~ **Eingetreten
  (22.09.2026): Lektion 12.** Der Anlass war ein echter Vorfall — freundliche Bitte an eine
  Bekannte, keine AfD-Wahlwerbung mehr per WhatsApp zu schicken; darauf ein unwahrer
  Vorwurf gegen den Mann des Nutzers, den sie gar nicht kennt. Die Lektion hat einen
  **eigenen Kursteil** bekommen (Teil IV · Die Übung), damit der Kurs nicht zum dritten Mal
  umnummeriert werden muss. Details im Record 0008.
- **Offen aus Lektion 10:** Ergebnis der Prüfung der Bundestagsverwaltung zum
  Musk-Weidel-Gespräch (am 16.09.2026 nicht auffindbar); Wirkungsforschung zu
  Desinformation und Wahlergebnissen (in der Lektion bewusst nicht behauptet).
- ~~Noch keine Lektion zu „Publikum am Stammtisch".~~ **Erledigt (16.09.2026): Lektion 2.**
- **Lektion 8, damals 7 (10.09.2026): das Muster war nicht inhaltlich, sondern formal.** Die vier
  Kapitel V, XI, XIII und XVII haben thematisch nichts gemeinsam — verbunden sind sie durch
  die Bauform: überwiegend kostenwirksame Forderungen, streckenweise ausgesprochen
  sympathische, und auf gut 40 Seiten genau **eine** Euro-Zahl (180 Mio. €, S. 237) und eine
  Prozentzahl (+20 % Medizinstudienplätze, S. 243). Dagegen steht die Selbstbindung in
  Kapitel XVI. Diese Zählung ist im Volltext nachprüfbar und ist das stärkste Argument der
  Lektion, weil sie ohne jede Wertung auskommt.
- **Vorsicht mit den eigenen Notizen in RESOURCES.md.** Beim Nachprüfen der
  BA-Pressemitteilung für Lektion 8 stand dort etwas anderes, als hier notiert war (keine
  Ärzteschaftszahlen, andere Beschäftigtenzahlen). Die Notiz ist korrigiert. Lehre: Auch die
  eigene Ressourcenliste ist Sekundärliteratur — vor der Verwendung einer Zahl die Quelle
  selbst öffnen, nicht die Zusammenfassung.
- **Offen aus Lektion 8:** Bologna-Ausstieg praktisch (keine Umsetzungsanalyse gefunden),
  Investitionsstau in Zahlen (KfW-Kommunalpanel 2025 ungeprüft), Kostenschätzung speziell
  für diese vier Kapitel (existiert nicht).
- **Erledigt aus Lektion 7 (10.09.2026):** Richterwahlausschuss und bilaterale
  Rückführungsverträge sind nachrecherchiert; beide Stellen in Lektion 7 sind präzisiert,
  die verbliebene Rest-Unsicherheit steht dort ausdrücklich. Details in RESOURCES.md
  unter `## Gaps`.
- **Zitierweise Bundesprogramme (seit 16.09.2026).** `GP 2016, 4.2, S. 30` und
  `WP 2025, S. 87 (PDF 88)`. **Die Seitenversätze sind unterschiedlich:** Beim
  Grundsatzprogramm ist gedruckte Seite = PDF-Seite, beim Wahlprogramm 2025 liegt die
  PDF-Seite um eins höher. Beim Landesprogramm sind sie wieder gleich. Vor jedem Zitat
  prüfen, sonst entstehen stillschweigend falsche Fundstellen.
- **Lehre aus der Nachprüfung von Lektion 5 (16.09.2026).** Der NATO-Abschnitt war aus
  zweiter Hand *nicht falsch* — der Satz zum europäischen Militärbündnis steht fast
  wörtlich im Programm. Aber die härteste Forderung fehlte ganz: „Abzug aller auf
  deutschem Boden stationierten alliierten Truppen und insbesondere ihrer Atomwaffen“
  (GP 2016, S. 31). **Zusammenfassungen sind nicht falsch, sie sind glatt.** Sie lassen
  weg, was sich schlecht zusammenfassen lässt — und das ist regelmäßig das Interessanteste.
- **Der Programm-Volltext liegt jetzt im Workspace** (`quellen/`). Vor jeder weiteren
  Programmaussage dort nachschlagen, nie aus zweiter Hand zitieren. Zitierweise, die sich
  eingebürgert hat: Kapitelnummer, Forderungsnummer, Seite — z. B. „XIV.7, S. 213".
  Gedruckte Seite = PDF-Seite.
- **Persönliche Betroffenheit beachten.** Der Nutzer ist evangelische Religionslehrkraft
  am Gymnasium (in Bayern). Das Programm greift den Beutelsbacher Konsens namentlich an
  (S. 77), nennt Bayern ausdrücklich als Vorbild (S. 78), streicht die Evangelische
  Akademie namentlich (S. 67) und bindet Staatsleistungen auf „geistliches Personal und
  Erhalt der Gebäude" (S. 66 f.). Lektion 6 macht das zum Thema — mit dem klaren Hinweis
  vorweg, dass für Bayern kein Wort davon gilt. **Nicht in Unterrichtsmaterial abgleiten**:
  Die Mission sagt eigene Urteilsbildung, nicht Schule. Wenn er Material fürs Kollegium
  will, ist das ein eigener Workspace mit eigenen didaktischen Regeln.

## Umnummerierung vom 16.09.2026 — unbedingt beachten

Die drei neuen Lektionen gehören inhaltlich an Position 2, 9 und 11. Weil Dateiname und
Lektionsnummer in diesem Workspace übereinstimmen, wurde der ganze Kurs umnummeriert:

| alt | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
|-----|---|---|---|---|---|---|---|---|
| neu | 1 | 3 | 4 | 5 | 6 | 7 | 8 | 10 |

**Die Learning Records 0001–0006 sind absichtlich nicht angepasst** — sie sind datierte
Protokolle. Wer dort „Lektion 3" liest, muss oben nachschlagen. Die Umrechnung steht auch
im Record 0007.

**Wenn wieder eine Lektion eingeschoben wird:** Das Skript dafür ist klein, aber der
Regex für „Lektion N" muss `Lektion(?:en)?` heißen — `Lektionen?` verlangt „Lektione" und
greift stillschweigend ins Leere. Genau das ist beim ersten Versuch passiert und fiel nur
auf, weil die Eyebrow-Zeilen danach noch die alten Nummern zeigten. Prüfen lässt es sich
hart: data-n, Dateinummer, `<title>` und Eyebrow müssen für jede Lektion dieselbe Zahl
ergeben; dieser Test läuft im Abschlussskript mit.

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

## Offene Fäden aus Lektion 12 (22.09.2026)

- **Volltext von Voit (2025)** zu Opferschaftsansprüchen in deutschen Wahlprogrammen ist
  paywalled (403). In der Lektion ist das mitten in der Quellenangabe gekennzeichnet.
  Wenn sich ein Zugang findet: Methode und Kodierregeln nachlesen, bevor die Aussage
  irgendwo mündlich verwendet wird.
- **Ungeklärt: Einzelchat oder Gruppe?** Bei einer Gruppe ändert sich der Rat erheblich —
  dann ist es der Publikumsfall aus Lektion 2. Die Lektion fragt das im Schlusskasten nach;
  bei der nächsten Sitzung daran denken.
- **Nachfassen:** Hat er geantwortet, und was kam zurück? Wenn die Geschichte
  weitergetragen wird, ist es keine Gesprächsführungsfrage mehr, sondern ein Fall für den
  Bundesverband Mobile Beratung (jetzt in RESOURCES.md und in der Lektion verlinkt).
- **Community-Lücke halb geschlossen.** Für den privaten Konfliktfall gibt es jetzt eine
  belastbare Anlaufstelle. Eine *Online*-Community für Gesprächstraining fehlt weiterhin —
  und Kommentarspalten bleiben dafür ungeeignet.
