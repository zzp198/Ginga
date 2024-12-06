package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	thirednet "github.com/shirou/gopsutil/net"
	"io"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

func main() {

	//State()

	//frontend.Server("")

	Web := gin.Default()

	Web.GET("/events", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		for {
			fmt.Fprintf(c.Writer, "data: %s\n\n", time.Now().Format(time.Stamp))
			c.Writer.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
	})

	Web.Run(":8080")

	//MailServer()

}

func State() {

	// CPU 信息
	cpuInfo, _ := cpu.Info()
	fmt.Println("CPU Info:", cpuInfo)

	// 内存信息
	memInfo, _ := mem.VirtualMemory()
	fmt.Println("Memory Info:", memInfo)

	// 磁盘信息
	diskInfo, _ := disk.Usage("/")
	fmt.Println("Disk Info:", diskInfo)

	// 网络信息
	netInfo, _ := thirednet.Interfaces()
	fmt.Println("Network Info:", netInfo)

}

func MailServer() {

	listener, err := net.Listen("tcp", "0.0.0.0:25")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("SMTP server is running on port %s\n", "0.0.0.0:25")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			continue
		}
		go handleConnection(conn) // 并发处理客户端连接
	}
}

// handleConnection 处理每个客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 向客户端发送欢迎消息
	fmt.Fprintf(conn, "220 Simple Go SMTP Server\r\n")

	var emailData strings.Builder
	scanner := bufio.NewScanner(conn)
	var readingData bool

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Client: %s\n", line)

		// 处理邮件内容
		if readingData {
			if line == "." { // 结束标志
				break
			}
			emailData.WriteString(line + "\r\n")
			continue
		}

		// 处理 SMTP 指令
		switch {
		case strings.HasPrefix(line, "HELO") || strings.HasPrefix(line, "EHLO"):
			fmt.Fprintf(conn, "250 Hello\r\n")

		case strings.HasPrefix(line, "MAIL FROM:"):
			fmt.Fprintf(conn, "250 OK\r\n")

		case strings.HasPrefix(line, "RCPT TO:"):
			fmt.Fprintf(conn, "250 OK\r\n")

		case strings.HasPrefix(line, "DATA"):
			fmt.Fprintf(conn, "354 End data with <CR><LF>.<CR><LF>\r\n")
			readingData = true

		case strings.HasPrefix(line, "QUIT"):
			fmt.Fprintf(conn, "221 Bye\r\n")
			return

		default:
			fmt.Fprintf(conn, "502 Command not implemented\r\n")
		}
	}

	// 解析邮件内容
	rawEmail := emailData.String()
	message, err := mail.ReadMessage(strings.NewReader(rawEmail))
	if err != nil {
		fmt.Printf("Failed to parse email: %v\n", err)
		fmt.Fprintf(conn, "550 Error parsing email\r\n")
		return
	}

	// 提取头部和正文
	headers := message.Header
	body, _ := io.ReadAll(message.Body)

	// 打印解析后的内容
	fmt.Println("==== Parsed Email ====")
	fmt.Printf("From: %s\n", headers.Get("From"))
	fmt.Printf("To: %s\n", headers.Get("To"))
	fmt.Printf("Subject: %s\n", headers.Get("Subject"))
	fmt.Printf("Body:\n%s\n", string(body))

	// 向客户端确认
	fmt.Fprintf(conn, "250 OK\r\n")
}
