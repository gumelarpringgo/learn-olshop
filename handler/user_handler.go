package handler

import (
	"encoding/json"
	"learn/model"
	"learn/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	Service  service.UserServive
	Validate *validator.Validate
}

// ChangePassword implements UserHandler
func (h *userHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req model.ChangePassReq

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

	passChange, err := h.Service.ChangePassword(id, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, passChange)
}

func NewUserHandler(srv service.UserServive, val *validator.Validate) UserHandler {
	return &userHandler{
		Service:  srv,
		Validate: val,
	}
}

// Register implements UserHandler
func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterReq

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

	userRgis, err := h.Service.Register(req)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, userRgis)
}

// Login implements UserHandler
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginReq

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.Validate.Struct(req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.Service.Login(req)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, token)
}

// Profile implements UserHandler
func (h *userHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["user_id"].(float64)
	id := int(user)

	userProfile, err := h.Service.Profile(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	WriteDataResponse(w, http.StatusOK, userProfile)
}
