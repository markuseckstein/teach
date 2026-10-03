// Shared retrieval-practice quiz widget for the "Wand verputzen" lessons.
//
// Choice quiz markup:
// <div class="quiz" data-quiz="choice">
//   <p class="q">Frage...</p>
//   <div class="choices">
//     <button class="choice" data-correct="false">Option A</button>
//     <button class="choice" data-correct="true">Option B</button>
//   </div>
//   <div class="feedback" data-correct-text="..." data-incorrect-text="..."></div>
// </div>
//
// Text quiz markup:
// <div class="quiz" data-quiz="text" data-answer="richtig|richtige antwort">
//   <p class="q">Frage...</p>
//   <input class="answer-input" type="text" />
//   <button data-action="check">Prüfen</button>
//   <button class="secondary" data-action="reveal">Antwort zeigen</button>
//   <div class="feedback" data-correct-text="..." data-incorrect-text="..."></div>
// </div>

function normalize(str) {
  return str
    .trim()
    .toLowerCase()
    .replace(/ß/g, "ss")
    .replace(/[.,;:!?]/g, "")
    .replace(/\s+/g, " ");
}

function setFeedback(el, ok) {
  const text = ok ? el.dataset.correctText : el.dataset.incorrectText;
  el.textContent = text || (ok ? "Richtig." : "Nicht ganz.");
  el.classList.remove("correct", "incorrect");
  el.classList.add(ok ? "correct" : "incorrect");
}

function initChoiceQuiz(quiz) {
  const feedback = quiz.querySelector(".feedback");
  const choices = quiz.querySelectorAll(".choice");
  let answered = false;
  choices.forEach((btn) => {
    btn.addEventListener("click", () => {
      if (answered) return;
      answered = true;
      const ok = btn.dataset.correct === "true";
      choices.forEach((b) => {
        b.disabled = true;
        if (b === btn) {
          b.classList.add(ok ? "correct-pick" : "incorrect-pick");
        } else if (b.dataset.correct === "true") {
          b.classList.add("correct-pick");
        }
      });
      if (feedback) setFeedback(feedback, ok);
    });
  });
}

function initTextQuiz(quiz) {
  const input = quiz.querySelector(".answer-input");
  const feedback = quiz.querySelector(".feedback");
  const answers = (quiz.dataset.answer || "")
    .split("|")
    .map(normalize)
    .filter(Boolean);
  const checkBtn = quiz.querySelector('[data-action="check"]');
  const revealBtn = quiz.querySelector('[data-action="reveal"]');

  if (checkBtn) {
    checkBtn.addEventListener("click", () => {
      const ok = answers.includes(normalize(input.value || ""));
      if (feedback) setFeedback(feedback, ok);
    });
  }
  if (revealBtn) {
    revealBtn.addEventListener("click", () => {
      if (feedback) {
        feedback.textContent = "Antwort: " + (quiz.dataset.answer || "").split("|")[0];
        feedback.classList.remove("incorrect");
        feedback.classList.add("correct");
      }
    });
  }
  if (input) {
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        e.preventDefault();
        checkBtn && checkBtn.click();
      }
    });
  }
}

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll('.quiz[data-quiz="choice"]').forEach(initChoiceQuiz);
  document.querySelectorAll('.quiz[data-quiz="text"]').forEach(initTextQuiz);
});
