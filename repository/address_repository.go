package repository

import (
	"errors"
	"fmt"
	"learn/common"
	"learn/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address model.Address) (model.Address, error)
	MarkAllAddressNonPrimary(addresId int) (bool, error)
	FindByUserId(userId int) ([]model.Address, error)
	FindByAddressId(addressId int) (model.Address, error)
	Update(address model.Address) (model.Address, error)
	Delete(addressId int) error
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
	emptyAddress   = model.Address{}
	emptyAddresses = []model.Address{}
)

// Create implements AddressRepository
func (r *addressRepository) Create(address model.Address) (model.Address, error) {
	err := r.DB.Create(&address).Error
	if err != nil {
		if errors.Is(err, common.ErrFailedCreateData) {
			return emptyAddress, fmt.Errorf("address: %w", common.ErrFailedCreateData)
		}
	}

	return address, nil
}

// MarkAllAddressNonPrimary implements AddressRepository
func (r *addressRepository) MarkAllAddressNonPrimary(UserId int) (bool, error) {
	err := r.DB.Model(&model.Address{}).Where("user_id = ?", UserId).Update("is_primary", "no").Error
	if err != nil {
		if errors.Is(err, common.ErrMustHavePrimary) {
			return false, fmt.Errorf("address primary: %w", common.ErrMustHavePrimary)
		}
	}

	return true, nil
}

// GetAllAddaress implements AddressRepository
func (r *addressRepository) FindByUserId(userId int) ([]model.Address, error) {
	addresses := []model.Address{}
	err := r.DB.Where("user_id = ?", userId).Find(&addresses).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return emptyAddresses, fmt.Errorf("address user id %d: %w", userId, common.ErrNotFound)
		}
	}

	return addresses, nil
}

// FindByAddressId implements AddressRepository
func (r *addressRepository) FindByAddressId(addressId int) (model.Address, error) {
	address := model.Address{}
	err := r.DB.Where("id = ?", addressId).Find(&address).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return emptyAddress, fmt.Errorf("address id %d: %w", addressId, common.ErrNotFound)
		}
	}

	return address, nil
}

// Create implements AddressRepository
func (r *addressRepository) Update(address model.Address) (model.Address, error) {
	err := r.DB.Save(&address).Error
	if err != nil {
		if errors.Is(err, common.ErrFailedUpdateData) {
			return address, fmt.Errorf("address : %w", common.ErrFailedUpdateData)
		}
	}

	return address, nil
}

// Delete implements AddressRepository
func (r *addressRepository) Delete(addressId int) error {
	err := r.DB.Delete(&model.Address{}, addressId).Error
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("address %d: %w", addressId, common.ErrNotFound)
		}
	}

	return nil
}
