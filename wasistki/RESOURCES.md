# Ressourcen: Was ist eigentlich „KI"?

Stand: August 2026. Alles unter „Landschaft" altert schnell — vor dem nächsten Workshop kurz gegenprüfen.

## Knowledge — Grundlagen

- [Agent Skills — offizielle Spezifikation (agentskills.io)](https://agentskills.io/home)
  Primärquelle für alles zum Thema Skills: Was eine `SKILL.md` ist, wie *progressive disclosure* in drei Stufen (Discovery → Activation → Execution) funktioniert, und welche Werkzeuge das Format unterstützen. Von Anthropic entwickelt, im Dezember 2025 als offener Standard freigegeben. **Nutzen für:** Lektion 7, Glossareintrag „Skill".

- [Context Rot: How Increasing Input Tokens Impacts LLM Performance — Chroma Research](https://www.trychroma.com/research/context-rot)
  Kontrollierte Studie über 18 aktuelle Modelle. Kernbefund im Originalwortlaut: *„models do not use their context uniformly; instead, their performance grows increasingly unreliable as input length grows."* Das ist der empirische Beleg dafür, dass ein voller Chat schlechter wird — nicht nur ein Bauchgefühl. **Nutzen für:** Lektion 3, die Begründung fürs Kontext-Leeren.

- [Lost in the Middle: How Language Models Use Long Contexts — Liu et al., Stanford (TACL 2024)](https://arxiv.org/abs/2307.03172)
  Die klassische Arbeit zum U-förmigen Verlauf: Modelle finden Informationen am Anfang und am Ende ihres Kontexts zuverlässig, in der Mitte deutlich schlechter. **Nutzen für:** das „Schreibtisch"-Bild in Lektion 3.

- [Attention Is All You Need — Vaswani et al., 2017](https://arxiv.org/abs/1706.03762)
  Die Ursprungsarbeit zur Transformer-Architektur. Für die Zielgruppe zu technisch, aber sie gehört als Fußnote in die Ressourcenliste — jemand wird danach fragen, wo das alles herkommt.

- [But what is a GPT? — 3Blue1Brown (Video, ~27 Min.)](https://www.youtube.com/watch?v=wjZofJX0v4M)
  Die mit Abstand beste visuelle Erklärung davon, was im Inneren passiert. Anschauliche Animationen, keine Mathematikkenntnisse nötig. **Nutzen für:** die Person im Workshop, die „aber wie funktioniert das *wirklich*" fragt.

- [Intro to Large Language Models — Andrej Karpathy (Video, ~60 Min.)](https://www.youtube.com/watch?v=zjkBMFhNj_g)
  Der beste Einstieg von einem der bekanntesten Praktiker. Erklärt insbesondere sauber, dass ein LLM erst einmal nur zwei Dateien ist. **Nutzen für:** Lektion 2, „das LLM ist zustandslos".

## Knowledge — Landschaft (altert schnell)

- [Anthropic — Claude Cowork](https://claude.com/product/cowork)
  Primärquelle zum Produkt. Leitspruch der Seite: *„Say what, not how."* Arbeitet in freigegebenen Ordnern und Anwendungen, zeigt jeden Schritt transparent an, läuft weiter, wenn man offline ist. Setzt einen kostenpflichtigen Plan voraus. **Nutzen für:** Lektion 8.

- [Claude Code — Dokumentation](https://code.claude.com/docs/)
  Primärquelle zum Agent-Harness für Entwickler. Auch relevant, weil hier Skills, Hooks und Unteragenten dokumentiert sind. **Nutzen für:** Lektionen 6–8.

- [OpenClaw — GitHub](https://github.com/openclaw/openclaw)
  Quelloffener, selbst gehosteter persönlicher Agent, der über die Kanäle erreichbar ist, die man ohnehin nutzt. Wichtig sind vor allem die Sicherheitshinweise der README: *„Treat inbound messages as untrusted input"*, und Werkzeuge laufen ohne zusätzliche Konfiguration direkt auf dem eigenen Rechner. **Nutzen für:** Lektion 8, Abschnitt „Was daran unbequem ist".

- [Matt Pococks Skills-Sammlung — GitHub](https://github.com/mattpocock/skills)
  Die Quelle der Skills `teach` und `grill-me`. Installation: `claude plugins install mattpocock-skills`. **Nutzen für:** Lektion 7, das durchgehende Beispiel.

- [LLM Stats — Leaderboard](https://llm-stats.com/)
  Fortlaufend aktualisierter Vergleich von Modellen nach Intelligenz, Geschwindigkeit und Preis. **Nutzen für:** die Modelltabelle auffrischen, bevor der Workshop gehalten wird. Nicht auswendig lernen — nachschlagen.

## Knowledge — Datenschutz und Schule

- [Datenschutzkonferenz (DSK) — Orientierungshilfe „Künstliche Intelligenz und Datenschutz"](https://www.datenschutzkonferenz-online.de/orientierungshilfen.html)
  Die abgestimmte Position der deutschen Datenschutzaufsichtsbehörden. Für Lehrkräfte die relevanteste offizielle Quelle. **Nutzen für:** Lektion 5. Vor dem Workshop auf die aktuelle Fassung schauen.

- Landesspezifische Handreichungen der jeweiligen Kultusministerien und Landesdatenschutzbeauftragten.
  **Wichtig:** Schulrecht ist Ländersache. Was in Bayern gilt, gilt nicht in NRW. Vor dem Workshop die Handreichung *des eigenen Bundeslandes* heraussuchen und verlinken — eine allgemeine Aussage reicht hier nicht.

## Wisdom — Communities

- [r/LocalLLaMA](https://www.reddit.com/r/LocalLLaMA/)
  Trotz des Namens die beste allgemeine Community zu Modellen und deren Fähigkeiten. Hohes Signal, technisch, aber Einsteigerfragen sind willkommen.

- [Anthropic Discord](https://discord.gg/anthropic)
  Für konkrete Fragen zu Claude Code, Cowork und Skills. Hier sitzen die Leute, die diese Werkzeuge täglich benutzen.

- **Das eigene Lehrerzimmer.** Die ehrlichste Rückmeldequelle, die es für diese Mission gibt: Wenn ein Rezept aus Lektion 9 einer Kollegin am Montag wirklich Zeit gespart hat, stimmt die Lektion. Wenn nicht, stimmt sie nicht — egal wie gut sie klingt.

- Fachnetzwerke für Lehrkräfte zu digitalen Medien auf Landesebene (z. B. Medienzentren, Landesinstitute).
  Hier laufen die schulrechtlich abgesicherten Empfehlungen zusammen — deutlich belastbarer als allgemeine KI-Ratgeber.

## Gaps

- **Keine belastbare deutschsprachige Primärquelle** gefunden, die die Ebenen Chatbot / Harness / autonomer Agent sauber trennt. Genau diese Lücke füllt dieser Workspace — die Begriffe in `GLOSSARY.md` sind daher teils selbst geprägt und als solche markiert.
- **Datenschutz:** Die Rechtslage zur Nutzung von KI-Diensten an Schulen ist in Bewegung und je nach Bundesland unterschiedlich. Die Lektion vermittelt bewusst eine *Faustregel*, keine Rechtsberatung, und verweist ausdrücklich auf die Schulleitung.
- **Kein guter Beleg** dafür gefunden, wie stark das Leeren des Kontexts die Ergebnisqualität in der Praxis verbessert. Die Empfehlung stützt sich auf die Context-Rot-Studie plus Erfahrungswerte — das wird in der Lektion offen gesagt.
