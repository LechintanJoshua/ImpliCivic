package main

import (
	"database/sql"
	"errors"
	"implicivic/internal/auth"
	"implicivic/internal/store"
	"net/http"
	"strings"
)

// registerHandler (POST api/auth/register) valideaza datele primite,
// hash-uieste parola si creeaza un cont nou de tip citizen
func registerHandler(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)

		if req.Username == "" || req.Email == "" ||
			len(req.Password) < 6 || len(req.Password) > 72 {
			writeError(w, http.StatusBadRequest, "username, email and password (6-72 chars) are required!")
			return
		}

		hash, err := auth.HashPassword(req.Password)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		user, err := store.CreateUser(r.Context(), pool, req.Username, req.Email, hash)

		if err != nil {
			if errors.Is(err, store.ErrDuplicateUser) {
				writeError(w, http.StatusConflict, "username or email already exists")
				return
			}

			writeError(w, http.StatusInternalServerError, "could not create user")
			return
		}

		writeJSON(w, http.StatusCreated, user)
	}
}

// loginHandler (POST /api/auth/login/) verifica email-ul si parola
// utilizatorului si returneaza un token JWT
func loginHandler(pool *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		user, err := store.GetUserByEmail(r.Context(), pool, req.Email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusUnauthorized, "invalid credentials")
				return
			}

			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role, jwtSecret)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
	}
}
