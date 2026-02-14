package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
)

var hexColourRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when trying to create a user that already exists
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserService handles user business logic
type UserService struct {
	repo *repository.Repository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetOrCreateUser gets an existing user or creates a new one based on provider ID
func (s *UserService) GetOrCreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	// Validate request
	if err := validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Try to find existing user by provider ID
	dbUser, err := s.repo.Users.GetByProviderID(ctx, req.ProviderID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if dbUser != nil {
		// User exists, return it
		return dbUserToDomain(dbUser), nil
	}

	// Check if this is the first user (make them admin)
	role := req.Role
	if role == "" {
		role = "member" // Default role
	}

	count, err := s.repo.Users.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	if count == 0 {
		// First user becomes admin
		role = "admin"
	}

	// Create new user
	dbUser, err = s.repo.Users.Create(ctx, repository.CreateUserParams{
		Name:       req.Name,
		Email:      req.Email,
		Provider:   req.Provider,
		ProviderID: req.ProviderID,
		Role:       role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return dbUserToDomain(dbUser), nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int32) (*domain.User, error) {
	dbUser, err := s.repo.Users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dbUserToDomain(dbUser), nil
}

// GetUserByProviderID retrieves a user by their OAuth provider ID
func (s *UserService) GetUserByProviderID(ctx context.Context, providerID string) (*domain.User, error) {
	dbUser, err := s.repo.Users.GetByProviderID(ctx, providerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dbUserToDomain(dbUser), nil
}

// SetBanned sets the banned status of a user (admin only; caller must enforce).
func (s *UserService) SetBanned(ctx context.Context, userID int32, banned bool) (*domain.User, error) {
	dbUser, err := s.repo.Users.SetBanned(ctx, userID, banned)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to set banned: %w", err)
	}
	return dbUserToDomain(dbUser), nil
}

// SetMaxFileSize sets the max file size override for a user (admin only; caller must enforce). nil = use site default.
func (s *UserService) SetMaxFileSize(ctx context.Context, userID int32, maxBytes *int64) (*domain.User, error) {
	dbUser, err := s.repo.Users.SetMaxFileSize(ctx, userID, maxBytes)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to set max file size: %w", err)
	}
	return dbUserToDomain(dbUser), nil
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	dbUsers, err := s.repo.Users.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = dbUserToDomain(dbUser)
	}

	return users, nil
}

// UpdateProfile updates the current user's display tag and colour.
// displayTag must be 1-3 chars (or empty to clear). colour must be hex #RRGGBB or empty to clear.
func (s *UserService) UpdateProfile(ctx context.Context, userID int32, displayTag, colour string) (*domain.User, error) {
	displayTag = strings.TrimSpace(displayTag)
	colour = strings.TrimSpace(colour)
	if displayTag != "" && (len([]rune(displayTag)) < 1 || len([]rune(displayTag)) > 3) {
		return nil, errors.New("display tag must be 1-3 characters")
	}
	if colour != "" && !hexColourRegex.MatchString(colour) {
		return nil, errors.New("colour must be hex #RRGGBB")
	}
	var tagPtr, colourPtr *string
	if displayTag != "" {
		tagPtr = &displayTag
	}
	if colour != "" {
		colourPtr = &colour
	}
	dbUser, err := s.repo.Users.UpdateProfile(ctx, repository.UpdateUserProfileParams{
		ID:         userID,
		DisplayTag: tagPtr,
		Colour:     colourPtr,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return dbUserToDomain(dbUser), nil
}

// Helper functions

func validateCreateUserRequest(req domain.CreateUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Provider == "" {
		return errors.New("provider is required")
	}
	if req.ProviderID == "" {
		return errors.New("provider ID is required")
	}
	if req.Role == "" {
		req.Role = "member" // Default role
	}
	return nil
}

func dbUserToDomain(u *repository.User) *domain.User {
	out := &domain.User{
		ID:                  u.ID,
		Name:                u.Name,
		Email:               u.Email,
		Provider:            u.Provider,
		ProviderID:          u.ProviderID,
		Role:                u.Role,
		Banned:              u.Banned,
		MaxFileSizeOverride: u.MaxFileSizeOverride,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
	if u.DisplayTag != nil {
		out.DisplayTag = *u.DisplayTag
	}
	if u.Colour != nil {
		out.Colour = *u.Colour
	}
	return out
}
