// Shared quiz component. Usage in a lesson:
//   <div class="quiz" id="quiz"></div>
//   <script src="../assets/quiz.js"></script>
//   <script>renderQuiz(document.getElementById('quiz'), [
//     { stem: "…?", options: ["a", "b", "c"], answer: 1, explain: "…" }
//   ]);</script>
// Options are shuffled on each render so position gives no clue.

function renderQuiz(container, questions) {
  container.innerHTML = '';
  questions.forEach(function (q, qi) {
    var box = document.createElement('div');
    box.className = 'quiz-q';

    var stem = document.createElement('p');
    stem.className = 'stem';
    stem.textContent = (qi + 1) + '. ' + q.stem;
    box.appendChild(stem);

    var explain = document.createElement('p');
    explain.className = 'explain';
    explain.textContent = q.explain || '';

    var order = q.options.map(function (_, i) { return i; });
    for (var i = order.length - 1; i > 0; i--) {
      var j = Math.floor(Math.random() * (i + 1));
      var t = order[i]; order[i] = order[j]; order[j] = t;
    }

    var buttons = [];
    order.forEach(function (optIndex) {
      var btn = document.createElement('button');
      btn.className = 'opt';
      btn.type = 'button';
      btn.textContent = q.options[optIndex];
      btn.addEventListener('click', function () {
        buttons.forEach(function (b) { b.disabled = true; });
        if (optIndex === q.answer) {
          btn.classList.add('correct');
        } else {
          btn.classList.add('wrong');
          buttons.forEach(function (b) {
            if (b.dataset.optIndex == q.answer) b.classList.add('correct');
          });
        }
        explain.classList.add('show');
      });
      btn.dataset.optIndex = optIndex;
      buttons.push(btn);
      box.appendChild(btn);
    });

    box.appendChild(explain);
    container.appendChild(box);
  });
}
