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
func (h *PortForwardingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	app := app.GetApp()
	id := r.URL.Path[len("/forwarding/"):]
	forwardingRecord, _ := app.FindRecordById("port_forwarding", id)
	if forwardingRecord == nil {
		log.Println("找不到对应的端口转发记录")
		http.Error(w, "端口转发记录不存在", http.StatusNotFound)
		return
	}

	if forwardingRecord.GetBool("running") {
		h.StopForwarding(id)
		w.WriteHeader(http.StatusAccepted) // 202 表示请求已接受但处理尚未完成
		return
	}

	sshRecord, _ := app.FindRecordById("ssh", forwardingRecord.GetString("sshId"))
	if sshRecord == nil {
		log.Println("找不到对应的 SSH 记录")
		http.Error(w, "SSH 记录不存在", http.StatusNotFound)
		return
	}

	// 立即返回响应
	w.WriteHeader(http.StatusAccepted) // 202 表示请求已接受但处理尚未完成
	fmt.Fprintf(w, "端口转发请求已接受: %s\n", id)

	// 异步启动端口转发
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.startPortForwarding(id, sshRecord, forwardingRecord)
	}()
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
	sshClient, err := ssh.Dial("tcp", sshRecord.GetString("host")+":"+sshRecord.GetString("port"), sshConfig)
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
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(remoteConn, localConn)
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
