package store

import (
	"context"
	"database/sql"
	"time"
)

type PublicDataset struct {
	ID          int64     `json:"id"`
	DatasetName string    `json:"dataset_name"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Address     *string   `json:"address"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	ImportedAt  time.Time `json:"imported_at"`
}

// CountDatasets returneaza numarul de randuri aflate in tabelele CSV
func CountDatasets(ctx context.Context, db *sql.DB) (int64, error) {
	query := `
	SELECT COUNT(*)
	FROM public_datasets`

	var count int64

	if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// InsertDataset introduce un rand al fisierului csv in baza de date
func InsertDataset(ctx context.Context, db *sql.DB, d PublicDataset) error {
	query := `
	INSERT INTO public_datasets (dataset_name, category, name, address, 
	latitude, longitude)
	VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, query, d.DatasetName, d.Category,
		d.Name, d.Address, d.Latitude, d.Longitude); err != nil {
		return err
	}

	return nil
}

// ListDatasets aduce din baza de date tseturile de date si le filtreaza
// daca este cazul, dupa categorie, si search-ul primit
// plimiteaza si numarul de seturi de date aduse in total
func ListDatasets(ctx context.Context, db *sql.DB, category, search string,
	limit, offset int) ([]PublicDataset, error) {
	query := `
	SELECT id, dataset_name, category, name, address, latitude, longitude, imported_at
	FROM public_datasets
	WHERE ($1 = '' OR category = $1)
	AND ($2 = '' OR name ILIKE "%" || $2 || "%" OR
	COALESCE (address, '') ILIKE "%" || $2 || "%")
	ORDER BY id ASC
	LIMIT $3 OFFSET $4`

	rows, err := db.QueryContext(ctx, query, category, search, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	datasets := make([]PublicDataset, 0)

	for rows.Next() {
		var pd PublicDataset

		if err := rows.Scan(&pd.ID, &pd.DatasetName, &pd.Category, &pd.Name,
			&pd.Address, &pd.Latitude, &pd.Longitude, &pd.ImportedAt); err != nil {
			return nil, err
		}

		datasets = append(datasets, pd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return datasets, nil
}
