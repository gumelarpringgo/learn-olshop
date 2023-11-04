package service

import (
	"fmt"
	"learn/common"
	"learn/model"
	"learn/repository"
)

type AddressService interface {
	AddAddress(req model.AddressReq, userId int) (model.AddressRes, error)
	GetAddresses(userId int) ([]model.AddressRes, error)
	UpdateAddress(req model.AddressReq, addressId int) (model.AddressRes, error)
	DeleteAddress(addressId int) (model.AddressResWithoutData, error)
}

type serviceAddress struct {
	Repo repository.AddressRepository
}

func NewAddressService(srv *repository.AddressRepository) AddressService {
	return &serviceAddress{
		Repo: *srv,
	}
}

var (
	emptyAddressRes         = model.AddressRes{}
	emptyAddressesRes       = []model.AddressRes{}
	emptyAddressWithoutData = model.AddressResWithoutData{}
)

// CreateAddress implements AddressService
func (s *serviceAddress) AddAddress(req model.AddressReq, userId int) (model.AddressRes, error) {
	address := model.Address{}
	isPrimary := "no"

	arrayAddress, err := s.Repo.FindByUserId(userId)
	if err != nil {
		return emptyAddressRes, fmt.Errorf("FindByUserId call failed: %w", err)
	}

	if len(arrayAddress) == 0 && req.IsPrimary != true {
		return emptyAddressRes, fmt.Errorf("address: %w", err)
	} else if len(arrayAddress) >= 1 && req.IsPrimary {
		isPrimary = "yes"

		_, err := s.Repo.MarkAllAddressNonPrimary(userId)
		if err != nil {
			return emptyAddressRes, fmt.Errorf("MarkAllAddressNonPrimary call failed: %w", err)
		}
	} else if req.IsPrimary {
		isPrimary = "yes"
	}

	address.Address = req.Address
	address.IsPrimary = isPrimary
	address.UserId = req.UserId

	addressDB, err := s.Repo.Create(address)
	if err != nil {
		return emptyAddressRes, fmt.Errorf("create call failed: %w", err)
	}

	var resIsPrimary bool
	if addressDB.IsPrimary == "yes" {
		resIsPrimary = true
	} else {
		resIsPrimary = false
	}

	response := model.AddressRes{
		Address:   addressDB.Address,
		IsPrimary: resIsPrimary,
		UserId:    addressDB.UserId,
	}

	return response, nil
}

// GetAddresses implements AddressService
func (s *serviceAddress) GetAddresses(userId int) ([]model.AddressRes, error) {
	formatAddresses := []model.AddressRes{}

	arrayAddress, err := s.Repo.FindByUserId(userId)
	if err != nil {
		return emptyAddressesRes, fmt.Errorf("address user id %d : %w", userId, common.ErrNotFound)
	}

	for _, addr := range arrayAddress {
		var resIsPrimary bool

		if addr.IsPrimary == "yes" {
			resIsPrimary = true
		} else {
			resIsPrimary = false
		}

		formatAddress := model.AddressRes{
			Address:   addr.Address,
			IsPrimary: resIsPrimary,
			UserId:    addr.UserId,
		}

		formatAddresses = append(formatAddresses, formatAddress)
	}

	return formatAddresses, nil
}

// UpdateAddress implements AddressService
func (s *serviceAddress) UpdateAddress(req model.AddressReq, addressId int) (model.AddressRes, error) {
	isPrimary := "no"

	address, err := s.Repo.FindByAddressId(addressId)
	if err != nil {
		return emptyAddressRes, fmt.Errorf("FindByAddressId call failed: %w", err)
	}

	if address.Id == 0 {
		return emptyAddressRes, fmt.Errorf("address %d : %w", address.Id, common.ErrNotFound)
	}

	if req.UserId != address.UserId {
		return emptyAddressRes, fmt.Errorf("address user %d : %w", req.UserId, common.ErrNotFound)
	}

	arrayAddress, err := s.Repo.FindByUserId(req.UserId)
	if err != nil {
		return emptyAddressRes, fmt.Errorf("FindByUserId call failed: %w", err)
	}

	for _, addrs := range arrayAddress {
		if req.IsPrimary {
			isPrimary = "yes"

			_, err := s.Repo.MarkAllAddressNonPrimary(addrs.UserId)
			if err != nil {
				return emptyAddressRes, fmt.Errorf("MarkAllAddressNonPrimary call failed: %w", err)
			} else {
				isPrimary = "no"
			}
		}
		if req.IsPrimary == false && addrs.IsPrimary == "yes" {
			return emptyAddressRes, fmt.Errorf("address: %w", err)
		}
	}

	address.Address = req.Address
	address.IsPrimary = isPrimary
	address.UserId = req.UserId

	updateAddress, err := s.Repo.Update(address)
	if err != nil {
		return emptyAddressRes, fmt.Errorf("update call failed: %w", err)
	}

	var resIsPrimary bool
	if updateAddress.IsPrimary == "yes" {
		resIsPrimary = true
	} else {
		resIsPrimary = false
	}

	response := model.AddressRes{
		Address:   updateAddress.Address,
		IsPrimary: resIsPrimary,
		UserId:    updateAddress.UserId,
	}

	return response, nil
}

// DeleteAddress implements AddressService
func (s *serviceAddress) DeleteAddress(addressId int) (model.AddressResWithoutData, error) {
	address, err := s.Repo.FindByAddressId(addressId)

	if err != nil {
		return emptyAddressWithoutData, fmt.Errorf("FindByAddressId call failed: %w", err)
	}

	if address.IsPrimary == "yes" {
		return emptyAddressWithoutData, fmt.Errorf("address : %w", common.ErrExists)
	}

	err = s.Repo.Delete(addressId)
	if err != nil {
		return emptyAddressWithoutData, fmt.Errorf("delete call failed: %w", err)
	}

	response := model.AddressResWithoutData{
		Message: "delete address successfully",
	}

	return response, nil
}
