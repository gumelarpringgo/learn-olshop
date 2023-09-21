package repository

import (
	"errors"
	"learn/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user model.User) (model.User, error)
	FindByID(id int) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByUsername(username string) (model.User, error)
	SaveNewPassword(user model.User) (model.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

var (
	SuccChangePass = "change password successfully"
)

var (
	ErrCreateUser     = errors.New("failed create user")
	ErrUserNotFound   = errors.New("user not found")
	ErrUpdatePassword = errors.New("failed update password")
)

var (
	emptyUser = model.User{}
)

// CreateUser implements UserRepository
func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.DB.Create(&user).Error
	if err != nil {
		return emptyUser, ErrCreateUser
	}

	return user, nil
}

// FindByID implements UserRepository
func (r *userRepository) FindByID(id int) (model.User, error) {
	dbUser := model.User{}

	if id == 0 {
		return emptyUser, ErrUserNotFound
	}

	err := r.DB.Where("id = ?", id).Find(&dbUser).Error
	if err != nil {
		return emptyUser, ErrUserNotFound
	}

	return dbUser, nil
}

// FindByEmail implements UserRepository
func (r *userRepository) FindByEmail(email string) (model.User, error) {
	dbUser := model.User{}

	err := r.DB.Where("email = ?", email).Find(&dbUser).Error
	if err != nil {
		return emptyUser, ErrUserNotFound
	}

	return dbUser, nil
}

// FindByUsername implements UserRepository
func (r *userRepository) FindByUsername(username string) (model.User, error) {
	dbUser := model.User{}

	err := r.DB.Where("username = ?", username).Find(&dbUser).Error
	if err != nil {
		return emptyUser, ErrUserNotFound
	}

	return dbUser, nil
}

// SaveNewPassword implements UserRepository
func (r *userRepository) SaveNewPassword(user model.User) (model.User, error) {
	err := r.DB.Save(&user).Error
	if err != nil {
		return emptyUser, ErrUpdatePassword
	}

	return user, nil
}
