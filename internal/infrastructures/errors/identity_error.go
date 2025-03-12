package errors

import app_errors "backend/pkg/errors"

type IdentityErrorValue int

const (
	EmailNotFound IdentityErrorValue = 1000 + iota
	CanNotHashPassword
	EmailExisted
	PasswordInvalid
	JWTError
	EmailAlreadyConfirmed
	UserNotFound
	OTPInvalid
	EmailNotConfirmed
)

func NewIdentityError(code IdentityErrorValue) app_errors.AppError {
	return &IdentityError{Code: code}
}

type IdentityError struct {
	Code IdentityErrorValue
}

func (e *IdentityError) Error() string {
	return IdentityMessage[e.Code]
}

func (e *IdentityError) GetCode() int {
	return int(e.Code)
}
func (e *IdentityError) GetMessage(code int) string {
	return IdentityMessage[IdentityErrorValue(code)]
}

var IdentityMessage = map[IdentityErrorValue]string{
	EmailNotFound:         "Email is not exists",
	CanNotHashPassword:    "Can't hash password",
	EmailExisted:          "Email is exsits",
	PasswordInvalid:       "Password is invalid",
	JWTError:              "Error generating JWT token",
	EmailAlreadyConfirmed: "Email already confirmed",
	OTPInvalid:            "OTP is invalid",
	EmailNotConfirmed:     "Email is not confirmed",
}
