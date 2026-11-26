package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDuplicateUser = errors.New("username or email already exists")

// User este structura care reprezinta datele din tabelul users
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateUser adauga un user nou in baza de date, verifica de
// erori si returneaza obiectul User adaugat
func CreateUser(ctx context.Context, db *sql.DB,
	username, email, passwordHash string) (User, error) {
	query :=
		`INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, role, created_at`

	var u User
	err := db.QueryRowContext(ctx, query, username, email, passwordHash).Scan(&u.ID,
		&u.Username, &u.Email, &u.Role, &u.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrDuplicateUser
		}

		return User{}, err
	}

	return u, nil
}

// GetUserByEmail verifica daca email-ul este inregistrat in db
// si il aduce in memorie si returneaza obiectul user
func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (User, error) {
	query :=
		`SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE email = $1`

	var u User
	err := db.QueryRowContext(ctx, query, email).Scan(&u.ID,
		&u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if err != nil {
		return User{}, err
	}

	return u, nil
}
