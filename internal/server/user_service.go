package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"simple-service/generated/auth"
	"simple-service/internal/config"
	"simple-service/internal/repo"
)

// UserServiceServer реализует интерфейс gRPC.
type AuthService struct {
	auth.UnimplementedAuthServiceServer
	repo   repo.UserRepository
	stdlog *log.Logger
	zaplog *zap.Logger
}

func NewAuthService(repo repo.UserRepository, logger *zap.SugaredLogger) *AuthService {
	return &AuthService{
		repo:   repo,
		stdlog: logger,
		zaplog: logger,
	}
}

// Метод регистрации нового пользователя
func (s *AuthService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.zaplog.Info("Register request received", zap.String("username", request.Username))

	exists, err := s.repo.CheckUserExists(ctx, request.Username)
	if err != nil {
		s.zaplog.Error("Failed to check user existence", zap.Error(err))
		return nil, errors.New("internal server error")
	}
	if exists {
		s.zaplog.Warn("User already exists", zap.String("username", request.Username))
		return nil, errors.New("user already exists")
	}
	//Хэширование пароля перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.zaplog.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("internal error: unable to hash password")
	}
	//создание нового пользователя
	user := config.User{
		Username: request.Username,
		Password: string(hashedPassword),
	}
	if _, err := s.repo.CreateUser(ctx, user); err != nil {
		s.zaplog.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("internal error: unable to create user")
	}

	s.zaplog.Info("User reqistered successfully", zap.String("username", request.Username))

	return &auth.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}
