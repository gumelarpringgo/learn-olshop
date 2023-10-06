package handler

import (
	"encoding/json"
	"errors"
	"learn/model"
	"learn/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AddressHandler interface {
	AddAddress(w http.ResponseWriter, r *http.Request)
	GetAddresses(w http.ResponseWriter, r *http.Request)
	UpdateAddress(w http.ResponseWriter, r *http.Request)
	DeleteAddress(w http.ResponseWriter, r *http.Request)
}

type addresshandler struct {
	Service  service.AddressService
	Validate *validator.Validate
}

func NewAddressHandler(srv service.AddressService, validate *validator.Validate) AddressHandler {
	return &addresshandler{
		Service:  srv,
		Validate: validate,
	}
}

// AddAddress implements AddressHandler
func (h *addresshandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	stringUserId := chi.URLParam(r, "user-id")
	intUserId, _ := strconv.Atoi(stringUserId)

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

	if intUserId != id {
		WriteErrorResponse(w, http.StatusInternalServerError, errors.New("user unauthorized"))
		return
	}

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
	stringUserId := chi.URLParam(r, "user-id")
	intUserId, _ := strconv.Atoi(stringUserId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["user_id"].(float64)
	id := int(user)

	if intUserId != id {
		WriteErrorResponse(w, http.StatusInternalServerError, errors.New("user unauthorized"))
		return
	}

	userAddresses, err := h.Service.GetAddresses(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, userAddresses)

}

// UpdateAddress implements AddressHandler
func (h *addresshandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var req model.AddressReq

	addressIdStr := chi.URLParam(r, "address-id")
	addressUserIdStr := chi.URLParam(r, "user-id")
	addressIdInt, _ := strconv.Atoi(addressIdStr)
	addressUserIdInt, _ := strconv.Atoi(addressUserIdStr)

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

	if addressUserIdInt != id {
		WriteErrorResponse(w, http.StatusInternalServerError, errors.New("user unauthorized"))
		return
	}

	req.UserId = id

	addressRes, err := h.Service.UpdateAddress(req, addressIdInt)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, addressRes)
}

// DeleteAddress implements AddressHandler
func (h *addresshandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	stringUserId := chi.URLParam(r, "user-id")
	intUserId, _ := strconv.Atoi(stringUserId)

	stringAddressId := chi.URLParam(r, "address-id")
	intAddressId, _ := strconv.Atoi(stringAddressId)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["user_id"].(float64)
	id := int(user)

	if intUserId != id {
		WriteErrorResponse(w, http.StatusInternalServerError, errors.New("user unauthorized"))
		return
	}

	response, err := h.Service.DeleteAddress(intAddressId)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, response)
}
