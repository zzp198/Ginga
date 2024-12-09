package backup

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {

	fmt.Println("临时邮件服务器")
	fmt.Println()

	listener, err := net.Listen("tcp", "0.0.0.0:25")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer func(listener net.Listener) { _ = listener.Close() }(listener)

	fmt.Printf("SMTP server is running on port %s\n", "0.0.0.0:25")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting client: %v\n", err)
			continue
		}

		go Parse(conn)
	}
}

func Parse1(conn net.Conn) {

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	// 向客户端发送欢迎消息
	_, _ = w.WriteString("220 Temporary Email Server by zzp198\r\n")
	_ = w.Flush()

	var data string
	var readingBody bool

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from client: %v\n", err)
			}
			break
		}

		//line = strings.TrimSpace(line)
		fmt.Println("C: " + line)

		if readingBody {
			if line == "." {
				readingBody = false
				_, _ = w.WriteString("250 OK\r\n")
				_ = w.Flush()

				fmt.Println(data)
				fmt.Println()

				//// 现在邮件发送时包含DKIM,net/mail不支持导致报错,还有各种转码问题,“\t代表换行内容”这样就好处理了
				//msg, err := mail.ReadMessage(strings.NewReader(data))
				//if err != nil {
				//	fmt.Printf("Error parsing mail: %v\n", err)
				//	return
				//}
				//
				//for key, values := range msg.Header {
				//	for _, value := range values {
				//		fmt.Printf("%s: %s", key, value)
				//	}
				//}
				//
				//body, err := io.ReadAll(msg.Body)
				//if err != nil {
				//	fmt.Printf("Error reading body: %v\n", err)
				//	return
				//}
				//
				//fmt.Println(string(body))

			} else {
				data += line + "\n"
			}
			continue
		}

		switch {
		case strings.HasPrefix(strings.ToUpper(line), "HELO"):
			_, _ = w.WriteString("250 Hello\r\n")
			fmt.Print("S: " + "250 Hello\r\n")
		case strings.HasPrefix(strings.ToUpper(line), "MAIL FROM"):
			_, _ = w.WriteString("250 OK\r\n")
			fmt.Print("S: " + "250 OK\r\n")
		case strings.HasPrefix(strings.ToUpper(line), "RCPT TO"):
			_, _ = w.WriteString("250 OK\r\n")
			fmt.Print("S: " + "250 OK\r\n")
		case strings.HasPrefix(strings.ToUpper(line), "DATA"):
			_, _ = w.WriteString("354 End data with <CR><LF>.<CR><LF>\r\n")
			fmt.Print("S: " + "354 End data with <CR><LF>.<CR><LF>\r\n")
			readingBody = true
		case strings.HasPrefix(strings.ToUpper(line), "QUIT"):
			_, _ = w.WriteString("221 Bye\r\n")
			_ = w.Flush()

			fmt.Print("S: " + "221 Bye\r\n")
			_ = conn.Close()
			return
		default:
			_, _ = w.WriteString("502 Command not implemented\r\n")
			fmt.Print("S: " + "502 Command not implemented\r\n")
		}
		_ = w.Flush()
	}
}
