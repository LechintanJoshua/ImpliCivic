package main

import (
	"context"
	"implicivic/internal/auth"
	"net/http"
	"strings"
)

type contextKey string

const claimsContextKey contextKey = "claims"

// corsMiddleware permite cererile cross-origin din frontend si opreste
// preflight-ul OPTIONS
func corsMiddleware(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requireAuth valideaza token-ul JWT din header-ul Authorization si
// ataseaza claims-urile in contextul cererii
func requireAuth(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr, ok := strings.CutPrefix(authHeader, "Bearer ")

		if !ok || tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		claims, err := auth.VerifyToken(tokenStr, jwtSecret)

		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// getClaims obtine claims-ul din contextul prmit
func getClaims(r *http.Request) (auth.Claims, bool) {
	claims, ok := r.Context().Value(claimsContextKey).(auth.Claims)
	return claims, ok
}
