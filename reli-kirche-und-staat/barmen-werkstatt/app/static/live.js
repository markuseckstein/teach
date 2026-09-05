// Progressive Enhancement: Ohne EventSource oder bei abgeschaltetem
// JavaScript bleibt die Seite bedienbar, nur ohne Live-Update — ein
// Neuladen von Hand zeigt immer den korrekten Stand.
if (window.EventSource) {
	var quelle = new EventSource("/gruppe/live");
	// Der Server schickt direkt nach Verbindungsaufbau ein erstes Signal
	// (sse.go), auch nach einem Wiederverbinden durch den Browser selbst —
	// das verhindert einen halb aktuellen Bildschirm nach einer
	// Funkloch-Unterbrechung. Bei einem ganz normalen Seitenaufruf ist dieses
	// erste Signal aber nur die Bestätigung "Verbindung steht", kein
	// tatsächlicher Zustandswechsel; würde es ein Neuladen auslösen, holte
	// sich die neu geladene Seite sofort ihr eigenes erstes Signal ab —
	// eine Endlosschleife, in der die Seite nie stillsteht. Deshalb zählt
	// nur das zweite und jedes weitere Signal als echte Änderung.
	var ersteNachricht = true;
	quelle.onmessage = function () {
		if (ersteNachricht) {
			ersteNachricht = false;
			return;
		}
		location.reload();
	};
}
