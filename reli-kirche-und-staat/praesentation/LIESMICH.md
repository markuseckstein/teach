# Beamer-Praesentation

`Kirche-und-Staat-Beamer.pptx` — 60 Folien, 16:9, alle sechs Doppelstunden.

## Was drin ist
Nur das, was im Unterricht tatsaechlich an die Wand kommt: Projektionsfolien,
Arbeitsauftraege und Merksaetze. Was ausgeteilt wird — Expertenblaetter,
Rollenkarten, Raster, Steckbriefe — steht bewusst nicht drin; das liegt als
Kopiervorlage in `word/` und `lessons/`.

## Die Referentenansicht ist der eigentliche Trick
**43 der 60 Folien haben Notizen**, und darin stehen die ausformulierten
Ueberleitungen, die Auflösungen und die Regieanweisungen. In der
Referentenansicht sehen Sie diesen Text, waehrend die Klasse nur die Folie
sieht.

- PowerPoint: Bildschirmpraesentation → **Referentenansicht verwenden**
- LibreOffice Impress: Bildschirmpraesentation → **Praesentationsmodus**

## Aufbau
Vor jeder Doppelstunde steht eine dunkelblaue Trennfolie. Sie koennen die
Praesentation also am Stueck durchlaufen lassen oder pro Stunde ab der
Trennfolie starten.

## Vor dem Einsatz pruefen
- **Folie „Was die Kirchen gesagt haben“ (DS 6)**: Die Zeile
  `SACHSTAND GEPRUEFT AM: ____` ausfuellen. Die Einstufung durch den
  Verfassungsschutz ist rechtlich in Bewegung.
- **Folie „Kirche im Sozialismus“ (DS 5)**: zwei belegte Zahlen ergaenzen.

## Bildnachweis
Die sieben Fotos stammen von Wikimedia Commons und sind frei lizenziert
(CC BY 2.0 / CC BY 3.0 / CC BY-SA 3.0 de / CC BY-SA 4.0). Die vollstaendigen
Angaben stehen am Ende der jeweiligen HTML-Kopiervorlage in `lessons/` und
in `assets/img/quellen.json` bzw. `quellen-ds2-6.json`. Bei Weitergabe der
Praesentation ausserhalb des Unterrichts muss der Nachweis mitlaufen.

## Neu erzeugen
    python3 build_pptx.py [ziel.pptx]
