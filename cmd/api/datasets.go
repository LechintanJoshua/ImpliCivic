package main

import (
	"database/sql"
	"implicivic/internal/store"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// listDatasetsHandler (GET /api/datasets) afiseaza dataset-urile
// salvate in baza de date si filtreaza unde este cazul
func listDatasetsHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		category := strings.TrimSpace(q.Get("category"))
		search := strings.TrimSpace(q.Get("search"))
		limit := 20
		offset := 0

		if limitUrlString := q.Get("limit"); limitUrlString != "" {
			v, err := strconv.Atoi(limitUrlString)

			if err != nil || v <= 0 || v > 100 {
				writeError(w, http.StatusBadRequest, "invalid limit")
				return
			}

			limit = v
		}

		if offsetUrlString := q.Get("offset"); offsetUrlString != "" {
			v, err := strconv.Atoi(offsetUrlString)

			if err != nil || v < 0 {
				writeError(w, http.StatusBadRequest, "invalid offset")
				return
			}

			offset = v
		}

		datasets, err := store.ListDatasets(r.Context(), pool, category, search, limit, offset)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "db error")
			log.Printf("ListDatasets error: %v", err)
			return
		}

		writeJSON(w, http.StatusOK, datasets)
	}
}
