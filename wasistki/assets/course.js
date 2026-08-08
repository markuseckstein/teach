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
  title: 'Was ist eigentlich „KI“?',
  subtitle: 'Ein Kurs für Menschen, die kein Informatikstudium anfangen wollen',

  lessons: [
    {
      n: 1,
      file: '0001-ki-ist-kein-ding.html',
      title: '„KI“ ist kein Ding',
      teaser: 'Warum ChatGPT keine KI ist, sondern ein Produkt mit einem Motor darin — und welche vier Ebenen man auseinanderhalten muss.',
      minutes: 8,
      part: 'Teil I · Die Bausteine'
    },
    {
      n: 2,
      file: '0002-das-llm-hat-kein-gedaechtnis.html',
      title: 'Das LLM hat kein Gedächtnis',
      teaser: 'Ein LLM ist zustandslos: Es erinnert sich an nichts, kennt Sie nicht und ist austauschbar wie ein Motor.',
      minutes: 9,
      part: 'Teil I · Die Bausteine'
    },
    {
      n: 3,
      file: '0003-das-kontextfenster.html',
      title: 'Der Schreibtisch: das Kontextfenster',
      teaser: 'Wenn das Modell nichts behält — woher weiß der Chat dann, worüber wir gerade sprechen? Und warum wird ein langer Chat schlechter?',
      minutes: 10,
      part: 'Teil I · Die Bausteine'
    },
    {
      n: 4,
      file: '0004-anbieter-und-modelle.html',
      title: 'Wer baut die Motoren?',
      teaser: 'Die Anbieterlandschaft, wofür sich welches Modell eignet — und warum Sie sich keine Modellnamen merken sollten.',
      minutes: 8,
      part: 'Teil II · Die Landschaft'
    },
    {
      n: 5,
      file: '0005-was-passiert-mit-meinen-daten.html',
      title: 'Was passiert mit meinen Daten?',
      teaser: 'Alles, was Sie hinschicken, verlässt Ihren Rechner. Eine Ampel für den Schulalltag — und wo Sie sich absichern müssen.',
      minutes: 9,
      part: 'Teil II · Die Landschaft'
    },
    {
      n: 6,
      file: '0006-vom-chatbot-zum-agenten.html',
      title: 'Vom Chatbot zum Agenten',
      teaser: 'Was eine Werkstatt um das Modell herum ändert: Werkzeuge, die Schleife, und ab wann etwas „autonom“ heißt.',
      minutes: 10,
      part: 'Teil III · Die Werkzeuge'
    },
    {
      n: 7,
      file: '0007-skills.html',
      title: 'Skills: Können zum Weitergeben',
      teaser: 'Der Unterschied zwischen einem guten Prompt und einem Skill — am Beispiel von „teach“ und „grill-me“.',
      minutes: 10,
      part: 'Teil III · Die Werkzeuge'
    },
    {
      n: 8,
      file: '0008-der-werkzeugkasten.html',
      title: 'Der Werkzeugkasten',
      teaser: 'Cowork, Claude Code, OpenClaw und die Alternativen — eine Landkarte statt einer Produktliste.',
      minutes: 9,
      part: 'Teil III · Die Werkzeuge'
    },
    {
      n: 9,
      file: '0009-ki-im-lehreralltag.html',
      title: 'KI im Lehreralltag',
      teaser: 'Der ehrliche Aufwand-Nutzen-Vergleich, und sieben Rezepte, die am Montag wirklich Zeit sparen.',
      minutes: 12,
      part: 'Teil IV · Anwenden'
    },
    {
      n: 10,
      file: '0010-workshop-fahrplan.html',
      title: 'Fahrplan für Ihren Workshop',
      teaser: 'Für Markus: 90 Minuten, sechs Demos, die typischen Einwände — und was Sie auf keinen Fall erklären sollten.',
      minutes: 11,
      part: 'Teil IV · Anwenden'
    }
  ],

  reference: [
    {
      file: 'glossar.html',
      title: 'Glossar',
      teaser: 'Alle Begriffe des Kurses, je in ein bis zwei Sätzen. Das Nachschlagewerk zum Ausdrucken.'
    },
    {
      file: 'spickzettel-ebenen.html',
      title: 'Spickzettel: die vier Ebenen',
      teaser: 'Eine Seite, ein Diagramm: LLM → Chatbot → Agent → autonomer Agent. Der Kern des ganzen Kurses.'
    },
    {
      file: 'anbieter-und-modelle.html',
      title: 'Anbieter & Modelle',
      teaser: 'Momentaufnahme der Landschaft mit Datum. Das Dokument, das am schnellsten altert — bewusst getrennt.'
    },
    {
      file: 'datenschutz-ampel.html',
      title: 'Datenschutz-Ampel',
      teaser: 'Rot, Gelb, Grün: Was darf in einen KI-Dienst, was nicht. Zum Aufhängen im Lehrerzimmer.'
    },
    {
      file: 'rezepte-lehreralltag.html',
      title: 'Rezeptkarten',
      teaser: 'Die sieben Anwendungen aus Lektion 9 als abtrennbare Karten mit fertigen Formulierungen.'
    }
  ]
};
