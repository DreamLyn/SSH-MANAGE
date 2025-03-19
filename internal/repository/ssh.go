package repository

import "context"

type SSHRepository struct{}

func NewSSHRepository() *SSHRepository {
	return &SSHRepository{}
}

func (s *SSHRepository) GetByID(ctx context.Context, id string) (map[string]any, error) {
	return nil, nil
}
