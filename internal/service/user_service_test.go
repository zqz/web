package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/tests"
)

func TestUserServiceGetOrCreateUserCreateNew(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Name:       "Alice",
		Email:      "alice@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-alice-123",
		Role:       testRoleMember,
	}

	user, err := svc.GetOrCreateUser(ctx, req)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, testProviderGoogle, user.Provider)
	assert.Equal(t, "google-alice-123", user.ProviderID)
	assert.Equal(t, "admin", user.Role) // First user becomes admin
	assert.False(t, user.Banned)
}

func TestUserServiceGetOrCreateUserIdempotentByProviderID(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Name:       "Bob",
		Email:      "bob@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-bob-456",
		Role:       testRoleMember,
	}

	user1, err := svc.GetOrCreateUser(ctx, req)
	require.NoError(t, err)
	user2, err := svc.GetOrCreateUser(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, user1.ID, user2.ID)
	assert.Equal(t, user1.ProviderID, user2.ProviderID)
}

func TestUserServiceGetOrCreateUserSecondUserIsMember(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	_, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "First",
		Email:      "first@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-first",
		Role:       testRoleMember,
	})
	require.NoError(t, err)

	second, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "Second",
		Email:      "second@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-second",
		Role:       testRoleMember,
	})
	require.NoError(t, err)
	assert.Equal(t, testRoleMember, second.Role)
}

func TestUserServiceGetOrCreateUserValidationErrors(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	_, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "",
		Email:      "a@b.com",
		Provider:   testProviderGoogle,
		ProviderID: "pid",
		Role:       testRoleMember,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	_, err = svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "A",
		Email:      "",
		Provider:   testProviderGoogle,
		ProviderID: "pid",
		Role:       testRoleMember,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")

	_, err = svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "A",
		Email:      "a@b.com",
		Provider:   "",
		ProviderID: "pid",
		Role:       testRoleMember,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provider is required")

	_, err = svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "A",
		Email:      "a@b.com",
		Provider:   testProviderGoogle,
		ProviderID: "",
		Role:       testRoleMember,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provider ID is required")
}

func TestUserServiceGetUserByID(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "Lookup",
		Email:      "lookup@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-lookup",
		Role:       testRoleMember,
	})
	require.NoError(t, err)

	user, err := svc.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Lookup", user.Name)

	_, err = svc.GetUserByID(ctx, 99999)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserServiceGetUserByProviderID(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "ProviderLookup",
		Email:      "pl@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-pl-789",
		Role:       testRoleMember,
	})
	require.NoError(t, err)

	user, err := svc.GetUserByProviderID(ctx, "google-pl-789")
	require.NoError(t, err)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "google-pl-789", user.ProviderID)

	_, err = svc.GetUserByProviderID(ctx, "nonexistent-provider-id")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserServiceListUsers(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	for i := 0; i < 3; i++ {
		_, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
			Name:       "User" + string(rune('A'+i)),
			Email:      "user" + string(rune('a'+i)) + "@example.com",
			Provider:   testProviderGoogle,
			ProviderID: "google-list-" + string(rune('0'+i)),
			Role:       testRoleMember,
		})
		require.NoError(t, err)
	}

	users, err := svc.ListUsers(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	users, err = svc.ListUsers(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, users, 2)

	users, err = svc.ListUsers(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUserServiceSetBanned(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "BanTarget",
		Email:      "ban@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-ban",
		Role:       testRoleMember,
	})
	require.NoError(t, err)
	assert.False(t, created.Banned)

	banned, err := svc.SetBanned(ctx, created.ID, true)
	require.NoError(t, err)
	assert.True(t, banned.Banned)

	unbanned, err := svc.SetBanned(ctx, created.ID, false)
	require.NoError(t, err)
	assert.False(t, unbanned.Banned)

	_, err = svc.SetBanned(ctx, 99999, true)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserServiceSetMaxFileSize(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "SizeUser",
		Email:      "size@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-size",
		Role:       testRoleMember,
	})
	require.NoError(t, err)
	assert.Nil(t, created.MaxFileSizeOverride)

	maxBytes := int64(50 * 1024 * 1024) // 50 MB
	updated, err := svc.SetMaxFileSize(ctx, created.ID, &maxBytes)
	require.NoError(t, err)
	require.NotNil(t, updated.MaxFileSizeOverride)
	assert.Equal(t, int64(50*1024*1024), *updated.MaxFileSizeOverride)

	cleared, err := svc.SetMaxFileSize(ctx, created.ID, nil)
	require.NoError(t, err)
	assert.Nil(t, cleared.MaxFileSizeOverride)

	_, err = svc.SetMaxFileSize(ctx, 99999, &maxBytes)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserServiceUpdateProfile(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "ProfileUser",
		Email:      "profile@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-profile",
		Role:       testRoleMember,
	})
	require.NoError(t, err)

	updated, err := svc.UpdateProfile(ctx, created.ID, "AB", "#ff0000")
	require.NoError(t, err)
	assert.Equal(t, "AB", updated.DisplayTag)
	assert.Equal(t, "#ff0000", updated.Colour)

	updated2, err := svc.UpdateProfile(ctx, created.ID, "X", "#00ff00")
	require.NoError(t, err)
	assert.Equal(t, "X", updated2.DisplayTag)
	assert.Equal(t, "#00ff00", updated2.Colour)
}

func TestUserServiceUpdateProfileValidationErrors(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	svc := NewUserService(repo)

	created, err := svc.GetOrCreateUser(ctx, domain.CreateUserRequest{
		Name:       "BadProfile",
		Email:      "bad@example.com",
		Provider:   testProviderGoogle,
		ProviderID: "google-bad",
		Role:       testRoleMember,
	})
	require.NoError(t, err)

	_, err = svc.UpdateProfile(ctx, created.ID, "toolong", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "display tag must be 1-3 characters")

	_, err = svc.UpdateProfile(ctx, created.ID, "A", "nothex")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "colour must be hex")

	_, err = svc.UpdateProfile(ctx, created.ID, "A", "#fff") // 3 chars, need 6
	require.Error(t, err)
	assert.Contains(t, err.Error(), "colour must be hex")

	_, err = svc.UpdateProfile(ctx, 99999, "A", "#ffffff")
	assert.ErrorIs(t, err, ErrUserNotFound)
}
