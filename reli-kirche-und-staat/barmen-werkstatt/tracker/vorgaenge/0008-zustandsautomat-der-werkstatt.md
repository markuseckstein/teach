---
id: 0008
titel: Zustandsautomat der Werkstatt
status: erledigt
label: blockiert
haengt-an: ["0006", "0007"]
---

# 0008 — Zustandsautomat der Werkstatt

## Worum es geht

Das Kernstück. Eine Gruppe wird durch die Werkstatt geführt: Bibelstelle wählen,
positiven Teil schreiben, Verwerfung schreiben, Vorschau, zwei Prüffragen, Freigabe.
Hier steckt die didaktische Logik im Code.

## Fertig, wenn

- Der Automat aus der Spezifikation ist vollständig umgesetzt:
  `BIBELSTELLE → POSITIV → VERWERFUNG → VORSCHAU → PRUEFUNG_1 → PRUEFUNG_2 → FREIGEGEBEN`.
- Ohne gewählte Bibelstelle ist kein Weiterkommen möglich — auch nicht durch eine
  direkt an einen späteren Schritt abgesetzte Anfrage.
- Der positive Teil hat zwei getrennte Felder („Weil …" / „gilt für uns: …"), die
  Verwerfung ein eigenes, vorbelegt mit „Wir verwerfen die falsche Lehre, als ob …".
- Ein „nein" bei einer Prüffrage führt zurück nach `VERWERFUNG` — **nie an den
  Anfang** —, der Text der Gruppe bleibt unverändert stehen, und `ueberarbeitungen`
  steigt um eins.
- Erst nach zwei Ja wird die These freigegeben und damit druckbar.
- Tests decken alle Pfade ab, beide Nein-Zweige einzeln.

## Gehört nicht dazu

Jede Form von Bewertung. Keine Wortlisten, keine Parteinamenerkennung, keine
automatische Ablehnung. Beide Prüffragen beantwortet die Gruppe selbst — ein Filter
würde genau das Überwältigungsverbot verletzen, das die Prüfung absichern soll.
Siehe Vorgang 0009.
