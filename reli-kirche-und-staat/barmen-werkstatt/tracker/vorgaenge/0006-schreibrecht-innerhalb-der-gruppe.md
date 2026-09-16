---
id: 0006
titel: Schreibrecht innerhalb der Gruppe
status: erledigt
label: blockiert
haengt-an: ["0005"]
---

# 0006 — Schreibrecht innerhalb der Gruppe

## Worum es geht

Genau ein Gerät je Gruppe schreibt. Die Übergabe soll eine sichtbare Handlung sein,
nicht davon abhängen, wer das Tablet festhält.

## Fertig, wenn

- Das erste Gerät einer Gruppe hält das Schreibrecht.
- Jedes andere Gerät der Gruppe kann es per Knopf an sich ziehen; die Übernahme wirkt
  sofort und ist für alle Geräte der Gruppe sichtbar.
- Ein Gerät ohne Schreibrecht kann keine Eingabe speichern, auch nicht durch eine
  direkt abgesetzte Anfrage.
- Der Zeitpunkt der letzten Übernahme wird festgehalten (Grundlage für die Warnung im
  Regiepult, Vorgang 0011).

## Gehört nicht dazu

Gleichzeitiges Bearbeiten. Es gibt bewusst kein CRDT und keine Konfliktauflösung — bei
einer 22-Minuten-Phase stünden die Kosten in keinem Verhältnis, und die sichtbare
Übergabe ist didaktisch das bessere Verhalten.
