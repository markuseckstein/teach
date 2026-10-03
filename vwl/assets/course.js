/* Kursverzeichnis + Navigation.
 *
 * Einzige Stelle, an der Lektionen und Referenzen registriert werden.
 * Neue Lektion: Eintrag in COURSE.lessons ergänzen (status "fertig").
 *
 * In Lektionen:
 *   <nav class="toc" data-auto></nav>          → Abschnitts-Inhaltsverzeichnis aus allen <h2 id="...">
 *   <nav class="lesson-nav" data-current="0001"></nav>  → Zurück · Übersicht · Weiter
 * Auf der Startseite:
 *   <div id="course-index"></div>              → komplettes Inhaltsverzeichnis
 *
 * Pfade sind relativ zum Workspace-Root; base ("" oder "../") wird automatisch bestimmt.
 */
const COURSE = {
  title: "Geld, Banken & Inflation",
  lessons: [
    { nr: "0001", file: "lessons/0001-woher-kommt-das-geld.html", title: "Woher kommt das Geld?", sub: "Guthaben sind Forderungen. Kredite schaffen Geld, Tilgung vernichtet es.", status: "fertig" },
    { nr: "0002", file: "lessons/0002-zentralbankgeld.html", title: "Zentralbankgeld", sub: "Was passiert, wenn Geld die Bank wechselt, und warum es eine Zentralbank braucht.", status: "fertig" },
    { nr: "0003", file: "lessons/0003-der-leitzins.html", title: "Der Leitzins", sub: "Wie die EZB die Kreditvergabe steuert, und was das für Bauzinsen und Immobilienpreise heißt.", status: "fertig" },
    { nr: "0004", file: "lessons/0004-inflation.html", title: "Inflation", sub: "Messen, Ursachen, Erwartungen. Wer gewinnt und wer verliert.", status: "fertig" },
    { nr: "0005", title: "Fall 1923: Hyperinflation", sub: "Kriegsfinanzierung, Notenpresse, Sachwerte gegen Geldwerte.", status: "geplant" },
    { nr: "0006", title: "Fall 1948: Währungsreform & Lastenausgleich", sub: "100 RM wurden zu 6,50 DM, Immobilienbesitzer zahlten mit.", status: "geplant" },
    { nr: "0007", title: "Bankpleite & Bank Run", sub: "Einlagensicherung, Depot als Sondervermögen, Zypern 2013.", status: "geplant" },
    { nr: "0008", title: "Synthese: Konto, Depot, Immobilie", sub: "Vier Krisentypen im Vergleich.", status: "geplant" },
  ],
  references: [
    { file: "reference/geldschoepfung-buchungen.html", title: "Spickzettel: Geldschöpfung in Buchungen", sub: "Die Grundbuchungen als T-Konten, inkl. Überweisung zwischen Banken." },
    { file: "reference/leitzins.html", title: "Spickzettel: Leitzins & Transmission", sub: "Die drei EZB-Zinsen, die Wirkungskette, die Kreditformel." },
    { file: "reference/inflation.html", title: "Spickzettel: Inflation", sub: "Warenkorb, vier Ursachen, Realzins, Gewinner und Verlierer." },
  ],
};

(function () {
  const css = `
  nav.toc { font-family: var(--sans); font-size: .8rem; border-left: 2px solid var(--rule); padding: .2rem 0 .2rem 1rem; margin: 0 0 2.5rem; }
  nav.toc .toc-h { font-size: .65rem; letter-spacing: .12em; text-transform: uppercase; color: var(--muted); margin-bottom: .3rem; }
  nav.toc ol { margin: 0; padding-left: 0; list-style: none; }
  nav.toc li { margin: .15rem 0; }
  nav.toc a { text-decoration: none; }
  nav.toc a:hover { text-decoration: underline; }
  .topbar { font-family: var(--sans); font-size: .75rem; margin-bottom: 2rem; }
  .topbar a { text-decoration: none; }
  #course-index .ci-section { margin: 2rem 0; }
  #course-index ol.ci { list-style: none; padding: 0; margin: 0; }
  #course-index ol.ci li { display: grid; grid-template-columns: 3.2rem 1fr; gap: .2rem .8rem; padding: .8rem 0; border-bottom: 1px solid var(--rule); }
  #course-index .ci-nr { font-family: var(--sans); font-size: .8rem; color: var(--muted); padding-top: .2rem; font-variant-numeric: tabular-nums; }
  #course-index .ci-t { font-size: 1.1rem; }
  #course-index .ci-t a { text-decoration: none; }
  #course-index .ci-t a:hover { text-decoration: underline; }
  #course-index .ci-s { grid-column: 2; color: var(--muted); font-size: .9rem; font-style: italic; }
  #course-index li.geplant .ci-t { color: var(--muted); }
  #course-index .badge { font-family: var(--sans); font-size: .62rem; letter-spacing: .08em; text-transform: uppercase; border: 1px solid var(--rule); border-radius: 2px; padding: .05rem .35rem; margin-left: .5rem; vertical-align: middle; color: var(--muted); }
  `;
  const style = document.createElement("style");
  style.textContent = css;
  document.head.appendChild(style);

  const inSub = /\/(lessons|reference)\/[^/]*$/.test(location.pathname);
  const base = inSub ? "../" : "";
  const href = (f) => base + f;

  document.addEventListener("DOMContentLoaded", () => {
    // Abschnitts-Inhaltsverzeichnis
    document.querySelectorAll("nav.toc[data-auto]").forEach((nav) => {
      const hs = [...document.querySelectorAll("main h2[id]")];
      nav.innerHTML = `<div class="toc-h">Inhalt</div><ol>${hs.map((h) => `<li><a href="#${h.id}">${h.textContent}</a></li>`).join("")}</ol>`;
    });

    // Leiste oben: zurück zur Übersicht
    document.querySelectorAll(".topbar[data-auto]").forEach((t) => {
      t.innerHTML = `<a href="${href("index.html")}">← ${COURSE.title}: Übersicht</a>`;
    });

    // Zurück / Übersicht / Weiter
    document.querySelectorAll("nav.lesson-nav[data-current]").forEach((nav) => {
      const ready = COURSE.lessons.filter((l) => l.status === "fertig");
      const i = ready.findIndex((l) => l.nr === nav.dataset.current);
      const prev = ready[i - 1], next = ready[i + 1];
      const upcoming = COURSE.lessons.find((l) => l.nr > nav.dataset.current && l.status !== "fertig");
      nav.innerHTML =
        (prev ? `<a href="${href(prev.file)}">← ${prev.nr} ${prev.title}</a>` : `<span class="disabled">← Anfang</span>`) +
        `<a href="${href("index.html")}">Übersicht</a>` +
        (next ? `<a href="${href(next.file)}">${next.nr} ${next.title} →</a>`
              : `<span class="disabled">${upcoming ? "Als Nächstes: " + upcoming.title : "Ende"} →</span>`);
    });

    // Startseite
    const idx = document.getElementById("course-index");
    if (idx) {
      const li = (l) => `<li class="${l.status || ""}"><span class="ci-nr">${l.nr || ""}</span>
        <span class="ci-t">${l.file && l.status !== "geplant" ? `<a href="${href(l.file)}">${l.title}</a>` : l.title}${l.status === "geplant" ? '<span class="badge">geplant</span>' : ""}</span>
        <span class="ci-s">${l.sub}</span></li>`;
      idx.innerHTML = `
        <section class="ci-section"><h2 id="lektionen">Lektionen</h2><ol class="ci">${COURSE.lessons.map(li).join("")}</ol></section>
        <section class="ci-section"><h2 id="referenz">Referenz</h2><ol class="ci">${COURSE.references.map((r) => li({ ...r, nr: "Ref" })).join("")}</ol></section>`;
    }
  });
})();
