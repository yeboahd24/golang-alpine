package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/yeboahd24/authentication/internal/api"
	"github.com/yeboahd24/authentication/internal/config"
	"github.com/yeboahd24/authentication/internal/db"
	sqlc "github.com/yeboahd24/authentication/internal/db/sqlc"
	"github.com/yeboahd24/authentication/internal/service"
)

func runDBMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetBaseFS(nil) // Clear any previously set filesystem

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}

func main() {
	dbConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewDB(&dbConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migrationsDir := filepath.Join("internal", "db", "migrations")
	if err := runDBMigrations(database, migrationsDir); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	queries := sqlc.New(database)

	jwtMaker := service.NewJWTMaker(dbConfig.JWT.SecretKey)

	emailService := service.NewEmailService(service.EmailConfig{
		Host:     dbConfig.Email.Host,
		Port:     dbConfig.Email.Port,
		Username: dbConfig.Email.Username,
		Password: dbConfig.Email.Password,
		From:     dbConfig.Email.From,
	})

	server := api.NewServer(queries, jwtMaker, emailService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Router(),
	}

	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
