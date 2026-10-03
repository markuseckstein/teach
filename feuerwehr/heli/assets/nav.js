/* Kursnavigation: eine zentrale Liste aller Seiten.
 *
 * Verwendung auf jeder Seite (Pfad zu assets/ relativ anpassen):
 *   <nav class="topnav" data-course-top></nav>      (optional, oben)
 *   <nav class="course" data-course-nav></nav>      (unten: zurück · Inhalt · weiter)
 *   <script src="../assets/nav.js"></script>
 *
 * Neue Lektion/Referenz: nur hier in COURSE eintragen.
 * Pfade sind relativ zum Kursverzeichnis (dem Ordner über assets/).
 */
(function () {
  var COURSE = [
    { href: "index.html", title: "Inhaltsverzeichnis", kind: "Übersicht" },
    { href: "lessons/0001-der-nachtlandeplatz.html", title: "Der Nachtlandeplatz in sechs Schritten", kind: "Lektion 1" },
    { href: "reference/nachtlandeplatz-taschenkarte.html", title: "Taschenkarte Nachtlandeplatz", kind: "Referenz" },
    { href: "lessons/0002-die-uebung-als-drehbuch.html", title: "Die Übung als Drehbuch", kind: "Lektion 2" },
    { href: "reference/uebung-2026-10-01-drehbuch.html", title: "Übungsdrehbuch Pattenhofen 1.10.", kind: "Referenz" },
    { href: "reference/melden-lagemeldung.html", title: "Meldekarte MELDEN", kind: "Referenz" },
    { href: "reference/briefing-leitfaden.html", title: "Briefing-Leitfaden (Übungsleiter)", kind: "Referenz" },
    { href: "reference/handout-mannschaft.html", title: "Handout für die Mannschaft", kind: "Handout" }
  ];

  var script = document.currentScript;
  var root = new URL("../", script.src);            // Kursverzeichnis
  var here = location.href.split("#")[0];

  var idx = -1;
  COURSE.forEach(function (p, i) {
    if (new URL(p.href, root).href === here) idx = i;
  });

  function link(p, label) {
    var a = document.createElement("a");
    a.href = new URL(p.href, root).href;
    a.textContent = label;
    return a;
  }

  var bottom = document.querySelector("[data-course-nav]");
  if (bottom) {
    bottom.innerHTML = "";
    var prev = idx > 0 ? COURSE[idx - 1] : null;
    var next = idx >= 0 && idx < COURSE.length - 1 ? COURSE[idx + 1] : null;
    bottom.appendChild(prev ? link(prev, "← " + prev.kind + ": " + prev.title) : document.createElement("span"));
    if (idx !== 0) bottom.appendChild(link(COURSE[0], "Inhaltsverzeichnis"));
    bottom.appendChild(next ? link(next, next.kind + ": " + next.title + " →") : document.createElement("span"));
  }

  var top = document.querySelector("[data-course-top]");
  if (top && idx !== 0) {
    top.innerHTML = "";
    top.appendChild(link(COURSE[0], "☰ Inhaltsverzeichnis"));
    if (idx > 0) {
      var s = document.createElement("span");
      s.textContent = "Seite " + idx + " von " + (COURSE.length - 1);
      top.appendChild(s);
    }
  }

  window.CourseNav = { pages: COURSE, root: root };
})();
