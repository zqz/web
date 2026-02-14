package repository

import "context"

type settingsRepository struct {
	queries *Queries
}

// NewSettingsRepository creates a new SettingsRepository
func NewSettingsRepository(queries *Queries) SettingsRepository {
	return &settingsRepository{queries: queries}
}

func (r *settingsRepository) Get(ctx context.Context, key string) (string, error) {
	return r.queries.GetSiteSetting(ctx, key)
}

func (r *settingsRepository) Set(ctx context.Context, key, value string) error {
	return r.queries.SetSiteSetting(ctx, SetSiteSettingParams{Key: key, Value: value})
}
