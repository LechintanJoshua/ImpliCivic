package main

import (
	"database/sql"
	"errors"
	"html"
	"implicivic/internal/store"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// sanitizeText este o functie helper anti-XSS
func sanitizeText(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}

// createIncidentHandler (POST api/incidents/) valideaza datele, le curata
// anti-XSS si creeaza o sesizare noua.
func createIncidentHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := getClaims(r)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req struct {
			Title       string   `json:"title"`
			Category    string   `json:"category"`
			Description string   `json:"description"`
			Address     string   `json:"address"`
			Latitude    *float64 `json:"latitude"`
			Longitude   *float64 `json:"longitude"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		req.Title = sanitizeText(req.Title)
		req.Category = sanitizeText(req.Category)
		req.Description = sanitizeText(req.Description)
		req.Address = sanitizeText(req.Address)

		if req.Title == "" || req.Category == "" || req.Description == "" ||
			req.Address == "" || utf8.RuneCountInString(req.Title) > 200 {
			writeError(w, http.StatusBadRequest, "empty fields or title too long")
			return
		}

		incident, err := store.CreateIncident(r.Context(), pool, claims.UserID,
			req.Title, req.Category, req.Description, req.Address,
			req.Latitude, req.Longitude)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create incident")
			return
		}

		writeJSON(w, http.StatusCreated, incident)
	}
}

// listIncidentsHandler (GET /api/incidents) listeaza incidentele salvate
// in baza de date si filtreaza, daca este cazul, dupa status si categorie
func listIncidentsHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		status := strings.TrimSpace(q.Get("status"))
		category := strings.TrimSpace(q.Get("category"))

		if status != "" && status != "INREGISTRAT" && status != "IN_LUCRU" &&
			status != "REZOLVAT" && status != "RESPINS" {
			writeError(w, http.StatusBadRequest, "invalid status")
			return
		}

		incidents, err := store.ListIncidents(r.Context(), pool, status, category, 0)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "db error")
			return
		}

		writeJSON(w, http.StatusOK, incidents)
	}
}

// getIncidentHandler (GET /api/incidents/{id) obtine o sesizare
// din baza de date
func getIncidentHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		incident, err := store.GetIncidentByID(r.Context(), pool, id)

		if err != nil {
			if errors.Is(err, store.ErrIncidentNotFound) {
				writeError(w, http.StatusNotFound, "no incident related to given id")
				return
			}

			writeError(w, http.StatusInternalServerError, "db error")
			return
		}

		writeJSON(w, http.StatusOK, incident)
	}
}

// myIncidentshandler (GET api/users/me/incidents) obtine toate sesizarile
// user-ului curent din baza de date
func myIncidentsHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := getClaims(r)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		incidents, err := store.ListIncidents(r.Context(), pool, "", "", claims.UserID)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "db error")
			return
		}

		writeJSON(w, http.StatusOK, incidents)
	}
}

// deleteIncidentHandler (DELETE /api/incidents/{id}) sterge o sesizare
// din baza de date
func deleteIncidentHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := getClaims(r)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		incident, err := store.GetIncidentByID(r.Context(), pool, id)

		if err != nil {
			if errors.Is(err, store.ErrIncidentNotFound) {
				writeError(w, http.StatusNotFound, "no incident related to given id")
				return
			}

			writeError(w, http.StatusInternalServerError, "db error")
			return
		}

		if claims.Role != "admin" &&
			(incident.UserId != claims.UserID || incident.Status != "INREGISTRAT") {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}

		if err := store.DeleteIncident(r.Context(), pool, id); err != nil {
			writeError(w, http.StatusInternalServerError, "db error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
