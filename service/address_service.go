package service

import (
	"errors"
	"learn/model"
	"learn/repository"
)

type AddressService interface {
	AddAddress(req model.AddressReq, userId int) (model.AddressRes, error)
	GetAddresses(userId int) (model.AddressesRes, error)
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
	address = model.Address{}

	emptyAddressRes   = model.AddressRes{}
	emptyAddressesRes = model.AddressesRes{}
)

var (
	errAddMustHavePrimary = errors.New("address must have primary")
)

// CreateAddress implements AddressService
func (s *serviceAddress) AddAddress(req model.AddressReq, userId int) (model.AddressRes, error) {
	isPrimary := "no"

	arrayAddress, err := s.Repo.GetAllAddaresses(userId)
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
func (s *serviceAddress) GetAddresses(userId int) (model.AddressesRes, error) {
	arrayAddress, err := s.Repo.GetAllAddaresses(userId)
	if err != nil {
		return emptyAddressesRes, err
	}

	return model.FormatAddresses(arrayAddress), nil
}
