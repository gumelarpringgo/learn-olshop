package repository

import (
	"errors"
	"learn/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address model.Address) (model.Address, error)
	MarkAllAddressNonPrimary(addresId int) (bool, error)
	GetAllAddaresses(userId int) ([]model.Address, error)
}

type addressRepository struct {
	DB *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{
		DB: db,
	}
}

var (
	errCreateAddress = errors.New("failed create address")
	errArrayAddress  = errors.New("failed get addresses")
)

// Create implements AddressRepository
func (r *addressRepository) Create(address model.Address) (model.Address, error) {
	err := r.DB.Create(&address).Error
	if err != nil {
		return address, errCreateAddress
	}

	return address, nil
}

// MarkAllAddressNonPrimary implements AddressRepository
func (r *addressRepository) MarkAllAddressNonPrimary(UserId int) (bool, error) {
	err := r.DB.Model(&model.Address{}).Where("user_id = ?", UserId).Update("is_primary", "no").Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetAllAddaress implements AddressRepository
func (r *addressRepository) GetAllAddaresses(userId int) ([]model.Address, error) {
	addresses := []model.Address{}

	err := r.DB.Where("user_id = ?", userId).Find(&addresses).Error
	if err != nil {
		return addresses, errArrayAddress
	}

	return addresses, nil
}
