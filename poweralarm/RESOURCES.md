# PowerAlarm im Feuerwehrhaus — Resources

## Knowledge

- [Anleitung PowerAlarm (PDF, Version 2.21 vom 15.06.2026) — FITT GmbH](https://www.fitt-gmbh.de/wp-content/uploads/hilfe_poweralarm.pdf)
  **Die** Primärquelle: 73 Seiten Herstellerhandbuch. Kapitel 4 = PA-Monitor
  (Windows), Kapitel 5 = PA Web-Monitor, Punkt 2.18.d = Zugangsdaten und URL für
  den Web-Monitor, Kapitel 3 = PAWinS. Bei jeder Frage zu Funktionsumfang zuerst
  hier nachsehen — die Marketing-Seiten sagen dazu praktisch nichts.

- [Produktseite PowerAlarm — FITT GmbH](https://www.fitt-gmbh.de/produkte/poweralarm/)
  Herstellerübersicht der Module (Alarmvisualisierung, PowerWarn, PowerSilent …).
  Use for: Was es überhaupt gibt und wen man anruft. Kein technisches Detail.

- [poweralarm.de](https://www.poweralarm.de/)
  Das Portal selbst, plus News-Bereich. Use for: Login, Änderungen am System.

- [Alarm Dispatcher: Konfiguration PowerAlarm (API-Key anlegen)](https://support.alarm-dispatcher.de/help/de-de/19-poweralarm/69-fur-den-endkunden-konfiguration-poweralarm)
  Fremddokumentation, zeigt aber den Weg zum API-Schlüssel im Portal
  (Administration → Verwaltung → API-Schlüssel) mit Screenshots.
  Use for: Orientierung im Portal, wenn das FITT-Handbuch zu knapp ist.

- [FITTcom PowerAlarm — HiOrg-Server Wiki](https://wiki.hiorg-server.de/admin/fitt)
  Drittanbieter-Doku zur PowerAlarm-Schnittstelle. Use for: zweite Meinung zur
  API, wenn das Handbuch unklar ist.

## Wisdom (Communities)

- [Deutsches Raspberry Pi Forum — Thread "Feuerwehr Display"](https://forum-raspberrypi.de/forum/thread/34064-feuerwehr-display/)
  Feuerwehrleute, die genau dieses Problem gelöst haben: Pi am Monitor,
  Alarmanzeige, Bildschirm ein/aus. Use for: Kiosk-Praxis, Bildschirmsteuerung,
  Fehler, die andere schon gemacht haben.

- [Deutsches Raspberry Pi Forum — Thread "DIVERA Monitor-App (Feuerwehr)"](https://forum-raspberrypi.de/forum/thread/61024-divera-monitor-app-feuerwehr/)
  Gleiches Muster mit anderer Plattform. Use for: übertragbare Kiosk- und
  Autostart-Lösungen.

- [Raspberry Pi Forums — "Alarmanzeige mit Steuerung für FFW"](https://forums.raspberrypi.com/viewtopic.php?t=288638)
  Use for: Bildschirm per HDMI-CEC oder schaltbarer Steckdose wecken.

- [Wachalarm-Kiosk — Robert-112 auf GitHub](https://github.com/Robert-112/Wachalarm-Kiosk)
  Fertiges Raspberry-Pi-Image, das eine Alarm-Webseite im Vollbild anzeigt,
  inklusive Sound und Standby. Andere Plattform, aber genau unser Bauteil.
  Use for: abschauen, wie man es sauber baut — oder direkt als Basis.

- FITT GmbH Support — service@fitt-gmbh.de, 0911 – 36 66 928
  Kein Community, aber die einzige Quelle für Lizenz- und Funktionsfragen.
  Use for: Freischaltungen, "kann der Web-Monitor X?".

## Gaps

- **Funktionsumfang des Web-Monitors ist unterdokumentiert.** Das Handbuch
  beschreibt in Kapitel 5 nur drei Anzeigefelder und erwähnt "Einstellungen
  (Optionen)", ohne sie aufzulisten. Ob Karte, Alarmfax, Uhrzeit, News oder
  **Alarmton** vorhanden sind, ist offen. → Muss bei FITT erfragt werden.
- **Keine öffentliche Doku zur Web-Monitor-URL.** Sie wird laut Handbuch nur im
  Portal angezeigt. Nicht erratbar, nirgends im Netz referenziert.
- **Kein Erfahrungsbericht zu PA Web-Monitor im Kiosk-Modus gefunden.** Für
  Divera und BlaulichtSMS gibt es Anleitungen, für PowerAlarm nicht. Wir sind
  hier Pioniere — Ergebnisse also selbst dokumentieren.
