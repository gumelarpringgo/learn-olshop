package repository

import (
	"errors"
	"fmt"
	"learn/common"
	"learn/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product model.Product) (model.Product, error)
	FindProductById(productId int) (model.Product, error)
	FindAllProductImagesByProductId(productId int) ([]model.ProductImage, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

var (
	emptyProduct       = model.Product{}
	emptyProductImages = []model.ProductImage{}
)

// CreateProduct implements ProductRepository
func (r *productRepository) CreateProduct(product model.Product) (model.Product, error) {
	err := r.DB.Create(&product).Error
	if err != nil {
		if errors.Is(err, common.ErrFailCreateData) {
			return emptyProduct, fmt.Errorf("product: %w", common.ErrFailCreateData)
		}
	}

	return product, nil
}

// FindProductById implements ProductRepository
func (r *productRepository) FindProductById(productId int) (model.Product, error) {
	product := model.Product{}

	err := r.DB.Where("id = ?", productId).Find(&product).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return emptyProduct, fmt.Errorf("product %d: %w", productId, common.ErrNotFound)
		}
	}

	return product, nil
}

// FindAllProductImagesByProductId implements ProductRepository
func (r *productRepository) FindAllProductImagesByProductId(productId int) ([]model.ProductImage, error) {
	productImage := []model.ProductImage{}

	err := r.DB.Where("product_id = ?", productId).Find(&productImage).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return emptyProductImages, fmt.Errorf("product %d: %w", productId, common.ErrNotFound)
		}
	}

	return productImage, nil
}
