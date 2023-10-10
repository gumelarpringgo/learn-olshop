package handler

import (
	"encoding/json"
	"learn/common"
	"learn/model"
	"learn/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type ProductHandler interface {
	AddProduct(w http.ResponseWriter, r *http.Request)
	FindProductById(w http.ResponseWriter, r *http.Request)
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
	var req model.AddProductReq

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
