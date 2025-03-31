package db_test

import (
	"paint-api/internal/db"
	"testing"
)

func createTestBrand() *db.PaintBrands {
	brand, err := db.PaintBrands{}.CreateBrand(testDB, &db.CreateBrandInput{
		Name: "Test Brand",
	})
	if err != nil {
		return nil
	}
	return brand
}

func TestListBrandsImplementations(t *testing.T) {
	t.Run("List brands with empty list", func(t *testing.T) {
		brands, err := db.PaintBrands{}.ListBrands(testDB)
		if err != nil {
			t.Errorf("Error listing brands: %v", err)
		}
		if len(brands) != 0 {
			t.Errorf("Expected empty list, got %v", brands)
		}
	})
	createTestBrand()
	t.Run("List brands", func(t *testing.T) {
		brands, err := db.PaintBrands{}.ListBrands(testDB)
		if err != nil {
			t.Errorf("Error listing brands: %v", err)
		}
		t.Log(brands)
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		_, err := db.PaintBrands{}.ListBrands(connection)
		if err == nil {
			t.Errorf("Expected error listing brands with nil db, got nil")
		}
	})
}

func TestCreateBrandImplementations(t *testing.T) {
	brandInput := &db.CreateBrandInput{
		Name: "Test Brand",
	}

	t.Run("Create brand", func(t *testing.T) {
		createdBrand, err := db.PaintBrands{}.CreateBrand(testDB, brandInput)
		if err != nil {
			t.Errorf("Error creating brand: %v", err)
		}
		if createdBrand == nil {
			t.Errorf("Expected created brand, got nil")
		}
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		_, err := db.PaintBrands{}.CreateBrand(connection, brandInput)
		if err == nil {
			t.Errorf("Expected error listing brands with nil db, got nil")
		}
	})
}

func TestGetBrandImplementations(t *testing.T) {
	testBrand := createTestBrand()

	t.Run("Get brand", func(t *testing.T) {
		foundBrand, err := db.PaintBrands{}.GetBrand(testDB, testBrand.ID)
		if err != nil {
			t.Errorf("Error getting brand: %v", err)
		}
		if foundBrand == nil {
			t.Errorf("Expected found brand, got nil")
		}
	})
	t.Run("Get brand with invalid ID", func(t *testing.T) {
		_, err := db.PaintBrands{}.GetBrand(testDB, 0)
		if err == nil {
			t.Errorf("Expected error getting brand with invalid ID, got nil")
		}
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		_, err := db.PaintBrands{}.GetBrand(connection, testBrand.ID)
		if err == nil {
			t.Errorf("Expected error listing brands with nil db, got nil")
		}
	})
}

func TestUpdateBrandImplementations(t *testing.T) {
	testBrand := createTestBrand()

	t.Run("Update brand", func(t *testing.T) {
		updateInput := &db.UpdateBrandInput{
			Name: "Updated Test Brand",
		}
		updatedBrand, err := db.PaintBrands{}.UpdateBrand(testDB, testBrand.ID, updateInput)
		if err != nil {
			t.Errorf("Error updating brand: %v", err)
		}
		if updatedBrand == nil {
			t.Errorf("Expected updated brand, got nil")
		}
		if updatedBrand.Name != updateInput.Name {
			t.Errorf("Expected updated brand name to be %s, got %s", updateInput.Name, updatedBrand.Name)
		}
		if testBrand.UpdatedAt == updatedBrand.UpdatedAt {
			t.Errorf("Expected updated brand to have updated at timestamp")
		}
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		updateInput := &db.UpdateBrandInput{
			Name: "Updated Test Brand",
		}

		_, err := db.PaintBrands{}.UpdateBrand(connection, testBrand.ID, updateInput)
		if err == nil {
			t.Errorf("Expected error listing brands with nil db, got nil")
		}
	})
	t.Run("Try updating brand with invalid ID", func(t *testing.T) {
		updateInput := &db.UpdateBrandInput{
			Name: "Updated Test Brand",
		}
		_, err := db.PaintBrands{}.UpdateBrand(testDB, 0, updateInput)
		if err == nil {
			t.Errorf("Expected error updating brand with invalid ID, got nil")
		}
	})
}

func TestDeleteBrandImplementations(t *testing.T) {
	t.Run("Delete brand", func(t *testing.T) {
		brandInput := &db.CreateBrandInput{
			Name: "Test Brand",
		}
		newBrand, err := db.PaintBrands{}.CreateBrand(testDB, brandInput)
		if err != nil {
			t.Errorf("Error creating brand: %v", err)
		}
		err = db.PaintBrands{}.DeleteBrand(testDB, newBrand.ID)
		if err != nil {
			t.Errorf("Error deleting brand: %v", err)
		}
		_, err = db.PaintBrands{}.GetBrand(testDB, newBrand.ID)
		if err == nil {
			t.Errorf("Expected error getting deleted brand, got nil")
		}
	})
	t.Run("transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()

		brand, err := db.PaintBrands{}.CreateBrand(connection, &db.CreateBrandInput{
			Name: "Test Brand",
		})
		if err != nil {
			t.Errorf("Error creating brand: %v", err)
		}

		sql, _ := connection.DB()
		sql.Close()

		err = db.PaintBrands{}.DeleteBrand(connection, brand.ID)
		if err == nil {
			t.Errorf("Expected transaction faliure, it did not fail")
		}
	})
	t.Run("Delete non-existing brand", func(t *testing.T) {
		err := db.PaintBrands{}.DeleteBrand(testDB, 0)
		if err == nil {
			t.Errorf("Expected error deleting non-existing brand, got nil")
		}
	})
}
