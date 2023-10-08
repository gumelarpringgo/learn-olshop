package model

import "time"

// DATABASE
type (
	Product struct {
		Id            int
		Name          string
		Description   string
		Quantity      int
		Price         int
		ProductImages []ProductImage
		CreatedAt     time.Time
		UpdatedAt     time.Time
		DeletedAt     time.Time
	}

	ProductImage struct {
		Id        int
		ProductId int
		FileName  string
		IsPrimary string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

// REQUEST
type (
	AddProductReq struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
		Quantity    int    `json:"quantity" validate:"required"`
		Price       int    `json:"price" validate:"required"`
	}
)

// RESPONSE
type (
	ProductImageRes struct {
		ProductId int    `json:"product_id"`
		FileName  string `json:"file_name"`
		IsPrimary string `json:"is_primary"`
	}

	ProductRes struct {
		Name          string            `json:"name"`
		Description   string            `json:"description"`
		Quantity      int               `json:"quantity"`
		Price         int               `json:"price"`
		ProductImages []ProductImageRes `json:"product_images"`
	}
)

// Formatter Response
func ProductFormatRes(product Product) ProductRes {
	response := ProductRes{
		Name:          product.Name,
		Description:   product.Description,
		Quantity:      product.Quantity,
		Price:         product.Price,
		ProductImages: ProductImagesFormatRes(product.ProductImages),
	}
	return response
}

func ProductImageFormatRes(pi ProductImage) ProductImageRes {
	return ProductImageRes{
		ProductId: pi.ProductId,
		FileName:  pi.FileName,
		IsPrimary: pi.IsPrimary,
	}
}

func ProductImagesFormatRes(images []ProductImage) []ProductImageRes {
	productImagesFormatRes := []ProductImageRes{}

	for _, productImage := range images {
		productImageFormatRes := ProductImageFormatRes(productImage)

		productImagesFormatRes = append(productImagesFormatRes, productImageFormatRes)
	}
	return productImagesFormatRes
}
