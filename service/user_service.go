package service

import (
	"fmt"
	"learn/common"
	"learn/config"
	"learn/model"
	"learn/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserServive interface {
	// PUBLIC
	Register(req model.RegisterReq) (model.RegisterRes, error)
	Login(req model.LoginReq) (model.LoginRes, error)
	Profile(id int) (model.ProfileRes, error)
	ChangePassword(id int, req model.ChangePassReq) (model.ChangePassRes, error)
	// ADMIN
	RegisterAdmin(req model.RegisterAdminReq) (model.RegisterAdminRes, error)
}

type userService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserServive {
	return &userService{
		Repo: repo,
	}
}

var (
	dbUser = model.User{}

	// USER
	emptyRegisRes      = model.RegisterRes{}
	emptyLoginRes      = model.LoginRes{}
	emptyProfileRes    = model.ProfileRes{}
	emptyChangePassRes = model.ChangePassRes{}
	// ADMIN
	emptyRegisAdminRes = model.RegisterAdminRes{}
)

// Register implements UserServive
func (s *userService) Register(req model.RegisterReq) (model.RegisterRes, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return emptyRegisRes, fmt.Errorf("GenerateFromPassword call failed: %w", err)
	}

	username, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return emptyRegisRes, fmt.Errorf("FindByUsername call failed: %w", err)
	}

	if username.Id != 0 {
		return emptyRegisRes, fmt.Errorf("user id %d : %w", username.Id, common.ErrNotFound)
	}

	userEmail, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return emptyRegisRes, fmt.Errorf("FindByEmail call failed: %w", err)
	}

	if userEmail.Id != 0 {
		return emptyRegisRes, fmt.Errorf("user email %d : %w", userEmail.Id, common.ErrNotFound)
	}

	dbUser.Username = req.Username
	dbUser.Email = req.Email
	dbUser.Password = string(passHash)
	dbUser.Role = "user"

	user, err := s.Repo.CreateUser(dbUser)
	if err != nil {
		return emptyRegisRes, fmt.Errorf("CreateUser call failed: %w", err)
	}

	response := model.RegisterRes{
		Username: user.Username,
	}

	return response, nil
}

// Login implements UserServive
func (s *userService) Login(req model.LoginReq) (model.LoginRes, error) {
	user, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return emptyLoginRes, fmt.Errorf("FindByUsername call failed: %w", err)
	}
	if user.Id == 0 {
		return emptyLoginRes, fmt.Errorf("user id %d : %w", user.Id, common.ErrNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return emptyLoginRes, fmt.Errorf("CompareHashAndPassword call failed: %w", err)
	}

	token, err := config.CreateToken(user.Id, user.Role)
	if err != nil {
		return emptyLoginRes, fmt.Errorf("CreateToken call failed: %w", err)
	}

	response := model.LoginRes{
		Token: token,
	}

	return response, nil
}

// Profile implements UserServive
func (s *userService) Profile(id int) (model.ProfileRes, error) {
	user, err := s.Repo.FindByID(id)
	if err != nil {
		return emptyProfileRes, fmt.Errorf("FindByID call failed: %w", err)
	}

	if user.Id == 0 {
		return emptyProfileRes, fmt.Errorf("user id %d : %w", user.Id, common.ErrNotFound)
	}

	response := model.ProfileRes{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	return response, nil
}

// ChangePassword implements UserServive
func (s *userService) ChangePassword(id int, req model.ChangePassReq) (model.ChangePassRes, error) {
	newPass, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return emptyChangePassRes, fmt.Errorf("GenerateFromPassword call failed: %w", err)
	}

	user, err := s.Repo.FindByID(id)
	if err != nil {
		return emptyChangePassRes, fmt.Errorf("FindByID call failed: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return emptyChangePassRes, fmt.Errorf("CompareHashAndPassword call failed: %w", err)
	}

	if user.Id == 0 {
		return emptyChangePassRes, fmt.Errorf("user id %d : %w", user.Id, common.ErrNotFound)
	}
	if req.NewPassword != req.ConfirmPassword {
		return emptyChangePassRes, fmt.Errorf("user passwored : %w", common.ErrNotMatch)
	}

	if req.Password == req.NewPassword {
		return emptyChangePassRes, fmt.Errorf("user passwored : %w", common.ErrNotMatch)
	}

	user.Password = string(newPass)

	_, err = s.Repo.SaveNewPassword(user)
	if err != nil {
		return emptyChangePassRes, fmt.Errorf("SaveNewPassword call failed: %w", err)
	}

	response := model.ChangePassRes{
		ChangePassword: "changed password successfully",
	}

	return response, nil
}

// RegisterAdmin implements UserServive
func (s *userService) RegisterAdmin(req model.RegisterAdminReq) (model.RegisterAdminRes, error) {
	codeAdmin := 181910

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return emptyRegisAdminRes, fmt.Errorf("GenerateFromPassword call failed: %w", err)
	}

	userUsername, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return emptyRegisAdminRes, fmt.Errorf("FindByUsername call failed: %w", err)
	}

	if userUsername.Id != 0 {
		return emptyRegisAdminRes, fmt.Errorf("admin email : %w", common.ErrExists)
	}

	if req.CodeAdmin != codeAdmin {
		return emptyRegisAdminRes, fmt.Errorf("admin code : %w", common.ErrNotMatch)
	}

	userEmail, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return emptyRegisAdminRes, fmt.Errorf("FindByEmail call failed: %w", err)
	}

	if userEmail.Id != 0 {
		return emptyRegisAdminRes, fmt.Errorf("admin email : %w", common.ErrExists)
	}

	dbUser.Username = req.Username
	dbUser.Email = req.Email
	dbUser.Password = string(passHash)
	dbUser.Role = "admin"

	newUser, err := s.Repo.CreateUser(dbUser)
	if err != nil {
		return emptyRegisAdminRes, err
	}

	response := model.RegisterAdminRes{
		Username: newUser.Username,
	}

	return response, nil
}
