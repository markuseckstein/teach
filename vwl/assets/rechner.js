/* Schieberegler-Rechner.
 *
 *   Rechner.mount(el, {
 *     inputs: [{ id: "zins", label: "Zins", min: 0.5, max: 7, step: 0.05, value: 4, fmt: (v) => v + " %" }, ...],
 *     compute: (v) => [{ label: "Darlehen", value: "300.000 €", big: true }, { label: "...", value: "..." }],
 *     presets: [{ label: "Dez. 2021", set: { zins: 1.37 } }],   // optional: Knöpfe, die Werte setzen
 *   });
 */
(function () {
  const css = `
  .rechner .r-in { display: grid; gap: .7rem; margin-bottom: 1rem; }
  .rechner label { display: grid; grid-template-columns: 11rem 1fr 5.5rem; align-items: center; gap: .8rem; font-family: var(--sans); font-size: .8rem; }
  .rechner label output { text-align: right; font-variant-numeric: tabular-nums; font-weight: 600; }
  .rechner input[type=range] { width: 100%; accent-color: var(--accent); }
  .rechner .r-presets { display: flex; flex-wrap: wrap; gap: .5rem; margin-bottom: 1rem; }
  .rechner .r-out { display: flex; flex-wrap: wrap; gap: 1rem 2.2rem; border-top: 1px solid var(--rule); padding-top: .8rem; }
  .rechner .r-o { font-family: var(--sans); font-size: .75rem; color: var(--muted); }
  .rechner .r-o b { display: block; font-family: var(--serif); font-weight: 400; font-size: 1.15rem; color: var(--fg); font-variant-numeric: tabular-nums; }
  .rechner .r-o.big b { font-size: 1.7rem; }
  @media (max-width: 600px) { .rechner label { grid-template-columns: 1fr 5rem; } .rechner label input { grid-column: 1 / -1; grid-row: 2; } }
  `;
  const style = document.createElement("style");
  style.textContent = css;
  document.head.appendChild(style);

  function mount(el, cfg) {
    el.classList.add("rechner");
    el.innerHTML = `
      ${cfg.presets ? `<div class="r-presets">${cfg.presets.map((p, i) => `<button data-p="${i}">${p.label}</button>`).join("")}</div>` : ""}
      <div class="r-in">${cfg.inputs.map((inp) => `
        <label><span>${inp.label}</span>
          <input type="range" data-id="${inp.id}" min="${inp.min}" max="${inp.max}" step="${inp.step}" value="${inp.value}">
          <output data-for="${inp.id}"></output></label>`).join("")}</div>
      <div class="r-out"></div>`;
    const ranges = [...el.querySelectorAll("input[type=range]")];
    const outEl = el.querySelector(".r-out");
    function update() {
      const v = {};
      ranges.forEach((r) => (v[r.dataset.id] = parseFloat(r.value)));
      cfg.inputs.forEach((inp) => (el.querySelector(`output[data-for="${inp.id}"]`).textContent = inp.fmt ? inp.fmt(v[inp.id]) : v[inp.id]));
      outEl.innerHTML = cfg.compute(v).map((o) => `<div class="r-o ${o.big ? "big" : ""}">${o.label}<b>${o.value}</b></div>`).join("");
    }
    ranges.forEach((r) => r.addEventListener("input", update));
    el.querySelectorAll("button[data-p]").forEach((b) => b.addEventListener("click", () => {
      const p = cfg.presets[+b.dataset.p];
      Object.entries(p.set).forEach(([k, val]) => { const r = ranges.find((x) => x.dataset.id === k); if (r) r.value = val; });
      update();
    }));
    update();
  }

  const eur = (n) => Math.round(n).toLocaleString("de-DE") + " €";
  window.Rechner = { mount, eur };
})();
