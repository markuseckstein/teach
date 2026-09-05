// Progressive Enhancement: Ohne EventSource oder bei abgeschaltetem
// JavaScript bleibt die Seite bedienbar, nur ohne Live-Update — ein
// Neuladen von Hand zeigt immer den korrekten Stand.
if (window.EventSource) {
	var quelle = new EventSource("/gruppe/live");
	quelle.onmessage = function () {
		location.reload();
	};
}
