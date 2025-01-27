package grpc

import (
	"context"
	"errors"

	"github.com/Melikhov-p/sso-grpc-go/protos/gen/go/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validErr error

var EmailEmpty validErr = errors.New("email is required")
var PasswordEmpty validErr = errors.New("password is required")

type Server struct {
	log *zap.Logger
	RPC *grpc.Server
}

type ForAuth struct {
	sso.UnimplementedAuthServer
	AuthService AuthService
}

type AuthService interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int32,
	) (string, error)
	Register(ctx context.Context,
		email string,
		password string,
	) (int64, error)
}

func NewServer(log *zap.Logger) *Server {
	srv := grpc.NewServer()

	return &Server{
		log: log,
		RPC: srv,
	}
}

func (s *Server) RegisterAuthService(service AuthService) {
	sso.RegisterAuthServer(s.RPC, ForAuth{AuthService: service})
}

func (a ForAuth) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	vErr := validateEmailAndPassword(req.GetEmail(), req.GetPassword())
	if vErr != nil {
		return nil, status.Error(codes.InvalidArgument, vErr.Error())
	}
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "appID is required")
	}

	token, err := a.AuthService.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to login")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (a ForAuth) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	vErr := validateEmailAndPassword(req.GetEmail(), req.GetPassword())
	if vErr != nil {
		return nil, status.Error(codes.InvalidArgument, vErr.Error())
	}

	userID, err := a.AuthService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to register new user")
	}

	return &sso.RegisterResponse{UserId: userID}, nil
}

func validateEmailAndPassword(email string, password string) validErr {
	if email == "" {
		return EmailEmpty
	}
	if password == "" {
		return PasswordEmpty
	}

	return nil
}
