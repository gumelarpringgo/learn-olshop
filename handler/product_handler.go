package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"learn/common"
	"learn/model"
	"learn/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type ProductHandler interface {
	AddProduct(w http.ResponseWriter, r *http.Request)
	FindProductById(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)

	GetAllProductImagesByProductId(w http.ResponseWriter, r *http.Request)
	UploadProductImage(w http.ResponseWriter, r *http.Request)
	DeleteProductImage(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	Service  service.ProductService
	Validate *validator.Validate
}

func NewProductHandler(service service.ProductService, validate *validator.Validate) ProductHandler {
	return &productHandler{
		Service:  service,
		Validate: validate,
	}
}

// AddProduct implements ProductHandler
func (h *productHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	urlRole := chi.URLParam(r, "role")
	var req model.ProductReq

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if urlRole != role || urlRole != "admin" {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	err = h.Validate.Struct(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	responseProduct, err := h.Service.AddProduct(req)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, responseProduct)
}

// FindProductById implements ProductHandler
func (h *productHandler) FindProductById(w http.ResponseWriter, r *http.Request) {
	urlRole := chi.URLParam(r, "role")
	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	responseProduct, err := h.Service.FindProductById(productIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, responseProduct)
}

// UpdateProduct implements ProductHandler
func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req model.ProductReq

	urlRole := chi.URLParam(r, "role")
	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.Validate.Struct(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.Service.UpdateProduct(req, productIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)
}

// DeleteProduct implements ProductHandler
func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	urlRole := chi.URLParam(r, "role")
	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	response, err := h.Service.DeleteProduct(productIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)
}

// GetAllProductImagesByProductId implements ProductHandler
func (h *productHandler) GetAllProductImagesByProductId(w http.ResponseWriter, r *http.Request) {
	urlRole := chi.URLParam(r, "role")
	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	response, err := h.Service.FindAllProductImagesByProductId(productIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)
}

// UploadProductImage implements ProductHandler
func (h *productHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	var req model.ProductImagesUploadReq

	urlRole := chi.URLParam(r, "role")
	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	isPrimary := r.FormValue("is_primary")

	uploadedFile, header, err := r.FormFile("file-image")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := "pringgodigdo.com/" + filename

	req.IsPrimary = isPrimary

	response, err := h.Service.UploadProductImages(req, productIdInt, filePath)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	out, err := os.Create("files/" + filename)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, uploadedFile)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)

}

// DeleteProductImage implements ProductHandler
func (h *productHandler) DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	urlRole := chi.URLParam(r, "role")

	productId := chi.URLParam(r, "product-id")
	productIdInt, _ := strconv.Atoi(productId)

	productImageId := chi.URLParam(r, "product-image-id")
	productImageIdInt, _ := strconv.Atoi(productImageId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	role := userInfo["role"].(string)

	if urlRole != role {
		WriteErrorResponse(w, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	response, err := h.Service.DeleteProductImageId(productImageIdInt, productIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)
}
