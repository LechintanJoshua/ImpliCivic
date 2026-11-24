package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// writeJSON trimite un raspuns HTTP cu status code si body JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// healthHandler verifica daca serverul si baza de date raspund
func healthHandler(pool *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.PingContext(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable,
				map[string]string{"status": "error", "db": "unreachable"})

			return
		}

		writeJSON(w, http.StatusOK,
			map[string]string{"status": "ok", "db": "connected"})
	})
}
