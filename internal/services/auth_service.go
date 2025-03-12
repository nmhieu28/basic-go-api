package services

import (
	"context"
	"errors"
	"fmt"

	"backend/email_template"
	"backend/internal/infrastructures/entities"
	identity_errors "backend/internal/infrastructures/errors"
	"backend/internal/infrastructures/repositories"
	"backend/internal/models/requests"
	"backend/internal/models/responses"
	"backend/pkg/cache"
	configs "backend/pkg/config"
	"backend/pkg/entity"
	app_errors "backend/pkg/errors"
	"backend/pkg/jwt_generate"
	"backend/pkg/logger"
	"backend/pkg/mailer"
	"backend/pkg/response"
	"backend/pkg/utils"

	"gorm.io/gorm"
)

type IdentityService struct {
	identityRepo repositories.UserRepository
	redisCache   cache.Cache `name:"redis_identity"`
	logger       logger.Logger
	mailer       mailer.Mailer
	appSetting   *configs.AppConfig
	jwtGen       jwt_generate.JwtGenerate
}

var (
	max_time_verify_otp = 3600
)

func NewIdentityService(identityRepo repositories.UserRepository,
	redisCache cache.Cache,
	logger logger.Logger,
	mailer mailer.Mailer,
	appSetting *configs.AppConfig,
	jwtGen jwt_generate.JwtGenerate,
) *IdentityService {

	return &IdentityService{identityRepo: identityRepo, redisCache: redisCache, logger: logger, mailer: mailer, appSetting: appSetting, jwtGen: jwtGen}
}

func (s *IdentityService) Register(ctx context.Context, request requests.CreateUserRequest) (bool, error) {
	user, err := s.identityRepo.FindByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, app_errors.NewGeneralError(app_errors.DatabaseError)
	}

	if user != nil {
		return false, identity_errors.NewIdentityError(identity_errors.EmailExisted)
	}

	passwordHash, err := utils.HashPassword(request.Password)

	if err != nil {
		return false, identity_errors.NewIdentityError(identity_errors.CanNotHashPassword)
	}

	newUser := &entities.User{
		BaseAuditTrackingEntity: entity.NewSQLModel(),
		Email:                   request.Email,
		PasswordHash:            passwordHash,
		EmailConfirm:            false,
		FirstName:               request.FirstName,
		LastName:                request.LastName,
	}
	user, err = s.identityRepo.Create(newUser, ctx)

	if err != nil {
		return false, app_errors.NewGeneralError(app_errors.DatabaseError)
	}
	token, _ := s.jwtGen.GenerateVerifyEmailToken(&jwt_generate.TokenPayload{
		UserId: user.Id,
		Email:  user.Email,
	})
	confirmUrl := fmt.Sprintf("%s/account/verify-account?token=%s", s.appSetting.ServiceUrl.Frontend, token)
	template, err := email_template.LoadTemplate(email_template.CONFIRM_ACCOUNT, &email_template.ConfirmAccountData{
		ConfirmationURL: confirmUrl,
	})

	if err != nil {
		s.logger.WithContext(ctx).Error("Cant not load Email Template")
	}

	if err := s.mailer.SendHTML(ctx, user.Email, "Welcome to AppName - Verify Your Account", template); err != nil {
		s.logger.WithContext(ctx).Error("Cant not send email")
	}
	return true, nil
}

func (s *IdentityService) Login(ctx context.Context, request requests.LoginRequest) *response.Response[*responses.AuthenResponse] {
	s.logger.WithContext(ctx).Info("Login", request)
	user, err := s.identityRepo.FindByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.FailureWithData[*responses.AuthenResponse](nil, app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	if user == nil {
		return response.FailureWithData[*responses.AuthenResponse](nil, identity_errors.NewIdentityError(identity_errors.EmailNotFound))
	}

	if err := utils.VerifyPassword(user.PasswordHash, request.Password); err != nil {
		return response.FailureWithData[*responses.AuthenResponse](nil, identity_errors.NewIdentityError(identity_errors.PasswordInvalid))
	}

	if !user.EmailConfirm {
		return response.FailureWithData[*responses.AuthenResponse](nil, identity_errors.NewIdentityError(identity_errors.EmailNotConfirmed))
	}

	token, err := s.jwtGen.GenerateToken(&jwt_generate.TokenPayload{
		UserId: user.Id,
		Email:  user.Email,
	})
	if err != nil {
		return response.FailureWithData[*responses.AuthenResponse](nil, identity_errors.NewIdentityError(identity_errors.JWTError))
	}

	refreshToken, err := s.jwtGen.GenerateRefreshToken(&jwt_generate.TokenPayload{
		UserId: user.Id,
		Email:  user.Email,
	})
	if err != nil {
		return response.FailureWithData[*responses.AuthenResponse](nil, identity_errors.NewIdentityError(identity_errors.JWTError))
	}

	return response.Success(&responses.AuthenResponse{
		AccessToken:        token,
		RefreshToken:       refreshToken,
		TokenExpire:        int64(s.appSetting.Jwt.TokenExpire),
		RefreshTokenExpire: int64(s.appSetting.Jwt.RefreshTokenExpire),
	})
}

func (s *IdentityService) VerifyEmail(ctx context.Context, token string) *response.Response[bool] {
	payload, err := s.jwtGen.VerifyToken(token, s.appSetting.Jwt.VerifyEmailSecretKey)

	if err != nil {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.JWTError))
	}
	user, err := s.identityRepo.FindByEmail(ctx, payload.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	if user == nil {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.EmailNotFound))
	}

	if user.EmailConfirm {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.EmailAlreadyConfirmed))
	}
	user.EmailConfirm = true
	if err = s.identityRepo.Update(user, ctx); err != nil {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	return response.Success(true)
}

func (s *IdentityService) ResetPassword(ctx context.Context, request requests.ResetPasswordRequest) *response.Response[bool] {
	user, err := s.identityRepo.FindByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	if user == nil {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.EmailNotFound))
	}
	otp, err := s.redisCache.Get(ctx, cache.ForgotPasswordKey(user.Id))
	if err != nil {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}
	if otp != request.Code {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.OTPInvalid))
	}

	user.PasswordHash, err = utils.HashPassword(request.Password)

	if err != nil {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	if err = s.identityRepo.Update(user, ctx); err != nil {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}
	return response.Success(true)
}

func (s *IdentityService) ForgotPassword(ctx context.Context, request requests.ForgotPasswordRequest) *response.Response[bool] {
	user, err := s.identityRepo.FindByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Failure(app_errors.NewGeneralError(app_errors.DatabaseError))
	}

	if user == nil {
		return response.Failure(identity_errors.NewIdentityError(identity_errors.EmailNotFound))
	}

	otp := utils.GenerateSecureOTP()

	if err := s.redisCache.Set(ctx, cache.ForgotPasswordKey(user.Id), otp, max_time_verify_otp); err != nil {
		s.logger.WithContext(ctx).Error("Cant not set otp to redis")
	}

	template, err := email_template.LoadTemplate(email_template.FORGOT_PASSWORD, &email_template.ForgotPasswordData{
		Token: otp,
	})

	if err != nil {
		s.logger.WithContext(ctx).Error("Cant not load Email Template")
	}

	if err := s.mailer.SendHTML(ctx, user.Email, "AppName - Reset Password", template); err != nil {
		s.logger.WithContext(ctx).Error("Cant not send email")
	}

	return response.Success(true)
}
