/* Reusable quiz widget.
 *
 * Markup contract:
 *   <div class="quiz" data-answer="2" data-explain="Why the answer is right.">
 *     <div class="quiz-q">Question text?</div>
 *     <button class="quiz-opt">Option 0</button>
 *     <button class="quiz-opt">Option 1</button>
 *     <button class="quiz-opt">Option 2</button>
 *     <div class="quiz-feedback"></div>
 *   </div>
 *
 * data-answer is the zero-based index of the correct option.
 * Options are shuffled on load so position carries no signal.
 */
document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('.quiz').forEach(function (quiz) {
    var opts = Array.from(quiz.querySelectorAll('.quiz-opt'));
    var answerIdx = parseInt(quiz.dataset.answer, 10);
    var correctText = opts[answerIdx].textContent;
    var feedback = quiz.querySelector('.quiz-feedback');

    // Shuffle option order (Fisher–Yates on the DOM nodes).
    for (var i = opts.length - 1; i > 0; i--) {
      var j = Math.floor(Math.random() * (i + 1));
      quiz.insertBefore(opts[j], opts[i].nextSibling);
      var tmp = opts[i]; opts[i] = opts[j]; opts[j] = tmp;
    }

    quiz.addEventListener('click', function (e) {
      var btn = e.target.closest('.quiz-opt');
      if (!btn || quiz.dataset.done) return;
      quiz.dataset.done = '1';
      var right = btn.textContent === correctText;
      btn.classList.add(right ? 'correct' : 'wrong');
      if (!right) {
        opts.forEach(function (o) {
          if (o.textContent === correctText) o.classList.add('correct');
        });
      }
      feedback.textContent = (right ? '✓ Correct. ' : '✗ Not quite. ') + (quiz.dataset.explain || '');
      feedback.classList.add(right ? 'good' : 'bad');
    });
  });
});
