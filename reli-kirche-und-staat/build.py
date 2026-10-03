#!/usr/bin/env python3
"""Baut die Kopiervorlagen: setzt Kopfblock (Schriften + Stil) und Bilder ein.

Jede Quelldatei in lessons-src/ enthaelt nur den Inhalt und zwei Platzhalter:
  <!--KOPF-->        wird durch assets/kopf-vorlage.html ersetzt
  BILD:<schluessel>  wird durch das Bild aus assets/img/<schluessel>.jpg
                     als data:-URI ersetzt (die Artifact-CSP blockt externe Bilder)
"""
import base64, pathlib, re, sys

ROOT = pathlib.Path(__file__).parent
KOPF = (ROOT/"assets/kopf-vorlage.html").read_text(encoding="utf-8")

def bild(key):
    p = ROOT/f"assets/img/{key}.jpg"
    if not p.exists():
        sys.exit(f"Bild fehlt: {p}")
    return "data:image/jpeg;base64," + base64.b64encode(p.read_bytes()).decode()

for src in sorted((ROOT/"lessons-src").glob("*.html")):
    s = src.read_text(encoding="utf-8")
    s = s.replace("<!--KOPF-->", '<meta charset="utf-8">\n' + KOPF)
    s = re.sub(r'BILD:([a-z0-9_]+)', lambda m: bild(m.group(1)), s)
    rest = re.findall(r'BILD:\S+|<!--KOPF-->', s)
    if rest:
        sys.exit(f"{src.name}: unaufgeloeste Platzhalter {rest}")
    out = ROOT/"lessons"/src.name
    out.write_text(s, encoding="utf-8")
    print(f"{out.name:34s} {round(len(s.encode())/1024):5d} KB")
