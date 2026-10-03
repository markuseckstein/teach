#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Erzeugt für jede Doppelstunde ein einseitiges Word-Dokument für SuS,
die krank zu Hause waren: Grundidee, vier Merkpunkte, ein Merksatz,
Anschluss an die nächste Stunde und ein freiwilliger Bett-Auftrag.

    python3 nacharbeit.py

Zielordner: word/nacharbeit/
Alle Links wurden am 5. September 2026 auf Erreichbarkeit geprüft.
"""
import pathlib
from docx import Document
from docx.enum.table import WD_TABLE_ALIGNMENT
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.oxml import OxmlElement
from docx.oxml.ns import qn
from docx.shared import Cm, Pt, RGBColor

ROOT = pathlib.Path(__file__).resolve().parent
ZIEL = ROOT / "word" / "nacharbeit"

TINTE   = RGBColor(0x15, 0x1A, 0x21)
AKZENT  = RGBColor(0x1E, 0x3A, 0x5F)
EPOCHE  = RGBColor(0x6A, 0x59, 0x36)
JETZT   = RGBColor(0x26, 0x5F, 0x53)
GRAU    = RGBColor(0x5A, 0x63, 0x6E)
SERIF   = "Cambria"
GROTESK = "Calibri"


def schrift(lauf, name=GROTESK, groesse=10.2, fett=False, kursiv=False, farbe=TINTE):
    lauf.font.name = name
    lauf.font.size = Pt(groesse)
    lauf.font.bold = fett
    lauf.font.italic = kursiv
    lauf.font.color.rgb = farbe
    rpr = lauf._element.get_or_add_rPr()
    rfonts = rpr.find(qn("w:rFonts"))
    if rfonts is None:
        rfonts = OxmlElement("w:rFonts")
        rpr.append(rfonts)
    for attr in ("w:ascii", "w:hAnsi", "w:cs"):
        rfonts.set(qn(attr), name)
    return lauf


def absatz(dok_oder_zelle, text="", vor=0, nach=4, **kw):
    p = dok_oder_zelle.add_paragraph()
    p.paragraph_format.space_before = Pt(vor)
    p.paragraph_format.space_after = Pt(nach)
    p.paragraph_format.line_spacing = 1.05
    if text:
        schrift(p.add_run(text), **kw)
    return p


def faerben(zelle, hexfarbe):
    shd = OxmlElement("w:shd")
    shd.set(qn("w:val"), "clear")
    shd.set(qn("w:fill"), hexfarbe)
    zelle._tc.get_or_add_tcPr().append(shd)


def rahmenlos(tabelle):
    tbl_pr = tabelle._tbl.tblPr
    borders = OxmlElement("w:tblBorders")
    for kante in ("top", "left", "bottom", "right", "insideH", "insideV"):
        el = OxmlElement(f"w:{kante}")
        el.set(qn("w:val"), "none")
        el.set(qn("w:sz"), "0")
        borders.append(el)
    tbl_pr.append(borders)


def kasten(dok, hexfarbe="EEF1F4"):
    """Einzellige Tabelle als farbiger Kasten. Gibt die Zelle zurück."""
    t = dok.add_table(rows=1, cols=1)
    t.alignment = WD_TABLE_ALIGNMENT.LEFT
    rahmenlos(t)
    z = t.cell(0, 0)
    faerben(z, hexfarbe)
    z.paragraphs[0]._element.getparent().remove(z.paragraphs[0]._element)
    return z


def hyperlink(p, url, text):
    teil = p.part
    r_id = teil.relate_to(
        url,
        "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink",
        is_external=True,
    )
    link = OxmlElement("w:hyperlink")
    link.set(qn("r:id"), r_id)
    r = OxmlElement("w:r")
    rpr = OxmlElement("w:rPr")
    u = OxmlElement("w:u"); u.set(qn("w:val"), "single"); rpr.append(u)
    c = OxmlElement("w:color"); c.set(qn("w:val"), "1E3A5F"); rpr.append(c)
    f = OxmlElement("w:rFonts")
    for a in ("w:ascii", "w:hAnsi", "w:cs"):
        f.set(qn(a), GROTESK)
    rpr.append(f)
    sz = OxmlElement("w:sz"); sz.set(qn("w:val"), "19"); rpr.append(sz)
    r.append(rpr)
    t = OxmlElement("w:t")
    t.text = text
    r.append(t)
    link.append(r)
    p._p.append(link)
    return p


DOPPELSTUNDEN = [
    dict(
        nr=1,
        titel="Wo steckt Kirche im Staat?",
        unter="Eine Spurensuche im eigenen Alltag",
        worum=(
            "Wir haben angeschaut, wo dir die Kirche im ganz normalen staatlichen Alltag begegnet, "
            "ohne dass ein Kirchturm in Sicht ist: auf der Lohnabrechnung deiner Eltern, in deinem "
            "Stundenplan, im Rettungswagen, bei der Vereidigung im Bundestag. Dahinter steckt eine "
            "Merkwürdigkeit: Das Grundgesetz sagt klar, dass es in Deutschland keine Staatskirche gibt "
            "— und trotzdem arbeiten Staat und Kirche an vielen Stellen eng zusammen. Man nennt das "
            "die „hinkende Trennung“."
        ),
        punkte=[
            ("Kein Staat ohne Grenze zur Kirche.",
             "Artikel 137 der Weimarer Reichsverfassung gilt über Artikel 140 Grundgesetz weiter: "
             "„Es besteht keine Staatskirche.“ Keine Religion regiert mit, niemand muss glauben."),
            ("Trotzdem vier feste Berührungspunkte.",
             "Religionsunterricht (ordentliches Lehrfach nach Art. 7 Abs. 3 GG — deshalb sitzt du hier), "
             "Kirchensteuer (das Finanzamt zieht sie für die Kirchen ein), Militärseelsorge "
             "(Pfarrerinnen und Pfarrer bei der Bundeswehr) und Diakonie (Kitas, Kliniken, Pflege "
             "— große Teile davon zahlt der Staat)."),
            ("Schon der erste Satz des Grundgesetzes.",
             "Die Präambel beginnt mit „Im Bewusstsein seiner Verantwortung vor Gott und den "
             "Menschen …“ — in einem Staat, der keine Staatskirche hat."),
            ("Darüber kann man streiten.",
             "Genau das haben wir am Ende gemacht: Sollte der Staat die Kirchensteuer einziehen? "
             "Gehört Religionsunterricht an eine staatliche Schule? Es gibt gute Argumente auf "
             "beiden Seiten — dein eigenes Urteil zählt hier mehr als das richtige."),
        ],
        merksatz="Keine Staatskirche — aber auch keine Mauer. Deutschland trennt Staat und Kirche hinkend: getrennt in der Macht, verbunden in der Arbeit.",
        weiter=(
            "In der nächsten Stunde springen wir ins Jahr 1933 und schauen, was passiert, wenn dieses "
            "Naheverhältnis von Kirche und Staat auf eine Diktatur trifft. Du brauchst dafür nur die "
            "vier Berührungspunkte von oben im Kopf — sonst nichts."),
        freiwillig=[
            ("Video · 4 Minuten", "„Kirchensteuer einfach erklärt“ von explainity. Danach weißt du, "
             "wer sie zahlt, wie hoch sie ist und wohin das Geld geht.",
             "https://www.youtube.com/watch?v=6y4wRMtxl1U"),
            ("Nachlesen · 3 Minuten", "Das Junge Politik-Lexikon der Bundeszentrale für politische "
             "Bildung, Stichwort Kirchensteuer.",
             "https://www.bpb.de/kurz-knapp/lexika/das-junge-politik-lexikon/320628/kirchensteuer/"),
        ],
        farbe="EEF2F6",
    ),
    dict(
        nr=2,
        titel="1933: Die Kirche jubelt mit",
        unter="Warum die Frage nicht lautet „Wie konnten die nur?“",
        worum=(
            "Wir haben Quellen aus dem Jahr 1933 gelesen — und die sind unbequem: Große Teile der "
            "evangelischen Kirche haben die Machtübernahme der Nationalsozialisten nicht nur "
            "hingenommen, sondern begeistert begrüßt. Gottesdienste mit Hakenkreuzfahnen waren "
            "keine Ausnahme. Die eigentliche Frage der Stunde war deshalb nicht „Wie konnten die "
            "nur?“, sondern die schwierigere: „Was war daran für die Menschen damals so "
            "überzeugend?“"
        ),
        punkte=[
            ("Die „Deutschen Christen“.",
             "Eine Bewegung innerhalb der evangelischen Kirche, die Christentum und "
             "nationalsozialistische Weltanschauung verschmelzen wollte. Bei der Kirchenwahl im Juli "
             "1933 gewann sie rund zwei Drittel der Sitze. Ludwig Müller wurde Reichsbischof."),
            ("Der Arierparagraph in der Kirche.",
             "Getaufte Christen jüdischer Herkunft sollten aus dem Pfarramt entfernt werden. Damit "
             "war der Streit kein politischer mehr, sondern ein Streit darüber, was die Taufe wert ist."),
            ("Vier Gründe, die wir herausgearbeitet haben.",
             "Angst vor dem Kommunismus, Sehnsucht nach Ordnung nach den unruhigen Weimarer Jahren, "
             "eine lange Tradition von „Thron und Altar“ — und ein Judenhass, den es in der "
             "Kirche lange vor 1933 gab."),
            ("Zwei Bibelverse, zwei Welten.",
             "Römer 13,1: „Jedermann sei untertan der Obrigkeit.“ Apostelgeschichte 5,29: „Man "
             "muss Gott mehr gehorchen als den Menschen.“ Beide Seiten konnten sich auf die Bibel "
             "berufen. Das ist der Kern der ganzen Einheit."),
        ],
        merksatz="Beide Seiten hatten einen Bibelvers. Die Frage war nicht, was in der Bibel steht — sondern welchen Vers man lesen wollte.",
        weiter=(
            "Nächste Stunde geht es um das Jahr 1934 und um den Satz, mit dem eine Minderheit "
            "widersprochen hat. Wenn du Römer 13 und Apostelgeschichte 5 parat hast, kommst du "
            "mühelos mit."),
        freiwillig=[
            ("Nachlesen · 10 Minuten", "„Kirchen im NS-Regime“ beim Lebendigen Museum Online des "
             "Deutschen Historischen Museums. Lies nur den Abschnitt zur evangelischen Kirche — der "
             "Rest kommt später dran.",
             "https://www.dhm.de/lemo/kapitel/ns-regime/innenpolitik/kirchen.html"),
            ("Wenn du mehr willst", "Der Wikipedia-Artikel „Deutsche Christen“ zeigt Plakate und "
             "Wahlergebnisse von 1933.",
             "https://de.wikipedia.org/wiki/Deutsche_Christen"),
        ],
        farbe="F2F0E8",
    ),
    dict(
        nr=3,
        titel="Barmen gegen Ansbach",
        unter="Zwei Texte, dreizehn Tage Abstand, dieselbe Kirche",
        worum=(
            "Diese Stunde war ein Rollenspiel: eine Streitrunde im Juni 1934. Der Hintergrund sind zwei "
            "Texte, die kurz nacheinander erschienen und sich gegenseitig widersprechen. Der eine wurde "
            "berühmt, der andere fast vergessen — obwohl beide von evangelischen Theologen stammen. "
            "Am Ende ging es um einen bayerischen Landesbischof, an dem sich zeigt, dass historische "
            "Urteile selten eindeutig ausfallen."
        ),
        punkte=[
            ("Barmer Theologische Erklärung, 31. Mai 1934.",
             "Das Gründungsdokument der Bekennenden Kirche, maßgeblich von Karl Barth formuliert. "
             "These 1: Jesus Christus ist das eine Wort Gottes, dem die Kirche zu gehorchen hat — "
             "und keiner anderen Macht, keinem Führer, keinem Volk."),
            ("Die Sprachform, die du dir merken solltest.",
             "Jede These hat zwei Teile: eine positive Aussage, woran man sich bindet, und dann die "
             "Verwerfung — „Wir verwerfen die falsche Lehre, als ob …“. Diesen Satzanfang "
             "brauchst du in der letzten Doppelstunde wieder."),
            ("Ansbacher Ratschlag, 11. Juni 1934.",
             "Die Gegenposition, unter anderem von Paul Althaus und Werner Elert: Volk, Familie und "
             "Führer seien von Gott gegebene Ordnungen; Hitler wird darin als frommer Oberherr "
             "bezeichnet. Dreizehn Tage nach Barmen, dieselbe Kirche."),
            ("Der Fall Hans Meiser.",
             "Der bayerische Landesbischof verhinderte, dass seine Landeskirche in die Reichskirche "
             "eingegliedert wurde — dafür stand er 1934 unter Hausarrest, und Tausende protestierten "
             "in Franken. Zugleich hat er sich judenfeindlich geäußert und gegen die Judenverfolgung "
             "nie öffentlich protestiert. Deshalb wurden Meiser-Straßen umbenannt, in München 2010."),
        ],
        merksatz="Die Grenze verlief 1934 nicht zwischen Kirche und Staat. Sie verlief mitten durch die Kirche.",
        weiter=(
            "In der nächsten Stunde hängen wir fünf Menschen an eine Wäscheleine zwischen "
            "Anpassung und Widerstand — Meiser ist einer davon. Es hilft, wenn du seine zwei Seiten "
            "im Kopf hast."),
        freiwillig=[
            ("Nachlesen · 5 Minuten", "Die Barmer Theologische Erklärung bei der EKD, mit dem "
             "Wortlaut der ersten These.",
             "https://www.ekd.de/barmer-theologische-erklarung-11292.htm"),
            ("Zum Anschauen", "Eine virtuelle Ausstellung der Deutschen Digitalen Bibliothek zu Barmen "
             "— Originalfotos und Dokumente, gut zum Durchklicken.",
             "https://ausstellungen.deutsche-digitale-bibliothek.de/jubilaeum-bekenntnis-barmen"),
        ],
        farbe="F2F0E8",
    ),
    dict(
        nr=4,
        titel="Was hätte ich getan?",
        unter="Fünf Biografien zwischen Anpassung und Widerstand",
        worum=(
            "Wir haben fünf evangelische Christen angeschaut, die alle dieselbe Zeit erlebt haben und "
            "sich sehr unterschiedlich verhalten haben. Ihre Karten hängen an einer Leine im "
            "Klassenzimmer, zwischen „Anpassung“ und „Widerstand“. Danach haben wir "
            "aufgehört, einzelne Menschen zu benoten, und stattdessen gefragt: Was hatten die Leute "
            "rechts an der Leine, was die Leute links nicht hatten?"
        ),
        punkte=[
            ("Dietrich Bonhoeffer (1906–1945).",
             "Theologe, Bekennende Kirche, Kontakt zur Widerstandsgruppe um das Attentat vom 20. Juli. "
             "1943 verhaftet, am 9. April 1945 im KZ Flossenbürg hingerichtet — keine zwei "
             "Autostunden von hier."),
            ("Martin Niemöller (1892–1984) und Paul Schneider (1897–1939).",
             "Niemöller war zuerst deutschnational, gründete dann den Pfarrernotbund gegen den "
             "Arierparagraphen und saß ab 1937 in Sachsenhausen und Dachau. Schneider, der "
             "„Prediger von Buchenwald“, verweigerte den Hitlergruß und predigte im KZ aus dem "
             "Fenster seiner Arrestzelle, bis er 1939 zu Tode gequält wurde."),
            ("Elisabeth Schmitz (1893–1977) und Hans Meiser (1881–1956).",
             "Schmitz schrieb 1935 die Denkschrift „Zur Lage der deutschen Nichtarier“ — eine "
             "der schärfsten kirchlichen Anklagen gegen die Judenverfolgung — und wurde jahrzehntelang "
             "vergessen. Meiser schützte seine Landeskirche und schwieg zur Judenverfolgung."),
            ("Das Ergebnis, das wir aufgehoben haben.",
             "In Vierergruppen haben wir gesammelt, was Widerspruch wahrscheinlicher macht. Oft "
             "genannt: früh anfangen, nicht allein sein, eine Überzeugung haben, die nicht vom "
             "Erfolg abhängt, gut informiert sein, schon einmal geübt haben, Nein zu sagen. Diese "
             "Liste kommt in der letzten Doppelstunde zurück."),
        ],
        merksatz="Menschen sind nicht von Anfang an Helden oder Feiglinge. Sie werden es — Entscheidung für Entscheidung.",
        weiter=(
            "Nächste Stunde springen wir vierzig Jahre weiter: gleiches Land, andere Diktatur, "
            "DDR. Bring dein Handy mit, wir bauen eine Ausstellung."),
        freiwillig=[
            ("Nachlesen · 5 Minuten", "Die Kurzbiografie Dietrich Bonhoeffers beim Lebendigen Museum "
             "Online des Deutschen Historischen Museums.",
             "https://www.dhm.de/lemo/biografie/dietrich-bonhoeffer"),
            ("Zum Nachdenken · 2 Minuten", "Such den Text „Als die Nazis die Kommunisten holten …“ "
             "von Martin Niemöller und lies ihn einmal laut. Und dann schreib dir einen Satz auf, den "
             "du niemandem zeigen musst: Wobei würde ich heute wegschauen?", None),
        ],
        farbe="F2F0E8",
    ),
    dict(
        nr=5,
        titel="Kerzen statt Steine",
        unter="Kirche in der DDR — und wie aus Gebeten eine Revolution wurde",
        worum=(
            "In dieser Stunde haben wir eine Ausstellung gebaut: fünf Tafeln zur Rolle der Kirche in "
            "der DDR, die jetzt im Klassenzimmer hängen. Der Ausgangspunkt war ein Aufnäher — ein "
            "Stück Stoff mit dem Bild eines Schmieds und dem Satz „Schwerter zu Pflugscharen“. "
            "Wer ihn in den 1980er Jahren auf der Jacke trug, konnte von der Schule fliegen. Ein Staat, "
            "der Angst vor einem Friedenssymbol hat: Damit war das Thema gesetzt."
        ),
        punkte=[
            ("Die Kirche war die einzige Nische.",
             "Der Staat war atheistisch, die Kirche geduldet. Sie war der einzige Ort in der DDR, an "
             "dem man Sätze sagen durfte, die sonst nirgends erlaubt waren — über Frieden, Umwelt, "
             "Reisefreiheit."),
            ("Der Preis dafür war hoch.",
             "Wer sich in der Kirche engagierte, kam oft nicht zum Abitur. Die Jugendweihe stand als "
             "gesellschaftlicher Zwang gegen die Konfirmation. Und die Staatssicherheit schleuste "
             "Inoffizielle Mitarbeiter bis in kirchliche Leitungsämter ein — Schuld gab es also "
             "auch innerhalb der Kirche."),
            ("„Schwerter zu Pflugscharen“.",
             "Der Satz stammt aus Micha 4,3. Das Bild geht auf eine Skulptur zurück, die die "
             "Sowjetunion den Vereinten Nationen geschenkt hatte — der Staat konnte das Motiv also "
             "schlecht verbieten und tat es doch."),
            ("Von den Friedensgebeten zur Friedlichen Revolution.",
             "Ab 1982 gab es montags Friedensgebete in der Leipziger Nikolaikirche. Am 4. September 1989 "
             "wurde daraus die erste große Montagsdemonstration; am 9. Oktober 1989 gingen rund "
             "70.000 Menschen auf die Straße — ohne einen einzigen Stein."),
        ],
        merksatz="Aus einem Gebetsraum wurde ein Schutzraum, aus dem Schutzraum eine Revolution. Und keine Kirche hatte das geplant.",
        weiter=(
            "Die letzte Doppelstunde bringt alles in die Gegenwart: Wo verläuft heute für Christen "
            "die Grenze? Du schreibst dann eine eigene These in der Barmer Sprachform. Handy "
            "mitbringen."),
        freiwillig=[
            ("Videoclips und Fotos · 10 Minuten", "„Friedensgebete und Montagsdemonstrationen“ auf "
             "jugendopposition.de — ein Projekt der Bundeszentrale für politische Bildung mit "
             "Zeitzeugenfilmen. Ideal fürs Bett.",
             "https://www.jugendopposition.de/jugendopposition/themen/565097/friedensgebete-und-montagsdemonstrationen/"),
            ("Kleiner Rechercheauftrag", "Finde heraus, was am 9. Oktober 1989 in Leipzig anders lief "
             "als in Peking im Juni desselben Jahres. Ein Satz genügt — sag ihn in der nächsten "
             "Stunde.", None),
        ],
        farbe="EBF2EF",
    ),
    dict(
        nr=6,
        titel="Barmen 2026?",
        unter="Wo verläuft heute die Grenze?",
        worum=(
            "Die letzte Doppelstunde bringt die Geschichte in die Gegenwart. Wir haben gelesen, was "
            "die Kirchen heute zum völkischen Nationalismus sagen, und ebenso ernst genommen, was "
            "Christen dagegen einwenden, die finden, die Kirche solle sich aus der Politik heraushalten. "
            "Danach hat jede Gruppe eine eigene These in der Barmer Sprachform geschrieben und an die "
            "Klassenzimmertür gehängt. Damit es kein Missverständnis gibt: Wir bewerten keine "
            "Parteien und geben keine Wahlempfehlung. Wir schauen an, was die Kirchen sagen, "
            "warum sie es sagen — und ob sie es sagen dürfen."
        ),
        punkte=[
            ("Was die Kirchen sagen.",
             "Die EKD-Synode 2023 und die deutschen Bischöfe 2024 haben erklärt, völkischer "
             "Nationalismus und Christentum seien unvereinbar. Begründung ist kein Parteiprogramm, "
             "sondern das Menschenbild: Jeder Mensch ist Gottes Ebenbild, nicht Mitglied eines Volkes "
             "erster oder zweiter Klasse."),
            ("Was dagegen eingewendet wird.",
             "Die Kirche sei zu politischer Neutralität verpflichtet; in ihren Gemeinden sitzen "
             "Menschen aller Parteien; wer sie ausschließt, verliert sie. Dieses Argument haben wir "
             "in seiner stärksten Form behandelt, nicht als Strohmann."),
            ("Die Barmer Sprachform, jetzt selbst benutzt.",
             "Erst der positive Teil: „Weil …, gilt für uns: …“. Dann die Verwerfung: "
             "„Wir verwerfen die falsche Lehre, als ob …“. Jede These musste sich an einer "
             "Bibelstelle festmachen lassen — sonst ist es nur Meinung."),
            ("Die Härteprüfung, an der sich alles entscheidet.",
             "Zwei Fragen musste jede These bestehen. Erstens: Richtet sie sich gegen einen Satz und "
             "nicht gegen Menschen? Zweitens: Würde sie auch dann noch gelten, wenn morgen eine ganz "
             "andere Partei regiert? Wenn nein, wurde überarbeitet. Genau das unterscheidet eine "
             "Bekenntnisthese von einem Wahlplakat."),
        ],
        merksatz="1933 war „Die Kirche soll sich raushalten“ die Mehrheitsmeinung. 1934 hat eine Minderheit gesagt: Es gibt eine Grenze. Wo sie verläuft, musst du selbst beantworten — aber du musst es begründen können.",
        weiter=(
            "Damit ist die Einheit zu Ende. Wenn du magst, schreib deine eigene These nach dem Muster "
            "oben und bring sie mit — sie kommt dann noch an die Tür. Und wenn du eine Frage hast, "
            "die du in der Klasse nicht stellen wolltest: Frag mich einfach direkt."),
        freiwillig=[
            ("Nachlesen · 5 Minuten", "Die Themenseite der EKD „Kirche gegen Rassismus, "
             "Rechtspopulismus und Rechtsextremismus“ — die Originalquelle, nicht eine "
             "Zusammenfassung.",
             "https://www.ekd.de/kirche-gegen-rechtspopulismus-und-rechtsextremismus-49866.htm"),
            ("Rechercheauftrag · 10 Minuten", "Such eine Gruppe oder Initiative bei dir vor Ort — "
             "Gemeinde, Verein, Schule —, bei der man sich für eines der sechs Themen der Stunde "
             "einsetzen kann. Notier den Namen und was man dort konkret tun kann.", None),
        ],
        farbe="EBF2EF",
    ),
]


def dokument(ds):
    dok = Document()
    for s in dok.sections:
        s.top_margin = Cm(1.9)
        s.bottom_margin = Cm(1.6)
        s.left_margin = Cm(2.1)
        s.right_margin = Cm(2.1)

    normal = dok.styles["Normal"]
    normal.font.name = GROTESK
    normal.font.size = Pt(10.2)

    kopf = absatz(dok, "Kirche und Staat  ·  Doppelstunde %d  ·  Wenn du krank warst"
                  % ds["nr"], nach=1, groesse=8.5, fett=True, farbe=GRAU)
    kopf.runs[0].font.all_caps = True
    kopf.runs[0].font.spacing = Pt(0.6)

    absatz(dok, ds["titel"], nach=0, name=SERIF, groesse=21, fett=True, farbe=AKZENT)
    absatz(dok, ds["unter"], nach=10, name=SERIF, groesse=11.5, kursiv=True, farbe=GRAU)

    z = kasten(dok, ds["farbe"])
    p = absatz(z, "Fünf Minuten lesen — mehr ist das nicht.", vor=6, nach=3,
               groesse=9.5, fett=True, farbe=AKZENT)
    p.runs[0].font.all_caps = True
    absatz(z, ds["worum"], nach=7, groesse=10.2)

    absatz(dok, "Das musst du wissen", vor=12, nach=4, name=SERIF, groesse=12.5,
           fett=True, farbe=AKZENT)
    for i, (fett, rest) in enumerate(ds["punkte"], 1):
        p = absatz(dok, nach=5)
        p.paragraph_format.left_indent = Cm(0.75)
        p.paragraph_format.first_line_indent = Cm(-0.75)
        schrift(p.add_run("%d  " % i), groesse=10.2, fett=True, farbe=AKZENT)
        schrift(p.add_run(fett + " "), groesse=10.2, fett=True)
        schrift(p.add_run(rest), groesse=10.2)

    z = kasten(dok, "E8E2CF")
    p = absatz(z, "Der eine Satz", vor=6, nach=2, groesse=9, fett=True, farbe=EPOCHE)
    p.runs[0].font.all_caps = True
    absatz(z, ds["merksatz"], nach=7, name=SERIF, groesse=12, kursiv=True, farbe=TINTE)

    absatz(dok, "So geht es weiter", vor=12, nach=3, name=SERIF, groesse=12.5,
           fett=True, farbe=AKZENT)
    absatz(dok, ds["weiter"], nach=8)

    absatz(dok, "Freiwillig — wenn dir langweilig ist", vor=6, nach=3, name=SERIF,
           groesse=12.5, fett=True, farbe=JETZT)
    for label, text, url in ds["freiwillig"]:
        p = absatz(dok, nach=3)
        p.paragraph_format.left_indent = Cm(0.4)
        schrift(p.add_run(label + "  "), groesse=10, fett=True, farbe=JETZT)
        schrift(p.add_run(text), groesse=10.2)
        if url:
            p = absatz(dok, nach=6)
            p.paragraph_format.left_indent = Cm(0.4)
            hyperlink(p, url, url)

    p = absatz(dok, "Wenn etwas unklar ist: Frag mich — vor der Stunde, in der Stunde "
               "oder per Nachricht. Das ist keine Hausaufgabe und wird nicht eingesammelt. "
               "Du sollst nur nicht das Gefühl haben, dass du etwas verpasst hast.",
               vor=12, nach=0, groesse=9.5, kursiv=True, farbe=GRAU)
    p.alignment = WD_ALIGN_PARAGRAPH.LEFT
    return dok


def main():
    ZIEL.mkdir(parents=True, exist_ok=True)
    for ds in DOPPELSTUNDEN:
        name = "DS%d-Wenn-du-krank-warst.docx" % ds["nr"]
        dokument(ds).save(ZIEL / name)
        print("geschrieben:", (ZIEL / name).relative_to(ROOT))


if __name__ == "__main__":
    main()
