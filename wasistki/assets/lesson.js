/* ---------------------------------------------------------------------------
   lesson.js — baut Navigation und interaktive Bausteine aus COURSE (course.js).

   Erwartet auf <body> ein Attribut:
     data-page="lesson"     + data-n="3"     → Lektion Nummer 3
     data-page="reference"  + data-file="glossar.html"
     data-page="index"

   Bausteine, die jede Lektion nutzen kann:
     .quiz      — Mehrfachauswahl mit sofortiger Rückmeldung
     .reveal    — Aufklappen (reines <details>, braucht kein JS)
   --------------------------------------------------------------------------- */

(function () {
  'use strict';

  var body = document.body;
  var page = body.getAttribute('data-page') || 'index';
  var root = page === 'index' ? '' : '../';

  function el(tag, cls, html) {
    var e = document.createElement(tag);
    if (cls) e.className = cls;
    if (html !== undefined) e.innerHTML = html;
    return e;
  }

  function lessonHref(l) { return root + 'lessons/' + l.file; }
  function refHref(r) { return root + 'reference/' + r.file; }
  function indexHref() { return root + 'index.html'; }

  /* ---------------------------------------------------------------- Kopfzeile */

  function buildTopBar(current) {
    var bar = el('nav', 'topbar');

    var home = el('a', 'topbar-home',
      '<span class="topbar-mark" aria-hidden="true">◆</span> ' + COURSE.title);
    home.href = indexHref();
    bar.appendChild(home);

    var spacer = el('span', 'topbar-spacer');
    bar.appendChild(spacer);

    if (current) {
      var pos = el('span', 'topbar-pos',
        'Lektion ' + current.n + ' <span class="of">von</span> ' + COURSE.lessons.length);
      bar.appendChild(pos);
    }

    // Sprungliste über alle Lektionen — immer vollständig, immer aktuell.
    var jump = el('div', 'topbar-jump');
    var toggle = el('button', 'jump-toggle', 'Alle Lektionen <span class="caret">▾</span>');
    toggle.setAttribute('aria-expanded', 'false');
    var list = el('div', 'jump-list');

    var lastPart = null;
    COURSE.lessons.forEach(function (l) {
      if (l.part !== lastPart) {
        list.appendChild(el('div', 'jump-part', l.part));
        lastPart = l.part;
      }
      var a = el('a', 'jump-item' + (current && l.n === current.n ? ' is-current' : ''),
        '<span class="jump-n">' + l.n + '</span> ' + l.title);
      a.href = lessonHref(l);
      list.appendChild(a);
    });

    list.appendChild(el('div', 'jump-part', 'Zum Nachschlagen'));
    COURSE.reference.forEach(function (r) {
      var a = el('a', 'jump-item jump-ref', '<span class="jump-n">◇</span> ' + r.title);
      a.href = refHref(r);
      list.appendChild(a);
    });

    toggle.addEventListener('click', function (ev) {
      ev.stopPropagation();
      var open = jump.classList.toggle('is-open');
      toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
    });
    document.addEventListener('click', function () {
      jump.classList.remove('is-open');
      toggle.setAttribute('aria-expanded', 'false');
    });

    jump.appendChild(toggle);
    jump.appendChild(list);
    bar.appendChild(jump);

    return bar;
  }

  /* ---------------------------------------------------------------- Fußzeile */

  function buildBottomNav(current) {
    var wrap = el('nav', 'pagenav');
    var idx = COURSE.lessons.indexOf(current);
    var prev = idx > 0 ? COURSE.lessons[idx - 1] : null;
    var next = idx > -1 && idx < COURSE.lessons.length - 1 ? COURSE.lessons[idx + 1] : null;

    if (prev) {
      var p = el('a', 'pagenav-card pagenav-prev',
        '<span class="pagenav-dir">← Zurück</span>' +
        '<span class="pagenav-title">' + prev.n + '. ' + prev.title + '</span>');
      p.href = lessonHref(prev);
      wrap.appendChild(p);
    } else {
      wrap.appendChild(el('span', 'pagenav-card is-empty'));
    }

    if (next) {
      var n = el('a', 'pagenav-card pagenav-next',
        '<span class="pagenav-dir">Weiter →</span>' +
        '<span class="pagenav-title">' + next.n + '. ' + next.title + '</span>' +
        '<span class="pagenav-teaser">' + next.teaser + '</span>');
      n.href = lessonHref(next);
      wrap.appendChild(n);
    } else {
      var done = el('a', 'pagenav-card pagenav-next',
        '<span class="pagenav-dir">Geschafft</span>' +
        '<span class="pagenav-title">Zurück zur Übersicht</span>' +
        '<span class="pagenav-teaser">Das war die letzte Lektion. Die Referenzdokumente bleiben Ihr Nachschlagewerk.</span>');
      done.href = indexHref();
      wrap.appendChild(done);
    }
    return wrap;
  }

  /* ------------------------------------------------------------ Fortschritt */

  function buildProgress() {
    var bar = el('div', 'readbar');
    var fill = el('div', 'readbar-fill');
    bar.appendChild(fill);
    function update() {
      var h = document.documentElement;
      var max = h.scrollHeight - h.clientHeight;
      var pct = max > 0 ? (h.scrollTop / max) * 100 : 0;
      fill.style.width = pct.toFixed(1) + '%';
    }
    window.addEventListener('scroll', update, { passive: true });
    window.addEventListener('resize', update);
    update();
    return bar;
  }

  /* ------------------------------------------------------------------ Quiz */

  function initQuizzes() {
    var quizzes = document.querySelectorAll('.quiz');
    Array.prototype.forEach.call(quizzes, function (q) {
      var opts = q.querySelectorAll('.quiz-opts li');
      var explain = q.querySelector('.quiz-explain');
      if (explain) explain.hidden = true;
      var answered = false;

      Array.prototype.forEach.call(opts, function (li) {
        var isCorrect = li.hasAttribute('data-correct');
        li.setAttribute('role', 'button');
        li.setAttribute('tabindex', '0');

        function choose() {
          if (answered) return;
          answered = true;
          q.classList.add('is-answered');
          li.classList.add(isCorrect ? 'is-picked-right' : 'is-picked-wrong');
          Array.prototype.forEach.call(opts, function (o) {
            if (o.hasAttribute('data-correct')) o.classList.add('is-right');
            o.setAttribute('tabindex', '-1');
          });
          if (explain) {
            explain.hidden = false;
            explain.classList.add('is-shown');
          }
        }
        li.addEventListener('click', choose);
        li.addEventListener('keydown', function (ev) {
          if (ev.key === 'Enter' || ev.key === ' ') { ev.preventDefault(); choose(); }
        });
      });
    });
  }

  /* ------------------------------------------------- Inhaltsverzeichnis */

  // Baut das Verzeichnis auf index.html. Reihenfolge, Titel und Teaser kommen
  // ausschliesslich aus COURSE — dieselbe Quelle wie die Nav-Links.
  function buildToc() {
    var host = document.getElementById('toc');
    if (!host) return;

    var lastPart = null;
    COURSE.lessons.forEach(function (l) {
      if (l.part !== lastPart) {
        host.appendChild(el('h2', 'toc-part', l.part));
        lastPart = l.part;
      }
      var a = el('a', 'toc-item',
        '<span class="toc-n">' + String(l.n).padStart(2, '0') + '</span>' +
        '<span class="toc-title">' + l.title + '</span>' +
        '<span class="toc-min">' + l.minutes + ' Min.</span>' +
        '<span class="toc-teaser">' + l.teaser + '</span>');
      a.href = lessonHref(l);
      host.appendChild(a);
    });

    var refHost = document.getElementById('ref-toc');
    if (!refHost) return;
    COURSE.reference.forEach(function (r) {
      var a = el('a', 'toc-item is-ref',
        '<span class="toc-n">◇</span>' +
        '<span class="toc-title">' + r.title + '</span>' +
        '<span class="toc-min">Referenz</span>' +
        '<span class="toc-teaser">' + r.teaser + '</span>');
      a.href = refHref(r);
      refHost.appendChild(a);
    });
  }

  /* --------------------------------------------------------------- Aufbau */

  document.addEventListener('DOMContentLoaded', function () {
    var current = null;
    if (page === 'lesson') {
      var n = parseInt(body.getAttribute('data-n'), 10);
      COURSE.lessons.forEach(function (l) { if (l.n === n) current = l; });
    }

    body.insertBefore(buildProgress(), body.firstChild);
    body.insertBefore(buildTopBar(current), body.children[1] || null);

    if (current) {
      var main = document.querySelector('main') || body;
      main.appendChild(buildBottomNav(current));
    }

    buildToc();
    initQuizzes();
  });
})();
