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
	UpdateProduct(product model.Product) (model.Product, error)
	DeleteProduct(productId int) error

	FindAllProductImagesByProductId(productId int) ([]model.ProductImage, error)
	CreateProductImages(productImages model.ProductImage) (model.ProductImage, error)
	MarkAllProductImagesNonPrimary(productId int) (bool, error)
	DeleteProductImageById(prodImgId int) error
	UpdateProductImageById(productImage model.ProductImage) (model.ProductImage, error)
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
	emptyProductImage  = model.ProductImage{}
)

// CreateProduct implements ProductRepository
func (r *productRepository) CreateProduct(product model.Product) (model.Product, error) {
	err := r.DB.Create(&product).Error
	if err != nil {
		if errors.Is(err, common.ErrFailedCreateData) {
			return emptyProduct, fmt.Errorf("product: %w", common.ErrFailedCreateData)
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

// UpdateProduct implements ProductRepository
func (r *productRepository) UpdateProduct(product model.Product) (model.Product, error) {
	err := r.DB.Save(&product).Error
	if err != nil {
		if errors.Is(err, common.ErrFailedUpdateData) {
			return emptyProduct, fmt.Errorf("product : %w", common.ErrNotFound)
		}
	}

	return product, nil
}

// DeleteProduct implements ProductRepository
func (r *productRepository) DeleteProduct(productId int) error {
	var product model.Product
	var productImage model.ProductImage

	productImages := []model.ProductImage{}
	err := r.DB.Where("product_id = ?", productId).Find(&productImages).Delete(&productImage).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("product %d: %w", productId, common.ErrDeleteData)
		}
	}

	err = r.DB.Delete(&product, productId).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("product %d: %w", productId, common.ErrNotFound)
		}
	}

	return nil
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

// CreateProductImages implements ProductRepository
func (r *productRepository) CreateProductImages(productImages model.ProductImage) (model.ProductImage, error) {
	var productImage model.ProductImage

	err := r.DB.Create(&productImages).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return emptyProductImage, fmt.Errorf("product %w", common.ErrNotFound)
		}
	}

	return productImage, nil
}

// MarkAllProductImagesNonPrimary implements ProductRepository
func (r *productRepository) MarkAllProductImagesNonPrimary(productId int) (bool, error) {
	err := r.DB.Model(&model.ProductImage{}).Where("product_id = ?", productId).Update("is_primary", "no").Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// FindProductImageById implements ProductRepository
func (r *productRepository) DeleteProductImageById(prodImgId int) error {
	productImage := model.ProductImage{}

	err := r.DB.Delete(&productImage, prodImgId).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("product %d: %w", prodImgId, common.ErrNotFound)
		}
	}

	return nil
}

// UpdateProductImageById implements ProductRepository
func (r *productRepository) UpdateProductImageById(productImage model.ProductImage) (model.ProductImage, error) {
	err := r.DB.Save(&productImage).Error
	if err != nil {
		if errors.Is(err, common.ErrFailedUpdateData) {
			return emptyProductImage, fmt.Errorf("product : %w", common.ErrFailedUpdateData)
		}
	}

	return productImage, nil
}
