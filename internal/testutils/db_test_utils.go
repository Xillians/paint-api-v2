package testutils

import (
	"log"
	"paint-api/internal/config"
	"paint-api/internal/db"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenTestConnection opens an in-memory database connection for testing purposes.
// This runs a migration on the database to ensure the schema is up to date.
//
// Returns a gorm.DB connection and a cleanup function.
func OpenTestConnection() (*gorm.DB, func()) {
	cfg := &config.DbConfig{
		DatabaseUrl: "file::memory:?cache=shared",
	}

	output, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        cfg.DatabaseUrl,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	err = output.AutoMigrate(&db.PaintBrands{}, &db.PaintCollection{}, &db.Users{}, &db.Paints{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Enable foreign key constraints
	output.Exec("PRAGMA foreign_keys = ON")

	// Return the database connection and a cleanup function
	cleanup := func() {
		sqlDB, _ := output.DB()
		sqlDB.Close()
	}

	return output, cleanup
}

type TestData struct {
	User    *db.Users
	Brand   *db.PaintBrands
	Paint   *db.Paints
	Entry   *db.CollectionPaintDetails
	Cleanup func()
}

func MakeTestData(connection *gorm.DB) (*TestData, error) {
	userInput := db.RegisterUserInput{
		GoogleUserId: "123456",
		Email:        "asd@fgh.io",
	}
	user, err := db.Users{}.RegisterUser(connection, userInput, "administrator")
	if err != nil {
		log.Fatalf("Failed to create test user: %v", err)
		return nil, err
	}

	brandInput := db.CreateBrandInput{
		Name: "Test Brand",
	}
	brand, err := db.PaintBrands{}.CreateBrand(connection, &brandInput)
	if err != nil {
		log.Fatalf("Failed to create test brand: %v", err)
		return nil, err
	}

	paintInput := db.CreatePaintInput{
		Name:        "Chaos black",
		BrandId:     brand.ID,
		ColorCode:   "#000000",
		Description: "Test description",
	}
	paint, err := db.Paints{}.CreatePaint(connection, &paintInput)
	if err != nil {
		log.Fatalf("Failed to create test paint: %v", err)
		return nil, err
	}

	collectionInput := db.CreateCollectionEntryInput{
		Quantity: 1,
		PaintID:  paint.Id,
		UserId:   user.ID,
	}
	entry, err := db.CollectionPaintDetails{}.CreateEntry(connection, collectionInput)
	if err != nil {
		log.Fatalf("Failed to create test collection entry: %v", err)
		return nil, err
	}

	return &TestData{
		User:  user,
		Brand: brand,
		Paint: paint,
		Entry: entry,
		Cleanup: func() {
			err := deleteTestData(connection, user, brand, paint, entry)
			if err != nil {
				log.Fatalf("Failed to cleanup test data: %v", err)
			}
		},
	}, nil
}

func deleteTestData(connection *gorm.DB, user *db.Users, brand *db.PaintBrands, paint *db.Paints, entry *db.CollectionPaintDetails) error {
	err := db.CollectionPaintDetails{}.DeleteEntry(connection, entry.ID)
	if err != nil {
		log.Fatalf("Failed to delete test collection entry: %v", err)
		return err
	}

	err = db.Paints{}.DeletePaint(connection, paint.Id)
	if err != nil {
		log.Fatalf("Failed to delete test paint: %v", err)
		return err
	}

	err = db.PaintBrands{}.DeleteBrand(connection, brand.ID)
	if err != nil {
		log.Fatalf("Failed to delete test brand: %v", err)
		return err
	}

	err = db.Users{}.DeleteUserByGoogleId(connection, user.GoogleUserId)
	if err != nil {
		log.Fatalf("Failed to delete test user: %v", err)
		return err
	}

	return nil
}
