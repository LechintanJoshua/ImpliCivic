package csvproc

import (
	"encoding/csv"
	"implicivic/internal/store"
	"io"
	"strconv"
	"time"
)

// WriteIncidentsCSV scrie in format csv tabelul de sesizari primit
func WriteIncidentsCSV(w io.Writer, incidents []store.Incident) error {
	writer := csv.NewWriter(w)
	header := []string{"id", "user_id", "title", "category", "description", "address",
		"latitude", "longitude", "status", "admin_notes", "created_at", "updated_at"}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, inc := range incidents {
		row := []string{
			strconv.FormatInt(inc.ID, 10),
			strconv.FormatInt(inc.UserId, 10),
			inc.Title,
			inc.Category,
			inc.Description,
			inc.Address,
			formatOptionalFloat(inc.Latitude),
			formatOptionalFloat(inc.Longitude),
			inc.Status,
			formatOptionalString(inc.AdminNotes),
			inc.CreatedAt.Format(time.RFC3339),
			inc.UpdatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

// formatOptionalFloat converteste un numar real la string, sau "" daca
// acesta nu exista
func formatOptionalFloat(p *float64) string {
	if p == nil {
		return ""
	}

	return strconv.FormatFloat(*p, 'f', 6, 64)
}

// formatOptionalString verifica daca pointerul pasat e nil
// si returneaza empty-string sau valoarea deferintiata
func formatOptionalString(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}
