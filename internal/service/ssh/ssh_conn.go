package ssh_service

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/ssh"
)

// SSHConnection 封装 SSH 客户端和会话
type SSHConnection struct {
	Client  *ssh.Client
	Session *ssh.Session
	Stdin   io.WriteCloser // 用于写入命令
	Stdout  io.Reader      // 用于读取输出
	mu      sync.Mutex
}

// NewSSHConnection 创建并初始化 SSHConnection
func NewSSHConnection(host, port, username, password string, data string) (*SSHConnection, error) {
	service := &SSHConnection{}

	// 配置 SSH 客户端
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境需验证主机密钥
	}

	// 连接 SSH 服务器
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}

	// 获取 Stdin 和 Stdout 管道
	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	// 请求伪终端 (TTY)
	var size struct {
		Cols int `json:"cols"`
		Rows int `json:"rows"`
	}
	if err := json.Unmarshal([]byte(data), &size); err != nil {

	}
	err = session.RequestPty("xterm", size.Rows, size.Cols, ssh.TerminalModes{})
	if err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	// 启动 Shell
	err = session.Shell()
	if err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	service.Client = client
	service.Session = session
	service.Stdin = stdin
	service.Stdout = stdout

	return service, nil
}

// WriteCommand 向 Shell 输入命令
func (s *SSHConnection) WriteCommand(command string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Stdin == nil {
		return fmt.Errorf("stdin is nil")
	}
	_, err := s.Stdin.Write([]byte(command)) // 添加换行符以模拟回车
	return err
}

// ReadOutput 异步读取 Shell 输出并通过回调函数发送
func (s *SSHConnection) ReadOutput(sendOutput func(string), sendError func(string)) {
	go func() {
		if s.Stdout == nil {
			sendError("stdout is nil")
			return
		}
		buf := make([]byte, 1024)
		for {
			n, err := s.Stdout.Read(buf)
			if err != nil {
				sendError("读取 Shell 输出失败: " + err.Error())
				return
			}
			if n > 0 {
				sendOutput(string(buf[:n]))
			}
		}
	}()
}

func (s *SSHConnection) Resize(cols, rows int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Session.WindowChange(rows, cols) // 调整终端大小
}

// Close 关闭 SSH 客户端和会话
func (s *SSHConnection) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Stdin != nil {
		s.Stdin.Close()
		s.Stdin = nil
	}
	if s.Session != nil {
		s.Session.Close()
		s.Session = nil
	}
	if s.Client != nil {
		s.Client.Close()
		s.Client = nil
	}
}
