#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Erzeugt tracker/INDEX.md aus den Frontmattern in tracker/vorgaenge/."""
import pathlib
import re

HIER = pathlib.Path(__file__).resolve().parent
SYMBOL = {"offen": "○", "in-arbeit": "◐", "erledigt": "●", "verworfen": "✕"}


def kopf(text):
    treffer = re.match(r"^---\n(.*?)\n---\n", text, re.S)
    if not treffer:
        raise SystemExit("Frontmatter fehlt")
    daten = {}
    for zeile in treffer.group(1).splitlines():
        if ":" not in zeile:
            continue
        schluessel, wert = zeile.split(":", 1)
        wert = wert.strip()
        if wert.startswith("["):
            wert = [t.strip().strip('"') for t in wert.strip("[]").split(",") if t.strip()]
        daten[schluessel.strip()] = wert
    return daten


def main():
    zeilen = [
        "# Übersicht",
        "",
        "Erzeugt mit `python3 tracker/index.py` — nicht von Hand ändern.",
        "",
        "| | Nr. | Vorgang | Status | Label | hängt an |",
        "|---|---|---|---|---|---|",
    ]
    offen = arbeit = fertig = 0
    for datei in sorted((HIER / "vorgaenge").glob("*.md")):
        d = kopf(datei.read_text(encoding="utf-8"))
        status = d.get("status", "offen")
        offen += status == "offen"
        arbeit += status == "in-arbeit"
        fertig += status == "erledigt"
        haengt = d.get("haengt-an") or []
        zeilen.append(
            "| %s | %s | [%s](vorgaenge/%s) | %s | `%s` | %s |"
            % (
                SYMBOL.get(status, "?"),
                d.get("id", "?"),
                d.get("titel", datei.stem),
                datei.name,
                status,
                d.get("label", ""),
                ", ".join(haengt) if haengt else "—",
            )
        )
    zeilen += ["", "%d offen · %d in Arbeit · %d erledigt" % (offen, arbeit, fertig), ""]
    (HIER / "INDEX.md").write_text("\n".join(zeilen), encoding="utf-8")
    print("INDEX.md erzeugt:", offen + arbeit + fertig, "Vorgänge")


if __name__ == "__main__":
    main()
