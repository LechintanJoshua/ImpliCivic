package csvproc

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"implicivic/internal/store"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type csvJob struct {
	datasetName string
	rowNum      int
	record      []string
}

// ImportDir citeste in streaming fisierele CSV din dir si le proceseaza
// in paralel folosind 4 goroutines (threads)
func ImportDir(ctx context.Context, db *sql.DB, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.csv"))

	if err != nil {
		return err
	}

	jobs := make(chan csvJob, 100)
	var wg sync.WaitGroup

	for range 4 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for job := range jobs {
				processRow(ctx, db, job)
			}
		}()
	}

	for _, file := range files {
		if err := feedFile(ctx, file, jobs); err != nil {
			log.Printf("Eroare la citirea fisierului %s: %v", file, err)
		}
	}

	close(jobs)
	wg.Wait()

	return nil
}

// feedFile citeste rand cu rand din fisierul csv, verifica de erori, si
// il trimite pe canal pentru a fi procesat de goroutines
func feedFile(ctx context.Context, path string, jobs chan<- csvJob) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()

	fileName := filepath.Base(path)

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	if _, err := reader.Read(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}

	rowNum := 2

	for {
		record, err := reader.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Printf("rand stricat: %v", err)
			rowNum++
			continue
		}

		select {
		case jobs <- csvJob{datasetName: fileName, rowNum: rowNum, record: record}:
		case <-ctx.Done():
			return ctx.Err()
		}

		rowNum++
	}

	return nil
}

// processRow parseaza un rand din fisierul csv, verifica daca a fost
// corupt si il salveaza in baza de date
func processRow(ctx context.Context, db *sql.DB, job csvJob) {
	if len(job.record) < 5 {
		log.Printf("skip %s row %d: expected 5 columns", job.datasetName, job.rowNum)
		return
	}

	category := strings.TrimSpace(job.record[0])
	name := strings.TrimSpace(job.record[1])
	address := strings.TrimSpace(job.record[2])
	latitude := strings.TrimSpace(job.record[3])
	longitude := strings.TrimSpace(job.record[4])

	if category == "" || name == "" {
		log.Printf("skip %s row %d: no name or category", job.datasetName, job.rowNum)
		return
	}

	var addrPtr *string

	if address != "" {
		addrPtr = &address
	}

	latptr, err := parseOptionalFloat(latitude)

	if err != nil {
		log.Printf("skip %s row %d: not_a_lat: %s", job.datasetName, job.rowNum, latitude)
		return
	}

	lngPtr, err := parseOptionalFloat(longitude)

	if err != nil {
		log.Printf("skip %s row %d: not_a_long: %s", job.datasetName, job.rowNum, longitude)
		return
	}

	if err := store.InsertDataset(ctx, db, store.PublicDataset{
		DatasetName: job.datasetName,
		Category:    category,
		Name:        name,
		Address:     addrPtr,
		Latitude:    latptr,
		Longitude:   lngPtr,
	}); err != nil {
		log.Printf("skip %s row %d: db insert error: %v", job.datasetName, job.rowNum, err)
		return
	}
}

// parseOptionalFloat retruneaza un pointer la un float64 parsat
// din string
func parseOptionalFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return nil, err
	}

	return &v, nil
}
