#!/usr/bin/env python3
"""Stellt die Schueleransprache von "Sie" auf "du/ihr" um.

Nicht angetastet wird bewusst:
  * Text, der die LEHRKRAFT anspricht (Regie, Erwartungshorizont, Hinweise)
  * das "Sie" in der Talkrunde von DS 3 — dort siezen die Schueler die
    Rollen (Barth, Althaus, Meiser); das ist inhaltlich richtig so
  * "sie/ihr" in der 3. Person ("Sie war offiziell freiwillig")
  * Bibel- und Quellenzitate

Beim Ansprechen der ganzen Klasse steht "ihr" (so duzt man eine Gruppe),
bei Aufgaben an Einzelne "du".

Aufruf:  python3 duzen.py [--probe]
"""
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).parent
DATEIEN = [
    ROOT / "lessons/0001-kirche-und-staat.html",
    ROOT / "lessons/0002-ds1-kopiervorlagen.html",
    *sorted((ROOT / "lessons-src").glob("*.html")),
    ROOT / "build_pptx.py",
]

# (alt, neu) — Reihenfolge zaehlt: Laengeres zuerst, damit sich kuerzere
# Muster nicht in bereits ersetzten Text hineinfressen.
ERSETZUNGEN = [
    # ---------------------------------------------------- DS 1
    ("Sie haben recht — in jedem", "Ihr habt recht — in jedem"),
    ("Aber schauen Sie noch einmal genau hin", "Aber schaut noch einmal genau hin"),
    ("Aber schauen Sie genau hin", "Aber schaut genau hin"),
    ("Diese Frage kann ich Ihnen nicht beantworten, ohne dass Sie wissen",
     "Diese Frage kann ich euch nicht beantworten, ohne dass ihr wisst"),
    ("Sie recherchieren jetzt selbst.", "Ihr recherchiert jetzt selbst."),
    ("Sie recherchieren selbst.", "Ihr recherchiert selbst."),
    ("Danach müssen Sie es den anderen erklären", "Danach müsst ihr es den anderen erklären"),
    ("Nur bei Ihnen.", "Nur bei euch."),
    ("„Stellen Sie sich vor, jemand ist aus der Kirche ausgetreten.",
     "„Stellt euch vor, jemand ist aus der Kirche ausgetreten."),
    ("Was stört diese Person an Ihrem Thema?", "Was stört diese Person an eurem Thema?"),
    ("Ihr Blatt bleibt liegen, Sie nehmen nur Ihren Kopf mit.",
     "Euer Blatt bleibt liegen, ihr nehmt nur euren Kopf mit."),
    ("Ihr Blatt bleibt liegen, Sie nehmen ", "Euer Blatt bleibt liegen, ihr nehmt "),
    ("nur Ihren Kopf mit.", "nur euren Kopf mit."),
    ("Sie haben dann viereinhalb Minuten pro Thema.", "Ihr habt dann viereinhalb Minuten pro Thema."),
    ("Wenn Ihre Zuhörer am Ende Ihr Feld", "Wenn eure Zuhörer am Ende euer Feld"),
    ("nicht ausfüllen können, haben Sie schlecht erklärt", "nicht ausfüllen können, habt ihr schlecht erklärt"),
    ("„Sie wissen jetzt, wie es ist.", "„Ihr wisst jetzt, wie es ist."),
    ("Stehen Sie bitte auf und schieben Sie die Stühle an die Tische.",
     "Steht bitte auf und schiebt die Stühle an die Tische."),
    ("Ich lese Ihnen gleich vier Sätze vor.", "Ich lese euch gleich vier Sätze vor."),
    ("Zu jedem Satz stellen Sie sich in die Ecke, die zu Ihrer Meinung passt.",
     "Zu jedem Satz stellt ihr euch in die Ecke, die zu eurer Meinung passt."),
    ("Stellen Sie sich in die Ecke, die zu Ihrer Meinung passt.",
     "Stellt euch in die Ecke, die zu eurer Meinung passt."),
    ("Wechseln Sie die Ecke", "Wechselt die Ecke"),
    ("Ich zeige Ihnen zum Schluss ein einziges Bild.", "Ich zeige euch zum Schluss ein einziges Bild."),
    ("Und die Frage, die Sie mitnehmen", "Und die Frage, die ihr mitnehmt"),
    ("Fragen Sie eine erwachsene Person in Ihrem Umfeld",
     "Frag eine erwachsene Person in deinem Umfeld"),
    ("Notieren Sie die Antwort in zwei Sätzen.", "Notiere die Antwort in zwei Sätzen."),

    # ---------------------------------------------------- DS 2
    ("Notieren Sie fünf Dinge, die Sie ", "Notiere fünf Dinge, die du "),
    ("tatsächlich sehen</em>", "tatsächlich siehst</em>"),
    ("Notieren Sie fünf Dinge, die Sie auf dem Bild tatsächlich sehen.",
     "Notiere fünf Dinge, die du auf dem Bild tatsächlich siehst."),
    ("Notieren Sie eine Sache, die Sie in einer Kirche erwarten würden",
     "Notiere eine Sache, die du in einer Kirche erwarten würdest"),
    ("Formulieren Sie die eine Frage, die Sie den Menschen auf diesem Bild stellen würden.",
     "Formuliere die eine Frage, die du den Menschen auf diesem Bild stellen würdest."),
    ("Die meisten Ihrer Fragen laufen auf dasselbe hinaus",
     "Die meisten eurer Fragen laufen auf dasselbe hinaus"),
    ("Ich möchte Ihnen diese Frage heute wegnehmen.", "Ich möchte euch diese Frage heute wegnehmen."),
    ("Ich möchte Ihnen diese Frage wegnehmen.", "Ich möchte euch diese Frage wegnehmen."),
    ("bei der Sie das Kirchensteuer-Interview gemacht haben",
     "bei der ihr das Kirchensteuer-Interview gemacht habt"),
    ("brauchen Sie acht Minuten Vorgeschichte", "braucht ihr acht Minuten Vorgeschichte"),
    ("Sie bekommen gleich Originaltexte", "Ihr bekommt gleich Originaltexte"),
    ("Ich möchte trotzdem, dass Sie sie nicht sofort verurteilen",
     "Ich möchte trotzdem, dass ihr sie nicht sofort verurteilt"),
    ("Unterstreichen Sie im Text die stärkste Forderung.", "Unterstreicht im Text die stärkste Forderung."),
    ("Achten Sie besonders darauf", "Achtet besonders darauf"),
    ("Beschreiben Sie eine konkrete Person", "Beschreibt eine konkrete Person"),
    ("Schreiben Sie Ihre Quelle in", "Schreibt eure Quelle in"),
    ("Ich gebe sie Ihnen nicht, weil sie eine Meinung", "Ich gebe sie euch nicht, weil sie eine Meinung"),
    ("Ich gebe sie Ihnen, ", "Ich gebe sie euch, "),
    ("Wenn Ihnen ein Satz zu nahe geht, sagen Sie es mir",
     "Wenn euch ein Satz zu nahe geht, sagt es mir"),
    ("„Hängen Sie Ihre Sätze an die Leine, in der Reihenfolge, in der ich Sie aufrufe.",
     "„Hängt eure Sätze an die Leine, in der Reihenfolge, in der ich euch aufrufe."),
    ("Schauen Sie sich das jetzt bitte als Ganzes an.", "Schaut euch das jetzt bitte als Ganzes an."),
    ("Ich zeige Ihnen beide — und Sie werden sehen", "Ich zeige euch beide — und ihr werdet sehen"),
    ("Formulieren Sie jeden Vers als", "Formuliert jeden Vers als"),
    ("Suchen Sie zu jedem Vers eine Person", "Sucht zu jedem Vers eine Person"),
    ("Beantworten Sie schriftlich", "Beantwortet schriftlich"),
    ("„Ich lasse Sie heute mit einem einzigen Satz gehen.",
     "„Ich lasse euch heute mit einem einzigen Satz gehen."),
    ("Überlegen Sie bis nächste Woche", "Überlegt bis nächste Woche"),
    ("Und Sie sitzen nicht im Publikum.", "Und ihr sitzt nicht im Publikum."),
    ("Sie sind die beiden Seiten.", "Ihr seid die beiden Seiten."),

    # ---------------------------------------------------- DS 3
    ("Zeigen Sie mit dem Finger: links oder rechts.", "Zeigt mit dem Finger: links oder rechts."),
    ("Woran haben Sie das festgemacht?", "Woran habt ihr das festgemacht?"),
    ("„Sie haben richtig getippt.", "„Ihr habt richtig getippt."),
    ("und Sie führen ihn.", "und ihr führt ihn."),
    ("Ab dem Moment, in dem Sie sich setzen, sind Sie nicht mehr Sie selbst.",
     "Ab dem Moment, in dem ihr euch hinsetzt, seid ihr nicht mehr ihr selbst."),
    ("Ich siezee Sie mit Ihrem Rollennamen, und Sie antworten in der Rolle.",
     "Ich rede euch mit eurem Rollennamen an und sieze euch — ihr antwortet in der Rolle."),
    ("Sie greifen Positionen an, nicht Personen.", "Ihr greift Positionen an, nicht Personen."),
    ("Bitte stehen Sie kurz auf, schütteln Sie sich einmal aus und setzen Sie sich auf einen anderen Platz.",
     "Steht bitte kurz auf, schüttelt euch einmal aus und setzt euch auf einen anderen Platz."),
    ("ab jetzt sind Sie wieder Sie selbst", "ab jetzt seid ihr wieder ihr selbst"),
    ("„Wie hat sich Ihre Rolle angefühlt?“", "„Wie hat sich eure Rolle angefühlt?“"),
    ("„Sie haben jetzt zweiundzwanzig Minuten lang über einen Text gestritten, den Sie noch gar nicht richtig gelesen haben.",
     "„Ihr habt jetzt zweiundzwanzig Minuten lang über einen Text gestritten, den ihr noch gar nicht richtig gelesen habt."),
    ("Ihre Aufgabe ist es, ihn so zu übersetzen", "Eure Aufgabe ist es, ihn so zu übersetzen"),
    ("Ihre Aufgabe ist, ", "Eure Aufgabe ist, "),
    ("und Sie hängen diese Menschen daran auf", "und ihr hängt diese Menschen daran auf"),
    ("vor der Sie sich heute noch drücken dürfen", "vor der ihr euch heute noch drücken dürft"),
    ("Wo hätten Sie gehangen?", "Wo hättet ihr gehangen?"),
    ("Ihr begegnet sie in der nächsten Doppelstunde wieder.",
     "Ihr begegnet ihr in der nächsten Doppelstunde wieder."),

    # ---------------------------------------------------- DS 4
    ("Schreiben Sie einen Satz auf, den Sie niemandem zeigen müssen",
     "Schreib einen Satz auf, den du niemandem zeigen musst"),
    ("„Behalten Sie diesen Zettel.", "„Behaltet diesen Zettel."),
    ("Der Zettel gehört Ihnen.", "Der Zettel gehört dir."),
    ("Fünf solcher Wege liegen jetzt auf Ihren Tischen.",
     "Fünf solcher Wege liegen jetzt auf euren Tischen."),
    ("Findet in Ihrem Text den Satz Ihrer Person", "Findet in eurem Text den Satz eurer Person"),
    ("Den Zwischenraum füllen Sie.", "Den Zwischenraum füllt ihr."),
    ("Und ich verspreche Ihnen", "Und ich verspreche euch"),
    ("Halten Sie das bitte aus, bis alle fünf hängen.", "Haltet das bitte aus, bis alle fünf hängen."),
    ("Ihre Chance:", "Eure Chance:"),
    ("„Ein Satz noch, dann gehen Sie.", "„Ein Satz noch, dann geht ihr."),
    ("diesmal bringen Sie Ihre Handys mit — Sie bauen eine Ausstellung",
     "diesmal bringt ihr eure Handys mit — ihr baut eine Ausstellung"),
    ("diesmal bringen Sie Ihre Handys mit, Sie bauen eine Ausstellung",
     "diesmal bringt ihr eure Handys mit, ihr baut eine Ausstellung"),

    # ---------------------------------------------------- DS 5
    ("Beschreiben Sie, was auf dem Bild passiert.", "Beschreibt, was auf dem Bild passiert."),
    ("Lesen Sie den Bibelvers.", "Lest den Bibelvers."),
    ("Zehn Minuten Hintergrund, dann übernehmen Sie.", "Zehn Minuten Hintergrund, dann übernehmt ihr."),
    ("und die suchen Sie jetzt selbst", "und die sucht ihr jetzt selbst"),
    ("Zu jeder Behauptung auf Ihrer Tafel muss ich fragen dürfen",
     "Zu jeder Behauptung auf eurer Tafel muss ich fragen dürfen"),
    ("und Sie müssen es mir zeigen können", "und ihr müsst es mir zeigen können"),
    ("„Sie gehen jetzt schweigend durch die Ausstellung", "„Ihr geht jetzt schweigend durch die Ausstellung"),
    ("Sie gehen schweigend durch die Ausstellung", "Ihr geht schweigend durch die Ausstellung"),
    ("und kleben Ihre Punkte auf die", "und klebt eure Punkte auf die"),
    ("die Sie am meisten beschäftigen", "die euch am meisten beschäftigen"),
    ("Ich möchte, dass Sie den Unterschied nicht bloß spüren, sondern benennen können.",
     "Ich möchte, dass ihr den Unterschied nicht bloß spürt, sondern benennen könnt."),
    ("Sie bringen Ihre Handys mit.", "Ihr bringt eure Handys mit."),
    ("über die in Ihren Familien vermutlich unterschiedlich gedacht wird, und Sie werden am Ende etwas tun",
     "über die in euren Familien vermutlich unterschiedlich gedacht wird, und ihr werdet am Ende etwas tun"),
    ("Und Sie werden am Ende etwas tun", "Und ihr werdet am Ende etwas tun"),
    ("Sie schreiben eigene Bekenntnissätze.", "Ihr schreibt eigene Bekenntnissätze."),

    # ---------------------------------------------------- DS 6
    ("Was fällt Ihnen auf, wenn Sie diese drei Schlagzeilen und das Plakat nebeneinander sehen?",
     "Was fällt euch auf, wenn ihr diese drei Schlagzeilen und das Plakat nebeneinander seht?"),
    ("„Sie merken, dass ich vorsichtig formuliere, und ich sage Ihnen auch warum.",
     "„Ihr merkt, dass ich vorsichtig formuliere, und ich sage euch auch warum."),
    ("die manche Menschen in Ihren Familien wählen", "die manche Menschen in euren Familien wählen"),
    ("Ich werde Ihnen nicht sagen, was Sie wählen sollen",
     "Ich werde euch nicht sagen, was ihr wählen sollt"),
    ("Und dann prüfen Sie das mit dem Werkzeug, das Sie sich in den letzten Wochen selbst gebaut haben.",
     "Und dann prüft ihr das mit dem Werkzeug, das ihr euch in den letzten Wochen selbst gebaut habt."),
    ("und ich möchte, dass Sie merken, warum", "und ich möchte, dass ihr merkt, warum"),
    ("Röm 13, Sie erinnern sich.", "Röm 13, ihr erinnert euch."),
    ("Römer 13, Sie erinnern sich.", "Römer 13, ihr erinnert euch."),
    ("Wo verläuft die Ihrer Meinung nach?", "Wo verläuft die eurer Meinung nach?"),
    ("„Sie haben in der dritten Doppelstunde einen Satzanfang übersetzt, den ich Ihnen jetzt zurückgebe",
     "„Ihr habt in der dritten Doppelstunde einen Satzanfang übersetzt, den ich euch jetzt zurückgebe"),
    ("Sie sind jetzt dran.", "Ihr seid jetzt dran."),
    ("jede Ihrer Thesen muss sich an einer Bibelstelle festmachen lassen",
     "jede eurer Thesen muss sich an einer Bibelstelle festmachen lassen"),
    ("Und jetzt schauen Sie bitte einmal zur Tür.", "Und jetzt schaut bitte einmal zur Tür."),
    ("einige von Ihnen haben mich gefragt, warum", "einige von euch haben mich gefragt, warum"),
    ("„Das haben Sie vor drei Wochen aufgeschrieben.", "„Das habt ihr vor drei Wochen aufgeschrieben."),
    ("Es war Ihre Antwort auf die", "Es war eure Antwort auf die"),
    ("Lesen Sie es noch einmal — und lesen Sie es diesmal nicht als Geschichte",
     "Lest es noch einmal — und lest es diesmal nicht als Geschichte"),
    ("Genau das haben Sie heute getan.", "Genau das habt ihr heute getan."),
    ("Kleben Sie bitte noch einmal auf das Plakat", "Klebt bitte noch einmal auf das Plakat"),
    ("Kleben Sie noch einmal auf das Plakat", "Klebt noch einmal auf das Plakat"),
    ("Ich sage Ihnen nicht, wo Sie kleben sollen.", "Ich sage euch nicht, wo ihr kleben sollt."),
    ("Ich möchte nur, dass Sie sehen, ob dieser Vormittag etwas verschoben hat.",
     "Ich möchte nur, dass ihr seht, ob dieser Vormittag etwas verschoben hat."),
    ("Kleben Sie einen Punkt auf die Linie", "Klebt einen Punkt auf die Linie"),
    ("Ihre eigene Antwort aus Doppelstunde 4", "Eure eigene Antwort aus Doppelstunde 4"),

    # ------------------------------------------- Leistungsnachweis (Einzelarbeit)
    ("Nennen Sie zwei Bereiche", "Nenne zwei Bereiche"),
    ("und geben Sie jeweils die rechtliche Grundlage an", "und gib jeweils die rechtliche Grundlage an"),
    ("Erläutern Sie, warum die Barmer", "Erläutere, warum die Barmer"),
    ("Beschreiben Sie an einem Beispiel", "Beschreibe an einem Beispiel"),
    ("und erklären Sie, warum die Kirche", "und erkläre, warum die Kirche"),
    ("Nehmen Sie begründet Stellung.", "Nimm begründet Stellung."),
    ("Beziehen Sie ein historisches und ein gegenwärtiges Beispiel ein.",
     "Beziehe ein historisches und ein gegenwärtiges Beispiel ein."),

    # ------------------------------------------- Praesentation, Folientexte
    ("SEHEN — Notieren Sie fünf Dinge, die Sie auf dem Bild tatsächlich sehen.",
     "SEHEN — Notiere fünf Dinge, die du auf dem Bild tatsächlich siehst."),
    ("VERMISSEN — Notieren Sie eine Sache, die Sie in einer Kirche erwarten würden",
     "VERMISSEN — Notiere eine Sache, die du in einer Kirche erwarten würdest"),
    ("FRAGEN — Formulieren Sie die eine Frage, die Sie den Menschen auf diesem Bild stellen würden.",
     "FRAGEN — Formuliere die eine Frage, die du den Menschen auf diesem Bild stellen würdest."),
    # ------------------------- Nachtrag: durch Tags/Zeilenumbrueche zerschnitten
    ("Findet in Ihrem Text den ", "Findet in eurem Text den "),
    ("Satz Ihrer Person</strong>", "Satz eurer Person</strong>"),
    ("Sie gehen jetzt schweigend durch die Ausstellung", "Ihr geht jetzt schweigend durch die Ausstellung"),
    ("du auf dem Bild tatsächlich sehen.", "du auf dem Bild tatsächlich siehst."),
    ('"Sie schreiben eigene "', '"Ihr schreibt eigene "'),
    ('„Sie haben in der dritten Doppelstunde einen Satzanfang übersetzt, den ich "',
     '„Ihr habt in der dritten Doppelstunde einen Satzanfang übersetzt, den ich "'),
    ('"Ihnen jetzt zurückgebe:', '"euch jetzt zurückgebe:'),
    ('Lesen Sie es noch einmal — und "', 'Lest es noch einmal — und "'),
    ('"lesen Sie es diesmal nicht als Geschichte', '"lest es diesmal nicht als Geschichte'),

    # ------------------------------- Regie- und Hinweistext an die LEHRKRAFT
    # Die Lehrerin hat ausdruecklich erlaubt, auch sie zu duzen.
    ("Fragen Sie zurück, wenn Sie zu einer einzelnen Doppelstunde die Materialien als druckfertige Blätter brauchen",
     "Frag zurück, wenn du zu einer einzelnen Doppelstunde die Materialien als druckfertige Blätter brauchst"),
    ("wenn Sie den Gesetzeswortlaut zeigen möchten, setzen Sie ihn von",
     "wenn du den Gesetzeswortlaut zeigen möchtest, setz ihn von"),
    ("Wenn Sie mit dem vollständigen Originalwortlaut arbeiten möchten, finden Sie die Dokumente",
     "Wenn du mit dem vollständigen Originalwortlaut arbeiten möchtest, findest du die Dokumente"),
    ("Die Rolle von Charlotte S. besetzen Sie bewusst.", "Die Rolle von Charlotte S. besetzt du bewusst."),
    ("Geben Sie diese Rolle jemandem, der sich behaupten kann, und greifen Sie als Moderation ein",
     "Gib diese Rolle jemandem, der sich behaupten kann, und greif als Moderation ein"),
    ("Was Sie vorbereiten und was Sie aushalten müssen", "Was du vorbereitest und was du aushalten musst"),
    ("Ihre Aufgabe ist es, die Begründungen einzufordern", "Deine Aufgabe ist es, die Begründungen einzufordern"),
    ("Fragen Sie die Klasse, ob der Text dadurch weniger wahr wird.",
     "Frag die Klasse, ob der Text dadurch weniger wahr wird."),
    ("Sagen Sie das ausdrücklich:", "Sag das ausdrücklich:"),
    ("Damit belohnen Sie Nachdenken statt Layout", "Damit belohnst du Nachdenken statt Layout"),
    ("geben Sie es einer Gruppe, die Ambivalenz aushält", "gib es einer Gruppe, die Ambivalenz aushält"),
    ("Setzen Sie diese Zahlen aus einer belegbaren Quelle ein und nennen Sie die Quelle auf der Folie",
     "Setz diese Zahlen aus einer belegbaren Quelle ein und nenne die Quelle auf der Folie"),
    ("Er zeigt der Klasse — und Ihnen —", "Er zeigt der Klasse — und dir —"),
    ("Prüfen Sie den Stand vor jedem Einsatz und tragen Sie das Datum ein.",
     "Prüf den Stand vor jedem Einsatz und trag das Datum ein."),
    ("Setzen Sie hier vor der Stunde eine aktuelle, belegte Zahl ein oder lassen Sie die Gruppe sie selbst suchen.",
     "Setz hier vor der Stunde eine aktuelle, belegte Zahl ein oder lass die Gruppe sie selbst suchen."),
    ("die Sie dazuschreiben", "die du dazuschreibst"),
    # ---------------------- Letzter Durchgang: Restliche Einzelstellen
    ("Bevor Sie M8b und M8c austeilen", "Bevor du M8b und M8c austeilst"),
    ("Was folgt daraus für Sie, hier, in dieser Stadt?", "Was folgt daraus für euch, hier, in dieser Stadt?"),
    ('"können, haben Sie schlecht erklärt', '"können, habt ihr schlecht erklärt'),
    ("diesmal bringen Sie Ihre Handys mit.", "diesmal bringt ihr eure Handys mit."),
    ('sondern dieses. Und Sie "', 'sondern dieses. Und ihr "'),
    ('"werden etwas tun, das seit 1934 kaum jemand mehr gemacht hat: Sie schreiben eigene "',
     '"werdet etwas tun, das seit 1934 kaum jemand mehr gemacht hat: Ihr schreibt eigene "'),
    ('"nur, dass Sie sehen, ob dieser Vormittag', '"nur, dass ihr seht, ob dieser Vormittag'),
]


def main():
    probe = "--probe" in sys.argv
    ungenutzt = [alt for alt, _ in ERSETZUNGEN]
    gesamt = 0
    for pfad in DATEIEN:
        inhalt = original = pfad.read_text(encoding="utf-8")
        treffer = 0
        for alt, neu in ERSETZUNGEN:
            if alt in inhalt:
                treffer += inhalt.count(alt)
                inhalt = inhalt.replace(alt, neu)
                if alt in ungenutzt:
                    ungenutzt.remove(alt)
        if inhalt != original and not probe:
            pfad.write_text(inhalt, encoding="utf-8")
        gesamt += treffer
        print(f"{pfad.name:38s} {treffer:4d} Ersetzungen")
    print(f"\ngesamt: {gesamt} Ersetzungen")
    if ungenutzt:
        print(f"\nNICHT GEFUNDEN ({len(ungenutzt)}) — pruefen:")
        for alt in ungenutzt:
            print("   ", alt[:95])
    return len(ungenutzt)


if __name__ == "__main__":
    sys.exit(1 if main() else 0)
