package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func WriteString(w *bufio.Writer, str string) {
	w.WriteString(str)
	w.Flush()
	fmt.Println("\x1b[32m" + str + "\x1b[0m")
}

func main() {
	fmt.Println("临时邮件服务器")

	listener, err := net.Listen("tcp", "0.0.0.0:25")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("运行于0.0.0.0:25")

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go fmt.Println(Process(conn))
	}

}

func Process(conn net.Conn) (data string, err error) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	WriteString(w, "220 Temporary Email Server\r\n")

	var reading bool

	for {
		line, err := r.ReadString('\n')

		if err != nil {
			//正常应该在QUIT命令后同时停止发送并断开连接,不会冗余读取
			if err == io.EOF {
				fmt.Printf("Error reading from client: %v\n", err)
				return data, err
			}
			fmt.Println(err)
			return "", err
		}

		fmt.Print(line)

		// 内容部分需要保留开头的\t
		if reading {
			if strings.TrimSpace(line) == "." {
				reading = false
				WriteString(w, "250 OK\r\n")
			} else {
				data += line
			}
			continue
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToUpper(line), "HELO") {
			WriteString(w, "250 OK\r\n")
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "MAIL FROM") {
			WriteString(w, "250 OK\r\n")
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "RCPT TO") {
			WriteString(w, "250 OK\r\n")
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "DATA") {
			reading = true
			WriteString(w, "354 <CR><LF>.<CR><LF>\r\n")
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "QUIT") {
			WriteString(w, "221 OK\r\n")
			return data, nil
		}

		WriteString(w, "502 Command not ready\r\n")
	}
}
