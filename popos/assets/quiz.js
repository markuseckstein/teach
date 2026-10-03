// Shared quiz widget logic for the PopOS/COSMIC teaching workspace.
// Two patterns supported via markup + data attributes:
//   .quiz[data-answer]        -> free-text input, .quiz-input + .quiz-check
//   .quiz-mc[data-answer]     -> multiple choice, .quiz-option[data-index] buttons

function normalize(s) {
  return s.trim().toLowerCase().replace(/\s+/g, " ");
}

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".quiz").forEach((quiz) => {
    const input = quiz.querySelector(".quiz-input");
    const button = quiz.querySelector(".quiz-check");
    const feedback = quiz.querySelector(".quiz-feedback");
    const answer = quiz.dataset.answer || "";

    const check = () => {
      const ok = normalize(input.value) === normalize(answer);
      feedback.textContent = ok
        ? "Richtig."
        : `Nicht ganz — richtig wäre: ${answer}`;
      feedback.className = "quiz-feedback " + (ok ? "correct" : "incorrect");
    };

    button.addEventListener("click", check);
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter") check();
    });
  });

  document.querySelectorAll(".quiz-mc").forEach((quiz) => {
    const correctIndex = quiz.dataset.answer;
    const feedback = quiz.querySelector(".quiz-feedback");
    const options = quiz.querySelectorAll(".quiz-option");

    options.forEach((opt) => {
      opt.addEventListener("click", () => {
        options.forEach((o) => o.classList.remove("correct", "incorrect"));
        const isCorrect = opt.dataset.index === correctIndex;
        opt.classList.add(isCorrect ? "correct" : "incorrect");
        if (!isCorrect) {
          const correctEl = quiz.querySelector(
            `.quiz-option[data-index="${correctIndex}"]`
          );
          if (correctEl) correctEl.classList.add("correct");
        }
        feedback.textContent = isCorrect ? "Richtig." : "Nicht ganz.";
        feedback.className = "quiz-feedback " + (isCorrect ? "correct" : "incorrect");
      });
    });
  });
});
