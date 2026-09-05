package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// geraetCookieName ist der Name des Cookies, über den ein Tablet
// wiedererkannt wird. Der Wert ist ein zufälliges, nicht rückführbares
// Token — kein Name, keine sonstige Identität.
const geraetCookieName = "geraet"

// geraetCookieMaxAge deckt die mehrwöchige Klammer der Unterrichtseinheit ab
// (Placemat aus DS 4, Ausstellung aus DS 5 kommen erst in DS 6 zurück).
const geraetCookieMaxAge = 60 * 24 * 60 * 60 // 60 Tage

func neuesGeraetToken() (string, error) {
	rohbytes := make([]byte, 16)
	if _, err := rand.Read(rohbytes); err != nil {
		return "", fmt.Errorf("zufallstoken erzeugen: %w", err)
	}
	return hex.EncodeToString(rohbytes), nil
}

func geraetCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     geraetCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   geraetCookieMaxAge,
	}
}

// aktuellesGeraet löst das Geräte-Cookie einer Anfrage auf. ok ist false,
// wenn kein Cookie mitgeschickt wurde oder das Token keinem Gerät (mehr)
// zugeordnet ist — dann muss das Gerät erneut beitreten.
func aktuellesGeraet(r *http.Request, database *sql.DB) (geraetID, gruppeID int64, ok bool) {
	cookie, err := r.Cookie(geraetCookieName)
	if err != nil || cookie.Value == "" {
		return 0, 0, false
	}

	err = database.QueryRow(`SELECT id, gruppe_id FROM geraet WHERE token = ?`, cookie.Value).
		Scan(&geraetID, &gruppeID)
	if err != nil {
		return 0, 0, false
	}

	// Wiedereinstieg aktualisiert "zuletzt gesehen"; ein Fehler dabei ist
	// kein Grund, die eigentliche Anfrage scheitern zu lassen.
	_, err = database.Exec(`UPDATE geraet SET zuletzt_gesehen = ? WHERE id = ?`, jetzt(), geraetID)
	if err != nil {
		log.Printf("zuletzt_gesehen aktualisieren: %v", err)
	}

	return geraetID, gruppeID, true
}

func jetzt() string {
	return time.Now().UTC().Format(time.RFC3339)
}

type beitretenAnsicht struct {
	Code   string
	Fehler string
}

var beitretenTmpl = template.Must(template.ParseFS(templatesFS, "templates/beitreten.html"))

// handleBeitreten nimmt den Kurscode aus der Query entgegen (auch vom
// gescannten QR-Code) und leitet bei Erfolg zur Gruppenwahl weiter. Ohne
// (gültigen) Code zeigt es das Eingabeformular.
func handleBeitreten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code")))
		if code == "" {
			renderTemplate(w, beitretenTmpl, beitretenAnsicht{})
			return
		}

		var kursID int64
		err := database.QueryRow(`SELECT id FROM kurs WHERE beitrittscode = ?`, code).Scan(&kursID)
		if errors.Is(err, sql.ErrNoRows) {
			renderTemplate(w, beitretenTmpl, beitretenAnsicht{Code: code, Fehler: "Diesen Kurscode gibt es nicht."})
			return
		}
		if err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("kurscode prüfen: %v", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/kurse/%d/gruppe-waehlen", kursID), http.StatusSeeOther)
	}
}

type gruppeWaehlenAnsicht struct {
	KursID  int64
	Gruppen []int
}

var gruppeWaehlenTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-waehlen.html"))

// handleGruppeWaehlen zeigt die Gruppen eines Kurses zur Auswahl.
func handleGruppeWaehlen(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		rows, err := database.Query(`SELECT nummer FROM gruppe WHERE kurs_id = ? ORDER BY nummer`, kursID)
		if err != nil {
			http.Error(w, "Gruppen konnten nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppen laden: %v", err)
			return
		}
		defer rows.Close()

		var gruppen []int
		for rows.Next() {
			var nummer int
			if err := rows.Scan(&nummer); err != nil {
				http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
				log.Printf("gruppe lesen: %v", err)
				return
			}
			gruppen = append(gruppen, nummer)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Gruppen konnten nicht gelesen werden", http.StatusInternalServerError)
			log.Printf("gruppen lesen: %v", err)
			return
		}
		if len(gruppen) == 0 {
			http.NotFound(w, r)
			return
		}

		renderTemplate(w, gruppeWaehlenTmpl, gruppeWaehlenAnsicht{KursID: kursID, Gruppen: gruppen})
	}
}

// handleGruppeBeitreten registriert das anfragende Gerät bei der gewählten
// Gruppe, setzt das Geräte-Cookie und legt — falls noch keine vorhanden
// ist — die zugehörige These an, damit "Gruppe" einen aktuellen Schritt hat.
func handleGruppeBeitreten(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		nummer, err := strconv.Atoi(r.PathValue("nummer"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var gruppeID int64
		err = database.QueryRow(`SELECT id FROM gruppe WHERE kurs_id = ? AND nummer = ?`, kursID, nummer).Scan(&gruppeID)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("gruppe laden: %v", err)
			return
		}

		token, err := neuesGeraetToken()
		if err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("token erzeugen: %v", err)
			return
		}

		res, err := database.Exec(
			`INSERT INTO geraet (gruppe_id, token, zuletzt_gesehen) VALUES (?, ?, ?)`,
			gruppeID, token, jetzt(),
		)
		if err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("geraet anlegen: %v", err)
			return
		}
		geraetID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("geraet-id lesen: %v", err)
			return
		}

		if _, err := database.Exec(
			`INSERT INTO these (gruppe_id) SELECT ? WHERE NOT EXISTS (SELECT 1 FROM these WHERE gruppe_id = ?)`,
			gruppeID, gruppeID,
		); err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("these anlegen: %v", err)
			return
		}

		// Das erste Gerät einer Gruppe hält automatisch das Schreibrecht. Die
		// Bedingung "schreibrecht_geraet IS NULL" im selben UPDATE macht das
		// Zuweisen atomar — kein separates SELECT, das mit einem zweiten,
		// gleichzeitig beitretenden Gerät wettlaufen könnte.
		if _, err := database.Exec(
			`UPDATE gruppe SET schreibrecht_geraet = ?, schreibrecht_seit = ? WHERE id = ? AND schreibrecht_geraet IS NULL`,
			geraetID, jetzt(), gruppeID,
		); err != nil {
			http.Error(w, "Beitritt fehlgeschlagen", http.StatusInternalServerError)
			log.Printf("schreibrecht zuweisen: %v", err)
			return
		}

		http.SetCookie(w, geraetCookie(token))
		http.Redirect(w, r, "/gruppe", http.StatusSeeOther)
	}
}

type gruppeStatusAnsicht struct {
	Nummer          int
	Schritt         string
	HatSchreibrecht bool
}

var gruppeStatusTmpl = template.Must(template.ParseFS(templatesFS, "templates/gruppe-status.html"))

// handleGruppeStatus ist die Landeseite eines Geräts: Gruppennummer,
// aktueller Schritt und ob dieses Gerät gerade das Schreibrecht hält —
// gelesen allein über das Geräte-Cookie. So landet ein Gerät nach einem
// Verbindungsabriss ohne erneutes Beitreten wieder hier.
func handleGruppeStatus(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geraetID, gruppeID, ok := aktuellesGeraet(r, database)
		if !ok {
			http.Redirect(w, r, "/beitreten", http.StatusSeeOther)
			return
		}

		var nummer int
		if err := database.QueryRow(`SELECT nummer FROM gruppe WHERE id = ?`, gruppeID).Scan(&nummer); err != nil {
			http.Error(w, "Gruppe konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("gruppe laden: %v", err)
			return
		}

		var schritt string
		if err := database.QueryRow(`SELECT schritt FROM these WHERE gruppe_id = ?`, gruppeID).Scan(&schritt); err != nil {
			http.Error(w, "These konnte nicht geladen werden", http.StatusInternalServerError)
			log.Printf("these laden: %v", err)
			return
		}

		hatSchreibrecht, err := geraetHatSchreibrecht(database, gruppeID, geraetID)
		if err != nil {
			http.Error(w, "Schreibrecht konnte nicht geprüft werden", http.StatusInternalServerError)
			log.Printf("schreibrecht prüfen: %v", err)
			return
		}

		renderTemplate(w, gruppeStatusTmpl, gruppeStatusAnsicht{
			Nummer:          nummer,
			Schritt:         schritt,
			HatSchreibrecht: hatSchreibrecht,
		})
	}
}

// handleKursQR liefert einen QR-Code, der auf den Beitritts-Link dieses
// Kurses zeigt — für den Beamer, damit die Klasse mit der Tabletkamera
// beitreten kann, statt den Code abzutippen.
func handleKursQR(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kursID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var beitrittscode string
		err = database.QueryRow(`SELECT beitrittscode FROM kurs WHERE id = ?`, kursID).Scan(&beitrittscode)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "QR-Code konnte nicht erzeugt werden", http.StatusInternalServerError)
			log.Printf("kurs laden: %v", err)
			return
		}

		// X-Forwarded-Proto zuerst: Vorgang 0018 setzt später Caddy als
		// TLS-Terminierung davor, dann kommt die Anfrage hier intern als
		// Klartext-HTTP an, auch wenn die Klasse https sieht.
		schema := r.Header.Get("X-Forwarded-Proto")
		if schema == "" {
			schema = "http"
			if r.TLS != nil {
				schema = "https"
			}
		}
		ziel := fmt.Sprintf("%s://%s/beitreten?code=%s", schema, r.Host, beitrittscode)

		png, err := qrcode.Encode(ziel, qrcode.Medium, 320)
		if err != nil {
			http.Error(w, "QR-Code konnte nicht erzeugt werden", http.StatusInternalServerError)
			log.Printf("qr-code erzeugen: %v", err)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(png)
	}
}
