# Spezifikation — Barmen-Werkstatt (Modul 7)

Stand: 5. September 2026 · Ergebnis der Grilling-Sitzung vom selben Tag
Umfang dieser Spezifikation: **Modul 7 vollständig**, Regiepult und Beamer-Ansicht
so weit, wie Doppelstunde 6 sie braucht. Alles andere ist bewusst ausgeklammert.

## Problemstellung

In Doppelstunde 6 schreiben sechs Gruppen je eine eigene Bekenntnisthese in der
Sprachform der Barmer Theologischen Erklärung. Diese Aufgabe ist das Kernstück der
gesamten Unterrichtseinheit — und die Stelle, an der auf Papier regelmäßig zweierlei
schiefgeht.

Erstens **zerfällt die Form**. Eine Barmer These hat zwei Teile: eine positive Bindung
("Weil …, gilt für uns: …") und eine Verwerfung ("Wir verwerfen die falsche Lehre,
als ob …"). Auf einem leeren Blatt schreiben Gruppen unter Zeitdruck nur den zweiten
Teil, oder sie schreiben eine Meinung ohne biblische Bindung. Was dabei herauskommt,
ist eine Parole, kein Bekenntnis.

Zweitens **hängt die didaktische Absicherung an der Lehrkraft**. Die Stunde berührt
eine im Bundestag vertretene Partei. Der Beutelsbacher Konsens verlangt, dass die
Lehrkraft nicht überwältigt. Die Einheit löst das über eine Härteprüfung: Eine These
muss sich gegen einen *Satz* richten, nicht gegen Menschen, und sie muss auch dann
noch gelten, wenn morgen eine andere Partei regiert. Auf Papier stellt diese beiden
Fragen die Lehrkraft — und damit sieht es aus, als zensiere sie. Genau das soll sie
nicht müssen.

Drittens, praktisch: Die Ergebnisse aus DS 4 (Placemat) und DS 5 (Ausstellung) sollen
in DS 6 zurückkommen. Auf Papier bedeutet das einen Ordner, der sechs Wochen überleben
muss.

## Lösung

Eine kleine Webanwendung, die die Gruppen an ihren Tablets durch die Werkstatt führt
und dabei genau drei Dinge tut, die Papier nicht kann:

1. Sie **erzwingt die Form**, indem sie die These in getrennte Felder zerlegt und ohne
   gewählte Bibelstelle nicht weitergeht.
2. Sie **stellt die Härteprüfung im richtigen Moment** — nach dem Schreiben, vor der
   Freigabe — und schickt bei "nein" die Gruppe mit ihrer eigenen Formulierung zurück
   in die Überarbeitung. Sie bewertet nichts, erkennt keine Parteinamen, filtert keine
   Wörter. Die Einsicht bleibt bei den Schülerinnen und Schülern; die Lehrkraft muss
   nicht eingreifen.
3. Sie **zeigt am Beamer, was entsteht** — welche Gruppe wie weit ist, und am Ende die
   sechs fertigen Thesen nebeneinander.

Der Thesenanschlag selbst bleibt analog: Nach bestandener Prüfung wird die These
gedruckt, die Gruppe geht zur Tür, liest laut vor und hängt auf. Die Anwendung endet
bewusst vor dem Ritual.

## User Stories

### Lehrkraft am Regiepult

1. Als Lehrkraft möchte ich vor der Stunde einen Kurs anlegen und einen Beitrittscode
   bekommen, damit ich die Klasse ohne Namenslisten und ohne Accounts hereinholen kann.
2. Als Lehrkraft möchte ich die Zahl der Gruppen festlegen, damit die Aufteilung zur
   tatsächlichen Klassenstärke passt.
3. Als Lehrkraft möchte ich jeder Gruppe eines der sechs Themenfelder zuweisen, damit
   die Bibelstellen-Vorräte stimmen und kein Themenfeld doppelt vergeben wird.
4. Als Lehrkraft möchte ich die Phase der Doppelstunde freischalten, damit alle
   gleichzeitig arbeiten und niemand vorausliest, was die Dramaturgie zerstören würde.
5. Als Lehrkraft möchte ich einen Timer starten und am Beamer sichtbar machen, damit
   die 22 Minuten der Werkstatt spürbar sind, ohne dass ich sie ansage.
6. Als Lehrkraft möchte ich auf einen Blick sehen, in welchem Schritt jede Gruppe
   steckt, damit ich weiß, wo ich hingehen muss.
7. Als Lehrkraft möchte ich sehen, wenn in einer Gruppe seit vielen Minuten dasselbe
   Gerät schreibt, damit ich die stille Übernahme durch eine Person bemerke — eine
   Diagnose, die ich auf Papier nie hatte.
8. Als Lehrkraft möchte ich eine Gruppe im Notfall manuell weiterschalten können, damit
   ein technischer Hänger die Stunde nicht anhält.
9. Als Lehrkraft möchte ich die fertigen Thesen drucken können, damit sie an die Tür
   kommen.
10. Als Lehrkraft möchte ich am Stundenende sehen, wie oft insgesamt wegen der
    Härteprüfung überarbeitet wurde, damit ich weiß, ob das Werkzeug seinen Zweck
    erfüllt hat.
11. Als Lehrkraft möchte ich den Kurs mit einem Knopf beenden und alle Daten löschen,
    damit nach der Einheit nichts liegen bleibt.
12. Als Lehrkraft möchte ich, dass die Anwendung ohne Vorbereitung am Stundenanfang
    startet, damit ich nicht vor der Klasse konfiguriere.

### Gruppe am Tablet

13. Als Schülerin möchte ich per QR-Code oder kurzem Code beitreten, ohne mich
    anzumelden, damit die ersten fünf Minuten nicht mit Passwörtern vergehen.
14. Als Schüler möchte ich meiner Gruppe beitreten und sehen, wer sonst noch verbunden
    ist, damit klar ist, dass wir am selben Text arbeiten.
15. Als Gruppe möchten wir genau ein Schreibgerät haben, damit nicht sechs Leute
    gleichzeitig in dasselbe Feld tippen.
16. Als Schülerin möchte ich das Schreibrecht per Knopf an mich ziehen können, damit
    das Weitergeben eine sichtbare Handlung ist und nicht davon abhängt, wer das Gerät
    festhält.
17. Als Schüler möchte ich auf meinem Gerät live sehen, was die schreibende Person
    tippt, damit ich mitdenken kann, auch wenn ich gerade nicht schreibe.
18. Als Gruppe möchten wir unser Themenfeld und den zugehörigen Bibelstellen-Vorrat
    angezeigt bekommen, damit wir nicht suchen müssen.
19. Als Gruppe möchten wir eine Bibelstelle auswählen müssen, bevor es weitergeht,
    damit unsere These eine Bindung hat und nicht bloß eine Meinung ist.
20. Als Gruppe möchten wir die gewählte Bibelstelle im Volltext lesen können, damit wir
    nicht nach einer Stellenangabe formulieren, die wir nicht kennen.
21. Als Gruppe möchten wir den positiven Teil in zwei getrennten Feldern schreiben
    ("Weil …" / "gilt für uns: …"), damit die Struktur nicht verwischt.
22. Als Gruppe möchten wir die Verwerfung in einem eigenen Feld schreiben, das mit
    "Wir verwerfen die falsche Lehre, als ob …" beginnt, damit die Sprachform
    unübersehbar ist.
23. Als Gruppe möchten wir sehen, wie unsere These zusammengesetzt aussieht, bevor wir
    sie einreichen, damit wir sie als Ganzes beurteilen können.
24. Als Gruppe möchten wir nach dem Schreiben die erste Prüffrage gestellt bekommen —
    steht da eine Aussage oder ein Vorwurf gegen Menschen? —, damit wir selbst merken,
    wenn wir über eine Gruppe statt über einen Satz geschrieben haben.
25. Als Gruppe möchten wir die zweite Prüffrage gestellt bekommen — gilt das auch unter
    einer anderen Regierung? —, damit wir merken, wenn unsere These an einer
    Tagespolitik hängt.
26. Als Gruppe möchten wir bei "nein" mit unserem eigenen Text zurück in die
    Überarbeitung geschickt werden, damit wir nicht von vorn anfangen müssen.
27. Als Gruppe möchten wir, dass die Anwendung unsere These nicht bewertet und keine
    Wörter verbietet, damit die Entscheidung unsere bleibt.
28. Als Gruppe möchten wir unsere These freigeben und dann als Blatt gedruckt bekommen,
    damit wir damit zur Tür gehen können.
29. Als Gruppe möchten wir mit unseren Vornamen unterschreiben können — freiwillig —,
    damit die Unterschrift eine Entscheidung bleibt und keine Pflicht.
30. Als Schüler möchte ich, dass mein Gerät nach einem kurzen Verbindungsabriss wieder
    dort einsteigt, wo wir waren, damit ein Funkloch nicht die Gruppenarbeit kostet.
31. Als Schülerin möchte ich, dass die Seite auf meinem alten Android-Tablet lesbar und
    bedienbar ist, damit ich nicht die Person bin, bei der es nicht geht.

### Klasse am Beamer

32. Als Klasse möchten wir am Beamer sehen, welche Gruppen wie weit sind, damit ein
    Arbeitsstand im Raum sichtbar ist, ohne dass jemand vorgelesen wird.
33. Als Klasse möchten wir am Beamer den Timer sehen, damit wir unser Tempo selbst
    steuern.
34. Als Klasse möchten wir am Ende alle sechs Thesen nebeneinander sehen, damit der
    Vergleich möglich ist.
35. Als Klasse möchten wir am Beamer *nicht* sehen, wie oft eine Gruppe überarbeiten
    musste, damit Überarbeitung nicht als Versagen erscheint.
36. Als Klasse möchten wir, dass die Beamer-Ansicht ohne unsere Namen auskommt, damit
    niemand vorgeführt wird.

### Zeitkapsel und Anschluss

37. Als Lehrkraft möchte ich die Fotos der Placemat-Mittelfelder aus DS 4 hochladen und
    in DS 6 wieder anzeigen können, damit die Klammer der Einheit hält.
38. Als Lehrkraft möchte ich das Foto der fertigen Thesentür archivieren, damit die
    Stunde einen Abschluss hat, der über den Tag hinausreicht.

## Umsetzungsentscheidungen

### Architektur

- **Ein einziges Go-Binary.** `html/template` serverseitig gerendert, Assets über
  `embed.FS` eingebettet, keine Build-Kette, kein Frontend-Framework.
- **SQLite als Datei** über einen reinen Go-Treiber (kein cgo), damit Cross-Compile und
  Backup trivial bleiben — Backup ist `cp`.
- **SSE für Regiepult → Geräte**, normale Formular-POSTs in der Gegenrichtung. Es gibt
  keinen Kanal, der bidirektional sein müsste; damit entfällt WebSocket-Zustand.
- **Betrieb:** ein VPS, systemd, Caddy als TLS-Terminierung davor.
- **Zwei Codebasen, eine Naht.** Die Anwendung liefert die interaktiven Module; die
  Inhaltsseiten und die Nacharbeit-Seiten bleiben statische Dateien, die dasselbe
  Binary mit ausliefert. Fällt die Anwendung aus, bleiben die Inhalte erreichbar.
- **Modernes CSS und HTML5 sind erlaubt** (Container Queries, `:has()`, `dialog`,
  View Transitions), aber als *progressive enhancement*: Die Werkstatt muss auf einem
  Gerät funktionieren, das nichts davon kennt. Kriterium ist nicht Schönheit, sondern
  dass keine Gruppe an ihrem Gerät scheitert.

### Identität und Datenhaltung

- **Keine Namen, keine Accounts, keine Anmeldung.** Kurs → Beitrittscode → Gruppe.
  Gruppen heißen nach Nummer ("Gruppe 3"). Ein Gerät wird über ein zufälliges,
  nicht rückführbares Token im Cookie wiedererkannt; das dient allein dem
  Wiedereinstieg nach Verbindungsabriss und dem Schreibrecht.
- **Freiwillige Vornamen** stehen im Thesentext, weil sie dort hingehören — sie sind
  Inhalt der Unterschrift, kein Identitätsmerkmal des Systems.
- **Daten bleiben bis "Kurs beenden".** Das ist eine bewusste Abweichung von der
  ursprünglichen Idee, am Stundenende zu löschen: Die Einheit braucht die Klammer über
  Wochen (Placemat aus DS 4, Ausstellung aus DS 5). "Kurs beenden" löscht alles.

### Schreibrecht

Genau ein Gerät je Gruppe hält das Schreibrecht. Jedes andere Gerät der Gruppe kann es
per Knopf an sich ziehen; die Übernahme ist sofort wirksam und für alle sichtbar. Es
gibt bewusst **kein** gleichzeitiges Bearbeiten und damit keine Konfliktauflösung —
die Kosten eines CRDT stehen in keinem Verhältnis zu einer 22-Minuten-Phase, und die
sichtbare Übergabe ist didaktisch das bessere Verhalten.

### Zustandsautomat der Werkstatt

Der Kern der Anwendung. Jede Gruppe durchläuft ihn genau einmal je These:

```
BIBELSTELLE ──(Stelle gewählt)──▶ POSITIV ──(beide Felder gefüllt)──▶ VERWERFUNG
                                                                          │
                                                        (Feld gefüllt)    ▼
                            ┌──────────────────────────────────── VORSCHAU
                            │                                            │
                            │                              (einreichen)  ▼
                            │                                      PRUEFUNG_1
                            │                          ┌────────────────┴───────┐
                            │                     nein │                        │ ja
                            └─── ueberarbeitungen++ ◀──┘                        ▼
                            │                                            PRUEFUNG_2
                            │                          ┌────────────────┴───────┐
                            │                     nein │                        │ ja
                            └─── ueberarbeitungen++ ◀──┘                        ▼
                                                                          FREIGEGEBEN
                                                                        (druckbar)
```

Zwei Festlegungen, die nicht verhandelbar sind, weil an ihnen die didaktische
Absicherung hängt:

- **Rückwärts geht es nur nach VERWERFUNG**, nie an den Anfang. Die Gruppe überarbeitet
  ihre eigene Formulierung, sie schreibt nicht neu.
- **Die Anwendung wertet nicht.** Keine Wortlisten, keine Parteinamenerkennung, keine
  automatische Ablehnung. Beide Prüffragen beantwortet die Gruppe selbst. Ein Filter
  würde exakt das Überwältigungsverbot verletzen, das die Prüfung absichern soll.

### Datenmodell

```
kurs(id, name, beitrittscode, aktive_phase, beendet_am)
gruppe(id, kurs_id, nummer, themenfeld, schreibrecht_geraet, schreibrecht_seit)
geraet(id, gruppe_id, token, zuletzt_gesehen)
these(id, gruppe_id, bibelstelle, weil, gilt, verwerfung, schritt,
      ueberarbeitungen, unterschriften, freigegeben_am)
kapsel(id, kurs_id, art, beschriftung, bild, erstellt_am)
```

`ueberarbeitungen` ist zugleich die Messgröße des Erfolgskriteriums und deshalb ein
fachlicher Wert, kein Protokoll-Nebenprodukt.

### Inhalte

Die sechs Themenfelder und ihre Bibelstellen-Vorräte stammen unverändert aus M25 der
bestehenden Einheit (Umgang mit Geflüchteten · Minderheiten · Wahrheit in der
Öffentlichkeit · Schöpfung und Klima · Nation und Volkszugehörigkeit · Freiheit der
Kirche gegenüber jeder Regierung). Sie werden als Daten geführt, nicht als Code, damit
die Lehrkraft sie ohne Neuübersetzung ändern kann.

### Medien

Eingebettet, nicht selbst ausgeliefert; YouTube ausschließlich im
`youtube-nocookie`-Modus, und der Player fällt nach dem Ende in die Beamer-Ansicht
zurück statt in Empfehlungen. Videos laufen über den Beamer, Recherchelinks auf die
Geräte — auch deshalb, damit die Lehrkraft sieht, was nach dem Clip kommt.

### Ausfallsicherheit

Die bestehenden Kopiervorlagen (`word/`, `lessons/`) sind der Plan B und bleiben
gepflegt. DS 5 und DS 6 einmal gedruckt im Schrank, alle sechs als PDF offline auf dem
Laptop der Lehrkraft. "Papierlos" heißt: kein Papier im Regelbetrieb, ein Satz Papier
im Schrank.

## Testentscheidungen

**Ein einziger Testzugang: HTTP.** Die Tests starten die Anwendung wie im Betrieb —
echter Router, echte SQLite-Datei in einem temporären Verzeichnis — und bedienen sie
über HTTP-Anfragen, so wie ein Tablet und das Regiepult es tun. Getestet wird das
äußere Verhalten: Was bekommt ein Gerät zu sehen, was verändert sich für die anderen
Geräte derselben Gruppe, was steht danach in der Datenbank. Interne Funktionen,
Template-Namen und Feldbezeichner sind ausdrücklich **nicht** Gegenstand der Tests;
sie dürfen sich ändern, ohne dass ein Test bricht.

Was geprüft wird:

- Der Zustandsautomat in allen Pfaden, insbesondere die beiden Nein-Zweige: Nach einem
  "nein" steht die Gruppe in VERWERFUNG, ihr Text ist unverändert vorhanden, und
  `ueberarbeitungen` ist um eins gestiegen.
- Dass ohne gewählte Bibelstelle kein Weiterkommen möglich ist, auch nicht durch eine
  direkt abgesetzte Anfrage an einen späteren Schritt.
- Dass eine These erst nach zwei Ja druckbar wird.
- Dass die Anwendung inhaltlich nicht wertet: Zwei Thesen mit unterschiedlichem Text,
  aber gleichem Antwortverhalten, kommen zum gleichen Ergebnis. Dieser Test ist die
  ausführbare Fassung der didaktischen Zusage und darf nie gelöscht werden.
- Schreibrecht: Übernahme durch ein zweites Gerät, Schreibversuch eines Geräts ohne
  Recht, Wiedereinstieg nach Verbindungsabriss.
- Isolation: Eine Gruppe sieht die Thesen anderer Gruppen erst in der Schlussphase.
- "Kurs beenden" hinterlässt keine Zeilen.

**Prior Art gibt es nicht** — dieser Workspace enthält bisher nur Python-Bauskripte und
keine Tests. Diese Suite legt die Konvention fest: Go-Standardbibliothek,
`net/http/httptest`, Testnamen als deutsche Verhaltenssätze passend zum übrigen
Vokabular des Projekts, keine Mocking-Bibliothek.

## Nicht im Umfang

- **Module 6, 8, 9** (Ausstellungsbaukasten, Zeitkapsel-Vollausbau, digitale
  Nacharbeit-Seiten). Sie werden erst entschieden, wenn Modul 7 lauffähig ist.
- **Module 2–5** (Bilderrätsel, Quellen-Lupe, Rollen-Set, Personenkarten). Das sind
  Inhalte, keine Interaktion; sie werden statische Seiten aus dem bestehenden HTML.
- **Die drei analogen Rituale**: Wäscheleine (DS 4), Klebepunkt-Skala (DS 6),
  Thesenanschlag an der Tür (DS 6). Sie bleiben körperlich, weil ihre Wirkung an der
  Sichtbarkeit der eigenen Festlegung hängt. Die Anwendung endet vor der Tür.
- **Bewertung und Leistungsnachweis.** Die Werkstatt erzeugt keine Noten und keine
  personenbezogenen Ergebnisse.
- **Mehrbenutzerbetrieb, Kollegiumsnutzung, Mandantenfähigkeit.** Das Werkzeug ist für
  eine Lehrkraft gebaut und darf ohne sie nicht laufen müssen.
- **Datenschutzkonformer Regelbetrieb.** Für dieses Experiment bewusst ausgeklammert.
  Sobald jemand außer der Autorin es einsetzt, ist das eine neue Entscheidung, keine
  Fortsetzung dieser.
- **Gleichzeitiges Bearbeiten** innerhalb einer Gruppe.

## Weitere Hinweise

**Erfolgskriterium, hart gemessen:** In DS 6 überarbeitet mindestens eine Gruppe ihre
These wegen der Härteprüfung, ohne dass die Lehrkraft etwas sagt. Der Zähler
`ueberarbeitungen` beantwortet das ohne Interpretation. Bleibt er bei null, war die
Anwendung ein hübsches Frontend für eine Aufgabe, die auf Papier genauso gut lief —
auch das ist ein Ergebnis des Experiments und kein Fehlschlag.

**Bekannte Risiken.** BYOD-Android heißt dreißig verschiedene Browser; die
Beamer-Ansicht und die Geräteansicht werden getrennt entwickelt und getestet, aber der
erste echte Klassendurchlauf wird etwas zutage fördern. Der Thesenanschlag braucht
einen funktionierenden Drucker im Haus; fällt er aus, schreibt die Gruppe ab.

**Verhältnis zum bestehenden Material.** Die Einheit bleibt ohne die Anwendung
vollständig unterrichtbar. Das ist keine Übergangslösung, sondern Bedingung: Ein
Werkzeug, das nur die Autorin betreiben kann, darf die Unterrichtseinheit nicht als
Geisel nehmen.
