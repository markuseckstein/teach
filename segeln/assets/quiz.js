/*
 * Quiz-Komponente für den Segelkurs.
 *
 * Verwendung in einer Lektion:
 *   <div class="quiz" id="quiz"></div>
 *   <script src="../assets/quiz.js"></script>
 *   <script>
 *     renderQuiz(document.getElementById('quiz'), [
 *       {
 *         q: 'Frage?',
 *         options: ['A', 'B', 'C'],   // gleiche Wortzahl je Option anstreben
 *         correct: 1,                  // Index der richtigen Antwort
 *         why: 'Kurze Erklärung, warum das stimmt.'
 *       },
 *     ]);
 *   </script>
 */
function renderQuiz(container, questions) {
  var answered = 0;
  var score = 0;

  var heading = document.createElement('h2');
  heading.textContent = 'Selbsttest — erst aus dem Kopf beantworten!';
  container.appendChild(heading);

  questions.forEach(function (item, qi) {
    var qDiv = document.createElement('div');
    qDiv.className = 'quiz-q';

    var qP = document.createElement('p');
    qP.className = 'question';
    qP.textContent = (qi + 1) + '. ' + item.q;
    qDiv.appendChild(qP);

    var opts = document.createElement('div');
    opts.className = 'quiz-options';

    var feedback = document.createElement('div');
    feedback.className = 'quiz-feedback';

    item.options.forEach(function (opt, oi) {
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.textContent = opt;
      btn.addEventListener('click', function () {
        if (qDiv.dataset.done) return;
        qDiv.dataset.done = '1';
        answered++;
        opts.querySelectorAll('button').forEach(function (b) { b.disabled = true; });
        opts.children[item.correct].classList.add('correct');
        if (oi === item.correct) {
          score++;
          feedback.textContent = 'Richtig. ' + item.why;
        } else {
          btn.classList.add('incorrect');
          feedback.textContent = 'Leider nein. ' + item.why;
        }
        feedback.classList.add('visible');
        if (answered === questions.length) showScore();
      });
      opts.appendChild(btn);
    });

    qDiv.appendChild(opts);
    qDiv.appendChild(feedback);
    container.appendChild(qDiv);
  });

  var scoreDiv = document.createElement('div');
  scoreDiv.className = 'quiz-score';
  container.appendChild(scoreDiv);

  function showScore() {
    scoreDiv.textContent = 'Ergebnis: ' + score + ' von ' + questions.length +
      (score === questions.length
        ? ' — sitzt! Weiter zur nächsten Lektion.'
        : ' — die verfehlten Themen oben nochmal lesen und den Test morgen wiederholen.');
    scoreDiv.classList.add('visible');
  }
}
