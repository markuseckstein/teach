#!/usr/bin/env python3
"""Baut die Beamer-Praesentation zur Einheit "Kirche und Staat".

Aufgenommen wird nur, was im Unterricht tatsaechlich an die Wand kommt:
Projektionsfolien, Arbeitsauftraege und Merksaetze. Kopiervorlagen, die
ausgeteilt werden (Expertenblaetter, Rollenkarten, Raster), fehlen bewusst.

Die Ueberleitungen stehen in den Notizen der jeweiligen Folie und sind damit
in der PowerPoint-Referentenansicht lesbar, waehrend die Klasse die Folie sieht.

Aufruf:  python3 build_pptx.py [ziel.pptx]
"""
import sys
import pathlib
from io import BytesIO

from pptx import Presentation
from pptx.util import Inches, Pt, Emu
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN, MSO_ANCHOR

from html2docx import normalize_jpeg

ROOT = pathlib.Path(__file__).parent
IMG = ROOT / "assets" / "img"

# ---------------------------------------------------------------- Gestaltung
PAPIER = RGBColor(0xF3, 0xF5, 0xF7)
BLATT = RGBColor(0xFF, 0xFF, 0xFF)
TINTE = RGBColor(0x15, 0x1A, 0x21)
GRAU = RGBColor(0x4B, 0x55, 0x63)
HELLGRAU = RGBColor(0x78, 0x82, 0x8F)
AKZENT = RGBColor(0x1E, 0x3A, 0x5F)
AKZENT_HELL = RGBColor(0xE5, 0xEC, 0xF4)
EPOCHE = RGBColor(0x6A, 0x59, 0x36)
EPOCHE_HELL = RGBColor(0xF0, 0xEB, 0xDF)
JETZT = RGBColor(0x26, 0x5F, 0x53)
JETZT_HELL = RGBColor(0xE2, 0xED, 0xEA)
SIGNAL = RGBColor(0x8C, 0x2A, 0x24)
SIGNAL_HELL = RGBColor(0xF6, 0xE7, 0xE5)
WEISS = RGBColor(0xFF, 0xFF, 0xFF)

# Auf einem Windows-Rechner sicher vorhandene Schriften. Die Schriften der
# HTML-Fassung (Zilla Slab, Source Sans 3, IBM Plex Mono) sind dort nicht
# installiert und wuerden stillschweigend ersetzt.
DISPLAY = "Georgia"
TEXT = "Calibri"
MONO = "Consolas"

BREITE = Inches(13.333)
HOEHE = Inches(7.5)
RAND = Inches(0.85)


# ------------------------------------------------------------------ Bausteine
def neue_folie(prs, hintergrund=PAPIER):
    folie = prs.slides.add_slide(prs.slide_layouts[6])  # leeres Layout
    fuellung = folie.background.fill
    fuellung.solid()
    fuellung.fore_color.rgb = hintergrund
    return folie


def kasten(folie, links, oben, breite, hoehe, farbe=None, rahmen=None,
           rahmenbreite=Pt(1)):
    form = folie.shapes.add_shape(1, links, oben, breite, hoehe)  # Rechteck
    form.shadow.inherit = False
    if farbe is None:
        form.fill.background()
    else:
        form.fill.solid()
        form.fill.fore_color.rgb = farbe
    if rahmen is None:
        form.line.fill.background()
    else:
        form.line.color.rgb = rahmen
        form.line.width = rahmenbreite
    form.text_frame.word_wrap = True
    return form


def text(folie, inhalt, links, oben, breite, hoehe, groesse=Pt(20),
         schrift=TEXT, farbe=TINTE, fett=False, kursiv=False,
         ausrichtung=PP_ALIGN.LEFT, anker=MSO_ANCHOR.TOP, zeilen=1.25,
         sperrung=None):
    feld = folie.shapes.add_textbox(links, oben, breite, hoehe)
    rahmen = feld.text_frame
    rahmen.word_wrap = True
    rahmen.vertical_anchor = anker
    for index, zeile in enumerate(inhalt.split("\n")):
        absatz = rahmen.paragraphs[0] if index == 0 else rahmen.add_paragraph()
        absatz.alignment = ausrichtung
        absatz.line_spacing = zeilen
        lauf = absatz.add_run()
        lauf.text = zeile
        zeichen = lauf.font
        zeichen.name = schrift
        zeichen.size = groesse
        zeichen.color.rgb = farbe
        zeichen.bold = fett
        zeichen.italic = kursiv
        if sperrung is not None:
            # Sperrung (letter-spacing) kennt python-pptx nicht direkt
            zeichen._rPr.set("spc", str(int(sperrung * 100)))
    return feld


def kennzeile(folie, inhalt, oben=Inches(0.45), farbe=AKZENT, links=RAND):
    return text(folie, inhalt.upper(), links, oben, BREITE - 2 * RAND,
                Inches(0.35), Pt(12), MONO, farbe, sperrung=1.6)


def notiz(folie, inhalt):
    folie.notes_slide.notes_text_frame.text = inhalt


def bild(folie, schluessel, links, oben, breite, hoehe):
    """Bild formatfuellend einsetzen (Mittenausschnitt, kein Verzerren)."""
    daten = normalize_jpeg((IMG / f"{schluessel}.jpg").read_bytes())
    form = folie.shapes.add_picture(BytesIO(daten), links, oben, breite, hoehe)
    eigen_b, eigen_h = form.image.size
    ziel = breite / hoehe
    quelle = eigen_b / eigen_h
    if quelle > ziel:                      # Bild ist breiter -> seitlich beschneiden
        anteil = (1 - ziel / quelle) / 2
        form.crop_left = form.crop_right = anteil
    else:                                  # Bild ist hoeher -> oben/unten beschneiden
        anteil = (1 - quelle / ziel) / 2
        form.crop_top = form.crop_bottom = anteil
    return form


# ----------------------------------------------------------------- Folientypen
def folie_titel(prs, obertitel, titel, untertitel, fakten):
    folie = neue_folie(prs, BLATT)
    kasten(folie, 0, 0, Inches(0.22), HOEHE, AKZENT)
    kennzeile(folie, obertitel, Inches(1.5), AKZENT, RAND)
    text(folie, titel, RAND, Inches(2.0), Inches(10.6), Inches(2.0),
         Pt(60), DISPLAY, TINTE, fett=True, zeilen=1.05)
    text(folie, untertitel, RAND, Inches(4.1), Inches(9.4), Inches(1.2),
         Pt(22), DISPLAY, GRAU, kursiv=True)
    text(folie, fakten, RAND, Inches(6.3), Inches(11), Inches(0.5),
         Pt(14), MONO, HELLGRAU, sperrung=1.2)
    return folie


def folie_trenner(prs, nummer, titel, untertitel):
    folie = neue_folie(prs, AKZENT)
    text(folie, f"DOPPELSTUNDE {nummer}", RAND, Inches(2.3), Inches(11),
         Inches(0.5), Pt(15), MONO, RGBColor(0x9F, 0xBA, 0xD8), sperrung=2.4)
    text(folie, titel, RAND, Inches(2.9), Inches(11.2), Inches(1.6),
         Pt(54), DISPLAY, WEISS, fett=True, zeilen=1.05)
    text(folie, untertitel, RAND, Inches(4.6), Inches(9.6), Inches(1.2),
         Pt(21), DISPLAY, RGBColor(0xC8, 0xD8, 0xEA), kursiv=True)
    return folie


def folie_aussage(prs, inhalt, quelle=None, kennung=None, hintergrund=PAPIER,
                  schriftfarbe=TINTE, hervor=None, notizen=None):
    """Grosses Zitat oder Merksatz, mittig — der haeufigste Folientyp."""
    folie = neue_folie(prs, hintergrund)
    if kennung:
        kennzeile(folie, kennung, Inches(0.5),
                  AKZENT if hintergrund != TINTE else RGBColor(0x8B, 0x95, 0xA3))
    laenge = len(inhalt)
    groesse = Pt(40) if laenge < 120 else Pt(32) if laenge < 260 else Pt(25)
    feld = text(folie, inhalt, RAND, Inches(1.5), BREITE - 2 * RAND, Inches(4.2),
                groesse, DISPLAY, schriftfarbe, anker=MSO_ANCHOR.MIDDLE,
                zeilen=1.32)
    if hervor:
        for absatz in feld.text_frame.paragraphs:
            for lauf in absatz.runs:
                if hervor in lauf.text:
                    vor, _, nach = lauf.text.partition(hervor)
                    lauf.text = vor
                    mitte = absatz.add_run()
                    mitte.text = hervor
                    mitte.font.name = DISPLAY
                    mitte.font.size = groesse
                    mitte.font.bold = True
                    mitte.font.color.rgb = SIGNAL if hintergrund != TINTE else WEISS
                    rest = absatz.add_run()
                    rest.text = nach
                    rest.font.name = DISPLAY
                    rest.font.size = groesse
                    rest.font.color.rgb = schriftfarbe
    if quelle:
        text(folie, quelle, RAND, Inches(6.0), BREITE - 2 * RAND, Inches(0.6),
             Pt(13), MONO, HELLGRAU, sperrung=0.8)
    if notizen:
        notiz(folie, notizen)
    return folie


def folie_bild(prs, schluessel, bildunterschrift, kennung=None, notizen=None):
    folie = neue_folie(prs, TINTE)
    bild(folie, schluessel, 0, 0, BREITE, Inches(6.05))
    if kennung:
        marke = kasten(folie, Inches(0.35), Inches(0.3), Inches(1.5),
                       Inches(0.42), RGBColor(0x15, 0x1A, 0x21))
        text(folie, kennung.upper(), Inches(0.45), Inches(0.34), Inches(1.4),
             Inches(0.35), Pt(12), MONO, WEISS, sperrung=1.4)
    text(folie, bildunterschrift, RAND, Inches(6.3), BREITE - 2 * RAND,
         Inches(1.0), Pt(17), TEXT, RGBColor(0xD8, 0xDD, 0xE4))
    if notizen:
        notiz(folie, notizen)
    return folie


def folie_auftrag(prs, titel, schritte, meta=None, kennung="Arbeitsauftrag",
                  farbe=AKZENT, hell=AKZENT_HELL, notizen=None):
    folie = neue_folie(prs, PAPIER)
    kasten(folie, RAND, Inches(0.9), BREITE - 2 * RAND, Inches(5.7), BLATT,
           farbe, Pt(1.5))
    kasten(folie, RAND, Inches(0.9), BREITE - 2 * RAND, Inches(0.72), hell)
    text(folie, kennung.upper(), RAND + Inches(0.3), Inches(1.06),
         Inches(6), Inches(0.4), Pt(13), MONO, farbe, fett=True, sperrung=1.6)
    if meta:
        text(folie, meta, BREITE - RAND - Inches(4.3), Inches(1.06), Inches(4),
             Inches(0.4), Pt(13), MONO, GRAU, ausrichtung=PP_ALIGN.RIGHT)
    text(folie, titel, RAND + Inches(0.3), Inches(1.85),
         BREITE - 2 * RAND - Inches(0.6), Inches(0.8), Pt(28), DISPLAY,
         TINTE, fett=True)
    zeilen = "\n".join(f"{i}.  {s}" for i, s in enumerate(schritte, 1))
    # Schriftgroesse nach Textmenge UND Schrittzahl — jeder Schritt braucht
    # mindestens eine eigene Zeile, sonst laeuft der Kasten unten aus.
    last = len(zeilen) + 60 * len(schritte)
    groesse = Pt(21) if last < 460 else Pt(18) if last < 620 else Pt(16)
    text(folie, zeilen, RAND + Inches(0.3), Inches(2.72),
         BREITE - 2 * RAND - Inches(0.6), Inches(3.7), groesse, TEXT, TINTE,
         zeilen=1.42)
    if notizen:
        notiz(folie, notizen)
    return folie


def folie_punkte(prs, titel, punkte, kennung=None, notizen=None, fuss=None):
    folie = neue_folie(prs, PAPIER)
    if kennung:
        kennzeile(folie, kennung)
    text(folie, titel, RAND, Inches(1.0), BREITE - 2 * RAND, Inches(0.9),
         Pt(38), DISPLAY, TINTE, fett=True)
    oben = Inches(2.1)
    hoehe_zeile = Inches(4.6) / max(len(punkte), 1)
    for kopf, rumpf in punkte:
        kasten(folie, RAND, oben + Inches(0.06), Inches(0.05),
               hoehe_zeile - Inches(0.3), AKZENT)
        text(folie, kopf, RAND + Inches(0.28), oben, Inches(3.4),
             hoehe_zeile, Pt(19), DISPLAY, AKZENT, fett=True)
        text(folie, rumpf, RAND + Inches(3.8), oben,
             BREITE - RAND - Inches(4.65), hoehe_zeile, Pt(18), TEXT, TINTE)
        oben += hoehe_zeile
    if fuss:
        text(folie, fuss, RAND, Inches(6.75), BREITE - 2 * RAND, Inches(0.5),
             Pt(14), MONO, HELLGRAU)
    if notizen:
        notiz(folie, notizen)
    return folie


def folie_tabelle(prs, titel, kopf, zeilen, kennung=None, spalten=None,
                  notizen=None, groesse=Pt(15)):
    folie = neue_folie(prs, PAPIER)
    if kennung:
        kennzeile(folie, kennung)
    text(folie, titel, RAND, Inches(0.95), BREITE - 2 * RAND, Inches(0.8),
         Pt(34), DISPLAY, TINTE, fett=True)
    breite = BREITE - 2 * RAND
    hoehe = Inches(4.6)
    form = folie.shapes.add_table(len(zeilen) + 1, len(kopf), RAND,
                                  Inches(1.95), breite, hoehe)
    tabelle = form.table
    tabelle.first_row = True
    if spalten:
        gesamt = sum(spalten)
        for index, anteil in enumerate(spalten):
            tabelle.columns[index].width = Emu(int(breite * anteil / gesamt))
    for index, titel_zelle in enumerate(kopf):
        zelle = tabelle.cell(0, index)
        zelle.text = titel_zelle
        zelle.fill.solid()
        zelle.fill.fore_color.rgb = AKZENT
        absatz = zelle.text_frame.paragraphs[0]
        absatz.runs[0].font.size = Pt(13)
        absatz.runs[0].font.name = MONO
        absatz.runs[0].font.bold = True
        absatz.runs[0].font.color.rgb = WEISS
    for r, zeile in enumerate(zeilen, 1):
        for c, wert in enumerate(zeile):
            zelle = tabelle.cell(r, c)
            zelle.text = wert
            zelle.fill.solid()
            zelle.fill.fore_color.rgb = BLATT if r % 2 else PAPIER
            for absatz in zelle.text_frame.paragraphs:
                for lauf in absatz.runs:
                    lauf.font.size = groesse
                    lauf.font.name = TEXT if c else MONO
                    lauf.font.color.rgb = TINTE
                    lauf.font.bold = (c == 0)
    if notizen:
        notiz(folie, notizen)
    return folie


def folie_zwei_spalten(prs, titel, links_titel, links_punkte, rechts_titel,
                       rechts_punkte, kennung=None, notizen=None):
    folie = neue_folie(prs, PAPIER)
    if kennung:
        kennzeile(folie, kennung)
    text(folie, titel, RAND, Inches(0.95), BREITE - 2 * RAND, Inches(0.8),
         Pt(34), DISPLAY, TINTE, fett=True)
    spalte = (BREITE - 2 * RAND - Inches(0.4)) / 2
    for versatz, (kopf, punkte, farbe, hell) in enumerate((
            (links_titel, links_punkte, JETZT, JETZT_HELL),
            (rechts_titel, rechts_punkte, SIGNAL, SIGNAL_HELL))):
        x = RAND + versatz * (spalte + Inches(0.4))
        kasten(folie, x, Inches(1.95), spalte, Inches(4.5), BLATT, farbe, Pt(1.5))
        kasten(folie, x, Inches(1.95), spalte, Inches(0.62), hell)
        text(folie, kopf, x + Inches(0.25), Inches(2.06), spalte - Inches(0.5),
             Inches(0.45), Pt(17), DISPLAY, farbe, fett=True)
        text(folie, "\n".join(f"•  {p}" for p in punkte), x + Inches(0.25),
             Inches(2.8), spalte - Inches(0.5), Inches(3.5), Pt(16), TEXT,
             TINTE, zeilen=1.45)
    if notizen:
        notiz(folie, notizen)
    return folie


# --------------------------------------------------------------- Der Foliensatz
def baue(ziel):
    prs = Presentation()
    prs.slide_width = BREITE
    prs.slide_height = HOEHE

    folie_titel(
        prs,
        "LehrplanPLUS Bayern · Gymnasium · Jgst. 9 · Ev. Religionslehre",
        "In Verantwortung vor Gott",
        "Sechs Doppelstunden über die Frage, wem man gehorcht,\nwenn Staat und Gewissen auseinandergehen.",
        "6 DOPPELSTUNDEN   ·   12 STUNDEN   ·   LERNBEREICH 4")

    # ================================================== DOPPELSTUNDE 1
    folie_trenner(prs, 1, "Wo steckt Kirche im Staat?",
                  "Eine Spurensuche im eigenen Alltag — mit dem Lohnzettel als Tatort.")

    folie_tabelle(
        prs, "Entgeltabrechnung · Muster · September",
        ["Position", "Betrag"],
        [["Bruttoentgelt", "3.200,00 €"],
         ["Lohnsteuer", "− 430,25 €"],
         ["Kirchensteuer ev. (8 %)", "− 34,42 €"],
         ["Sozialversicherung", "− 665,60 €"],
         ["Auszahlungsbetrag", "2.069,73 €"]],
        kennung="M1 · Bild 1 von 7", spalten=[2, 1], groesse=Pt(17),
        notizen="Zehn Sekunden zeigen. Nichts sagen. Steuerklasse I, Kirche: ev.")

    folie_tabelle(
        prs, "Stundenplan · Klasse 9",
        ["Std.", "Mo", "Di", "Mi", "Do", "Fr"],
        [["1.", "D", "M", "E", "Ph", "D"],
         ["2.", "M", "E", "Ge", "M", "Sp"],
         ["3.", "Ev / K / Eth", "Ch", "D", "E", "Ku"],
         ["4.", "Bio", "D", "Ev / K / Eth", "Ge", "M"],
         ["5.", "Sp", "Geo", "M", "Inf", "E"]],
        kennung="M1 · Bild 2 von 7", spalten=[1, 2, 1, 2, 1, 1],
        groesse=Pt(15),
        notizen="Drei Fächer teilen sich dieselbe Stunde — an einer staatlichen Schule.")

    folie_bild(prs, "bild3",
               "Ein evangelischer und ein katholischer Militärgeistlicher bei einer Gedenkfeier im Auslandseinsatz. Stola über dem Feldanzug — und keine Waffe.",
               "M1 · 3 von 7",
               notizen="Häufiger Irrtum: Militärpfarrer seien Soldaten. Sind sie nicht. Jetzt nicht auflösen.")

    folie_bild(prs, "bild4",
               "Geschäftsstelle der Diakonie in München. Das Zeichen an der Wand ist das Kronenkreuz.",
               "M1 · 4 von 7")

    folie_aussage(
        prs,
        "An stillen Tagen — dazu gehört in Bayern der Karfreitag ganztägig — sind öffentliche Unterhaltungsveranstaltungen untersagt, sofern sie dem ernsten Charakter dieser Tage widersprechen.",
        "Zusammenfassung der Regelung des Bayerischen Feiertagsgesetzes",
        "M1 · Bild 5 von 7", hervor="stillen Tagen",
        notizen="Ein staatliches Gesetz schützt einen kirchlichen Feiertag — und verbietet allen etwas, auch denen, die nicht in der Kirche sind.")

    folie_aussage(
        prs,
        "„Ich schwöre, dass ich meine Kraft dem Wohle des deutschen Volkes widmen … werde. So wahr mir Gott helfe.“\n\nDer Eid kann auch ohne religiöse Beteuerung geleistet werden.",
        "Amtseid nach Artikel 56 des Grundgesetzes",
        "M1 · Bild 6 von 7", hervor="So wahr mir Gott helfe.")

    folie_bild(prs, "bild7",
               "Ein evangelischer Kindergarten. Träger ist die Kirchengemeinde — bezahlt überwiegend aus öffentlichen Mitteln und Elternbeiträgen.",
               "M1 · 7 von 7")

    folie_auftrag(
        prs, "Was haben diese sieben Bilder gemeinsam?",
        ["Notiert zu zweit in einem Satz, was die sieben Bilder verbindet.",
         "Einigt euch auf ein Bild, bei dem ihr euch am wenigsten sicher seid, warum es dabei ist."],
        meta="Partnerarbeit · 3 min", kennung="Arbeitsauftrag 1.1",
        notizen="ÜBERLEITUNG: „Ich höre viermal dasselbe Wort: Kirche. Sie haben recht — in jedem "
                "dieser Bilder steckt Kirche. Aber schauen Sie genau hin: In keinem einzigen ist eine "
                "Kirche zu sehen. Das hier ist eine Lohnabrechnung, das ein Rettungswagen, das ein "
                "Stundenplan. Das sind lauter staatliche Angelegenheiten. Und trotzdem ist Kirche darin. "
                "Wie kommt sie da hinein?“")

    folie_aussage(
        prs,
        "„Im Bewußtsein seiner Verantwortung vor Gott und den Menschen … hat sich das Deutsche Volk kraft seiner verfassungsgebenden Gewalt dieses Grundgesetz gegeben.“",
        "Präambel des Grundgesetzes", "M2 · Erster Satz", hervor="vor Gott")

    folie_aussage(
        prs, "„Es besteht keine Staatskirche.“",
        "Artikel 137 Abs. 1 der Weimarer Reichsverfassung — nach Art. 140 GG Bestandteil des Grundgesetzes",
        "M2 · Zweiter Satz",
        notizen="Zuruf: Welcher Satz passt zu welchem Bild von vorhin?")

    folie_aussage(
        prs,
        "Deutschland hat keine Staatskirche — und trotzdem sind Kirche und Staat überall miteinander verflochten.\n\nWie viel Nähe ist gut? Und wann wird sie gefährlich?",
        None, "Leitfrage der ganzen Einheit", hintergrund=AKZENT_HELL,
        notizen="An die Tafel schreiben. Bleibt sechs Wochen lang stehen.\n\n"
                "ÜBERLEITUNG: „Diese Frage kann ich Ihnen nicht beantworten, ohne dass Sie wissen, "
                "worüber wir reden. Also machen wir es umgekehrt: Sie recherchieren selbst. Jede Gruppe "
                "wird für zwanzig Minuten die einzige Person im Raum, die sich mit ihrem Thema auskennt.“")

    folie_auftrag(
        prs, "Expertengruppen: vier Steckbriefe",
        ["Lest euer Materialblatt und klärt alle Begriffe, die keiner erklären kann — bevor ihr das Handy anfasst.",
         "Recherchiert die beiden Aufträge. Regel: Was ihr aufschreibt, muss auf einer Seite stehen, die ihr nennen könnt.",
         "Füllt den A3-Steckbrief aus: WAS · WOHER · WEM NÜTZT ES · STREITPUNKT. Je Feld höchstens drei Sätze.",
         "Legt fest, wer den Steckbrief später in 90 Sekunden erklärt. Übt es einmal."],
        meta="Gruppenarbeit · 30 min", kennung="Arbeitsauftrag 1.2",
        notizen="ÜBERLEITUNG danach: „Zeit ist um. Ab jetzt gilt: Ihr Blatt bleibt liegen, Sie nehmen "
                "nur Ihren Kopf mit. Wenn Ihre Zuhörer am Ende Ihr Feld ‚Streitpunkt‘ nicht ausfüllen "
                "können, haben Sie schlecht erklärt — nicht die anderen schlecht zugehört.“")

    for nummer, satz in enumerate([
            "„Der Staat soll die Kirchensteuer nicht mehr für die Kirchen einziehen.“",
            "„Religionsunterricht gehört nicht an eine staatliche Schule.“",
            "„Eine Kirche, die Krankenhäuser betreibt, soll ihren Angestellten dieselben Streikrechte geben wie jeder andere Arbeitgeber.“",
            "„Solange die Kirche vom Staat Geld bekommt, sollte sie sich mit Kritik am Staat zurückhalten.“"], 1):
        folie_aussage(
            prs, satz, "Stellen Sie sich in die Ecke, die zu Ihrer Meinung passt.",
            f"M5 · Streitsatz {nummer} von 4",
            notizen=("Nach zwei Begründungen fragen: „Wer hat ein Argument gehört, das ihn ins Wanken "
                     "bringt? Wechseln Sie die Ecke — das ist ausdrücklich erlaubt.“"
                     + ("\n\nSatz 4 ist der wichtigste. Verteilung im Raum notieren — sie kommt in "
                        "Doppelstunde 6 zurück." if nummer == 4 else "")))

    folie_bild(prs, "ds2_kirchenwahl",
               "Eine evangelische Kirche in Deutschland. Das Bild ist ungefähr neunzig Jahre alt.",
               "Stundenschluss",
               notizen="Zehn Sekunden stehen lassen, nichts sagen.\n\n"
                       "ÜBERLEITUNG: „Damals war die Nähe zwischen Kirche und Staat so eng wie nie — und "
                       "für manche Menschen wurde sie lebensgefährlich. Nächste Stunde fangen wir genau "
                       "hier an. Und die Frage, die Sie mitnehmen, ist nicht ‚Wie konnten die nur?‘, "
                       "sondern: ‚Was hätte ich in diesem Raum gemacht?‘“\n\n"
                       "HAUSAUFGABE: Eine erwachsene Person fragen: Zahlst du Kirchensteuer — und warum "
                       "eigentlich (nicht)? Zwei Sätze, ohne Namen.")

    # ================================================== DOPPELSTUNDE 2
    folie_trenner(prs, 2, "1933: Die Kirche jubelt mit",
                  "Warum die Frage nicht lautet „Wie konnten die nur?“, sondern „Was war daran so überzeugend?“")

    folie_bild(prs, "ds2_kirchenwahl", "", "M6a",
               notizen="Hängt schon an der Wand, wenn die Klasse hereinkommt. Ohne Titel, ohne Jahreszahl.")

    folie_auftrag(
        prs, "Bildlesen in drei Schritten",
        ["SEHEN — Notieren Sie fünf Dinge, die Sie auf dem Bild tatsächlich sehen. Keine Deutung, nur Gegenstände.",
         "VERMISSEN — Notieren Sie eine Sache, die Sie in einer Kirche erwarten würden und die hier fehlt oder überdeckt ist.",
         "FRAGEN — Formulieren Sie die eine Frage, die Sie den Menschen auf diesem Bild stellen würden."],
        meta="Einzelarbeit · 8 min", kennung="Arbeitsauftrag 2.1",
        notizen="Auf dem Transparent hat das Kreuz nur die Hälfte des Platzes — die andere gehört dem "
                "Hakenkreuz. Wenn niemand darauf kommt, nur fragen: „Wer wirbt hier eigentlich für wen?“\n\n"
                "ÜBERLEITUNG: „Die meisten Ihrer Fragen laufen auf dasselbe hinaus: Wie konnten die nur? "
                "Ich möchte Ihnen diese Frage wegnehmen. Sie ist bequem — sie stellt uns auf die richtige "
                "Seite, bevor wir irgendetwas verstanden haben. Die schwierigere Frage lautet: Was war "
                "daran so überzeugend?“")

    folie_tabelle(
        prs, "Das Jahr, in dem sich die evangelische Kirche entschied",
        ["Datum", "Ereignis"],
        [["30. Jan. 1933", "Hitler wird Reichskanzler — in vielen Gemeinden Dankgottesdienste"],
         ["März 1933", "Ermächtigungsgesetz — Widerspruch wird gefährlich"],
         ["Frühjahr 1933", "Die „Deutschen Christen“ werden Massenbewegung"],
         ["23. Juli 1933", "KIRCHENWAHL — deutliche Mehrheit für die Deutschen Christen"],
         ["Sept. 1933", "Ludwig Müller wird Reichsbischof · Arierparagraph erreicht die Kirche"],
         ["13. Nov. 1933", "Sportpalast: Forderung, das Alte Testament abzuschaffen"]],
        kennung="M7 · Zeitstrahl", spalten=[1, 3],
        notizen="Die Kirchenwahl ist der wichtigste Tag: Die Anpassung wurde GEWÄHLT, nicht befohlen.")

    folie_bild(prs, "ds2_reichsbischof",
               "Berliner Dom, September 1933: die Einführung des Reichsbischofs Ludwig Müller. Talare und Uniformen auf denselben Stufen.",
               "M6b")

    folie_auftrag(
        prs, "Quellen-Autopsie",
        ["Wer spricht hier, wann, zu wem? Ein Satz.",
         "Was wird gefordert? Unterstreicht im Text die stärkste Forderung.",
         "Womit wird es begründet? Achtet darauf, wo religiöse Sprache Politisches fordert.",
         "Was wäre die Folge? Beschreibt eine konkrete Person, für die dieser Text gefährlich wird.",
         "Ein Satz in heutiger Alltagssprache auf den Papierstreifen — für die Wäscheleine."],
        meta="Gruppenarbeit · 25 min", kennung="Arbeitsauftrag 2.2",
        notizen="Vor dem Austeilen ansagen: „Zwei dieser Texte sind judenfeindlich. Ich gebe sie Ihnen, "
                "weil man nicht verstehen kann, was passiert ist, ohne zu sehen, wie normal diese Sätze "
                "damals klangen. Wenn Ihnen ein Satz zu nahe geht, sagen Sie es mir.“")

    folie_aussage(
        prs,
        "„Jedermann sei untertan der Obrigkeit, die Gewalt über ihn hat. Denn es ist keine Obrigkeit außer von Gott.“",
        "Römer 13,1 · Lutherbibel 2017", "M9 · Vers 1 von 2")

    folie_aussage(
        prs, "„Man muss Gott mehr gehorchen als den Menschen.“",
        "Apostelgeschichte 5,29 · Lutherbibel 2017", "M9 · Vers 2 von 2",
        notizen="ÜBERLEITUNG: „Die Menschen, die 1933 gejubelt haben, waren keine schlechten Christen in "
                "ihren eigenen Augen. Sie hatten einen Bibelvers dafür. Und die, die widersprochen haben, "
                "hatten auch einen.“")

    folie_aussage(
        prs,
        "Der ganze Kirchenkampf ist ein Streit darüber, wo die Grenze zwischen diesen beiden Versen verläuft.",
        None, "Merksatz — trägt bis Doppelstunde 6", hintergrund=AKZENT_HELL)

    folie_aussage(
        prs, "„Jesus Christus … ist das eine Wort Gottes, das wir zu hören … haben.“",
        "Mai 1934", "Cliffhanger", hervor="das eine Wort Gottes",
        notizen="ÜBERLEITUNG: „Das EINE Wort. Nicht: eines von mehreren. Überlegen Sie bis nächste Woche, "
                "wem dieser Satz wehtut. Denn dreizehn Tage später haben acht evangelische Theologen — "
                "sechs davon aus Franken — eine Gegenschrift veröffentlicht. Nächste Stunde lassen wir die "
                "beiden Seiten aufeinandertreffen. Und Sie sitzen nicht im Publikum. Sie sind die beiden Seiten.“")

    # ================================================== DOPPELSTUNDE 3
    folie_trenner(prs, 3, "Barmen gegen Ansbach",
                  "Zwei Texte, dreizehn Tage Abstand, dieselbe Kirche — und ein bayerischer Landesbischof mittendrin.")

    folie_aussage(
        prs,
        "„Jesus Christus, wie er uns in der Heiligen Schrift bezeugt wird, ist das eine Wort Gottes, das wir zu hören, dem wir im Leben und im Sterben zu vertrauen und zu gehorchen haben.“",
        None, "M10 · Text A")

    folie_aussage(
        prs,
        "Gott begegne uns auch in den natürlichen Ordnungen des Lebens — in Familie, Volk und Rasse. Der Christ sei zum Gehorsam gegenüber diesen Ordnungen und gegenüber der von Gott gegebenen Obrigkeit verpflichtet.",
        "sinngemäße Zusammenfassung, kein Zitat", "M10 · Text B",
        hervor="Familie, Volk und Rasse",
        notizen="Frage: Welcher Text stammt von Gegnern des Nationalsozialismus? Zeigen mit dem Finger.\n\n"
                "AUFLÖSUNG: Beide sind von 1934, beide von evangelischen Theologen, dreizehn Tage "
                "auseinander. Text A: Barmer Erklärung, 31. Mai. Text B: Ansbacher Ratschlag, 11. Juni.")

    folie_tabelle(
        prs, "„Streit am Sonntagabend“ · Ansbach, Juni 1934",
        ["Dauer", "Teil", "Regie"],
        [["5 min", "Eröffnung", "Je 60 Sekunden Statement: Wagner, Althaus, Barth, Meiser, Charlotte S."],
         ["10 min", "Freie Runde", "Moderation greift nur ein, wenn jemand nicht zu Wort kommt"],
         ["4 min", "Publikumsfrage", "Eine Frage aus dem Halbkreis"],
         ["3 min", "Schlusswort", "Je 30 Sekunden. Charlotte S. spricht zuletzt"]],
        kennung="M12 · Ablauf der Talkrunde", spalten=[1, 2, 5],
        notizen="Zwei Regeln laut ansagen: Angegriffen werden Positionen, nie Personen. Auf „Stopp“ "
                "frieren alle ein.\n\nNACH DER RUNDE: „Bitte aufstehen, einmal ausschütteln, auf einen "
                "anderen Platz setzen. Damit sind die Rollen abgelegt.“")

    folie_aussage(
        prs,
        "„Wir verwerfen die falsche Lehre, als könne und müsse die Kirche als Quelle ihrer Verkündigung außer und neben diesem einen Worte Gottes auch noch andere Ereignisse und Mächte, Gestalten und Wahrheiten als Gottes Offenbarung anerkennen.“",
        "Barmer Theologische Erklärung, These 1, 31. Mai 1934", "M13")

    folie_auftrag(
        prs, "Barmen in eigener Sprache",
        ["Lest These 1 laut vor — einmal ganz, ohne zu stocken.",
         "Übersetzt sie in höchstens 25 Wörter heutiger Alltagssprache. Kein Wort aus dem Original außer „Gott“ und „Jesus Christus“.",
         "Schreibt euren eigenen Verwerfungssatz für 1934: „Wir verwerfen die falsche Lehre, als ob …“",
         "Haltet fest: Wem tut dieser Satz weh? Nennt zwei Gruppen."],
        meta="Partnerarbeit · 14 min", kennung="Arbeitsauftrag 3.3",
        notizen="ÜBERLEITUNG: „Der Originaltext ist sperrig — Theologendeutsch von 1934. Ihre Aufgabe ist, "
                "ihn so zu übersetzen, dass ihn eine Siebtklässlerin versteht. Das ist die härteste Prüfung "
                "für Verstehen, die ich kenne: Wer etwas nicht einfach sagen kann, hat es nicht verstanden.“\n\n"
                "Diese Satzform kommt in Doppelstunde 6 zurück. Die besten Übersetzungen einsammeln.")

    folie_zwei_spalten(
        prs, "Wie umgehen mit halben Helden?",
        "Das spricht für ihn",
        ["Verhinderte die Eingliederung der bayerischen Landeskirche in die Reichskirche",
         "Oktober 1934 Hausarrest — Tausende demonstrierten in Franken für ihn",
         "Die bayerische Landeskirche blieb „intakt“",
         "Baute die Kirche nach 1945 wieder auf"],
        "Das spricht gegen ihn",
        ["Äußerte sich antijudaistisch",
         "Schwieg öffentlich zur Verfolgung und Deportation der Juden",
         "War nachsichtig mit belasteten Pfarrern nach 1945",
         "Straßen in München (2007) und Nürnberg (2010) wurden umbenannt"],
        kennung="M14 · Landesbischof Hans Meiser",
        notizen="Blitzumfrage: Umbenennen — ja oder nein? Dann die dritte Möglichkeit ins Spiel bringen: "
                "Name behalten, Zusatzschild anbringen. Was müsste darauf stehen? Einen Satz gemeinsam "
                "formulieren.\n\nÜBERLEITUNG: „Meiser hat seine Kirche gerettet und zu den Deportationen "
                "geschwiegen. Bonhoeffer hat an einem Attentat mitgewirkt und ist dafür gehängt worden. "
                "Niemöller hat Hitler zuerst gewählt und ist dann im KZ gelandet. Nächste Stunde spannen "
                "wir eine Leine quer durch diesen Raum.“")

    # ================================================== DOPPELSTUNDE 4
    folie_trenner(prs, 4, "Was hätte ich getan?",
                  "Fünf Biografien an einer Wäscheleine — zwischen Anpassung und Widerstand.")

    folie_aussage(
        prs,
        "Als die Nationalsozialisten die Kommunisten holten, habe ich geschwiegen; ich war ja kein Kommunist.\n\nAls sie die Sozialdemokraten einsperrten, habe ich geschwiegen; ich war ja kein Sozialdemokrat.\n\nAls sie die Gewerkschafter holten, habe ich geschwiegen; ich war ja kein Gewerkschafter.\n\nAls sie mich holten, gab es keinen mehr, der protestieren konnte.",
        "Martin Niemöller · keine autorisierte Urfassung — er hat den Gedanken mehrfach anders formuliert",
        "M15 · Chorlesung",
        notizen="Vier Gruppen lesen nacheinander je einen Abschnitt, die letzte Zeile liest die ganze "
                "Klasse gemeinsam. Danach zehn Sekunden Stille — nicht auflösen.\n\n"
                "Die fehlende Urfassung gehört in die Auswertung: Auch ein berühmtes Zitat hat eine "
                "Geschichte. Wird der Text dadurch weniger wahr?")

    folie_aussage(
        prs, "Wobei würde ich heute wegschauen?",
        "Ein Satz. Wird nicht eingesammelt und nicht vorgelesen. Der Zettel gehört Ihnen.",
        "Stille Arbeit · 3 min", hintergrund=EPOCHE_HELL,
        notizen="ÜBERLEITUNG: „Behalten Sie diesen Zettel. Der Mann, der den Text geschrieben hat, war "
                "1924 stolzer Nationalist und U-Boot-Kommandant. Später saß er sieben Jahre in "
                "Sachsenhausen und Dachau. Menschen sind nicht von Anfang an Helden oder Feiglinge — "
                "sie werden es, Entscheidung für Entscheidung.“")

    folie_auftrag(
        prs, "Fünf Biografien, fünf Personenkarten",
        ["Gestaltet die Personenkarte: Name · Lebensdaten · EINE Entscheidung, die alles veränderte · was sie gekostet hat.",
         "Sucht den Satz eurer Person im Text und schreibt ihn groß auf die Rückseite.",
         "Ordnet eure Person auf der Skala ein: 1 = Anpassung, 5 = Widerstand unter Lebensgefahr. Ein Satz, warum nicht eine Stufe höher oder tiefer.",
         "Beantwortet für euch, ohne es aufzuschreiben: Hätte diese Person mehr tun können? Was?"],
        meta="Gruppenarbeit · 26 min", kennung="Arbeitsauftrag 4.2",
        notizen="ÜBERLEITUNG: „Ich habe eine Leine gespannt. Links: Anpassung. Rechts: Widerstand. "
                "Dazwischen ist nichts — kein Raster, keine Markierungen. Den Zwischenraum füllen Sie. "
                "Und ich verspreche Ihnen: Sobald die dritte Karte hängt, wird jemand widersprechen wollen.“")

    folie_auftrag(
        prs, "Was macht Widerstand wahrscheinlicher?",
        ["ALLEIN (4 min) — Jede und jeder schreibt ins eigene Außenfeld drei Bedingungen, die es Menschen leichter gemacht haben zu widersprechen.",
         "REIHUM (4 min) — Jede und jeder liest vor. Es wird nicht diskutiert.",
         "GEMEINSAM (4 min) — Einigt euch auf die drei wichtigsten Bedingungen und schreibt sie groß in die Mitte."],
        meta="Placemat · 12 min", kennung="Arbeitsauftrag 4.3", farbe=JETZT,
        hell=JETZT_HELL,
        notizen="ÜBERLEITUNG: „Wir hören auf, einzelne Menschen zu benoten. Keiner von uns weiß, wie er "
                "sich verhalten hätte — und wer behauptet, er wüsste es, hat meistens am wenigsten darüber "
                "nachgedacht. Die nützlichere Frage lautet: Was hatten die Leute rechts an der Leine, was "
                "die Leute links nicht hatten?“\n\nMITTELFELDER EINSAMMELN — sie hängen in Doppelstunde 6 "
                "wieder an der Wand.")

    folie_aussage(
        prs, "„Dies ist das Ende — für mich der Beginn des Lebens.“",
        "Dietrich Bonhoeffer, wenige Tage vor seiner Hinrichtung im KZ Flossenbürg, April 1945",
        "M18 · Stundenschluss", hintergrund=TINTE, schriftfarbe=RGBColor(0xF3, 0xF5, 0xF7),
        notizen="Vorlesen. Pause. Nichts erklären, nicht deuten, keine Hausaufgabe hinterherschieben. "
                "Projektion an lassen, bis der letzte Schüler den Raum verlassen hat.\n\n"
                "Erst beim Hinausgehen, im Vorbeigehen: „Nächste Woche springen wir vierzig Jahre weiter. "
                "Gleiches Land, andere Diktatur — und diesmal bringen Sie Ihre Handys mit.“")

    # ================================================== DOPPELSTUNDE 5
    folie_trenner(prs, 5, "Kerzen statt Steine",
                  "Kirche in der DDR — und wie aus Friedensgebeten eine Revolution wurde.")

    folie_bild(prs, "ds5_pflugscharen",
               "Das Zeichen der kirchlichen Friedensbewegung in der DDR. Als Aufnäher auf der Jacke konnte es das Abitur kosten.",
               "M19")

    folie_aussage(
        prs,
        "„Sie werden ihre Schwerter zu Pflugscharen und ihre Spieße zu Sicheln machen. Es wird kein Volk wider das andere das Schwert erheben, und sie werden hinfort nicht mehr lernen, Krieg zu führen.“",
        "Micha 4,3 · Lutherbibel 2017", "M19 · Der Bibelvers auf dem Emblem",
        notizen="Frage: Warum sollte ein Staat Angst vor einem Zeichen haben, das zum Frieden aufruft?\n\n"
                "DER WITZ AN DER SACHE: Das Motiv geht auf eine Plastik zurück, die die Sowjetunion 1959 "
                "den Vereinten Nationen geschenkt hat. Der Staat konnte es schlecht verbieten, ohne sich "
                "selbst zu widersprechen — und ging trotzdem gegen die Träger vor.")

    folie_punkte(
        prs, "Kirche im Sozialismus",
        [("1 · Geduldet", "Der Staat ist atheistisch. Religion gilt als Rest der Vergangenheit. Verboten wird die Kirche nicht — sie wird an den Rand gedrängt."),
         ("2 · Teuer", "Wer in der Kirche aktiv ist, wird beim Abitur benachteiligt. Der Preis für Überzeugung ist die eigene Zukunft."),
         ("3 · Jugendweihe", "Der Staat schafft ein eigenes Übergangsritual. Formal freiwillig, faktisch fast unumgänglich."),
         ("4 · Die Formel 1971", "„Kirche im Sozialismus“ — nicht gegen den Staat, nicht für ihn, sondern in ihm. Bis heute umstritten."),
         ("5 · Bespitzelt", "Die Staatssicherheit setzt Inoffizielle Mitarbeiter bis in kirchliche Leitungsämter ein.")],
        kennung="M20 · Input",
        fuss="Zwei belegte Zahlen ergänzen: Anteil Jugendweihe in den 1980ern · Kirchenmitgliedschaft in der DDR",
        notizen="Der eine Satz, der hängen bleiben soll: Die Kirche war in der DDR der einzige Ort, an dem "
                "man Sätze sagen durfte, die sonst nirgends erlaubt waren. Das machte sie kostbar — und für "
                "den Staat zu einem Problem, das er nicht zerschlagen konnte.")

    folie_auftrag(
        prs, "Wir bauen eine Ausstellung",
        ["Eine Überschrift, die neugierig macht und nichts verrät.",
         "Eine Jahreszahl oder Zeitspanne, groß.",
         "Fünf Sätze, die erklären, was passiert ist. Nicht sechs.",
         "Ein Originalzitat eines beteiligten Menschen, mit Namen.",
         "Eine Bildidee: Welches Foto müsste hier hängen — und wo habt ihr es gefunden?",
         "Eine offene Frage am unteren Rand, auf die ihr selbst keine sichere Antwort habt."],
        meta="Redaktionsteams · 30 min", kennung="Arbeitsauftrag 5.2",
        notizen="Regel laut ansagen: „Zu jeder Behauptung auf Ihrer Tafel muss ich fragen dürfen: Wo steht "
                "das? — und Sie müssen es mir zeigen können.“ Startquellen: LeMO (dhm.de/lemo) und "
                "Bundesstiftung Aufarbeitung.")

    folie_aussage(
        prs,
        "Sie gehen schweigend durch die Ausstellung — wie in einem echten Museum.\n\nDrei Klebepunkte, und die kommen auf die offenen Fragen, die Sie am meisten beschäftigen.\n\nNicht auf die schönste Tafel. Auf die beste Frage.",
        None, "Museumsrundgang · 7 min", hintergrund=JETZT_HELL,
        notizen="Die anschließende Besprechung der meistgeklebten Fragen ist die inhaltlich wertvollste "
                "Phase der Stunde. Sie darf überziehen.")

    folie_tabelle(
        prs, "Zwei Diktaturen, eine Kirche",
        ["Frage", "Nationalsozialismus", "DDR"],
        [["Was wollte der Staat von der Kirche?", "", ""],
         ["Wie groß war ihr Freiraum?", "", ""],
         ["Woran orientierten sich die, die widersprachen?", "", ""],
         ["Welchen Preis zahlte, wer widersprach?", "", ""],
         ["Was hatte die Kirche 1989, was ihr 1933 fehlte?", "", ""]],
        kennung="M22 · Partnerarbeit", spalten=[3, 2, 2], groesse=Pt(14),
        notizen="ERWARTUNGSHORIZONT für die letzte Zeile — nicht vorsagen: die Erfahrung von 1933. "
                "1989 wusste die Kirche, wie es ausgeht, wenn man sich vereinnahmen lässt. Barmen war "
                "inzwischen ein Text, auf den man sich berufen konnte.\n\n"
                "ÜBERLEITUNG: „Nächste Woche steht kein Jahr mehr an der Tafel, sondern dieses. Und Sie "
                "werden etwas tun, das seit 1934 kaum jemand mehr gemacht hat: Sie schreiben eigene "
                "Bekenntnissätze. Und die hängen wir an die Tür.“")

    # ================================================== DOPPELSTUNDE 6
    folie_trenner(prs, 6, "Barmen 2026?",
                  "Die Kirchen und das Erstarken der AfD — und die Frage, wo für Christen heute die Grenze verläuft.")

    folie_aussage(
        prs, "„Die Kirche soll sich aus der Politik raushalten.“",
        "Kleben Sie einen Punkt auf die Linie — bevor wir anfangen.",
        "Einstieg · Plakat an der Wand", hintergrund=SIGNAL_HELL,
        notizen="Das Plakat hängt schon, wenn die Klasse hereinkommt. Alle kleben beim Hereinkommen, "
                "bevor irgendetwas besprochen wird. Verteilung unkommentiert stehen lassen. Am Stundenende "
                "klebt jede und jeder einen zweiten Punkt in anderer Farbe.")

    folie_punkte(
        prs, "Was die Kirchen gesagt haben",
        [("Herbst 2023", "Die Synode der EKD stellt fest, die menschenverachtenden Haltungen und Äußerungen insbesondere der rechtsextremen Kräfte innerhalb der AfD seien mit den Grundsätzen des christlichen Glaubens in keiner Weise vereinbar."),
         ("22. Feb. 2024", "Die Deutsche Bischofskonferenz erklärt völkischen Nationalismus und Christentum für unvereinbar. Die EKD begrüßt die Erklärung ausdrücklich."),
         ("Mai 2025", "Das Bundesamt für Verfassungsschutz stuft die AfD als gesichert rechtsextremistisch ein. Gegen die Einstufung läuft ein Gerichtsverfahren.")],
        kennung="M23 · Drei Schlagzeilen",
        fuss="SACHSTAND GEPRÜFT AM: ______________________",
        notizen="VOR DER STUNDE PRÜFEN: Die dritte Meldung ist rechtlich in Bewegung. Datum eintragen — "
                "das ist Teil der Stunde: Die Klasse soll sehen, wie man mit einem laufenden Verfahren "
                "korrekt umgeht.\n\nFrage: Was heißt eigentlich „unvereinbar“? Und was heißt es nicht?\n\n"
                "ÜBERLEITUNG: „Ich werde Ihnen nicht sagen, was Sie wählen sollen — das steht mir nicht zu "
                "und wäre auch verboten. Wir schauen uns an, was die Kirchen sagen, warum sie es sagen, "
                "und ob sie es sagen dürfen.“")

    folie_auftrag(
        prs, "Vier Stimmen",
        ["Fasst die Position eures Textes in EINEM Satz zusammen, der mit „Wir sagen: …“ beginnt.",
         "Sucht die stärkste Stelle — die Formulierung, mit der ihr jemanden überzeugen könntet, der anderer Meinung ist.",
         "Sucht die schwächste Stelle. Wo ist euer eigener Text angreifbar? Diese Aufgabe ist Pflicht.",
         "Prüft eine Tatsachenbehauptung mit dem Handy nach — auf der Originalseite der Institution, nicht in einer Zusammenfassung."],
        meta="Gruppenarbeit · 20 min", kennung="Arbeitsauftrag 6.2",
        notizen="ÜBERLEITUNG: „Aufgabe 3 war die wichtigste. Wer nur die starken Stellen der eigenen Seite "
                "kennt, kann kein Gespräch führen, sondern nur einen Schlagabtausch. Jede Gruppe stellt "
                "deshalb NICHT ihre Position vor, sondern nennt zuerst die schwächste Stelle ihres eigenen "
                "Textes.“\n\nAn der Tafel zwei Spalten führen: für Zurückhaltung / für Einmischung. Beide "
                "müssen am Ende ähnlich voll sein.")

    folie_aussage(
        prs,
        "Genau das war 1933 die Mehrheitsmeinung in dieser Kirche. Römer 13, Sie erinnern sich.\n\nUnd 1934 hat eine Minderheit gesagt: Es gibt eine Grenze.\n\nWo verläuft die Ihrer Meinung nach?",
        None, "Die Frage, an der die Einheit zusammenläuft", hintergrund=AKZENT_HELL,
        notizen="Einblenden, wenn das Argument „Die Kirche soll sich raushalten“ fällt. Die Frage bleibt "
                "offen — sie wird in der Werkstatt bearbeitet, nicht im Gespräch entschieden.")

    folie_auftrag(
        prs, "Die Barmen-Werkstatt",
        ["Lest die Bibelstellen auf eurer Karte. Wählt EINE aus, die euch am meisten zu sagen hat.",
         "Erster Teil eurer These — woran ihr euch bindet:  „Weil …, gilt für uns: …“",
         "Zweiter Teil — die Verwerfung:  „Wir verwerfen die falsche Lehre, als ob …“",
         "Macht die Härteprüfung. Beide Fragen müssen mit Ja beantwortbar sein.",
         "Schreibt die These groß auf das Thesenblatt. Unterschreiben ist freiwillig."],
        meta="Gruppenarbeit · 22 min", kennung="Arbeitsauftrag 6.3 · Kernstück",
        notizen="ÜBERLEITUNG: „Sie haben in der dritten Doppelstunde einen Satzanfang übersetzt, den ich "
                "Ihnen jetzt zurückgebe: ‚Wir verwerfen die falsche Lehre, als ob …‘ 1934 haben acht Sätze "
                "in dieser Form eine Kirche verändert. Sie sind jetzt dran. Nicht als Übung, sondern im "
                "Ernst. Und weil das eine hohe Messlatte ist, gibt es eine Bedingung: Jede These muss sich "
                "an einer Bibelstelle festmachen lassen. Sonst ist es Meinung.“")

    folie_auftrag(
        prs, "Die Härteprüfung",
        ["Steht in eurer These eine AUSSAGE — oder nur ein Vorwurf gegen eine Gruppe?  Eine Verwerfung richtet sich gegen einen Satz, nie gegen Menschen.",
         "Könnte eure These auch dann noch gelten, wenn morgen eine ganz andere Partei regiert?  Wenn nein: überarbeiten."],
        meta="beide Fragen: Ja", kennung="Beide Fragen müssen mit Ja beantwortbar sein",
        farbe=SIGNAL, hell=SIGNAL_HELL,
        notizen="Das ist der Unterschied zwischen einer Bekenntnisthese und einem Wahlkampfplakat — und "
                "zugleich die didaktische Absicherung der Stunde. Eine These, die eine Partei beim Namen "
                "nennt, besteht Kriterium 2 nicht und wird von der Gruppe SELBST überarbeitet. Damit "
                "erledigt sich die Neutralitätsfrage nicht durch ein Verbot der Lehrkraft, sondern durch "
                "eine Einsicht der Klasse.\n\nBarmen war stark, weil es nicht an einer Person hing.")

    folie_aussage(
        prs,
        "Ein Theologieprofessor hat vor gut fünfhundert Jahren seine Thesen an eine Kirchentür geschlagen, weil er wollte, dass darüber gestritten wird — nicht, weil er sicher war, recht zu haben.\n\nIn diesem Geist hängen wir jetzt auf.",
        None, "Thesenanschlag", hintergrund=EPOCHE_HELL,
        notizen="Jede Gruppe geht zur Tür, liest ihre These LAUT vor und hängt sie auf. Kein Kommentar, "
                "keine Bewertung, kein Applaus — die Wirkung entsteht durch das Nacheinander. Erst wenn "
                "alle sechs hängen, tritt die Klasse zurück und liest still.")

    folie_auftrag(
        prs, "Nach dem Anschlag",
        ["Welche der sechs Thesen wäre 1934 am gefährlichsten gewesen?",
         "Welche ist heute am unbequemsten — und für wen?"],
        meta="Plenum · 4 min", kennung="Arbeitsauftrag 6.4",
        notizen="ÜBERLEITUNG: „Eine Sache fehlt noch, und sie ist die entscheidende. Es ist wohlfeil, "
                "Thesen an eine Tür zu hängen, wenn dahinter nichts folgt. Die Bekennende Kirche war genau "
                "deshalb schwach: Sie hat 1934 einen großartigen Text geschrieben — und dann zur "
                "Verfolgung ihrer jüdischen Nachbarn überwiegend geschwiegen. Ein Bekenntnis, aus dem kein "
                "Handeln folgt, ist ein Aufsatz.“")

    folie_auftrag(
        prs, "Vom Text zur Tat",
        ["Sucht auf der Seite einer Kirchengemeinde oder eines Dekanats HIER VOR ORT ein Beispiel für politisches oder soziales Engagement.",
         "Notiert: Was wird getan? Wer macht mit? Was würde hier fehlen, wenn es die Kirche nicht gäbe?",
         "Nennt eine Möglichkeit, wie eine Fünfzehnjährige dort mitmachen könnte — ohne fromm sein zu müssen. Findet ihr keine: Schreibt das auf. Auch das ist ein Ergebnis."],
        meta="Partnerarbeit · 10 min", kennung="Arbeitsauftrag 6.5", farbe=JETZT,
        hell=JETZT_HELL,
        notizen="Zwei bis drei Beispiele sammeln und neben die Thesen an die Tür hängen — mit der "
                "Überschrift, die Sie dazuschreiben: „Und das machen wir daraus.“")

    folie_aussage(
        prs,
        "Früh anfangen.\nNicht allein sein.\nVorher geübt haben.\n\nGenau das haben Sie heute getan.",
        "Ihre eigene Antwort aus Doppelstunde 4 — vor drei Wochen aufgeschrieben",
        "Schluss der Einheit", hintergrund=AKZENT_HELL,
        notizen="Die Placemat-Mittelfelder aus DS 4 werden ohne Vorankündigung wieder aufgehängt.\n\n"
                "WORTLAUT: „Das haben Sie vor drei Wochen aufgeschrieben. Es war Ihre Antwort auf die "
                "Frage, was Menschen 1934 geholfen hat zu widersprechen. Lesen Sie es noch einmal — und "
                "lesen Sie es diesmal nicht als Geschichte, sondern als Gebrauchsanweisung.“\n\n"
                "„Und zum Schluss: Kleben Sie noch einmal auf das Plakat von heute Morgen. Dieselbe Frage, "
                "zweiter Klebepunkt, andere Farbe. Ich sage Ihnen nicht, wo Sie kleben sollen. Ich möchte "
                "nur, dass Sie sehen, ob dieser Vormittag etwas verschoben hat.“")

    ziel.parent.mkdir(parents=True, exist_ok=True)
    prs.save(str(ziel))
    return len(prs.slides.__iter__.__self__._sldIdLst)


if __name__ == "__main__":
    pfad = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else \
        ROOT / "praesentation" / "Kirche-und-Staat-Beamer.pptx"
    anzahl = baue(pfad)
    print(f"geschrieben: {pfad}  ({anzahl} Folien)")
