/* Wiederverwendbares Quiz-Modul.
 *
 * Verwendung:
 *   <div class="quiz" data-quiz="meinQuiz"></div>
 *   <script src="../assets/quiz.js"></script>
 *   <script>
 *     Quiz.render("meinQuiz", [
 *       { q: "Frage?", options: ["A", "B", "C"], answer: 1, why: "Begründung mit Quelle." },
 *     ]);
 *   </script>
 *
 * - Antwortoptionen werden gemischt, damit die Position nichts verrät.
 * - Sofortige Rückmeldung mit Begründung nach dem ersten Klick.
 * - Kein Speichern: jede Wiederholung ist wieder eine echte Abrufübung.
 */
(function () {
  function shuffle(arr) {
    const a = arr.slice();
    for (let i = a.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [a[i], a[j]] = [a[j], a[i]];
    }
    return a;
  }

  function render(id, questions) {
    const root = document.querySelector('[data-quiz="' + id + '"]');
    if (!root) return;
    root.innerHTML = "";
    let answered = 0, correct = 0;

    const score = document.createElement("div");
    score.className = "quiz-score";

    function updateScore() {
      score.textContent = answered + " von " + questions.length + " beantwortet · " + correct + " richtig";
      if (answered === questions.length) {
        const reset = document.createElement("button");
        reset.className = "quiz-reset";
        reset.textContent = "Nochmal (neu gemischt)";
        reset.onclick = () => render(id, questions);
        score.appendChild(reset);
      }
    }

    questions.forEach((item, qi) => {
      const box = document.createElement("div");
      box.className = "q";
      const text = document.createElement("p");
      text.className = "q-text";
      text.textContent = (qi + 1) + ". " + item.q;
      box.appendChild(text);

      const opts = document.createElement("div");
      opts.className = "q-opts";
      const fb = document.createElement("div");
      fb.className = "q-fb";

      const order = shuffle(item.options.map((o, i) => i));
      const buttons = order.map((oi) => {
        const b = document.createElement("button");
        b.type = "button";
        b.textContent = item.options[oi];
        b.onclick = () => {
          buttons.forEach((x) => (x.disabled = true));
          const ok = oi === item.answer;
          answered++; if (ok) correct++;
          b.classList.add(ok ? "right" : "wrong");
          if (!ok) buttons[order.indexOf(item.answer)].classList.add("right");
          fb.className = "q-fb " + (ok ? "right" : "wrong");
          fb.innerHTML = "<strong>" + (ok ? "Richtig." : "Nicht ganz.") + "</strong>" +
            '<span class="why"></span>';
          fb.querySelector(".why").innerHTML = item.why || "";
          updateScore();
        };
        opts.appendChild(b);
        return b;
      });

      box.appendChild(opts);
      box.appendChild(fb);
      root.appendChild(box);
    });

    root.appendChild(score);
    updateScore();
  }

  window.Quiz = { render: render };
})();
