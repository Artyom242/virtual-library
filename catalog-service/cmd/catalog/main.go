package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"
	"virtual-library/catalog-service/internal/httpapi"
	"virtual-library/catalog-service/internal/postgres"
	"virtual-library/catalog-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Cannot connect to PostgreSQL:", err)
	}
	log.Println("Connected to PostgreSQL")

	tx := postgres.NewPgxTransactor(pool)
	// Author process
	authorRepo := postgres.NewAuthorRepository(pool)
	authorService := service.NewAuthorService(tx, authorRepo)
	// Book process
	bookRepo := postgres.NewBookRepository(pool)
	bookService := service.NewService(tx, bookRepo, authorRepo)
	router := httpapi.NewRouter(bookService, authorService)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
