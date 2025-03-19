package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/dreamlyn/ssh-manage/internal/app"
	ssh_service "github.com/dreamlyn/ssh-manage/internal/service/ssh"
	"github.com/gorilla/websocket"
)

// SSHHandler 封装 WebSocket 处理逻辑
type SSHHandler struct {
	upgrader websocket.Upgrader
	sessions map[*websocket.Conn]*Session
	mu       sync.Mutex
}

// Session 存储 WebSocket 和 SSHConnection
type Session struct {
	Conn    *websocket.Conn
	SSHConn *ssh_service.SSHConnection
	mu      sync.Mutex
}

// NewSSHHandler 创建新的 SSHHandler
func NewSSHHandler() *SSHHandler {
	return &SSHHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 生产环境需限制跨域
			},
		},
		sessions: make(map[*websocket.Conn]*Session),
	}
}

// Handle 处理 WebSocket 连接
func (h *SSHHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 从 URL 中获取参数 id
	app := app.GetApp()
	id := r.URL.Path[len("/ws/"):]
	sshRecord, _ := app.FindRecordById("ssh", id)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	defer func() {
		h.mu.Lock()
		if session, exists := h.sessions[conn]; exists {
			if session.SSHConn != nil {
				session.SSHConn.Close()
			}
			delete(h.sessions, conn)
		}
		h.mu.Unlock()
		conn.Close()
	}()
	log.Println("客户端已连接")
	session := &Session{Conn: conn}
	h.mu.Lock()
	h.sessions[conn] = session
	h.mu.Unlock()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}

		var message struct {
			Type string `json:"type"`
			Data string `json:"data,omitempty"`
			ID   string `json:"id,omitempty"`
		}
		if err := json.Unmarshal(msg, &message); err != nil {
			h.sendError(conn, "无效的消息格式")
			continue
		}
		switch message.Type {
		case "connect":
			port := sshRecord.GetString("port")
			if port == "" {
				port = "22"
			}
			sshConnection, err := ssh_service.NewSSHConnection(sshRecord.GetString("host"), port,
				sshRecord.GetString("username"), sshRecord.GetString("password"), message.Data)
			if err != nil {
				h.sendError(conn, err.Error())
			} else {
				session.SSHConn = sshConnection
				h.sendInfo(conn, "connected")
				sshConnection.ReadOutput(
					func(data string) {
						h.sendOutput(conn, data)
					},
					func(err string) {
						h.sendError(conn, err)
					},
				)
			}
		case "command":
			if session.SSHConn == nil {
				h.sendError(conn, "未建立 SSHConn 连接")
				continue
			}
			err := session.SSHConn.WriteCommand(message.Data)
			if err != nil {
				h.sendError(conn, err.Error())
			}
		case "resize":
			if session.SSHConn == nil {
				h.sendError(conn, "未建立 SSHConn 连接")
				continue
			}
			var size struct {
				Cols int `json:"cols"`
				Rows int `json:"rows"`
			}
			if err := json.Unmarshal([]byte(message.Data), &size); err != nil {
				h.sendError(conn, "无效的数据格式")
				continue
			}
			session.SSHConn.Resize(size.Cols, size.Rows)
		default:
			h.sendError(conn, "未知的消息类型")
		}
	}
}

// sendInfo 发送输出消息
func (h *SSHHandler) sendInfo(conn *websocket.Conn, data string) {
	msg := struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}{
		Type: "info",
		Data: data,
	}
	conn.WriteJSON(msg)
}

// sendOutput 发送输出消息
func (h *SSHHandler) sendOutput(conn *websocket.Conn, data string) {
	msg := struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}{
		Type: "output",
		Data: data,
	}
	conn.WriteJSON(msg)
}

// sendError 发送错误消息
func (h *SSHHandler) sendError(conn *websocket.Conn, data string) {
	msg := struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}{
		Type: "error",
		Data: data,
	}
	conn.WriteJSON(msg)
}
