package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"innogen-backend/auth_service/internal/dto"
	"innogen-backend/auth_service/internal/jwt"
	"innogen-backend/auth_service/internal/repository"
)

// Sentinel errors for the service layer.
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthService handles authentication business logic.
type AuthService struct {
	repo   *repository.UserRepository
	jwtSvc *jwt.Service
}

// New creates a new AuthService with the given dependencies.
func New(repo *repository.UserRepository, jwtSvc *jwt.Service) *AuthService {
	return &AuthService{repo: repo, jwtSvc: jwtSvc}
}

// Login authenticates a user by email and password, returning a JWT token and user info.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if msg := req.Validate(); msg != "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
		User:        dto.ToUserResponse(user),
	}, nil
}

// CurrentUser returns the user info for the given user ID.
func (s *AuthService) CurrentUser(ctx context.Context, userID int) (*dto.MeResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	return &dto.MeResponse{
		User: dto.ToUserResponse(user),
	}, nil
}
