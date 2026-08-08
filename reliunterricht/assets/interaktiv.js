/* ============================================================
   Interaktive Komponenten der Einheit „Todesstrafe"
   Alles rein clientseitig, offline-fähig (Beamer-Einsatz).

   Komponenten (per data-Attribut aktiviert):
   1. Aufdecken:   <button class="aufdecken" data-ziel="idDerAntwort">
                   Blendet das Element mit der id ein (Schätzquiz).
   2. Tally:       <div class="tally" data-tally="eindeutiger-schluessel">
                   Strichliste für Handzeichen-Abstimmungen. Stand wird
                   in localStorage gespeichert, damit eine spätere Stunde
                   (gleicher Browser!) den Vorher-Wert wieder zeigen kann.
   ============================================================ */

document.addEventListener('DOMContentLoaded', function () {

  /* ---------- 1. Aufdecken-Buttons ---------- */
  document.querySelectorAll('.aufdecken[data-ziel]').forEach(function (knopf) {
    knopf.addEventListener('click', function () {
      var ziel = document.getElementById(knopf.dataset.ziel);
      if (!ziel) return;
      ziel.classList.add('sichtbar');
      knopf.style.display = 'none';
    });
  });

  /* ---------- 2. Tally (Strichlisten-Abstimmung) ---------- */
  document.querySelectorAll('.tally[data-tally]').forEach(function (tally) {
    var schluessel = 'tally:' + tally.dataset.tally;
    var spalten = tally.querySelectorAll('.spalte');
    var stand;
    try { stand = JSON.parse(localStorage.getItem(schluessel)) || {}; }
    catch (e) { stand = {}; }

    function speichern() {
      try { localStorage.setItem(schluessel, JSON.stringify(stand)); } catch (e) {}
    }

    function zeichnen() {
      var summe = 0;
      spalten.forEach(function (s) { summe += stand[s.dataset.wert] || 0; });
      spalten.forEach(function (s) {
        var n = stand[s.dataset.wert] || 0;
        s.querySelector('.zahl').textContent = n;
        var balken = s.querySelector('.balken');
        if (balken) balken.style.width = summe ? (100 * n / summe) + '%' : '0';
      });
    }

    spalten.forEach(function (spalte) {
      var wert = spalte.dataset.wert;
      if (!(wert in stand)) stand[wert] = 0;
      spalte.querySelector('[data-aktion="plus"]').addEventListener('click', function () {
        stand[wert]++; speichern(); zeichnen();
      });
      var minus = spalte.querySelector('[data-aktion="minus"]');
      if (minus) minus.addEventListener('click', function () {
        if (stand[wert] > 0) stand[wert]--;
        speichern(); zeichnen();
      });
    });

    var reset = tally.parentElement.querySelector('[data-aktion="reset"][data-fuer="' + tally.dataset.tally + '"]');
    if (reset) reset.addEventListener('click', function () {
      if (!confirm('Diese Abstimmung wirklich auf null zurücksetzen?')) return;
      Object.keys(stand).forEach(function (k) { stand[k] = 0; });
      speichern(); zeichnen();
    });

    zeichnen();
  });

  /* ---------- 3. Umrechner (z.B. „Die Welt in unserer Klasse") ----------
     <div class="umrechner" data-umrechner> mit <input data-rolle="gesamt">;
     jedes Element mit data-prozent zeigt den gerundeten Anteil der Gesamtzahl. */
  document.querySelectorAll('[data-umrechner]').forEach(function (block) {
    var eingabe = block.querySelector('input[data-rolle="gesamt"]');
    if (!eingabe) return;
    function rechnen() {
      var gesamt = parseInt(eingabe.value, 10) || 0;
      block.querySelectorAll('[data-prozent]').forEach(function (ziel) {
        ziel.textContent = Math.round(gesamt * parseFloat(ziel.dataset.prozent) / 100);
      });
    }
    eingabe.addEventListener('input', rechnen);
    rechnen();
  });

  /* ---------- 4. Timer (Beamer-Countdown, z.B. für Debattenphasen) ----------
     <div class="timer" data-timer>
       <div class="timer-anzeige">0:00</div>
       <button data-minuten="1">1 Min</button> …
       <button data-timer-aktion="startpause">Start/Pause</button>
       <button data-timer-aktion="reset">Zurücksetzen</button>
     </div> */
  document.querySelectorAll('[data-timer]').forEach(function (timer) {
    var anzeige = timer.querySelector('.timer-anzeige');
    var gewaehlt = 0, rest = 0, intervall = null;

    function stoppen() {
      if (intervall) { clearInterval(intervall); intervall = null; }
    }
    function zeigen() {
      var m = Math.floor(rest / 60), s = rest % 60;
      anzeige.textContent = m + ':' + (s < 10 ? '0' : '') + s;
      timer.classList.toggle('abgelaufen', gewaehlt > 0 && rest === 0);
    }
    function tick() {
      if (rest > 0) rest--;
      if (rest === 0) stoppen();
      zeigen();
    }

    timer.querySelectorAll('button[data-minuten]').forEach(function (knopf) {
      knopf.addEventListener('click', function () {
        stoppen();
        gewaehlt = Math.round(parseFloat(knopf.dataset.minuten) * 60);
        rest = gewaehlt;
        zeigen();
      });
    });
    var startpause = timer.querySelector('[data-timer-aktion="startpause"]');
    if (startpause) startpause.addEventListener('click', function () {
      if (intervall) { stoppen(); }
      else if (rest > 0) { intervall = setInterval(tick, 1000); }
    });
    var reset = timer.querySelector('[data-timer-aktion="reset"]');
    if (reset) reset.addEventListener('click', function () {
      stoppen();
      rest = gewaehlt;
      zeigen();
    });
    zeigen();
  });

  /* ---------- 5. Tally-Vergleich (z.B. Positionierung vorher/nachher) ----------
     Zellen mit data-quelle="tally-schluessel" und data-wert="spalte" zeigen den
     gespeicherten Stand; ein Knopf [data-vergleich-aktion="aktualisieren"]
     liest neu ein (z.B. nachdem die Nachher-Abstimmung gezählt wurde). */
  function tallyLesen(schluessel) {
    try { return JSON.parse(localStorage.getItem('tally:' + schluessel)) || {}; }
    catch (e) { return {}; }
  }
  function vergleichFuellen() {
    document.querySelectorAll('[data-quelle][data-wert]').forEach(function (zelle) {
      zelle.textContent = tallyLesen(zelle.dataset.quelle)[zelle.dataset.wert] || 0;
    });
  }
  document.querySelectorAll('[data-vergleich-aktion="aktualisieren"]').forEach(function (knopf) {
    knopf.addEventListener('click', vergleichFuellen);
  });
  vergleichFuellen();

  /* Hilfsfunktion für spätere Stunden: gespeicherten Tally-Stand lesen */
  window.tallyStand = function (schluessel) {
    try { return JSON.parse(localStorage.getItem('tally:' + schluessel)) || {}; }
    catch (e) { return {}; }
  };
});
