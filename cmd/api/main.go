package main

import (
	"implicivic/internal/db"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Eroare la incarcarea variabilelor de mediu")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL nu este setat")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET nu este setat")
	}

	pool, err := db.Connect(databaseURL)

	if err != nil {
		log.Fatalf("Eroare la conectare DB: %v", err)
	}

	defer pool.Close()

	log.Printf("Conectat la PostgreSQL cu success!")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", healthHandler(pool))
	mux.HandleFunc("POST /api/auth/register", registerHandler(pool))
	mux.HandleFunc("POST /api/auth/login", loginHandler(pool, jwtSecret))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(allowedOrigin, mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server pornit pe http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}
