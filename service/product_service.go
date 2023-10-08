package service

import (
	"fmt"
	"learn/model"
	"learn/repository"
)

type ProductService interface {
	AddProduct(req model.AddProductReq) (model.ProductRes, error)
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
	DbProduct = model.Product{}
)

var (
	emptyAddProductRes = model.ProductRes{}
)

var ()

// AddProduct implements ProductService
func (s *productService) AddProduct(req model.AddProductReq) (model.ProductRes, error) {
	DbProduct.Name = req.Name
	DbProduct.Description = req.Description
	DbProduct.Quantity = req.Quantity
	DbProduct.Price = req.Price

	product, err := s.Repo.CreateProduct(DbProduct)
	if err != nil {
		return emptyAddProductRes, fmt.Errorf("CreateProduct call failed: %w", err)
	}

	response := model.ProductFormatRes(product)

	return response, nil
}
