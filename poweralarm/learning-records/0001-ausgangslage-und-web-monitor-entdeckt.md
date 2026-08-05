# Ausgangslage im Feuerwehrhaus, und: der PA Web-Monitor existiert

Markus betreibt in einem kleinen Feuerwehrhaus den PowerAlarm **PA-Monitor als
Windows-Version** auf einem Rechner, der auf Windows 10 festsitzt und kein
Upgrade mehr bekommt. Er kennt die Windows-Anzeige also aus dem Betrieb — das ist
vorhandenes Vorwissen und muss nicht erklärt werden. Was ihm fehlte, war das
**Komponentenmodell** von PowerAlarm (Portal / PAWinS / PA-Monitor / PA
Web-Monitor) und das Wissen, dass es überhaupt eine Web-Variante gibt.

In Lektion 1 etabliert: Der **PA Web-Monitor** (Handbuch Kap. 5) ist die Lösung
für Browser-Kiosk auf Linux/Pi. Er hängt an einem von FITT freigeschalteten
API-Key mit der Option „PA Web-Monitor“; die URL wird ausschließlich im Portal
unter *Benutzer → Web-Monitor* (Punkt 2.18.d) angezeigt und ist nirgends
öffentlich dokumentiert.

**Implications für kommende Sessions:**

- Der entscheidende noch offene Fakt ist, ob auf dem Windows-Rechner auch
  **PAWinS / FMS32 / POC32 / monitord** läuft. Ja → der Rechner hat einen zweiten
  Job (Funkdekodierung), den der Web-Monitor nicht ersetzt, und Windows kann nicht
  einfach weg. Nein → Umstieg ist reine Anzeigefrage. Bis das geklärt ist, keine
  Lektion planen, die einen der beiden Wege voraussetzt.
- Der zweite kritische Unbekannte ist der **Alarmton** im Web-Monitor. Fehlt er,
  kippt die ganze Mission-Strategie (dann eher: Windows-Monitor auf einem
  aktuellen, sparsamen Windows-Gerät, oder Ton separat lösen). Das ist keine
  Lernfrage, sondern eine Herstellerfrage — Lektion 1 endet daher bewusst mit
  einer konkreten Mail an FITT statt mit Theorie.
- Nächste Lektion, sobald geantwortet wurde: entweder **Kiosk-Betrieb sauber
  aufsetzen** (Autostart, Vollbild, Bildschirm bei Alarm wecken, Autoplay-Policy
  für Ton) oder **PAWinS ohne Windows-Arbeitsplatz**.
