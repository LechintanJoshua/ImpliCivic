package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrIncidentNotFound = errors.New("incident not found")

type Incident struct {
	ID          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Status      string    `json:"status"`
	AdminNotes  *string   `json:"admin_notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateIncidents creaza adauga in baza de date un incident nou si
// returneaza incidentul creat
func CreateIncident(ctx context.Context, db *sql.DB, userID int64,
	title, category, description, address string, lat, lng *float64) (Incident, error) {
	query := `
		INSERT INTO incidents (user_id, title, category, description, address, latitude, longitude)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, title, category, description, address, latitude, longitude,
		status, admin_notes, created_at, updated_at`

	var i Incident

	if err := db.QueryRowContext(ctx, query, userID, title, category,
		description, address, lat, lng).Scan(&i.ID, &i.UserId, &i.Title, &i.Category,
		&i.Description, &i.Address, &i.Latitude, &i.Longitude, &i.Status,
		&i.AdminNotes, &i.CreatedAt, &i.UpdatedAt); err != nil {
		return Incident{}, err
	}

	return i, nil
}

// GetIncidentById verifica daca id-ul exista, il aduce in memorie
// si returneaza obiectul Incident
func GetIncidentByID(ctx context.Context, db *sql.DB, id int64) (Incident, error) {
	query := `
	SELECT id, user_id, title, category, description, address, latitude, longitude, 
	status, admin_notes, created_at, updated_at
	FROM incidents
	WHERE id = $1`

	var i Incident

	err := db.QueryRowContext(ctx, query, id).Scan(&i.ID, &i.UserId, &i.Title,
		&i.Category, &i.Description, &i.Address, &i.Latitude, &i.Longitude,
		&i.Status, &i.AdminNotes, &i.CreatedAt, &i.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Incident{}, ErrIncidentNotFound
		}

		return Incident{}, err
	}

	return i, nil
}

// ListIncidents aduce din baza de date lista de incidente al unui user
// si le filtrazea, daca este cazul, dupa user, status si categorie
func ListIncidents(ctx context.Context, db *sql.DB, status,
	category string, userID int64) ([]Incident, error) {
	query := `
		SELECT id, user_id, title, category, description, address, 
		latitude, longitude, status, admin_notes, created_at, updated_at
		FROM incidents
		WHERE ($1 = '' OR status = $1)
		AND ($2 = '' OR category = $2)
		AND ($3 = 0 OR user_id = $3)
		ORDER BY created_at DESC`

	rows, err := db.QueryContext(ctx, query, status, category, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	incidents := make([]Incident, 0)

	for rows.Next() {
		var i Incident

		if err := rows.Scan(&i.ID, &i.UserId, &i.Title,
			&i.Category, &i.Description, &i.Address, &i.Latitude, &i.Longitude,
			&i.Status, &i.AdminNotes, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		incidents = append(incidents, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

// DeleteIncidents sterge un incident din baza de date dupa id
func DeleteIncident(ctx context.Context, db *sql.DB, id int64) error {
	res, err := db.ExecContext(ctx, "DELETE FROM INCIDENTS WHERE id = $1", id)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrIncidentNotFound
	}

	return nil
}
