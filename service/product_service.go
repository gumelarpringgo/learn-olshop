package service

import (
	"fmt"
	"learn/common"
	"learn/model"
	"learn/repository"
)

type ProductService interface {
	AddProduct(req model.AddProductReq) (model.ProductRes, error)
	FindProductById(productId int) (model.ProductRes, error)
}

type productService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		Repo: repo,
	}
}

var (
	emptyAddProductRes = model.ProductRes{}
)

// AddProduct implements ProductService
func (s *productService) AddProduct(req model.AddProductReq) (model.ProductRes, error) {
	dbProduct := model.Product{}
	dbProduct.Name = req.Name
	dbProduct.Description = req.Description
	dbProduct.Quantity = req.Quantity
	dbProduct.Price = req.Price

	product, err := s.Repo.CreateProduct(dbProduct)
	if err != nil {
		return emptyAddProductRes, fmt.Errorf("CreateProduct call failed: %w", err)
	}

	response := model.ProductFormatRes(product)

	return response, nil
}

// FindProductById implements ProductService
func (s *productService) FindProductById(productId int) (model.ProductRes, error) {
	product, err := s.Repo.FindProductById(productId)
	if err != nil {
		return emptyAddProductRes, fmt.Errorf("FindProductById call failed: %w", err)
	}

	if product.Id == 0 {
		return emptyAddProductRes, fmt.Errorf("product %d : %w", productId, common.ErrNotFound)
	}

	productImages, err := s.Repo.FindAllProductImagesByProductId(productId)
	if err != nil {
		return emptyAddProductRes, fmt.Errorf("FindAllProductImagesByProductId call failed: %w", err)
	}

	product.ProductImages = productImages

	response := model.ProductFormatRes(product)
	return response, nil
}
