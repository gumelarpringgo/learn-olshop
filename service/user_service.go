package service

import (
	"errors"
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

	// PUBLIC
	emptyRegisRes      = model.RegisterRes{}
	emptyLoginRes      = model.LoginRes{}
	emptyProfileRes    = model.ProfileRes{}
	emptyChangePassRes = model.ChangePassRes{}
	// ADMIN
	emptyRegisAdminRes = model.RegisterAdminRes{}
)

var (
	errorUserNotFound = errors.New("user not found")
)

// Register implements UserServive
func (s *userService) Register(req model.RegisterReq) (model.RegisterRes, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return emptyRegisRes, errors.New("failed generate pasword")
	}

	username, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return emptyRegisRes, err
	}

	if username.Id != 0 {
		return emptyRegisRes, errors.New("username already exists")
	}

	userEmail, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return emptyRegisRes, err
	}

	if userEmail.Id != 0 {
		return emptyRegisRes, errors.New("email already exists")
	}

	dbUser.Username = req.Username
	dbUser.Email = req.Email
	dbUser.Password = string(passHash)
	dbUser.Role = "user"

	user, err := s.Repo.CreateUser(dbUser)
	if err != nil {
		return emptyRegisRes, err
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
		return emptyLoginRes, err
	}

	if user.Id == 0 {
		return emptyLoginRes, errors.New("username incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return emptyLoginRes, errors.New("password incorrect")
	}

	token, err := config.CreateToken(user.Id, user.Role)
	if err != nil {
		return emptyLoginRes, err
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
		return emptyProfileRes, err
	}

	if user.Id == 0 {
		return emptyProfileRes, errorUserNotFound
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
		return emptyChangePassRes, errors.New("failed generate pasword")
	}

	user, err := s.Repo.FindByID(id)
	if err != nil {
		return emptyChangePassRes, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return emptyChangePassRes, errors.New("password incorrect")
	}

	if user.Id == 0 {
		return emptyChangePassRes, errorUserNotFound
	}

	if req.NewPassword != req.ConfirmPassword {
		return emptyChangePassRes, errors.New("new password do not match")
	}

	if req.Password == req.NewPassword {
		return emptyChangePassRes, errors.New("new password cannot be the same as old")
	}

	user.Password = string(newPass)

	_, err = s.Repo.SaveNewPassword(user)
	if err != nil {
		return emptyChangePassRes, err
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
		return emptyRegisAdminRes, errors.New("failed generate pasword")
	}

	userUsername, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return emptyRegisAdminRes, err
	}

	if userUsername.Id != 0 {
		return emptyRegisAdminRes, errors.New("email already exists")
	}

	if req.CodeAdmin != codeAdmin {
		return emptyRegisAdminRes, errors.New("code admin is wrong")
	}

	userEmail, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return emptyRegisAdminRes, err
	}

	if userEmail.Id != 0 {
		return emptyRegisAdminRes, errors.New("email already exists")
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
