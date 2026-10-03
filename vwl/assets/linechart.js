/* Liniendiagramm (SVG) mit Fadenkreuz-Tooltip, Legende, Endwert-Labels und Tabellenansicht.
 *
 *   LineChart.mount(el, {
 *     labels: ["2021-11", ...],             // x-Achse (gleich lang wie jede Serie)
 *     series: [{ name: "Bauzins", values: [1.36, ...], color: 1, step: false }, ...],  // color: Slot 1|2
 *     unit: " %", yMin: -1, yMax: 5, yStep: 1,
 *     xTick: (label, i) => label.endsWith("-01") ? label.slice(0, 4) : null,  // welche x-Ticks beschriftet werden
 *     fmtX: (label) => label,               // Format im Tooltip / in der Tabelle
 *     annotations: [{ at: "2022-07", text: "1. Zinserhöhung" }],
 *   });
 *
 * Farben (validiert mit dataviz/validate_palette.js, hell + dunkel): Slot 1 Ziegelrot, Slot 2 Blau.
 */
(function () {
  const css = `
  :root { --series-1: #b5452c; --series-2: #1f6fb8; --grid: #e6e3d6; }
  @media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) { --series-1: #d9694a; --series-2: #3d8fd6; --grid: #34322c; } }
  :root[data-theme="dark"] { --series-1: #d9694a; --series-2: #3d8fd6; --grid: #34322c; }
  .lc { position: relative; font-family: var(--sans); }
  .lc-legend { display: flex; gap: 1.2rem; flex-wrap: wrap; font-size: .78rem; color: var(--fg); margin-bottom: .4rem; }
  .lc-legend span::before { content: ""; display: inline-block; width: 14px; height: 2px; vertical-align: middle; margin-right: .4rem; background: var(--c); }
  .lc svg { width: 100%; height: auto; display: block; overflow: visible; }
  .lc .axis text { font-size: 11px; fill: var(--muted); }
  .lc .grid line { stroke: var(--grid); stroke-width: 1; }
  .lc .zero { stroke: var(--muted); stroke-width: 1; }
  .lc .line { fill: none; stroke-width: 2; stroke-linejoin: round; stroke-linecap: round; }
  .lc .endlabel { font-size: 11px; fill: var(--fg); }
  .lc .annot line { stroke: var(--muted); stroke-width: 1; }
  .lc .annot text { font-size: 10.5px; fill: var(--muted); }
  .lc .cross { stroke: var(--muted); stroke-width: 1; }
  .lc .dot { stroke: var(--panel); stroke-width: 2; }
  .lc-tip { position: absolute; pointer-events: none; background: var(--bg); border: 1px solid var(--rule); border-radius: 3px; padding: .35rem .55rem; font-size: .75rem; line-height: 1.4; white-space: nowrap; box-shadow: 0 2px 8px rgba(0,0,0,.08); display: none; }
  .lc-tip i { display: inline-block; width: 10px; height: 2px; vertical-align: middle; margin-right: .35rem; }
  .lc details { margin-top: .6rem; font-size: .8rem; }
  .lc table { border-collapse: collapse; margin-top: .4rem; font-variant-numeric: tabular-nums; }
  .lc td, .lc th { padding: .15rem .6rem; border-bottom: 1px solid var(--rule); text-align: right; }
  .lc th:first-child, .lc td:first-child { text-align: left; }
  `;
  const style = document.createElement("style");
  style.textContent = css;
  document.head.appendChild(style);

  const fmtN = (v) => v.toLocaleString("de-DE", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

  function mount(el, cfg) {
    const W = 680, H = 300, m = { t: 14, r: 56, b: 28, l: 36 };
    const iw = W - m.l - m.r, ih = H - m.t - m.b;
    const n = cfg.labels.length;
    const x = (i) => m.l + (i / (n - 1)) * iw;
    const y = (v) => m.t + (1 - (v - cfg.yMin) / (cfg.yMax - cfg.yMin)) * ih;
    const unit = cfg.unit || "";
    const fmtX = cfg.fmtX || ((l) => l);
    const col = (s) => `var(--series-${s.color || 1})`;

    let g = "";
    for (let v = cfg.yMin; v <= cfg.yMax + 1e-9; v += cfg.yStep) {
      g += `<g class="grid"><line x1="${m.l}" x2="${m.l + iw}" y1="${y(v)}" y2="${y(v)}"/></g>`;
      g += `<g class="axis"><text x="${m.l - 6}" y="${y(v) + 4}" text-anchor="end">${v.toLocaleString("de-DE")}</text></g>`;
    }
    if (cfg.yMin < 0) g += `<line class="zero" x1="${m.l}" x2="${m.l + iw}" y1="${y(0)}" y2="${y(0)}"/>`;
    cfg.labels.forEach((l, i) => {
      const t = cfg.xTick ? cfg.xTick(l, i) : null;
      if (t) g += `<g class="axis"><text x="${x(i)}" y="${H - 8}" text-anchor="middle">${t}</text></g>`;
    });
    (cfg.annotations || []).forEach((a) => {
      const i = cfg.labels.indexOf(a.at);
      if (i < 0) return;
      g += `<g class="annot"><line x1="${x(i)}" x2="${x(i)}" y1="${m.t}" y2="${m.t + ih}"/><text x="${x(i) + 4}" y="${m.t + 10}">${a.text}</text></g>`;
    });

    cfg.series.forEach((s) => {
      let d = "", prev = null;
      s.values.forEach((v, i) => {
        if (v == null) { prev = null; return; }
        if (prev == null) d += `M${x(i)},${y(v)}`;
        else d += s.step ? `H${x(i)}V${y(v)}` : `L${x(i)},${y(v)}`;
        prev = v;
      });
      g += `<path class="line" d="${d}" style="stroke:${col(s)}"/>`;
      let li = s.values.length - 1;
      while (li >= 0 && s.values[li] == null) li--;
      g += `<circle class="dot" cx="${x(li)}" cy="${y(s.values[li])}" r="4" style="fill:${col(s)}"/>`;
      g += `<text class="endlabel" x="${x(li) + 8}" y="${y(s.values[li]) + 4}">${fmtN(s.values[li])}${unit}</text>`;
    });

    el.classList.add("lc");
    el.innerHTML = `
      <div class="lc-legend">${cfg.series.map((s) => `<span style="--c:${col(s)}">${s.name}</span>`).join("")}</div>
      <svg viewBox="0 0 ${W} ${H}" role="img" aria-label="${cfg.series.map((s) => s.name).join(" und ")}">
        ${g}
        <line class="cross" y1="${m.t}" y2="${m.t + ih}" style="display:none"/>
        ${cfg.series.map((s) => `<circle class="dot hov" r="4" style="fill:${col(s)};display:none"/>`).join("")}
        <rect x="${m.l}" y="0" width="${iw}" height="${H}" fill="transparent" class="hit"/>
      </svg>
      <div class="lc-tip"></div>
      <details><summary>Als Tabelle anzeigen</summary>
        <table><tr><th>Zeitpunkt</th>${cfg.series.map((s) => `<th>${s.name}</th>`).join("")}</tr>
        ${cfg.labels.map((l, i) => `<tr><td>${fmtX(l)}</td>${cfg.series.map((s) => `<td>${s.values[i] == null ? "–" : fmtN(s.values[i]) + unit}</td>`).join("")}</tr>`).join("")}
        </table></details>`;

    const svg = el.querySelector("svg"), tip = el.querySelector(".lc-tip");
    const cross = el.querySelector(".cross"), hov = [...el.querySelectorAll(".hov")];
    function show(evt) {
      const r = svg.getBoundingClientRect();
      const sx = ((evt.clientX - r.left) / r.width) * W;
      const i = Math.max(0, Math.min(n - 1, Math.round(((sx - m.l) / iw) * (n - 1))));
      cross.setAttribute("x1", x(i)); cross.setAttribute("x2", x(i)); cross.style.display = "";
      cfg.series.forEach((s, k) => {
        const v = s.values[i];
        if (v == null) { hov[k].style.display = "none"; return; }
        hov[k].setAttribute("cx", x(i)); hov[k].setAttribute("cy", y(v)); hov[k].style.display = "";
      });
      tip.innerHTML = `<strong>${fmtX(cfg.labels[i])}</strong><br>` + cfg.series.map((s) =>
        `<i style="background:${col(s)}"></i>${s.name}: ${s.values[i] == null ? "–" : fmtN(s.values[i]) + unit}`).join("<br>");
      tip.style.display = "block";
      const px = (x(i) / W) * r.width;
      tip.style.left = Math.min(px + 12, r.width - tip.offsetWidth) + "px";
      tip.style.top = "2rem";
    }
    function hide() { tip.style.display = "none"; cross.style.display = "none"; hov.forEach((h) => (h.style.display = "none")); }
    const hit = el.querySelector(".hit");
    hit.addEventListener("pointermove", show);
    hit.addEventListener("pointerleave", hide);
  }

  window.LineChart = { mount };
})();
