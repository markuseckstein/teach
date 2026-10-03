# Barmen-Werkstatt

Begleitende Webanwendung zur Unterrichtseinheit **„Kirche und Staat"** (Ev. Religions-
lehre, Gymnasium Bayern, Jgst. 9). Sie führt in Doppelstunde 6 sechs Gruppen durch die
Barmen-Werkstatt: eigene Bekenntnisthese schreiben, Härteprüfung bestehen, drucken —
und dann analog an die Tür hängen.

Ein Experiment. Ein Betreiber. Die Unterrichtseinheit bleibt ohne diese Anwendung
vollständig unterrichtbar, und das ist Bedingung, nicht Übergangslösung.

## Wo was steht

| | |
|---|---|
| `SPEZIFIKATION.md` | **was** gebaut wird und warum — die verbindliche Fassung |
| `tracker/` | dateibasierte Vorgangsverwaltung; `tracker/INDEX.md` ist die Übersicht |
| `tracker/README.md` | die Spielregeln des Trackers |
| `app/` | das Go-Modul mit dem eigentlichen Server |

Der Unterrichtsstoff selbst liegt eine Ebene höher: `../lessons/`, `../word/`,
`../assets/`. Diese Anwendung liefert ihn später mit aus (Vorgang 0017), verändert ihn
aber nicht.

## Technik in einem Absatz

Ein Go-Binary. `html/template` serverseitig, Assets über `embed.FS`, SQLite als Datei
über einen reinen Go-Treiber, SSE für die Richtung Regiepult → Geräte, normale
Formular-POSTs zurück. Ein VPS, systemd, Caddy davor. Kein Framework, kein
Build-Schritt, kein Container.

## Stand

Vorgänge 0001 (Projektgerüst, HTTP-Server) und 0002 (Datenbankschema) sind erledigt.
Die Testnaht ist etabliert: **genau ein Zugang, HTTP** — echte Anwendung, echte
SQLite-Datei im Temp-Verzeichnis, geprüft wird äußeres Verhalten.

Lokal starten:

    cd app
    go run . -port 8080 -db barmen.db

Für die erste Version läuft die Anwendung nur lokal; Betrieb auf einem VPS
(Vorgang 0018) ist bewusst zurückgestellt.

Nächster Schritt: Vorgang 0003 (Testnaht über HTTP formalisieren) und 0004.
