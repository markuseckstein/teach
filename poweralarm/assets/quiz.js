/* Reusable retrieval-practice quiz widget.
 *
 * Usage in a lesson:
 *   <div class="quiz" data-quiz></div>
 *   <script src="../assets/quiz.js"></script>
 *   <script>
 *     PowerQuiz.render(document.querySelector('[data-quiz]'), [
 *       { q: 'Frage?', a: ['Antwort A', 'Antwort B'], correct: 0,
 *         why: 'Erklärung, warum A richtig ist.' },
 *     ]);
 *   </script>
 *
 * Design rules (see SKILL.md):
 *  - Answers must be of near-identical length so formatting leaks no clue.
 *  - Feedback is immediate and explains *why*, not just right/wrong.
 *  - Options are shuffled per render so repeat visits are real retrieval,
 *    not position recall.
 */
(function (global) {
  'use strict';

  var CSS = [
    '.quiz{margin:2em 0;padding:0}',
    '.quiz-item{margin-bottom:1.8em;padding:1rem 1.2rem;border:1px solid var(--rule);border-radius:5px}',
    '.quiz-q{font-weight:600;margin:0 0 .8em}',
    '.quiz-opts{list-style:none;padding:0;margin:0}',
    '.quiz-opts li{margin:0 0 .45em}',
    '.quiz-opt{display:block;width:100%;text-align:left;font:inherit;font-size:.98rem;',
    'background:var(--code-bg);color:var(--ink);border:1px solid var(--rule);',
    'border-radius:4px;padding:.55em .8em;cursor:pointer}',
    '.quiz-opt:hover:not(:disabled){border-color:var(--accent)}',
    '.quiz-opt:disabled{cursor:default;opacity:.75}',
    '.quiz-opt.is-right{background:var(--ok-soft);border-color:var(--ok);color:var(--ink);opacity:1}',
    '.quiz-opt.is-wrong{background:var(--accent-soft);border-color:var(--accent);opacity:1}',
    '.quiz-why{margin:.9em 0 0;font-size:.95rem;color:var(--ink-soft);display:none}',
    '.quiz-why.shown{display:block}',
    '.quiz-score{font-family:"Gill Sans","Gill Sans MT",Calibri,system-ui,sans-serif;',
    'font-size:.8rem;letter-spacing:.07em;text-transform:uppercase;color:var(--ink-soft)}',
    '.quiz-again{font:inherit;font-size:.9rem;background:none;border:0;border-bottom:1px solid var(--accent);',
    'color:var(--ink);cursor:pointer;padding:0}',
    '@media print{.quiz-opt{background:none}.quiz-why{display:block}.quiz-again{display:none}}'
  ].join('');

  function injectCSS() {
    if (document.getElementById('powerquiz-css')) return;
    var s = document.createElement('style');
    s.id = 'powerquiz-css';
    s.textContent = CSS;
    document.head.appendChild(s);
  }

  function shuffle(n) {
    var idx = [], i, j, t;
    for (i = 0; i < n; i++) idx.push(i);
    for (i = n - 1; i > 0; i--) {
      j = Math.floor(Math.random() * (i + 1));
      t = idx[i]; idx[i] = idx[j]; idx[j] = t;
    }
    return idx;
  }

  function render(root, items) {
    if (!root) return;
    injectCSS();
    root.innerHTML = '';

    var answered = 0, right = 0;

    var score = document.createElement('p');
    score.className = 'quiz-score';
    function updateScore() {
      score.textContent = answered
        ? right + ' von ' + answered + ' richtig · ' + items.length + ' Fragen insgesamt'
        : items.length + ' Fragen · aus dem Kopf, nicht nachblättern';
    }
    updateScore();
    root.appendChild(score);

    items.forEach(function (item) {
      var order = shuffle(item.a.length);

      var box = document.createElement('div');
      box.className = 'quiz-item';

      var q = document.createElement('p');
      q.className = 'quiz-q';
      q.textContent = item.q;
      box.appendChild(q);

      var ul = document.createElement('ul');
      ul.className = 'quiz-opts';

      var why = document.createElement('p');
      why.className = 'quiz-why';
      why.textContent = item.why || '';

      var buttons = [];
      var done = false;

      order.forEach(function (origIdx) {
        var li = document.createElement('li');
        var btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'quiz-opt';
        btn.textContent = item.a[origIdx];
        btn.addEventListener('click', function () {
          if (done) return;
          done = true;
          var isRight = origIdx === item.correct;
          answered++;
          if (isRight) right++;
          buttons.forEach(function (b) {
            b.disabled = true;
            if (b.dataset.orig === String(item.correct)) b.classList.add('is-right');
          });
          if (!isRight) btn.classList.add('is-wrong');
          why.classList.add('shown');
          updateScore();
        });
        btn.dataset.orig = String(origIdx);
        buttons.push(btn);
        li.appendChild(btn);
        ul.appendChild(li);
      });

      box.appendChild(ul);
      box.appendChild(why);
      root.appendChild(box);
    });

    var again = document.createElement('button');
    again.type = 'button';
    again.className = 'quiz-again';
    again.textContent = 'Nochmal — Reihenfolge neu mischen';
    again.addEventListener('click', function () { render(root, items); });
    root.appendChild(again);
  }

  global.PowerQuiz = { render: render };
})(window);
