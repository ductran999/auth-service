package account

import "errors"

var (
	ErrEmailExisted           = errors.New("email already registered")
	ErrAccountNotFound        = errors.New("account not found")
	ErrPasswordMismatch       = errors.New("password mismatch")
	ErrNewPasswordMustChanged = errors.New("new password must e different")
)
