package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"
	"virtual-library/catalog-service/internal/app/services"
	"virtual-library/catalog-service/internal/httpapi"
	"virtual-library/catalog-service/internal/infra/config"
	pgrs "virtual-library/catalog-service/internal/infra/db/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgrs.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	tx := pgrs.NewPgxTransactor(pool)
	// Author process
	authorRepo := pgrs.NewAuthorRepository(pool)
	authorService := services.NewAuthorService(tx, authorRepo)
	// Book process
	bookRepo := pgrs.NewBookRepository(pool)
	bookService := services.NewService(tx, bookRepo, authorRepo)
	router := httpapi.NewRouter(bookService, authorService)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
