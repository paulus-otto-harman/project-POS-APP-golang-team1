package repository

import (
	"errors"
	"fmt"
	"project/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewInventoryRepository(db *gorm.DB, log *zap.Logger) *InventoryRepository {
	return &InventoryRepository{db: db, log: log}
}

func (repo InventoryRepository) All(page, limit int, productStatus, categoryName, stock string, quantity int, minPrice, maxPrice float64) ([]*domain.Inventory, int64, error) {

	var inventory []*domain.Inventory
	var totalItems int64

	// Mulai query
	query := repo.db.Model(&domain.Inventory{}).
		Joins("JOIN categories ON categories.id = inventories.category_id")

	// Filter by Category Name
	if categoryName != "" {
		// Cek apakah category_name ada di database
		var count int64
		err := repo.db.Model(&domain.Category{}).Where("name = ?", categoryName).Count(&count).Error
		if err != nil {
			repo.log.Error("Failed to check category", zap.Error(err))
			return nil, 0, err
		}

		if count == 0 {
			repo.log.Warn("Category not found", zap.String("category", categoryName))
			return nil, 0, fmt.Errorf("category '%s' not found", categoryName)
		}

		// Tambahkan filter ke query jika category_name ada
		query = query.Where("categories.name = ?", categoryName)
	}

	// Filter by Product Status (Active, Inactive, All)
	if productStatus != "" && productStatus != "all" {
		query = query.Where("inventories.status = ?", productStatus)
	}

	// Filter by Stock (In Stock, Low Stock, Out Of Stock)
	if stock == "In Stock" {
		query = query.Where("inventories.quantity > ?", 5)
	} else if stock == "Low Stock" {
		query = query.Where("inventories.quantity > ?", 0).Where("inventories.quantity <= ?", 5)
	} else if stock == "Out Of Stock" {
		query = query.Where("inventories.quantity < ?", 1)
	}

	// Filter by Quantity
	if quantity > 0 {
		query = query.Where("inventories.quantity = ?", quantity)
	}

	// Filter by Price Range
	if minPrice > 0 {
		query = query.Where("inventories.price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("inventories.price <= ?", maxPrice)
	}

	// Count total items after applying filters
	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total inventory", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No inventory found with given filters")
		return []*domain.Inventory{}, 0, nil
	}

	// Fetch filtered inventory data
	err := query.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, created_at, updated_at")
	}).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&inventory).Error

	if err != nil {
		repo.log.Error("Failed to fetch inventory", zap.Error(err))
		return nil, 0, err
	}

	return inventory, totalItems, nil
}

func (repo InventoryRepository) Add(inventory *domain.Inventory, categoryName string) (*domain.Inventory, error) {
	var category domain.Category

	// Cari Category berdasarkan nama
	err := repo.db.Where("name = ?", categoryName).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			repo.log.Warn("Category not found", zap.String("category_name", categoryName))
			return nil, fmt.Errorf("category name '%s' not found", categoryName)
		}
		// Log error database lainnya
		repo.log.Error("Database error during category validation", zap.Error(err))
		return nil, fmt.Errorf("failed to validate category: %v", err)
	}

	// Assign CategoryID ke inventory
	inventory.CategoryID = category.ID

	// Simpan data inventory baru
	if err := repo.db.Create(inventory).Error; err != nil {
		repo.log.Error("Failed to add inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to add inventory: %v", err)
	}

	repo.log.Info("Inventory added successfully", zap.Uint("inventory_id", inventory.ID))
	return inventory, nil
}

func (repo InventoryRepository) Update(id uint, inventoryData *domain.Inventory, categoryName string) (*domain.Inventory, error) {
	var existingInventory domain.Inventory

	// Cek apakah inventory dengan ID tertentu ada di database
	if err := repo.db.First(&existingInventory, id).Error; err != nil {
		repo.log.Error("Inventory not found", zap.Uint("inventory_id", id), zap.Error(err))
		return nil, fmt.Errorf("inventory with ID %d not found", id)
	}

	// Validasi Category Name jika diberikan
	if categoryName != "" {
		var category domain.Category
		if err := repo.db.Where("name = ?", categoryName).First(&category).Error; err != nil {
			repo.log.Warn("Category not found", zap.String("category_name", categoryName))
			return nil, fmt.Errorf("category '%s' not found", categoryName) // Pastikan error ini terpropagasi dengan benar
		}
		inventoryData.CategoryID = category.ID
	} else {
		// Jika kategori kosong, biarkan ID kategori lama tetap dipakai
		inventoryData.CategoryID = existingInventory.CategoryID
	}

	// Cek nilai 'status' jika tidak ada, gunakan nilai lama
	if inventoryData.Status == "" {
		inventoryData.Status = existingInventory.Status // Pertahankan status lama
	}

	// Cek nilai 'name' jika tidak ada, gunakan nilai lama
	if inventoryData.Name == "" {
		inventoryData.Name = existingInventory.Name // Pertahankan name lama
	}

	// Cek nilai 'price' jika tidak ada, gunakan nilai lama
	if inventoryData.Price == 0 {
		inventoryData.Price = existingInventory.Price // Pertahankan price lama
	}

	// Cek nilai 'quantity' jika tidak ada, gunakan nilai lama
	if inventoryData.Quantity == 0 {
		inventoryData.Quantity = existingInventory.Quantity // Pertahankan quantity lama
	} else if inventoryData.Quantity > 0 && inventoryData.Quantity <= 5 {
		inventoryData.Stock = "Low Stock"
	} else if inventoryData.Quantity > 5 {
		inventoryData.Stock = "In Stock"
	} else if inventoryData.Quantity < 1 {
		inventoryData.Stock = "Out of Stock"
	}

	// Cek nilai 'image' jika tidak ada, gunakan nilai lama
	if inventoryData.Image == "" {
		inventoryData.Image = existingInventory.Image // Pertahankan image lama
	}

	// Cek nilai 'code_product' jika tidak ada, gunakan nilai lama
	if inventoryData.CodeProduct == "" {
		inventoryData.CodeProduct = existingInventory.CodeProduct // Pertahankan code_product lama
	}

	// Update field lainnya
	if err := repo.db.Model(&existingInventory).Select(
		"CategoryID", "Image", "Name", "CodeProduct", "Quantity", "Price", "Status", "Stock",
	).Updates(inventoryData).Error; err != nil {
		repo.log.Error("Failed to update inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to update inventory: %v", err)
	}

	repo.log.Info("Inventory updated successfully", zap.Uint("inventory_id", existingInventory.ID))
	return &existingInventory, nil
}

func (repo *InventoryRepository) Delete(id uint) error {
	var inventory domain.Inventory

	// Cek apakah inventory dengan ID tertentu ada di database
	if err := repo.db.First(&inventory, id).Error; err != nil {
		repo.log.Error("Inventory not found", zap.Uint("inventory_id", id), zap.Error(err))
		return fmt.Errorf("inventory with ID %d not found", id)
	}

	// Lakukan soft delete dengan mengisi field deleted_at
	currentTime := time.Now()
	inventory.DeletedAt = &currentTime

	// Update database dengan soft delete
	if err := repo.db.Save(&inventory).Error; err != nil {
		repo.log.Error("Failed to soft delete inventory", zap.Error(err))
		return fmt.Errorf("failed to soft delete inventory: %v", err)
	}

	repo.log.Info("Inventory soft deleted successfully", zap.Uint("inventory_id", inventory.ID))
	return nil
}
