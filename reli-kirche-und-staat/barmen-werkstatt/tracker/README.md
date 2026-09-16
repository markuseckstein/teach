# Vorgangsverwaltung

Der Issue-Tracker dieses Projekts ist ein Ordner. Kein GitHub, kein Jira — ein
Vorgang ist eine Markdown-Datei in `vorgaenge/`, und der Verlauf steht in der
Git-Historie. Das reicht für ein Ein-Personen-Experiment und hat den Vorteil, dass
Spezifikation, Vorgänge und Code in derselben Änderung wandern.

## Aufbau

    tracker/
      README.md            diese Datei — die Spielregeln
      INDEX.md             erzeugte Übersicht, nicht von Hand pflegen
      index.py             erzeugt INDEX.md aus den Frontmattern
      vorgaenge/
        0001-<kurzname>.md

## Ein Vorgang

Jede Datei beginnt mit einem Frontmatter-Block und enthält danach vier Abschnitte:
**Worum es geht** (das Problem, nicht die Lösung), **Fertig, wenn** (prüfbare
Akzeptanzkriterien), **Gehört nicht dazu** (die Abgrenzung, die verhindert, dass ein
Vorgang wächst) und optional **Hinweise**.

```yaml
---
id: 0008
titel: Zustandsautomat der Werkstatt
status: offen          # offen | in-arbeit | erledigt | verworfen
label: ready-for-agent # siehe unten
haengt-an: [0006, 0007]
---
```

## Statuswerte

| Status | Bedeutung |
|---|---|
| `offen` | beschrieben, noch niemand dran |
| `in-arbeit` | jemand arbeitet gerade daran |
| `erledigt` | Akzeptanzkriterien erfüllt, Tests grün |
| `verworfen` | bewusst nicht gebaut; die Begründung bleibt in der Datei stehen |

Ein verworfener Vorgang wird **nicht gelöscht**. Der Grund, etwas nicht zu bauen, ist
oft wertvoller als der Vorgang selbst.

## Labels

| Label | Bedeutung |
|---|---|
| `ready-for-agent` | vollständig beschrieben, alle Abhängigkeiten erledigt — kann sofort umgesetzt werden |
| `blockiert` | wartet auf einen Vorgang unter `haengt-an` |
| `spaeter` | gehört zu einem Modul außerhalb der aktuellen Spezifikation |
| `inhalt` | fachliche Arbeit am Unterrichtsmaterial, keine Programmierung |
| `betrieb` | Deployment, Sicherung, Betriebsfragen |

`ready-for-agent` ist eine Zusage: Wer den Vorgang aufmacht, muss nichts mehr klären.
Ist noch eine Entscheidung offen, gehört sie erst in die Datei — dann das Label.

## Übersicht neu erzeugen

    python3 tracker/index.py

Vor dem Commit ausführen, wenn ein Status oder Label sich geändert hat.

## Verhältnis zur Spezifikation

`SPEZIFIKATION.md` sagt, **was** gebaut wird und warum. Die Vorgänge sagen, **in
welcher Reihenfolge** und **woran man merkt, dass es fertig ist**. Widersprechen sich
beide, gewinnt die Spezifikation — und der Vorgang wird korrigiert, nicht umgekehrt.
