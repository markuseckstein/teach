/* Quiz-Komponente.
 *
 * Markup:
 * <div class="quiz">
 *   <p class="q">Frage?</p>
 *   <ul>
 *     <li data-correct data-feedback="Warum richtig">Antwort A</li>
 *     <li data-feedback="Warum falsch">Antwort B</li>
 *   </ul>
 * </div>
 * <div class="quiz-score"></div>   (optional, zeigt Gesamtstand)
 *
 * Optionen werden gemischt, damit die Position nichts verrät.
 * Erster Klick zählt für die Wertung; danach darf weiter probiert werden.
 */
(function () {
  const css = `
  .quiz { background: var(--panel); border: 1px solid var(--rule); border-radius: 4px; padding: 1rem 1.1rem; margin: 1.2rem 0; }
  .quiz .q { font-weight: 600; margin: 0 0 .7rem; }
  .quiz .q .qn { font-family: var(--sans); font-size: .7rem; color: var(--muted); font-weight: 400; display: block; letter-spacing: .1em; text-transform: uppercase; }
  .quiz ul { list-style: none; padding: 0; margin: 0; display: grid; gap: .45rem; }
  .quiz li button { width: 100%; font-family: var(--serif); font-size: .95rem; padding: .55rem .8rem; }
  .quiz li.right button { border-color: var(--good); background: var(--good-soft); }
  .quiz li.wrong button { border-color: var(--bad); background: var(--bad-soft); }
  .quiz .fb { font-size: .88rem; margin: .7rem 0 0; padding: .6rem .8rem; border-left: 3px solid var(--rule); display: none; }
  .quiz .fb.show { display: block; }
  .quiz .fb.ok { border-color: var(--good); }
  .quiz .fb.no { border-color: var(--bad); }
  .quiz-score { font-family: var(--sans); font-size: .85rem; color: var(--muted); margin: 1rem 0; }
  @media print { .quiz li button { display: block; } }
  `;
  const style = document.createElement("style");
  style.textContent = css;
  document.head.appendChild(style);

  function shuffle(a) {
    for (let i = a.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [a[i], a[j]] = [a[j], a[i]];
    }
    return a;
  }

  document.addEventListener("DOMContentLoaded", () => {
    const quizzes = [...document.querySelectorAll(".quiz")];
    const scoreEl = document.querySelector(".quiz-score");
    let answered = 0, correct = 0;
    const total = quizzes.length;

    function updateScore() {
      if (!scoreEl) return;
      scoreEl.textContent = answered === 0
        ? `${total} Fragen. Erst nachdenken, dann klicken: Die erste Antwort zählt.`
        : `Stand: ${correct} von ${answered} beim ersten Versuch richtig (${total} Fragen insgesamt).` +
          (answered === total ? (correct === total ? " Stark!" : " Lies die Erklärungen zu den falschen Antworten noch einmal.") : "");
    }

    quizzes.forEach((quiz, idx) => {
      const q = quiz.querySelector(".q");
      if (q) q.insertAdjacentHTML("afterbegin", `<span class="qn">Frage ${idx + 1} / ${total}</span>`);
      const ul = quiz.querySelector("ul");
      const items = shuffle([...ul.children]);
      items.forEach((li) => ul.appendChild(li));
      const fb = document.createElement("p");
      fb.className = "fb";
      quiz.appendChild(fb);
      let first = true;

      items.forEach((li) => {
        const btn = document.createElement("button");
        btn.innerHTML = li.innerHTML;
        li.innerHTML = "";
        li.appendChild(btn);
        btn.addEventListener("click", () => {
          const ok = li.hasAttribute("data-correct");
          items.forEach((x) => x.classList.remove("right", "wrong"));
          li.classList.add(ok ? "right" : "wrong");
          fb.className = "fb show " + (ok ? "ok" : "no");
          fb.innerHTML = (ok ? "<strong>Richtig.</strong> " : "<strong>Nicht ganz.</strong> ") + (li.dataset.feedback || "");
          if (first) {
            first = false;
            answered++;
            if (ok) correct++;
            updateScore();
          }
        });
      });
    });
    updateScore();
  });
})();
