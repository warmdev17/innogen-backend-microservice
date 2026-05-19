package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"innogen-backend/auth_service/internal/dto"
	"innogen-backend/auth_service/internal/jwt"
	"innogen-backend/auth_service/internal/repository"
	"innogen-backend/shared/config"
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
	cfg    *config.Config
}

// New creates a new AuthService with the given dependencies.
func New(repo *repository.UserRepository, jwtSvc *jwt.Service, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, jwtSvc: jwtSvc, cfg: cfg}
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

// FrontendURL returns the configured frontend URL.
func (s *AuthService) FrontendURL() string {
	return s.cfg.FrontendURL
}

// GithubConnectURL generates the GitHub App installation URL with a state JWT.
func (s *AuthService) GithubConnectURL(userID int, userEmail string) string {
	stateToken, _ := s.jwtSvc.GenerateToken(userID, userEmail, "user")
	return s.cfg.GitHubAppInstallURL + "?state=" + stateToken
}

// HandleGithubCallback processes the GitHub OAuth callback.
// Returns the frontend redirect URL.
func (s *AuthService) HandleGithubCallback(ctx context.Context, installationID string, stateToken string) (string, error) {
	// Parse state JWT to get userID
	claims, err := s.jwtSvc.ParseToken(stateToken)
	if err != nil {
		return s.cfg.FrontendURL + "/settings?github_status=error&message=invalid_state", nil
	}

	// Upsert github_account link
	if err := s.repo.UpsertGithubAccount(ctx, claims.UserID, installationID); err != nil {
		return "", err
	}

	// Try backfill owner from webhook data (if webhook arrived first)
	owner, ownerType, err := s.repo.GetGithubInstallationOwner(ctx, installationID)
	if err == nil && owner != "" {
		_ = s.repo.UpdateGithubAccountOwner(ctx, claims.UserID, owner, ownerType, "active")
	}

	return s.cfg.FrontendURL + "/settings?github_status=connected", nil
}

// GithubStatus returns the GitHub connection status for a user.
func (s *AuthService) GithubStatus(ctx context.Context, userID int) (*dto.GithubStatusResponse, error) {
	instID, owner, ownerType, status, err := s.repo.GetGithubAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if instID == "" {
		return &dto.GithubStatusResponse{Connected: false}, nil
	}
	return &dto.GithubStatusResponse{
		Connected:       true,
		InstallationID:  instID,
		GithubOwner:     owner,
		GithubOwnerType: ownerType,
		Status:          status,
	}, nil
}
