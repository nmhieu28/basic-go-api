package errors

type GeneralErrorValue int

const (
	Success       GeneralErrorValue = 0 + iota
	DatabaseError GeneralErrorValue = 501
	DataInvalid   GeneralErrorValue = 502
)

type AppError interface {
	Error() string
	GetCode() int
	GetMessage(code int) string
}

type GeneralError struct {
	Code GeneralErrorValue
}

func NewGeneralError(code GeneralErrorValue) AppError {
	return &GeneralError{Code: code}
}
func (e *GeneralError) Error() string {
	return GeneralMessage[e.Code]
}

func (e *GeneralError) GetCode() int {
	return int(e.Code)
}
func (e *GeneralError) GetMessage(code int) string {
	return GeneralMessage[GeneralErrorValue(code)]
}

var GeneralMessage = map[GeneralErrorValue]string{
	Success:       "Successfully!",
	DatabaseError: "DB Error",
	DataInvalid:   "Data is invalid",
}
