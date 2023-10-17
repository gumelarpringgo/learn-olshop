package service

import (
	"fmt"
	"learn/common"
	"learn/model"
	"learn/repository"
)

type ProductService interface {
	AddProduct(req model.ProductReq) (model.ProductRes, error)
	FindProductById(productId int) (model.ProductRes, error)
	UpdateProduct(req model.ProductReq, productId int) (model.ProductRes, error)
	DeleteProduct(productId int) (model.MessageResponse, error)

	FindAllProductImagesByProductId(productId int) (model.ProductImagesRes, error)
	UploadProductImages(req model.ProductImagesUploadReq, productId int, productName string) (model.MessageResponse, error)
	DeleteProductImageId(prodImgId int, roductId int) (model.MessageResponse, error)
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
	emptyMessageRes    = model.MessageResponse{}
	emptyProductImages = model.ProductImagesRes{}
)

// AddProduct implements ProductService
func (s *productService) AddProduct(req model.ProductReq) (model.ProductRes, error) {
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

// UpdateProductById implements ProductService
func (s *productService) UpdateProduct(req model.ProductReq, productId int) (model.ProductRes, error) {
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

	product.Name = req.Name
	product.Description = req.Description
	product.Quantity = req.Quantity
	product.Price = req.Price
	product.ProductImages = productImages

	productUpdate, err := s.Repo.UpdateProduct(product)
	if err != nil {
		return emptyAddProductRes, fmt.Errorf("UpdateProduct call failed: %w", err)
	}

	response := model.ProductFormatRes(productUpdate)
	return response, nil
}

// DeleteProduct implements ProductService
func (s *productService) DeleteProduct(productId int) (model.MessageResponse, error) {
	product, err := s.Repo.FindProductById(productId)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("FindProductById call failed: %w", err)
	}

	if product.Id == 0 {
		return emptyMessageRes, fmt.Errorf("product %d : %w", productId, common.ErrNotFound)
	}

	err = s.Repo.DeleteProduct(productId)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("DeleteProduct call failed: %w", err)
	}

	response := model.MessageResponse{
		Message: fmt.Sprintf("product id %d successfully deleted", productId),
	}

	return response, nil
}

// FindAllProductImagesByProductId implements ProductService
func (s *productService) FindAllProductImagesByProductId(productId int) (model.ProductImagesRes, error) {
	productImages, err := s.Repo.FindAllProductImagesByProductId(productId)
	if err != nil {
		return emptyProductImages, fmt.Errorf("FindAllProductImagesByProductId call failed: %w", err)
	}

	formatProductImagesRes := model.ProductImagesFormatRes(productImages)
	response := model.ProductImagesRes{
		ProductImages: formatProductImagesRes,
	}
	return response, nil
}

// UploadProductImages implements ProductService
func (s *productService) UploadProductImages(req model.ProductImagesUploadReq, productId int, productName string) (model.MessageResponse, error) {
	productImage := model.ProductImage{}

	prodImages, err := s.Repo.FindAllProductImagesByProductId(productId)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("FindAllProductImagesByProductId call failed: %w", err)
	}

	isPrimary := "no"
	if len(prodImages) == 0 && req.IsPrimary == "no" {
		return emptyMessageRes, common.ErrMustHavePrimary
	} else if len(prodImages) > 0 && req.IsPrimary == "yes" {
		isPrimary = "yes"

		_, err := s.Repo.MarkAllProductImagesNonPrimary(productId)
		if err != nil {
			return emptyMessageRes, fmt.Errorf("MarkAllProductImagesNonPrimary call failed: %w", err)
		}
	} else if req.IsPrimary == "yes" {
		isPrimary = "yes"
	}

	productImage.ProductId = productId
	productImage.FileName = productName
	productImage.IsPrimary = isPrimary

	_, err = s.Repo.CreateProductImages(productImage)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("CreateProductImages call failed: %w", err)
	}

	response := model.MessageResponse{
		Message: "upload product image successfully",
	}

	return response, nil
}

// DeleteProductImageId implements ProductService
func (s *productService) DeleteProductImageId(prodImgId int, productId int) (model.MessageResponse, error) {
	err := s.Repo.DeleteProductImageById(prodImgId)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("DeleteProductImageById call failed: %w", err)
	}

	productImages, err := s.Repo.FindAllProductImagesByProductId(productId)
	if err != nil {
		return emptyMessageRes, fmt.Errorf("FindAllProductImagesByProductId call failed: %w", err)
	}

	for _, productImage := range productImages {
		if productImage.IsPrimary != "yes" {
			prodImg := productImages[0]
			prodImg.IsPrimary = "yes"

			_, err := s.Repo.UpdateProductImageById(prodImg)
			if err != nil {
				return emptyMessageRes, fmt.Errorf("UpdateProductImageById call failed: %w", err)
			}
		}
	}

	response := model.MessageResponse{
		Message: fmt.Sprintf("product image id %d successfully deleted", prodImgId),
	}

	return response, nil
}
