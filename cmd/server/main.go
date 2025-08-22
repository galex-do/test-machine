package main

import (
        "log"
        "net/http"
        "os"

        "github.com/galex-do/test-machine/internal/config"
        "github.com/galex-do/test-machine/internal/database"
        "github.com/galex-do/test-machine/internal/handlers"
        "github.com/galex-do/test-machine/internal/repository"
        "github.com/galex-do/test-machine/internal/service"
        "github.com/pressly/goose/v3"
        _ "github.com/lib/pq"
)

func main() {
        // Load configuration
        cfg := config.Load()

        // Initialize database
        db, err := database.Connect(cfg.DatabaseURL)
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }
        defer db.Close()

        if err := db.Ping(); err != nil {
                log.Fatal("Failed to ping database:", err)
        }
        log.Println("Connected to PostgreSQL database")

        // Run database migrations with Goose (skip if SKIP_MIGRATIONS is set)
        if os.Getenv("SKIP_MIGRATIONS") == "true" {
                log.Println("Skipping database migrations (SKIP_MIGRATIONS=true)")
        } else {
                if err := goose.SetDialect("postgres"); err != nil {
                        log.Printf("Failed to set goose dialect: %v", err)
                } else {
                        // Check if migrations directory exists
                        migrationDir := "migrations"
                        if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
                                log.Printf("Migrations directory '%s' does not exist, skipping migrations", migrationDir)
                        } else {
                                log.Printf("Running database migrations from directory: %s", migrationDir)
                                if err := goose.Up(db, migrationDir); err != nil {
                                        log.Printf("Failed to run database migrations: %v", err)
                                        log.Printf("Continuing without migrations - they may already be applied")
                                } else {
                                        log.Println("Database migrations completed successfully")
                                }
                        }
                }
        }

        // Initialize repositories
        projectRepo := repository.NewProjectRepository(db)
        testSuiteRepo := repository.NewTestSuiteRepository(db)
        testCaseRepo := repository.NewTestCaseRepository(db)
        testRunRepo := repository.NewTestRunRepository(db)

        // Initialize services
        projectService := service.NewProjectService(projectRepo)
        testSuiteService := service.NewTestSuiteService(testSuiteRepo)
        testCaseService := service.NewTestCaseService(testCaseRepo)
        testRunService := service.NewTestRunService(testRunRepo)

        // Initialize handlers
        handler := handlers.NewHandler(projectService, testSuiteService, testCaseService, testRunService)

        // Setup routes
        mux := handler.SetupRoutes()

        log.Printf("Server starting on port %s...", cfg.Port)
        log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}