package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Melikhov-p/sso-grpc-go/internal/lib/jwt"
	"github.com/Melikhov-p/sso-grpc-go/internal/models"
	"github.com/Melikhov-p/sso-grpc-go/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserCreator interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AppProvider interface {
	GetAppByID(ctx context.Context, appID int32) (*models.App, error)
}

type AuthService struct {
	log          *zap.Logger
	userCreator  UserCreator
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func NewAuthService(log *zap.Logger, creator UserCreator, usrProvider UserProvider, appProvider AppProvider) *AuthService {
	return &AuthService{
		log:          log,
		userCreator:  creator,
		userProvider: usrProvider,
		appProvider:  appProvider,
	}
}

func (as *AuthService) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	op := "service.Auth.Login"

	var (
		app  *models.App
		user *models.User
		err  error
	)

	if user, err = as.userProvider.GetUserByEmail(ctx, email); err != nil {
		as.log.Error("error getting user by email", zap.Error(err), zap.String("EMAIL", email))

		if errors.Is(err, storage.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		as.log.Error("error compare pass hash for user", zap.Error(err), zap.Int("USERID", int(user.ID)))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if app, err = as.appProvider.GetAppByID(ctx, appID); err != nil {
		as.log.Error("error getting app with ID", zap.Error(err), zap.Int("appID", int(appID)))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.BuildJWTToken(user.ID, app, as.tokenTTL)
	if err != nil {
		as.log.Error("error building jwt Token", zap.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	as.log.Debug("user login in", zap.Int("userID", int(user.ID)), zap.String("Token", token))

	return token, nil
}

func (as *AuthService) Register(ctx context.Context, email string, password string) (int64, error) {
	op := "service.Auth.Register"

	var (
		userID   int64 // ID of new user
		passHash []byte
		err      error
	)

	passHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userID, err = as.userCreator.CreateUser(ctx, email, passHash)
	if err != nil {
		as.log.Error("error creating new user", zap.Error(err))

		return -1, fmt.Errorf("%s: %w", op, err)
	}

	as.log.Debug("created new user", zap.Int("userID", int(userID)))

	return userID, nil
}
