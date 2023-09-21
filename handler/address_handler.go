package handler

import (
	"encoding/json"
	"learn/model"
	"learn/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AddressHandler interface {
	AddAddress(w http.ResponseWriter, r *http.Request)
	GetAddresses(w http.ResponseWriter, r *http.Request)
}

type addresshandler struct {
	Service  service.AddressService
	Validate *validator.Validate

	// ChangePassword implements
}

func NewAddressHandler(srv service.AddressService, validate *validator.Validate) AddressHandler {
	return &addresshandler{
		Service:  srv,
		Validate: validate,
	}
}

// AddAddress implements AddressHandler
func (h *addresshandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	var req model.AddressReq

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

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["user_id"].(float64)
	id := int(user)

	req.UserId = id

	resAddress, err := h.Service.AddAddress(req, id)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, resAddress)
}

// GetAddresses implements AddressHandler
func (h *addresshandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["user_id"].(float64)
	id := int(user)

	userAddresses, err := h.Service.GetAddresses(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, userAddresses)

}
