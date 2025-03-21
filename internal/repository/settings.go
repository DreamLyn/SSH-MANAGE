package repository

import (
	"database/sql"
	"errors"

	"github.com/dreamlyn/ssh-manage/internal/app"
	"github.com/dreamlyn/ssh-manage/internal/domain"
	"github.com/pocketbase/dbx"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) GetByName(name string) (*domain.Settings, error) {
	record, err := app.GetApp().FindFirstRecordByFilter(
		domain.CollectionNameSettings,
		"name={:name}",
		dbx.Params{"name": name},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	settings := &domain.Settings{
		Meta: domain.Meta{
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		Name:    record.GetString("name"),
		Content: record.GetString("content"),
	}
	return settings, nil
}
