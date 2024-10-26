package main

import (
	"flag"
	"github.com/gin-gonic/gin"
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

func main() {
	ip := flag.String("ip", "0.0.0.0:8080", "")

	flag.Parse()

	r := gin.Default()

	r.GET("/xterm", func(c *gin.Context) {

	})

	r.GET("api/xterm", func(c *gin.Context) {
		// 升级 HTTP 请求为 WebSocket 连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		//addr, ok := c.GetQuery("addr")
		//if !ok {
		//	c.JSON(500, gin.H{"error": "no addr"})
		//	return
		//}
		//
		//user, ok := c.GetQuery("user")
		//if !ok {
		//	c.JSON(500, gin.H{"error": "no user"})
		//	return
		//}
		//
		//auth, ok := c.GetQuery("auth")
		//if !ok {
		//	c.JSON(500, gin.H{"error": "no auth"})
		//	return
		//}

		sshConfig := &ssh.ClientConfig{
			User:            "root",
			Auth:            []ssh.AuthMethod{ssh.Password("258237")},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         10 * time.Minute,
		}

		sshConn, err := ssh.Dial("tcp", "[2001:41d0:800:482:565f:2505:9913:9a20]:22", sshConfig)
		if err != nil {
			log.Println("Failed to connect to SSH server:", err)
			c.JSON(200, gin.H{"error": "Failed to connect to SSH server"})
			return
		}
		defer sshConn.Close()

		session, err := sshConn.NewSession()
		if err != nil {
			log.Println("Failed to create SSH session:", err)
			c.JSON(200, gin.H{"error": "Failed to create SSH session"})
			return
		}
		defer session.Close()

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err = session.RequestPty("xterm", 120, 40, modes); err != nil {
			log.Println("Failed to request pseudo terminal:", err)
			c.JSON(200, gin.H{"error": "Failed to request pseudo terminal"})
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
	})

	r.Run(*ip)
}
