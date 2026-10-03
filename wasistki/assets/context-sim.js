/* ---------------------------------------------------------------------------
   context-sim.js — Komponente: der Schreibtisch des Chatbots.

   Einbinden mit:  <div class="ctxsim"></div>

   Zeigt anschaulich, wie sich ein Kontextfenster füllt, was passiert, wenn es
   überläuft, und wie ein altes Thema in ein neues hineinredet.

   WICHTIG: Das ist eine Veranschaulichung, keine Messung. Die Zahlen sind
   plausibel gewählt, nicht gemessen. Der Zweck ist das Gefühl für den Verlauf —
   dass es nicht linear schlechter wird, sondern lange gut und dann schnell
   schlecht. Genau so steht es auch in der Beschriftung unter dem Baustein.
   --------------------------------------------------------------------------- */

(function () {
  'use strict';

  var POOL = [
    { id: 'sys',  label: 'Systemanweisung des Anbieters', cost: 4,  topic: 'setup', fixed: true },
    { id: 'brief',label: 'Elternbrief-Entwurf',           cost: 8,  topic: 'a' },
    { id: 'fb',   label: 'Ihre Korrekturwünsche dazu',    cost: 6,  topic: 'a' },
    { id: 'pdf',  label: 'Lehrplan-PDF, 40 Seiten',       cost: 46, topic: 'a' },
    { id: 'tab',  label: 'Notentabelle eingefügt',        cost: 14, topic: 'a' },
    { id: 'web',  label: 'Ergebnis einer Websuche',       cost: 12, topic: 'a' },
    { id: 'neu',  label: 'Neue Frage: Physik Klasse 9',   cost: 5,  topic: 'b' }
  ];

  function quality(fill) {
    if (fill <= 30) return { label: 'hoch',              cls: 'ok',   note: 'Alles im Blick.' };
    if (fill <= 60) return { label: 'noch gut',          cls: 'ok',   note: 'Die Mitte wird unschärfer.' };
    if (fill <= 85) return { label: 'sinkt merklich',    cls: 'warn', note: 'Details aus der Mitte gehen unter.' };
    return              { label: 'unzuverlässig',        cls: 'bad',  note: 'Anfang und Ende zählen, der Rest kaum noch.' };
  }

  function build(host) {
    var state = { items: [POOL[0]] };

    host.classList.add('ctxsim-ready');
    host.innerHTML =
      '<div class="ctxsim-head">' +
        '<span class="ctxsim-title">Der Schreibtisch</span>' +
        '<span class="ctxsim-fill" data-fill></span>' +
      '</div>' +
      '<div class="ctxsim-bar" data-bar></div>' +
      '<ul class="ctxsim-list" data-list></ul>' +
      '<div class="ctxsim-verdict" data-verdict></div>' +
      '<div class="ctxsim-actions" data-actions></div>';

    var barEl     = host.querySelector('[data-bar]');
    var listEl    = host.querySelector('[data-list]');
    var fillEl    = host.querySelector('[data-fill]');
    var verdictEl = host.querySelector('[data-verdict]');
    var actionsEl = host.querySelector('[data-actions]');

    function total() {
      return state.items.reduce(function (s, i) { return s + i.cost; }, 0);
    }

    function render() {
      var sum = total();
      var over = sum > 100;

      // Beim Überlauf fallen die ältesten beweglichen Einträge heraus.
      var dropped = [];
      if (over) {
        var running = sum;
        for (var i = 0; i < state.items.length && running > 100; i++) {
          if (state.items[i].fixed) continue;
          dropped.push(state.items[i].id);
          running -= state.items[i].cost;
        }
      }

      var shown = Math.min(sum, 100);
      fillEl.textContent = shown + ' % voll' + (over ? ' — übergelaufen' : '');
      fillEl.className = 'ctxsim-fill ' + quality(shown).cls;

      barEl.innerHTML = '';
      state.items.forEach(function (it) {
        var seg = document.createElement('span');
        seg.className = 'ctxsim-seg t-' + it.topic +
          (dropped.indexOf(it.id) > -1 ? ' is-dropped' : '');
        seg.style.flexGrow = it.cost;
        seg.title = it.label;
        barEl.appendChild(seg);
      });
      if (!over && shown < 100) {
        var rest = document.createElement('span');
        rest.className = 'ctxsim-seg is-free';
        rest.style.flexGrow = 100 - shown;
        barEl.appendChild(rest);
      }

      listEl.innerHTML = '';
      state.items.forEach(function (it) {
        var li = document.createElement('li');
        var gone = dropped.indexOf(it.id) > -1;
        li.className = 'ctxsim-item t-' + it.topic + (gone ? ' is-dropped' : '');
        li.innerHTML =
          '<span class="ctxsim-dot"></span>' +
          '<span class="ctxsim-label">' + it.label + '</span>' +
          '<span class="ctxsim-cost">' + it.cost + '</span>' +
          (gone ? '<span class="ctxsim-gone">herausgefallen</span>' : '');
        listEl.appendChild(li);
      });

      var q = quality(shown);
      var msg = '<b>Zuverlässigkeit: ' + q.label + '.</b> ' + q.note;

      if (over) {
        msg += ' <b>Und jetzt wird es unangenehm:</b> Der Anfang passt nicht mehr ' +
               'hinein und fällt weg. Der Chatbot sagt Ihnen das nicht. Er antwortet ' +
               'einfach weiter — nur eben ohne das, was Sie ganz am Anfang festgelegt haben.';
      }

      var hasA = state.items.some(function (i) { return i.topic === 'a'; });
      var hasB = state.items.some(function (i) { return i.topic === 'b'; });
      if (hasA && hasB) {
        msg += ' <b>Themenverschleppung:</b> Ihre Physikfrage liegt jetzt auf einem ' +
               'Schreibtisch voller Deutsch-Material. Das Modell sieht beides gleichzeitig ' +
               'und mischt — Sie bekommen Physik im Ton des Elternbriefs, mit Bezügen ' +
               'zum Lehrplan-PDF, nach dem Sie nicht gefragt haben.';
      }

      verdictEl.innerHTML = msg;
      verdictEl.className = 'ctxsim-verdict ' + q.cls;

      actionsEl.innerHTML = '';
      POOL.forEach(function (p) {
        if (p.fixed) return;
        if (state.items.indexOf(p) > -1) return;
        var b = document.createElement('button');
        b.className = 'ctxsim-btn t-' + p.topic;
        b.innerHTML = '+ ' + p.label + ' <span class="ctxsim-btn-cost">' + p.cost + '</span>';
        b.addEventListener('click', function () {
          state.items.push(p);
          render();
        });
        actionsEl.appendChild(b);
      });

      var reset = document.createElement('button');
      reset.className = 'ctxsim-btn is-reset';
      reset.textContent = '↺ Neuer Chat';
      reset.addEventListener('click', function () {
        state.items = [POOL[0]];
        render();
      });
      actionsEl.appendChild(reset);
    }

    render();
  }

  document.addEventListener('DOMContentLoaded', function () {
    Array.prototype.forEach.call(document.querySelectorAll('.ctxsim'), build);
  });
})();
