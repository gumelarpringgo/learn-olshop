package service

import (
	"errors"
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
	// address = model.Address{}

	emptyAddressRes         = model.AddressRes{}
	emptyAddressesRes       = []model.AddressRes{}
	emptyAddressWithoutData = model.AddressResWithoutData{}
)

var (
	errAddMustHavePrimary = errors.New("address must have primary")
	errAddressNotFound    = errors.New("address not found")
	errAddressNotOwner    = errors.New("address not owner")
)

// CreateAddress implements AddressService
func (s *serviceAddress) AddAddress(req model.AddressReq, userId int) (model.AddressRes, error) {
	address := model.Address{}
	isPrimary := "no"

	arrayAddress, err := s.Repo.FindByUserId(userId)
	if err != nil {
		return emptyAddressRes, err
	}

	if len(arrayAddress) == 0 && req.IsPrimary != true {
		return emptyAddressRes, errAddMustHavePrimary
	} else if len(arrayAddress) >= 1 && req.IsPrimary {
		isPrimary = "yes"

		_, err := s.Repo.MarkAllAddressNonPrimary(userId)
		if err != nil {
			return emptyAddressRes, err
		}
	} else if req.IsPrimary {
		isPrimary = "yes"
	}

	address.Address = req.Address
	address.IsPrimary = isPrimary
	address.UserId = req.UserId

	addressDB, err := s.Repo.Create(address)
	if err != nil {
		return emptyAddressRes, err
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
		return emptyAddressesRes, err
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
		return emptyAddressRes, err
	}

	if address.Id == 0 {
		return emptyAddressRes, errAddressNotFound
	}

	if req.UserId != address.UserId {
		return emptyAddressRes, errAddressNotOwner
	}

	arrayAddress, err := s.Repo.FindByUserId(req.UserId)
	if err != nil {
		return emptyAddressRes, err
	}

	for _, addrs := range arrayAddress {
		if req.IsPrimary {
			isPrimary = "yes"

			_, err := s.Repo.MarkAllAddressNonPrimary(addrs.UserId)
			if err != nil {
				return emptyAddressRes, err
			}
		} else {
			isPrimary = "no"
		}

		if req.IsPrimary == false && addrs.IsPrimary == "yes" {
			return emptyAddressRes, errAddMustHavePrimary
		}
	}

	address.Address = req.Address
	address.IsPrimary = isPrimary
	address.UserId = req.UserId

	updateAddress, err := s.Repo.Update(address)
	if err != nil {
		return emptyAddressRes, err
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
		return emptyAddressWithoutData, errAddressNotFound
	}

	if address.IsPrimary == "yes" {
		return emptyAddressWithoutData, errAddMustHavePrimary
	}

	err = s.Repo.Delete(addressId)
	if err != nil {
		return emptyAddressWithoutData, err
	}

	response := model.AddressResWithoutData{
		Message: "delete address successfully",
	}

	return response, nil
}
