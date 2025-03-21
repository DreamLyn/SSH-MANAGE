package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/dreamlyn/ssh-manage/internal/app"
	"github.com/dreamlyn/ssh-manage/internal/rest/resp"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/crypto/ssh"
)

// PortForwardingHandler 封装端口转发处理逻辑
type PortForwardingHandler struct {
	cancelFuncs sync.Map // 存储每个转发的取消函数
	wg          sync.WaitGroup
}

// NewPortForwardingHandler 创建新的 PortForwardingHandler
func NewPortForwardingHandler() *PortForwardingHandler {
	return &PortForwardingHandler{}
}

// Handle 处理 HTTP 请求并启动端口转发
func (h *PortForwardingHandler) Handle(e *core.RequestEvent) error {
	r := e.Request
	w := e.Response
	app := app.GetApp()
	id := r.URL.Path[len("/forwarding/"):]
	forwardingRecord, _ := app.FindRecordById("port_forwarding", id)
	if forwardingRecord == nil {
		log.Println("找不到对应的端口转发记录")
		return resp.Err(e, e.Error(500, "端口转发记录不存在", nil))
	}

	if forwardingRecord.GetBool("running") {
		h.StopForwarding(id)
		w.WriteHeader(http.StatusAccepted) // 202 表示请求已接受但处理尚未完成
		return resp.Ok(e, nil)
	}

	sshRecord, _ := app.FindRecordById("ssh", forwardingRecord.GetString("sshId"))
	if sshRecord == nil {
		log.Println("找不到对应的 SSH 记录")
		return resp.Err(e, e.Error(500, "SSH 记录不存在", nil))
	}
	// 异步启动端口转发
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.startPortForwarding(id, sshRecord, forwardingRecord)
	}()

	return resp.Ok(e, nil)
}

// startPortForwarding 启动端口转发
func (h *PortForwardingHandler) startPortForwarding(id string, sshRecord, forwardingRecord *core.Record) {
	ctx, cancel := context.WithCancel(context.Background())
	h.cancelFuncs.Store(id, cancel)
	defer h.cancelFuncs.Delete(id)

	// 配置 SSH 客户端
	sshConfig := &ssh.ClientConfig{
		User: sshRecord.GetString("username"),
		Auth: []ssh.AuthMethod{
			ssh.Password(sshRecord.GetString("password")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接 SSH 服务器
	port := sshRecord.GetString("port")
	if port == "" {
		port = "22"
	}
	sshClient, err := ssh.Dial("tcp", sshRecord.GetString("host")+":"+port, sshConfig)
	if err != nil {
		log.Printf("无法连接 SSH 服务器: %v", err)
		return
	}
	defer sshClient.Close()

	// 从 forwardingRecord 获取地址
	localAddr := fmt.Sprintf("%s:%s", forwardingRecord.GetString("localAddress"), forwardingRecord.GetString("localPort"))
	remoteAddr := fmt.Sprintf("%s:%s", forwardingRecord.GetString("remoteAddress"), forwardingRecord.GetString("remotePort"))

	// 启动本地监听
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("无法监听本地端口: %v", err)
		return
	}
	defer listener.Close()

	// 存储 listener，以便 StopForwarding 时关闭
	h.cancelFuncs.Store(id+"_listener", listener)

	log.Printf("SSH 端口转发已启动: %s -> %s", localAddr, remoteAddr)
	forwardingRecord.Set("running", true)
	app.GetApp().Save(forwardingRecord)

	// 监控 SSH 连接状态
	go func() {
		err := sshClient.Wait()
		if err != nil {
			log.Printf("SSH 连接断开: %v", err)
		}
		listener.Close()
		cancel()
	}()

	// 处理连接转发
	for {
		select {
		case <-ctx.Done():
			log.Printf("端口转发 %s 已停止", id)
			forwardingRecord.Set("running", false)
			app.GetApp().Save(forwardingRecord)
			return
		default:
			localConn, err := listener.Accept()
			if err != nil {
				log.Printf("接受本地连接失败: %v", err)
				return
			}
			go h.forward(localConn, sshClient, remoteAddr, ctx)
		}
	}
}

// forward 处理端口转发
func (h *PortForwardingHandler) forward(localConn net.Conn, sshClient *ssh.Client, remoteAddr string, ctx context.Context) {
	defer localConn.Close()

	remoteConn, err := sshClient.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("连接远程地址失败: %v", err)
		return
	}
	defer remoteConn.Close()

	errChan := make(chan error, 2)
	// go func() {
	// 	_, err := io.Copy(localConn, remoteConn)
	// 	errChan <- err
	// }()
	// go func() {
	// 	_, err := io.Copy(remoteConn, localConn)
	// 	errChan <- err
	// }()
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		if err != nil {
			log.Printf("本地到远程转发失败: %v", err)
		}
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("远程到本地转发失败: %v", err)
		}
		errChan <- err
	}()
	select {
	case <-ctx.Done():
		// 上下文取消时，主动关闭连接
		localConn.Close()
		remoteConn.Close()
		log.Printf("转发连接被取消")
	case err := <-errChan:
		if err != nil {
			log.Printf("转发连接断开: %v", err)
		}
	}
}

// StopForwarding 停止指定 ID 的端口转发
func (h *PortForwardingHandler) StopForwarding(id string) {
	if cancel, ok := h.cancelFuncs.Load(id); ok {
		cancel.(context.CancelFunc)()
		h.cancelFuncs.Delete(id)

		forwardingRecord, _ := app.GetApp().FindRecordById("port_forwarding", id)
		forwardingRecord.Set("running", false)
		app.GetApp().Save(forwardingRecord)

		// 关闭监听端口（存储 listener 的 map）
		if l, ok := h.cancelFuncs.Load(id + "_listener"); ok {
			l.(net.Listener).Close()
			h.cancelFuncs.Delete(id + "_listener")
		}
	}
}

// Shutdown 停止所有端口转发并清理资源
func (h *PortForwardingHandler) Shutdown() {
	// 停止所有端口转发
	h.cancelFuncs.Range(func(key, value interface{}) bool {
		if cancel, ok := value.(context.CancelFunc); ok {
			cancel() // 触发上下文取消
		} else if listener, ok := value.(net.Listener); ok {
			listener.Close() // 关闭监听器
		}
		h.cancelFuncs.Delete(key)
		return true
	})

	// 等待所有协程完成
	h.wg.Wait()
	log.Println("所有端口转发已停止，资源已清理")
}
