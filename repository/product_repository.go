package repository

import (
	"errors"
	"fmt"
	"project/domain"
	"project/helper"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewProductRepository(db *gorm.DB, log *zap.Logger) *ProductRepository {
	return &ProductRepository{db: db, log: log}
}

func (repo ProductRepository) All(
	page, limit int,
	productStatus, categoryName, availability string,
	quantity int,
	minPrice, maxPrice float64,
) ([]*domain.Product, int64, error) {

	var products []*domain.Product
	var totalItems int64

	// Mulai query
	query := repo.db.Model(&domain.Product{}).
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.deleted_at IS NULL")

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
		query = query.Where("products.status = ?", productStatus)
	}

	// Filter by Stock (In Stock, Low Stock, Out Of Stock)
	if availability == "In Stock" {
		query = query.Where("products.stock > ?", 5)
	} else if availability == "Low Stock" {
		query = query.Where("products.stock > ?", 0).Where("products.stock <= ?", 5)
	} else if availability == "Out Of Stock" {
		query = query.Where("products.stock < ?", 1)
	}

	// Filter by Quantity
	if quantity > 0 {
		query = query.Where("products.stock = ?", quantity)
	}

	// Filter by Price Range
	if minPrice > 0 {
		query = query.Where("products.price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("products.price <= ?", maxPrice)
	}

	// Count total items after applying filters
	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total Product", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No Product found with given filters")
		return []*domain.Product{}, 0, nil
	}

	err := query.Scopes(
		helper.Paginate(uint(page), uint(limit)),
	).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, created_at, updated_at")
	}).Find(&products).Error

	if err != nil {
		repo.log.Error("Failed to fetch product", zap.Error(err))
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (repo ProductRepository) Add(product *domain.Product, categoryName string) (*domain.Product, error) {
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

	// Assign CategoryID ke product
	product.CategoryID = int(category.ID)

	// Simpan data product baru
	if err := repo.db.Create(product).Error; err != nil {
		repo.log.Error("Failed to add product", zap.Error(err))
		return nil, fmt.Errorf("failed to add product: %v", err)
	}

	repo.log.Info("product added successfully", zap.Uint("product_id", product.ID))
	return product, nil
}

func (repo ProductRepository) Update(id uint, ProductData *domain.Product, categoryName string) (*domain.Product, error) {
	var existingProduct domain.Product

	// Cek apakah product dengan ID tertentu ada di database
	if err := repo.db.First(&existingProduct, id).Error; err != nil {
		repo.log.Error("product not found", zap.Uint("product_id", id), zap.Error(err))
		return nil, fmt.Errorf("product with ID %d not found", id)
	}

	// Validasi Category Name jika diberikan
	if categoryName != "" {
		var category domain.Category
		if err := repo.db.Where("name = ?", categoryName).First(&category).Error; err != nil {
			repo.log.Warn("Category not found", zap.String("category_name", categoryName))
			return nil, fmt.Errorf("category '%s' not found", categoryName) // Pastikan error ini terpropagasi dengan benar
		}
		ProductData.CategoryID = int(category.ID)
	}

	// Cek nilai 'quantity' jika tidak ada, gunakan nilai lama
	if ProductData.Stock == 0 {
		ProductData.Availability = existingProduct.Availability // Pertahankan Stock lama
	} else if ProductData.Stock > 0 && ProductData.Stock <= 5 {
		ProductData.Availability = "Low Stock"
	} else if ProductData.Stock > 5 {
		ProductData.Availability = "In Stock"
	} else if ProductData.Stock < 1 {
		ProductData.Availability = "Out of Stock"
	}

	// Update field lainnya
	if err := repo.db.Model(&existingProduct).Select(
		"CategoryID", "Image", "Name", "CodeProduct", "Stock", "Price", "Status", "Availability",
	).Updates(ProductData).Error; err != nil {
		repo.log.Error("Failed to update inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to update inventory: %v", err)
	}

	repo.log.Info("Inventory updated successfully", zap.Uint("inventory_id", existingProduct.ID))
	return &existingProduct, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	var Product domain.Product

	// Cek apakah Product dengan ID tertentu ada di database
	if err := repo.db.First(&Product, id).Error; err != nil {
		repo.log.Error("Product not found", zap.Uint("Product_id", id), zap.Error(err))
		return fmt.Errorf("product with id %d not found", id)
	}

	// Lakukan soft delete dengan mengisi field deleted_at
	currentTime := time.Now()
	Product.DeletedAt = &currentTime

	// Update database dengan soft delete
	if err := repo.db.Save(&Product).Error; err != nil {
		repo.log.Error("Failed to soft delete Product", zap.Error(err))
		return fmt.Errorf("failed to soft delete Product: %v", err)
	}

	repo.log.Info("product soft deleted successfully", zap.Uint("product_id", Product.ID))
	return nil
}
