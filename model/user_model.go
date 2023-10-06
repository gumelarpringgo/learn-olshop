package model

import "time"

// DATABASE
type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// REQUEST
type (
	RegisterReq struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	LoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ChangePassReq struct {
		Password        string `json:"password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}

	RegisterAdminReq struct {
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
		CodeAdmin int    `json:"code_admin"`
	}
)

// RESPONSE
type (
	RegisterRes struct {
		Username string `json:"username"`
	}

	LoginRes struct {
		Token string `json:"token"`
	}

	ProfileRes struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	ChangePassRes struct {
		ChangePassword string `json:"change_password"`
	}

	RegisterAdminRes struct {
		Username string `json:"username"`
	}
)
