# Glossar: Was ist eigentlich „KI"?

Die verbindliche Sprache dieses Workspaces. Jede Lektion benutzt diese Begriffe und keine anderen — auch dort, wo draußen andere Wörter kursieren. Die druckbare Fassung für Lernende liegt unter `reference/glossar.html`; diese Datei hier ist die Quelle.

## Die vier Ebenen

**LLM (Sprachmodell)**:
Eine Rechenmaschine, die aus Text neuen Text erzeugt. Zustandslos und austauschbar.
_Vermeiden_: „die KI", „der Algorithmus", „das neuronale Netz"

**Chatbot**:
Ein Produkt, das ein [[LLM]] mit einem Gesprächsverlauf verbindet und diesen bei jeder Anfrage vollständig mitschickt.
_Vermeiden_: „KI-Assistent", „Sprachassistent"

**Agent**:
Ein Chatbot, dem eine [[Werkstatt]] Werkzeuge gibt und der in einer Schleife läuft, bis eine Aufgabe erledigt ist.
_Vermeiden_: „KI-Assistent", „Bot"

**Autonomer Agent**:
Ein Agent, der handeln darf, ohne vorher nachzufragen — bis hin zum Weiterlaufen ohne anwesende Person. Autonomie ist eine Erlaubnis, keine Fähigkeit.
_Vermeiden_: „selbstständige KI", „starke KI"

## Bausteine

**Werkstatt** (englisch _agent harness_):
Das Programm um das [[LLM]] herum, das dessen Vorschläge tatsächlich ausführt und dabei die Rechte kontrolliert. Das Modell schlägt vor, die Werkstatt tut.
_Vermeiden_: „Framework", „Runtime", „Geschirr" — im Deutschen ist „Werkstatt" das tragfähigere Bild.

**Kontextfenster**:
Die maximale Textmenge, die in eine einzelne Anfrage passt. Im Kurs bildlich: der Schreibtisch. Was nicht darauf liegt, existiert für das Modell nicht.
_Vermeiden_: „Speicher", „Gedächtnis", „Arbeitsspeicher"

**Zustandslos** (englisch _stateless_):
Zwischen zwei Aufrufen bleibt nichts erhalten. Jede Anfrage trifft auf ein Modell, für das sie die erste ist.
_Vermeiden_: „vergesslich", „hat kein Langzeitgedächtnis"

**Werkzeug**:
Eine Fähigkeit, die die [[Werkstatt]] dem Modell anbietet — Dateien lesen, suchen, ein Programm ausführen. Das Modell benennt das Werkzeug, benutzen kann es nur die Werkstatt.
_Vermeiden_: „Funktion", „Plugin", „Skill" (letzteres ist etwas anderes!)

**Skill**:
Ein Ordner mit einer Datei `SKILL.md`, die in normaler Sprache beschreibt, wie eine wiederkehrende Aufgabe zu erledigen ist. Ein Verfahren zum Aufheben und Weitergeben.
_Vermeiden_: „Prompt-Vorlage", „Makro", „Plugin"

**Schrittweises Aufdecken** (englisch _progressive disclosure_):
Das Verfahren, mit dem ein Agent viele [[Skill]]s bereithalten kann, ohne das [[Kontextfenster]] zu füllen: erst nur Name und Beschreibung, dann bei Bedarf die vollständige Anweisung, dann bei Bedarf die zugehörigen Dateien.

## Erscheinungen

**Context Rot** (Kontextfäule):
Der gemessene Effekt, dass die Zuverlässigkeit eines Modells sinkt, je länger die Eingabe wird — auch bei einfachen Aufgaben und auch weit unterhalb der Kapazitätsgrenze.
_Vermeiden_: „das Modell wird müde", „Überlastung"

**Themenverschleppung**:
Der Effekt, dass ein früheres Thema im selben Chat die Antwort auf ein neues Thema einfärbt, weil beides gleichzeitig auf dem [[Kontextfenster|Schreibtisch]] liegt. Im Kurs selbst geprägter Begriff — im Englischen gibt es dafür keinen etablierten Namen.

**Erfinden** (englisch _hallucination_):
Wenn das Modell plausibel klingende, aber falsche Angaben erzeugt. Kein Fehlerzustand, sondern die normale Funktionsweise, angewandt auf eine Lücke.
_Vermeiden_: „Halluzination" (klingt nach Krankheit), „Lüge" (unterstellt Absicht), „Bug"

**Offenes Modell** (englisch _open weights_):
Ein [[LLM]], dessen trainierte Datei zum Herunterladen bereitsteht und das man deshalb auf eigener Hardware betreiben kann. Nicht gleichbedeutend mit Open Source.
_Vermeiden_: „Open-Source-KI", „freie KI"

## Anmerkung zur Begriffsbildung

„Werkstatt" und „Themenverschleppung" sind in diesem Workspace geprägt, weil im Deutschen keine brauchbaren etablierten Begriffe existieren. Beim Weitergeben an Dritte sollte der englische Fachbegriff einmal genannt werden — sonst finden die Lernenden draußen nichts wieder.
