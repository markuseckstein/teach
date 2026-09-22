/* ---------------------------------------------------------------------------
   course.js — die EINZIGE Quelle für Reihenfolge, Titel und Navigation.

   Wer eine Lektion hinzufügt, macht genau zwei Dinge:
     1. die HTML-Datei in ./lessons/ anlegen
     2. hier unten einen Eintrag in COURSE.lessons ergänzen

   Nav-Links in bestehenden Lektionen werden NIE von Hand angefasst —
   lesson.js baut sie zur Laufzeit aus dieser Datei.

   Klassisches Script mit globaler Variable, kein ES-Modul: Die Seiten müssen
   per Doppelklick über file:// laufen, und dort blockiert der Browser
   Modul-Importe und fetch().
   --------------------------------------------------------------------------- */

var COURSE = {
  title: 'Die AfD verstehen',
  subtitle: 'Ein Kurs über das Gespräch, das Programm und die Folgen',

  lessons: [
    {
      n: 1,
      file: '0001-protest-oder-ueberzeugung.html',
      title: 'Protest oder Überzeugung?',
      teaser: 'Drei Fragen, mit denen du im laufenden Gespräch erkennst, wen du vor dir hast — und warum die bequemste Antwort („die wählen doch nur aus Protest“) empirisch die schwächste ist.',
      minutes: 11,
      part: 'Teil I · Das Gespräch'
    },
    {
      n: 2,
      file: '0002-die-stillen-am-tisch.html',
      title: 'Die Stillen am Tisch',
      teaser: 'Sobald andere zuhören, ändert sich alles — nur nicht das, was die meisten denken. Nicht dein Gegenüber wird wichtiger, sondern dein Schweigen. Drei Sätze für die Mitte des Tisches.',
      minutes: 11,
      part: 'Teil I · Das Gespräch'
    },
    {
      n: 3,
      file: '0003-wer-regiert-eigentlich.html',
      title: 'Wer regiert eigentlich?',
      teaser: 'Drei Verfassungsartikel entscheiden alles: warum 39 von 83 Sitzen zum Blockieren reichen, zum Gestalten nicht — und warum eine Enthaltung im dritten Wahlgang keine neutrale Handlung ist.',
      minutes: 11,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 4,
      file: '0004-was-ginge-ueberhaupt.html',
      title: 'Die Schere: was ginge, was nicht',
      teaser: 'Die Programmanalyse Feld für Feld. Befund: Fast alles, worüber gestritten wird, ist rechtlich blockiert — und fast alles, was tief wirkt, wird kaum diskutiert.',
      minutes: 12,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 5,
      file: '0005-ausblick-deutschland.html',
      title: 'Ausblick: und wenn ganz Deutschland?',
      teaser: 'Euro, EU, NATO, Russland — jetzt mit beiden Bundesprogrammen im Wortlaut. Darunter die härteste Forderung von allen: Abzug aller alliierten Truppen und ihrer Atomwaffen.',
      minutes: 15,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 6,
      file: '0006-schule-kirche-kanzel.html',
      title: 'Schule, Kirche, Kanzel',
      teaser: 'Die einzige Lektion, die nicht vom Land handelt, sondern von deinem Beruf: Das Programm nennt den Beutelsbacher Konsens beim Namen, die Evangelische Akademie beim Namen — und Bayern als Vorbild.',
      minutes: 12,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 7,
      file: '0007-polizei-und-justiz.html',
      title: 'Polizei und Justiz',
      teaser: 'Aufsicht wird zurückgebaut, Strafrecht ausgeweitet — und eine Kriminalstatistik auf Seite 109 trägt mehr Gewicht, als sie tragen kann. Wie man sie richtig liest.',
      minutes: 15,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 8,
      file: '0008-bestellt-nicht-bezahlt.html',
      title: 'Bestellt, nicht bezahlt',
      teaser: 'Wissenschaft, Landwirtschaft, Verkehr, Gesundheit: vier Kapitel voller sympathischer Forderungen — und auf gut 40 Seiten genau eine Euro-Zahl. Die zweite Schere, diesmal nicht rechtlich, sondern haushalterisch.',
      minutes: 12,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 9,
      file: '0009-die-freiheit-die-sie-meinen.html',
      title: 'Die Freiheit, die sie meinen',
      teaser: 'Kapitel VIII fordert mehr Volksentscheide, mehr Transparenz, weniger Überwachung — und auf denselben acht Seiten nur noch eine Demo pro Ort. Beides ist Landesrecht, beides ginge sofort.',
      minutes: 15,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 10,
      file: '0010-moskau-und-mar-a-lago.html',
      title: 'Moskau und Mar-a-Lago',
      teaser: 'Russland und die US-Milliardäre: das Thema, bei dem man am schnellsten als Verschwörungstheoretiker dasteht. Wo die Aktenlage wirklich aufhört — und warum sie weiter reicht als gedacht.',
      minutes: 14,
      part: 'Teil II · Die Folgen'
    },
    {
      n: 11,
      file: '0011-wie-es-anderswo-ausging.html',
      title: 'Wie es anderswo ausging',
      teaser: 'Die Eskalationsfrage mit Quellen statt Bauchgefühl. Ungarn wurde im April 2026 abgewählt — Polen dreht seit 2023 nichts zurück. Was sich davon auf Sachsen-Anhalt überträgt und was nicht.',
      minutes: 15,
      part: 'Teil III · Die lange Sicht'
    },
    {
      n: 12,
      file: '0012-der-vorwurf-nach-der-grenze.html',
      title: 'Der Vorwurf nach der Grenze',
      teaser: 'Du hast freundlich eine Grenze gezogen — und bekommst einen Vorwurf zurück, der mit der Sache nichts zu tun hat. Warum gerade die Höflichkeit der Auslöser war, und welche eine Nachricht den Streit beendet.',
      minutes: 12,
      part: 'Teil IV · Die Übung'
    }
  ],

  reference: [
    {
      file: 'gespraechs-spickzettel.html',
      title: 'Spickzettel: Das Gespräch',
      teaser: 'Eine Seite zum Ausdrucken: die drei Diagnosefragen, die drei Lagen, was in jeder Lage hilft — und was sicher schiefgeht.'
    },
    {
      file: 'fundstellen-programm.html',
      title: 'Fundstellen im Programm',
      teaser: 'Wo steht was? Seitenzahlen und Wortlaut zu allen analysierten Forderungen — plus der Vergleich, was die Kurzfassung der Partei weglässt.'
    },
    {
      file: 'zahlen-sachsen-anhalt.html',
      title: 'Zahlen: Sachsen-Anhalt 2026',
      teaser: 'Das Wahlergebnis, die Sitzverteilung und die belastbaren Kostenzahlen zum Programm — mit amtlicher Quelle, zum Nachschlagen im Gespräch.'
    },
    {
      file: 'glossar.html',
      title: 'Glossar',
      teaser: 'Jeder Fachbegriff des Kurses in zwei Sätzen — mit der Fundstelle in der Lektion und, wo es zählt, dem Hinweis, womit er gern verwechselt wird.'
    },
    {
      file: 'wiederholungsplan.html',
      title: 'Wiederholungsplan',
      teaser: 'Fünf kurze Abrufe an Tag 1, 3, 7, 21 und 60 — Fragen aus dem Gedächtnis, Antworten zum Aufklappen, Datumsfelder zum Eintragen.'
    }
  ]
};
