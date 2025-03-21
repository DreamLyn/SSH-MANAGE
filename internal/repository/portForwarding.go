package repository

import (
	"context"
	"errors"

	"github.com/dreamlyn/ssh-manage/internal/app"
	"github.com/pocketbase/dbx"
)

type PortForwardingRepository struct{}

func NewPortForwardingRepository() *PortForwardingRepository {
	return &PortForwardingRepository{}
}

func (r *PortForwardingRepository) Initialize() error {
	_, err := app.GetApp().DB().NewQuery(
		"UPDATE {{port_forwarding}} SET running = {:running}",
	).Bind(dbx.Params{"running": false}).Execute()

	if err != nil {
		return errors.New("failed to initialize port forwarding records: " + err.Error())
	}
	return nil
}

func (s *PortForwardingRepository) GetByID(ctx context.Context, id string) (map[string]any, error) {
	return nil, nil
}
