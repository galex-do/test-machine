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
        _ "github.com/lib/pq"
        "github.com/pressly/goose/v3"
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

        // Run database migrations with Goose
        if err := goose.SetDialect("postgres"); err != nil {
                log.Fatal("Failed to set goose dialect:", err)
        }

        // Determine migrations directory (Docker vs local)
        migrationsDir := "./migrations"
        if _, err := os.Stat("/app/migrations"); err == nil {
                migrationsDir = "/app/migrations"
        }

        // Run migrations from the migrations directory
        if err := goose.Up(db, migrationsDir); err != nil {
                log.Fatal("Failed to run database migrations:", err)
        }
        log.Println("Database migrations completed successfully")

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
