package repository

import (
	"context"
	"database/sql"
)

type userRepository struct {
	queries *Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(queries *Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, params CreateUserParams) (*User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByProviderID(ctx context.Context, providerID string) (*User, error) {
	user, err := r.queries.GetUserByProviderID(ctx, providerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]*User, error) {
	users, err := r.queries.ListUsers(ctx, ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}

func (r *userRepository) Update(ctx context.Context, params UpdateUserParams) (*User, error) {
	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, params UpdateUserProfileParams) (*User, error) {
	user, err := r.queries.UpdateUserProfile(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SetBanned(ctx context.Context, userID int32, banned bool) (*User, error) {
	user, err := r.queries.SetUserBanned(ctx, SetUserBannedParams{ID: userID, Banned: banned})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SetMaxFileSize(ctx context.Context, userID int32, maxBytes *int64) (*User, error) {
	user, err := r.queries.SetUserMaxFileSize(ctx, SetUserMaxFileSizeParams{ID: userID, MaxFileSizeOverride: maxBytes})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) CountBanned(ctx context.Context) (int64, error) {
	return r.queries.CountBannedUsers(ctx)
}
