package main

import (
	"database/sql"
	"log"
	"os"
	"virtual-library/catalog-service/internal/catalog"
	"virtual-library/catalog-service/internal/httpapi"
	"virtual-library/catalog-service/internal/memory"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalf("Error closing DB connection: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to PostgreSQL:", err)
	}
	log.Println("Connected to PostgreSQL")

	repo := memory.NewMemoryBookRepository()
	//repo := postgres.NewBookRepository(db)
	service := catalog.NewService(repo)
	handler := httpapi.NewBookHandler(service)
	router := httpapi.NewRouter(handler)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
