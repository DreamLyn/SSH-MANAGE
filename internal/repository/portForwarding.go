package repository

import "context"

type PortForwardingRepository struct{}

func NewPortForwardingRepository() *PortForwardingRepository {
	return &PortForwardingRepository{}
}

func (s *PortForwardingRepository) GetByID(ctx context.Context, id string) (map[string]any, error) {
	return nil, nil
}
