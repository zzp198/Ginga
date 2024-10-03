package main

import (
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"time"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理 WebSSH 请求
func handleSSH(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// SSH 客户端配置
	sshConfig := &ssh.ClientConfig{
		User: "", // 替换为你的 SSH 用户名
		Auth: []ssh.AuthMethod{
			ssh.Password(""), // 替换为你的 SSH 密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         18 * time.Minute,
	}

	// 连接 SSH 服务器
	sshConn, err := ssh.Dial("tcp", "[]:22", sshConfig)
	if err != nil {
		log.Println("Failed to connect to SSH server:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to connect to SSH server"))
		return
	}
	defer sshConn.Close()

	// 创建新的 SSH 会话
	session, err := sshConn.NewSession()
	if err != nil {
		log.Println("Failed to create SSH session:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to create SSH session"))
		return
	}
	defer session.Close()

	// 请求伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 禁用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	if err := session.RequestPty("xterm", 120, 40, modes); err != nil {
		log.Println("Failed to request pseudo terminal:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to request pseudo terminal"))
		return
	}

	// 启动 shell
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Println("Unable to setup stdin for SSH session:", err)
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Println("Unable to setup stdout for SSH session:", err)
		return
	}

	if err := session.Start("/bin/bash"); err != nil {
		log.Println("Failed to start shell:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to start shell"))
		return
	}

	// Goroutine：处理从 SSH 到 WebSocket 的数据流
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading from SSH stdout:", err)
				}
				break
			}

			// 将 SSH 输出发送到 WebSocket
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Println("Error sending data to WebSocket:", err)
				break
			}
		}
	}()

	// 主线程：处理从 WebSocket 到 SSH 的数据流
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// 将 WebSocket 输入发送到 SSH
		if _, err := stdin.Write(message); err != nil {
			log.Println("Error writing to SSH stdin:", err)
			break
		}
	}

	// 关闭 SSH 会话
	if err := session.Wait(); err != nil {
		log.Println("SSH session ended with error:", err)
	}
}

func main() {
	http.HandleFunc("/ssh", handleSSH)
	log.Println("WebSSH server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
