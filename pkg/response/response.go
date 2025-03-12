package response

import "backend/pkg/errors"

type Response[T any] struct {
	Data      T      `json:"data"`
	Code      int    `json:"code"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

type ResponseWithPaging[T any, P any] struct {
	Data      T      `json:"data"`
	Paging    P      `json:"paging"`
	Code      string `json:"code"`
	IsSuccess bool   `json:"isSuccess"`
}

func generate[T any](data T, isSuccess bool, err errors.AppError) *Response[T] {
	code := err.GetCode()
	message := err.GetMessage(code)
	return &Response[T]{Data: data, Code: code, IsSuccess: isSuccess, Message: message}
}

func Success[T any](data T) *Response[T] {
	return generate(data, true, errors.NewGeneralError(errors.Success))
}

func Failure(err errors.AppError) *Response[bool] {
	return generate(false, false, err)
}

func FailureWithData[T any](data T, err errors.AppError) *Response[T] {
	return generate(data, false, err)
}
