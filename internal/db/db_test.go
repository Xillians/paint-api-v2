package db_test

import (
	"log"
	"os"
	"paint-api/internal/config"
	"paint-api/internal/db"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	testDB = OpenTestConnection()

	code := m.Run()

	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

// Opens a in memory database connection for testing purposes.
// This runs a migration on the database to ensure the schema is up to date.
//
// Returns a gorm.DB connection
func OpenTestConnection() *gorm.DB {
	cfg := &config.DbConfig{
		DatabseUrl: "file::memory:?cache=shared",
	}

	output, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        cfg.DatabseUrl,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	err = output.AutoMigrate(&db.PaintBrands{}, &db.PaintCollection{}, &db.Users{}, &db.Paints{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
	return output
}
