/* Bilanz-Simulator: T-Konten für mehrere Akteure.
 *
 * Verwendung:
 *   Bilanz.mount(element, {
 *     actors: [
 *       { id: "anna", name: "Anna (Haushalt)", aktiva: { "Bargeld": 0 }, passiva: { "Eigenkapital": 0 } },
 *       { id: "bankA", name: "Bank A", kind: "bank", aktiva: {...}, passiva: {...} },
 *     ],
 *     geldmenge: (s) => s.anna.aktiva["Guthaben bei Bank A"] + ...,   // optional
 *     geldmengeLabel: "Geldmenge (Bargeld + Guthaben der Nichtbanken)",
 *     // alternativ mehrere Zähler: counters: [{ label: "Geldmenge", fn: (s) => ... }, { label: "Zentralbankgeld", fn: ... }]
 *     actions: [
 *       { label: "Anna nimmt 1.000 € Kredit", run: (s, book) => { book("anna","aktiva","Guthaben",1000); ...; return "Erklärung"; },
 *         enabled: (s) => true }
 *     ]
 *   });
 *
 * book(actorId, "aktiva"|"passiva", posten, betrag) bucht einen Betrag (negativ = Abgang).
 * Posten mit Wert 0 werden ausgeblendet. Geänderte Zeilen werden kurz hervorgehoben.
 */
(function () {
  const css = `
  .bilanz-sim .accounts { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 1rem; margin-bottom: 1rem; }
  .bilanz-sim .tkonto { background: var(--bg); border: 1px solid var(--rule); border-radius: 3px; padding: .6rem .7rem; font-family: var(--sans); font-size: .78rem; }
  .bilanz-sim .tkonto h4 { margin: 0 0 .4rem; font-size: .8rem; font-weight: 600; text-align: center; }
  .bilanz-sim .tkonto.bank h4 { color: var(--accent); }
  .bilanz-sim .tgrid { display: grid; grid-template-columns: 1fr 1fr; border-top: 1.5px solid var(--fg); }
  .bilanz-sim .side { padding: .3rem .35rem; min-height: 3.2rem; }
  .bilanz-sim .side + .side { border-left: 1.5px solid var(--fg); }
  .bilanz-sim .side-h { font-size: .62rem; letter-spacing: .1em; text-transform: uppercase; color: var(--muted); margin-bottom: .2rem; }
  .bilanz-sim .row { display: flex; justify-content: space-between; gap: .4rem; padding: .08rem .15rem; border-radius: 2px; transition: background 1.2s; }
  .bilanz-sim .row .v { font-variant-numeric: tabular-nums; white-space: nowrap; }
  .bilanz-sim .row.up { background: var(--good-soft); transition: none; }
  .bilanz-sim .row.down { background: var(--bad-soft); transition: none; }
  .bilanz-sim .sum { display: grid; grid-template-columns: 1fr 1fr; border-top: 1px solid var(--rule); font-variant-numeric: tabular-nums; color: var(--muted); }
  .bilanz-sim .sum span { padding: .2rem .5rem; text-align: right; }
  .bilanz-sim .gms { display: flex; flex-wrap: wrap; gap: 0 2.2rem; }
  .bilanz-sim .gm { display: flex; align-items: baseline; gap: .8rem; flex-wrap: wrap; font-family: var(--sans); margin: .2rem 0 1rem; }
  .bilanz-sim .gm .gm-v { font-size: 1.6rem; font-variant-numeric: tabular-nums; font-family: var(--serif); }
  .bilanz-sim .gm .gm-d { font-size: .8rem; }
  .bilanz-sim .gm .gm-d.up { color: var(--good); } .bilanz-sim .gm .gm-d.down { color: var(--bad); }
  .bilanz-sim .gm .gm-l { font-size: .75rem; color: var(--muted); }
  .bilanz-sim .actions { display: flex; flex-wrap: wrap; gap: .5rem; }
  .bilanz-sim .log { margin-top: 1rem; font-size: .9rem; min-height: 2.6rem; padding: .6rem .8rem; background: var(--bg); border-left: 3px solid var(--accent); }
  `;
  const style = document.createElement("style");
  style.textContent = css;
  document.head.appendChild(style);

  const fmt = (n) => n.toLocaleString("de-DE") + " €";
  const clone = (o) => JSON.parse(JSON.stringify(o));

  function mount(el, cfg) {
    el.classList.add("bilanz-sim");
    const counters = cfg.counters || (cfg.geldmenge ? [{ label: cfg.geldmengeLabel || "Geldmenge", fn: cfg.geldmenge }] : []);
    const initial = {};
    cfg.actors.forEach((a) => (initial[a.id] = { aktiva: { ...a.aktiva }, passiva: { ...a.passiva } }));
    let state = clone(initial);
    let changed = {};
    let last = counters.map((c) => c.fn(state));
    let deltas = counters.map(() => 0);

    el.innerHTML = `
      <div class="gms">${counters.map((c) => `<div class="gm"><span class="gm-v"></span><span class="gm-d"></span><span class="gm-l">${c.label}</span></div>`).join("")}</div>
      <div class="accounts"></div>
      <div class="actions"></div>
      <div class="log">${cfg.intro || "Wähle eine Buchung und beobachte, welche Posten sich ändern."}</div>`;
    const accEl = el.querySelector(".accounts");
    const actEl = el.querySelector(".actions");
    const logEl = el.querySelector(".log");

    function book(id, side, posten, betrag) {
      const s = state[id][side];
      s[posten] = (s[posten] || 0) + betrag;
      changed[id + "|" + side + "|" + posten] = betrag > 0 ? "up" : "down";
    }

    function sideHtml(id, side, label) {
      const rows = Object.entries(state[id][side])
        .filter(([, v]) => v !== 0)
        .map(([k, v]) => `<div class="row ${changed[id + "|" + side + "|" + k] || ""}"><span>${k}</span><span class="v">${fmt(v)}</span></div>`)
        .join("");
      return `<div class="side"><div class="side-h">${label}</div>${rows || '<div class="row" style="color:var(--muted)">–</div>'}</div>`;
    }

    function render() {
      accEl.innerHTML = cfg.actors.map((a) => {
        const sA = Object.values(state[a.id].aktiva).reduce((x, y) => x + y, 0);
        const sP = Object.values(state[a.id].passiva).reduce((x, y) => x + y, 0);
        return `<div class="tkonto ${a.kind || ""}"><h4>${a.name}</h4>
          <div class="tgrid">${sideHtml(a.id, "aktiva", "Aktiva · was ich habe")}${sideHtml(a.id, "passiva", "Passiva · was ich schulde")}</div>
          <div class="sum"><span>${fmt(sA)}</span><span>${fmt(sP)}</span></div></div>`;
      }).join("");
      el.querySelectorAll(".gm").forEach((g, i) => {
        g.querySelector(".gm-v").textContent = fmt(counters[i].fn(state));
        const d = g.querySelector(".gm-d"), dv = deltas[i];
        d.textContent = dv === 0 ? "" : (dv > 0 ? "▲ +" : "▼ ") + fmt(dv);
        d.className = "gm-d " + (dv > 0 ? "up" : dv < 0 ? "down" : "");
      });
      actEl.querySelectorAll("button[data-i]").forEach((b) => {
        const a = cfg.actions[+b.dataset.i];
        b.disabled = a.enabled ? !a.enabled(state) : false;
      });
      const btns = [...actEl.querySelectorAll("button[data-i]")];
      if (btns.length && btns.every((b) => b.disabled) && !logEl.querySelector(".stuck"))
        logEl.insertAdjacentHTML("beforeend", ' <em class="stuck">Keine Buchung mehr möglich. Setz den Simulator zurück.</em>');
      // Hervorhebung nach dem Zeichnen ausblenden lassen
      requestAnimationFrame(() => setTimeout(() => accEl.querySelectorAll(".row.up,.row.down").forEach((r) => r.classList.remove("up", "down")), 900));
    }

    cfg.actions.forEach((a, i) => {
      const b = document.createElement("button");
      b.textContent = a.label;
      b.dataset.i = i;
      b.addEventListener("click", () => {
        changed = {};
        const msg = a.run(state, book);
        const now = counters.map((c) => c.fn(state));
        deltas = now.map((v, i) => v - last[i]);
        last = now;
        logEl.innerHTML = msg || "";
        render();
      });
      actEl.appendChild(b);
    });
    const reset = document.createElement("button");
    reset.textContent = "↺ Zurücksetzen";
    reset.addEventListener("click", () => {
      state = clone(initial); changed = {};
      last = counters.map((c) => c.fn(state)); deltas = counters.map(() => 0);
      logEl.innerHTML = cfg.intro || "Zurückgesetzt.";
      render();
    });
    actEl.appendChild(reset);
    render();
    return { get state() { return state; } };
  }

  window.Bilanz = { mount, fmt };
})();
