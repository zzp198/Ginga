package main

import (
	"errors"
	"fmt"
	"github.com/alash3al/go-smtpsrv/v3"
)

func main() {

	// 只有谷歌能正常接收和解析，163卡在解码，qq直接收不了
	cfg := smtpsrv.ServerConfig{
		BannerDomain:    "980520.xyz",
		ListenAddr:      ":25",
		MaxMessageBytes: 5 * 1024,
		Handler: smtpsrv.HandlerFunc(func(c *smtpsrv.Context) error {
			msg, err := c.Parse()
			if err != nil {
				return errors.New("Cannot read your message: " + err.Error())
			}

			fmt.Println(msg.From)
			fmt.Println(msg.To)
			fmt.Println(msg.Subject)
			fmt.Println(msg.TextBody)

			return nil
		}),
	}

	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}
