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
		UserId int `json:"user_id"`
	}
)

// Response
type (
	AddressRes struct {
		Address   string `json:"address"`
		IsPrimary bool   `json:"is_primary"`
		UserId    int    `json:"user_id"`
	}

	AddressesRes struct {
		Addresses []AddressRes `json:"addresses"`
	}
)

func FormatAddresses(arrayAddress []Address) AddressesRes {
	addressesRes := []AddressRes{}
	addrsRes := AddressRes{}
	for _, address := range arrayAddress {

		resAddr := addrsRes
		resAddr.Address = address.Address

		resAddsPrimary := false
		if address.IsPrimary == "yes" {
			resAddsPrimary = true
		} else {
			resAddsPrimary = false
		}

		resAddr.IsPrimary = resAddsPrimary
		resAddr.UserId = address.UserId

		addressesRes = append(addressesRes, resAddr)
	}
	return AddressesRes{}
}
