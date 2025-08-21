package database

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"test-management/models"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase() {
	var err error
	
	// Connect to SQLite database
	DB, err = gorm.Open(sqlite.Open("test_management.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = DB.AutoMigrate(
		&models.Project{},
		&models.TestSuite{},
		&models.TestCase{},
		&models.TestStep{},
		&models.TestRun{},
		&models.TestExecution{},
	)
	if err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	log.Println("Database initialized successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
