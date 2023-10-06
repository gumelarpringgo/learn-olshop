package model

import "time"

// DATABASE
type Address struct {
	Id        int
	Address   string
	IsPrimary string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}

// Request
type (
	AddressReq struct {
		Address   string `json:"address" validate:"required"`
		IsPrimary bool   `json:"is_primary"`
		UserId    int    `json:"user_id"`
	}

	GetAllAddress struct {
		UserId int `json:"user_id" validate:"required"`
	}
)

// Response
type (
	AddressResWithoutData struct {
		Message string `json:"message"`
	}

	AddressRes struct {
		Address   string `json:"address"`
		IsPrimary bool   `json:"is_primary"`
		UserId    int    `json:"user_id"`
	}
)
