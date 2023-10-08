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
	emptyProduct = model.Product{}
)

// CreateProduct implements ProductRepository
func (r *productRepository) CreateProduct(product model.Product) (model.Product, error) {
	err := r.DB.Create(&product).Error
	if err != nil {
		if errors.Is(err, common.ErrFailCreateData) {
			return emptyProduct, fmt.Errorf("product: %w", common.ErrFailCreateData)
		}

		return emptyProduct, err
	}

	return product, nil
}
