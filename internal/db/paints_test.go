package db_test

import (
	"paint-api/internal/db"
	"testing"
)

func createTestPaint() *db.Paints {
	brand := createTestBrand()
	input := db.CreatePaintInput{
		BrandId:     brand.ID,
		Name:        "Chaos black",
		ColorCode:   "#000000",
		Description: "A very dark black",
	}
	paint, err := db.Paints{}.CreatePaint(testDB, &input)
	if err != nil {
		return nil
	}
	return paint
}

func TestCreatePaint(t *testing.T) {
	brand := createTestBrand()
	testInput := db.CreatePaintInput{
		BrandId:     brand.ID,
		Name:        "Scar white",
		ColorCode:   "#FFFFFF",
		Description: "A very bright white",
	}

	t.Run("Create paint", func(t *testing.T) {
		paint, err := db.Paints{}.CreatePaint(testDB, &testInput)
		if err != nil {
			t.Errorf("Error creating paint: %v", err)
		}

		err = db.Paints{}.DeletePaint(testDB, paint.Id)
		if err != nil {
			t.Errorf("Error deleting paint by id: %v", err)
		}
	})
	t.Run("Create paint with invalid brand", func(t *testing.T) {
		testInput.BrandId = 0
		_, err := db.Paints{}.CreatePaint(testDB, &testInput)
		if err == nil {
			t.Errorf("Expected error creating paint with invalid brand")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting database connection: %v", err)
		}
		sql.Close()

		_, err = db.Paints{}.CreatePaint(connection, &testInput)
		if err == nil {
			t.Errorf("Expected error creating paint with invalid connection")
		}
	})
}

func TestGetPaint(t *testing.T) {
	testPaint := createTestPaint()
	t.Run("Get paint", func(t *testing.T) {
		_, err := db.Paints{}.GetPaint(testDB, testPaint.Id)
		if err != nil {
			t.Errorf("Error fetching paint: %v", err)
		}
	})
	t.Run("Get paint not found", func(t *testing.T) {
		_, err := db.Paints{}.GetPaint(testDB, 0)
		if err == nil {
			t.Errorf("Expected error fetching paint not found")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting database connection: %v", err)
		}
		sql.Close()

		_, err = db.Paints{}.GetPaint(connection, testPaint.Id)
		if err == nil {
			t.Errorf("Expected error fetching paint with invalid connection")
		}
	})
	t.Cleanup(func() {
		err := db.Paints{}.DeletePaint(testDB, testPaint.Id)
		if err != nil {
			t.Errorf("Error deleting paint: %v", err)
		}
		err = db.PaintBrands{}.DeleteBrand(testDB, testPaint.BrandId)
		if err != nil {
			t.Errorf("Error deleting brand: %v", err)
		}
	})
}

func TestListPaints(t *testing.T) {

	t.Run("List paints", func(t *testing.T) {
		_, err := db.Paints{}.ListPaints(testDB)
		if err != nil {
			t.Errorf("Error listing paints: %v", err)
		}
	})

	testPaint := createTestPaint()

	t.Run("List paints", func(t *testing.T) {
		_, err := db.Paints{}.ListPaints(testDB)
		if err != nil {
			t.Errorf("Error listing paints: %v", err)
		}
	})

	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting database connection: %v", err)
		}
		sql.Close()

		_, err = db.Paints{}.ListPaints(connection)
		if err == nil {
			t.Errorf("Expected error listing paints with invalid connection")
		}
	})

	t.Cleanup(func() {
		err := db.Paints{}.DeletePaint(testDB, testPaint.Id)
		if err != nil {
			t.Errorf("Error deleting paint: %v", err)
		}
		err = db.PaintBrands{}.DeleteBrand(testDB, testPaint.BrandId)
		if err != nil {
			t.Errorf("Error deleting brand: %v", err)
		}
	})
}

func TestUpdatePaints(t *testing.T) {
	testPaint := createTestPaint()

	t.Run("Update paint", func(t *testing.T) {
		input := db.UpdatePaintInput{
			Name:        "Scar white",
			BrandId:     testPaint.BrandId,
			ColorCode:   "#FFFFFF",
			Description: "A very bright white",
		}
		paint, err := db.Paints{}.UpdatePaint(testDB, testPaint.Id, &input)
		if err != nil {
			t.Errorf("Error updating paint: %v", err)
		}
		if paint.Name != input.Name {
			t.Errorf("Expected name to be %s, got %s", input.Name, paint.Name)
		}
		if paint.ColorCode != input.ColorCode {
			t.Errorf("Expected color code to be %s, got %s", input.ColorCode, paint.ColorCode)
		}
		if paint.Description != input.Description {
			t.Errorf("Expected description to be %s, got %s", input.Description, paint.Description)
		}
		if paint.BrandId != input.BrandId {
			t.Errorf("Expected brand id to be %d, got %d", input.BrandId, paint.BrandId)
		}
		if paint.UpdatedAt == testPaint.UpdatedAt {
			t.Errorf("Expected updated at to be different")
		}
	})
	t.Run("Update paint with invalid brand", func(t *testing.T) {
		input := db.UpdatePaintInput{
			Name:        "Scar white",
			BrandId:     testPaint.BrandId,
			ColorCode:   "#FFFFFF",
			Description: "A very bright white",
		}
		_, err := db.Paints{}.UpdatePaint(testDB, 0, &input)
		if err == nil {
			t.Errorf("Expected error updating paint with invalid brand")
		}
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting database connection: %v", err)
		}
		sql.Close()

		input := db.UpdatePaintInput{
			Name:        "Scar white",
			BrandId:     testPaint.BrandId,
			ColorCode:   "#FFFFFF",
			Description: "A very bright white",
		}
		_, err = db.Paints{}.UpdatePaint(connection, testPaint.Id, &input)
		if err == nil {
			t.Errorf("Expected error updating paint with invalid connection")
		}
	})
	t.Cleanup(func() {
		err := db.Paints{}.DeletePaint(testDB, testPaint.Id)
		if err != nil {
			t.Errorf("Error deleting paint: %v", err)
		}
		err = db.PaintBrands{}.DeleteBrand(testDB, testPaint.BrandId)
		if err != nil {
			t.Errorf("Error deleting brand: %v", err)
		}
	})
}

func TestDeletePaints(t *testing.T) {
	t.Run("Delete paint", func(t *testing.T) {
		testPaint := createTestPaint()
		err := db.Paints{}.DeletePaint(testDB, testPaint.Id)
		if err != nil {
			t.Errorf("Error deleting paint: %v", err)
		}
	})
	t.Run("Delete paint not found", func(t *testing.T) {
		err := db.Paints{}.DeletePaint(testDB, 0)
		if err == nil {
			t.Errorf("Expected error deleting paint not found")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting database connection: %v", err)
		}
		sql.Close()

		err = db.Paints{}.DeletePaint(connection, 0)
		if err == nil {
			t.Errorf("Expected error deleting paint with invalid connection")
		}
	})
}
