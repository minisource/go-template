package usecase

import (
	"context"
	"errors"

	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/domain/model"
	"github.com/minisource/template_go/domain/repository"
	auth "github.com/minisource/auth/service"
	"github.com/minisource/auth/service/models"
	"github.com/minisource/go-common/logging"
)

type UserUsecase struct {
	logger      logging.Logger
	cfg         *config.Config
	authService *auth.AuthService
	repository  repository.UserRepository
}

func NewUserUsecase(cfg *config.Config, repository repository.UserRepository) *UserUsecase {
	logger := logging.NewLogger(&cfg.Logger)
	return &UserUsecase{
		cfg:         cfg,
		repository:  repository,
		logger:      logger,
		authService: auth.GetAuthService(),
	}
}

func (u UserUsecase) SendOtpByMobileNumber(countryCode, mobileNumber string) error {
	if countryCode == "" {
		countryCode = "+98"
	}
	phone := countryCode + mobileNumber
	err := u.authService.SendOTP(phone)
	if err != nil {
		return err
	}
	return nil
}

// Register/login by mobile number
func (u *UserUsecase) RegisterAndLoginByMobileNumber(ctx context.Context, countryCode, mobileNumber string, otp string) (*models.AccessTokenResponse, error) {
	if countryCode == "" {
		countryCode = "+98"
	}
	phone := countryCode + mobileNumber
	
	// verify otp
	result, err := u.authService.VerifyCode(phone, otp)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("invalid OTP")
	}

	// register and get user
	user, err := u.authService.GetUserInfoByPhone(phone)
	if err != nil {
		return nil, err
	}

	userDb := model.User{UserId: user.Id}
	exists, err := u.repository.ExistsUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err = u.repository.CreateUser(ctx, userDb)
		if err != nil {
			return nil, err
		}
	}

	return u.authService.GenerateJWT(user.Name)
}
