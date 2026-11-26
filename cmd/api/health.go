package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

// healthHandler (GET /api/health/) verifica daca serverul si baza de date raspund
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
