package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"strings"
)

func WriteString(w *bufio.Writer, str string) {
	w.WriteString(str)
	w.Flush()
	fmt.Print("\x1b[32m" + str + "\x1b[0m")
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

		go func() {
			data, err := Process(conn)
			if err != nil {
				fmt.Println(err)
				return
			}

			if data == "" { // 对面服务器EHLO不通过直接QUIT,data都不给...
				return
			}

			msg, err := mail.ReadMessage(strings.NewReader(data))
			if err != nil {
				fmt.Println(err)
				return
			}

			// TODO, From和To也需要处理,To还可能包含多个.
			// 预处理Subject的格式,163,QQ等默认GBK需要转UTF8.
			wd := mime.WordDecoder{CharsetReader: charset.NewReaderLabel}
			subject, err := wd.Decode(msg.Header.Get("Subject"))
			if err == nil { // 也可能不需要解码
				msg.Header["Subject"] = []string{subject}
			}

			fmt.Println("Date: ", msg.Header.Get("Date"))
			fmt.Println("From: ", msg.Header.Get("From"))
			fmt.Println("To: ", msg.Header.Get("To"))
			fmt.Println("Subject: ", msg.Header.Get("Subject"))
			fmt.Println("Content-Type: ", msg.Header.Get("Content-Type"))
			fmt.Println()

			if msg.Header.Get("Content-Type") == "" {
				bin, _ := io.ReadAll(msg.Body)
				fmt.Println(string(bin))
				return
			}

			mimeType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
			if err != nil {
				fmt.Println(fmt.Errorf("failed to parse media type: %v", err))
				return
			}

			if strings.HasPrefix(mimeType, "multipart/") {
				mr := multipart.NewReader(msg.Body, params["boundary"])

				var htmlContent string
				var attachmentsInfo string

				for {
					part, err := mr.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						fmt.Println(err)
						return
					}

					content, attachmentInfo, err := ParsePart(part)
					if err != nil {
						fmt.Println(err)
						continue
					}

					if content != "" {
						htmlContent += content
					}

					if attachmentInfo != "" {
						attachmentsInfo += attachmentInfo
					}

					fmt.Println("HTML Content:", htmlContent)
					fmt.Println("Attachments Info:\n", attachmentsInfo)
				}
			}
		}()
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
			//正常会在QUIT命令后停止读取并断开连接,不会冗余读取
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

func ParsePart(part *multipart.Part) (string, string, error) {
	mimeType, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse media type: %v", err)
	}

	switch {
	case strings.HasPrefix(mimeType, "multipart/"):
		// Recursive parsing for multipart
		mr := multipart.NewReader(part, params["boundary"])
		var result string
		var attachmentsInfo string
		for {
			subPart, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", "", fmt.Errorf("error reading sub-part: %v", err)
			}

			content, attachmentInfo, err := ParsePart(subPart)
			if err != nil {
				fmt.Printf("Error parsing sub-part: %v\n", err)
				continue
			}
			if content != "" {
				result += content
			}
			if attachmentInfo != "" {
				attachmentsInfo += attachmentInfo
			}
		}
		return result, attachmentsInfo, nil

	case mimeType == "text/plain", mimeType == "text/html":
		// Decode text or HTML content
		encoding := part.Header.Get("Content-Transfer-Encoding")
		content, err := io.ReadAll(part)
		if err != nil {
			return "", "", fmt.Errorf("failed to read part content: %v", err)
		}
		return decodeContent(content, encoding, params["charset"]), "", nil

	case part.FileName() != "":
		// Handle attachments
		filename := part.FileName()
		if decodedFilename, err := decodeHeader(filename); err == nil {
			filename = decodedFilename
		}
		content, err := io.ReadAll(part)
		if err != nil {
			return "", "", fmt.Errorf("failed to read attachment: %v", err)
		}
		attachmentInfo := fmt.Sprintf("Attachment: %s, Content Length: %d\n", filename, len(content))
		return "", attachmentInfo, nil

	default:
		// Ignore other content types
		return "", "", nil
	}
}

func decodeContent(content []byte, encoding string, charset string) string {
	switch strings.ToLower(encoding) {
	case "base64":
		decoded, _ := io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(content)))
		content = decoded
	case "quoted-printable":
		decoded, _ := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(content)))
		content = decoded
	}
	return convertToUTF8(charset, content)
}

func convertToUTF8(label string, content []byte) string {
	if strings.ToLower(label) == "utf-8" || label == "" {
		return string(content)
	}
	reader, err := charset.NewReaderLabel(label, bytes.NewReader(content))
	if err != nil {
		return string(content) // Fallback to raw content
	}
	utf8Content, _ := io.ReadAll(reader)
	return string(utf8Content)
}

// decodeHeader decodes a MIME header to a readable string.
func decodeHeader(header string) (string, error) {
	wd := mime.WordDecoder{
		CharsetReader: charset.NewReaderLabel,
	}
	return wd.Decode(header)
}
