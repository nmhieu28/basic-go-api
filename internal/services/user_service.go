package services

import (
	identity_errors "backend/internal/infrastructures/errors"
	"backend/internal/infrastructures/repositories"
	"backend/internal/models/responses"
	app_errors "backend/pkg/errors"
	"backend/pkg/logger"
	"backend/pkg/response"
	"context"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
	logger   logger.Logger
}

func NewUserService(userRepo repositories.UserRepository, logger logger.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}

func (s *UserService) GetUser(id uuid.UUID, ctx context.Context) *response.Response[*responses.UserResponse] {
	user, err := s.userRepo.GetByID(id, ctx)
	if err != nil {
		return response.FailureWithData[*responses.UserResponse](nil, app_errors.NewGeneralError(app_errors.DatabaseError))
	}
	if user == nil {
		return response.FailureWithData[*responses.UserResponse](nil, identity_errors.NewIdentityError(identity_errors.UserNotFound))
	}
	return response.Success(&responses.UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		FullName:    user.FullName(),
		Avatar:      user.Avatar,
		DateOfBirth: user.DateOfBirth,
	})
}
