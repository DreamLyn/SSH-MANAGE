package ssh_service

import (
	"context"

	"github.com/dreamlyn/ssh-manage/internal/domain"
)

type sshRepository interface {
	GetByID(ctx context.Context, id string) (*domain.SSH, error)
}

type SSHService struct {
	sshRepo sshRepository
}

// 创建对象
func NewSSHService(sshRepo sshRepository) *SSHService {
	return &SSHService{
		sshRepo: sshRepo,
	}
}

// func (n *SSHService) ConnectSSH(ctx context.Context, req *domain.SSH) error {

// 	client, err := goph.New("root", "192.168.31.2", goph.Password("PbVQ#fk&F!4LZK6!"))

// }
