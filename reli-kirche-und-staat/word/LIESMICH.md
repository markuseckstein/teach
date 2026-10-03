# Word-Dokumente

Erzeugt aus den HTML-Vorlagen mit `python3 html2docx.py <quelle.html> <ziel.docx>`
(Konverter liegt im Projektwurzelverzeichnis, benoetigt python-docx).

| Datei | Inhalt |
|---|---|
| `DS1-Kopiervorlagen-Kirche-und-Staat.docx` | Wo steckt Kirche im Staat? — M1–M5, Steckbrief, Eckenschilder |
| `DS2-Kopiervorlagen-1933-Die-Kirche-jubelt-mit.docx` | Bildlesen, Zeitstrahl, Quellen-Autopsie M6–M9 |
| `DS3-Kopiervorlagen-Barmen-gegen-Ansbach.docx` | Rollenkarten, Beobachtungsbogen, Barmer Thesen M10–M14 |
| `DS4-Kopiervorlagen-Was-haette-ich-getan.docx` | Chorlesung, fuenf Biografien, Placemat M15–M18 |
| `DS5-Kopiervorlagen-Kerzen-statt-Steine.docx` | Kirche in der DDR, Themenkarten, Ausstellungstafel M19–M22 |
| `DS6-Kopiervorlagen-Barmen-2026.docx` | Beutelsbach, vier Stimmen, Barmen-Werkstatt M23–M26 |
| `Unterrichtseinheit-Kirche-und-Staat.docx` | Die komplette Einheit mit allen Ueberleitungen |

## Bekannte Abweichungen zum HTML
- Schriften: Cambria / Calibri / Consolas statt Zilla Slab / Source Sans 3 / IBM Plex Mono
  (die Originalschriften sind auf diesem Rechner nicht installiert).
- Farbige Kaesten sind einzellige Tabellen mit Hintergrundfarbe statt CSS-Kaesten.
- Die Projektionsfolien verlieren das 16:9-Format. **Zum Projizieren die HTML-Fassung
  benutzen**, zum Kopieren und Anpassen die Word-Fassung.
- In Tabellenzellen sind Nummerierungen von Hand gesetzt und zaehlen beim Umsortieren
  nicht automatisch neu.
- Versalien-Kennzeilen schreiben "ss" statt scharfes s: DENKANSTOSS statt Denkanstoß.

## Neu erzeugen
    python3 build.py                      # HTML aus lessons-src/ bauen
    python3 html2docx.py lessons/0003-ds2-kopiervorlagen.html word/DS2-....docx

# Nacharbeit fuer kranke SuS — `word/nacharbeit/`

Sechs einseitige Blaetter, eines je Doppelstunde, erzeugt mit `python3 nacharbeit.py`.
Gedacht zum Verschicken an SuS, die gefehlt haben — nicht als Ersatz fuer die Stunde,
sondern damit sie in der Folgestunde mitreden koennen.

| Datei | Kern |
|---|---|
| `DS1-Wenn-du-krank-warst.docx` | Hinkende Trennung, vier Beruehrungspunkte |
| `DS2-Wenn-du-krank-warst.docx` | Deutsche Christen 1933, Roem 13 gegen Apg 5 |
| `DS3-Wenn-du-krank-warst.docx` | Barmen, Ansbach, der Fall Meiser |
| `DS4-Wenn-du-krank-warst.docx` | Fuenf Biografien, Bedingungen fuer Widerstand |
| `DS5-Wenn-du-krank-warst.docx` | Kirche in der DDR, Friedensgebete |
| `DS6-Wenn-du-krank-warst.docx` | Kirchliche Positionen, Barmer Sprachform, Haerteprufung |

Aufbau jedes Blattes: Kasten „Fuenf Minuten lesen“ · vier nummerierte Merkpunkte ·
Merksatz · „So geht es weiter“ (Anschluss an die Folgestunde) · freiwilliger
Bett-Auftrag mit Link · Schlusszeile „Frag mich“.

## Hinweise
- **Alle Links sind am 5. September 2026 auf Erreichbarkeit geprueft** (EKD, bpb,
  DHM/LeMO, jugendopposition.de, Deutsche Digitale Bibliothek, ein explainity-Video
  bei YouTube). Vor dem Verschicken trotzdem kurz anklicken.
- Die freiwilligen Auftraege sind bewusst als *freiwillig* gekennzeichnet und werden
  nicht eingesammelt — kranke SuS sollen keinen Leistungsdruck bekommen.
- **DS 6 ist das vollste Blatt** (knapp eine Seite). Wenn es bei dir auf zwei Seiten
  umbricht, den Einleitungskasten um einen Satz kuerzen.
- DS 6 haelt den Beutelsbacher Konsens ausdruecklich fest ("keine Wahlempfehlung"),
  weil dieses Blatt anders als das Unterrichtsmaterial in Elternhaende geraet.
- Gruppenprodukte (Steckbriefe, Personenkarten, Ausstellungstafeln) muessen nicht
  nachgeholt werden; die Blaetter sagen das implizit, indem sie nur das Ergebnis nennen.
