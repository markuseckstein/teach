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
