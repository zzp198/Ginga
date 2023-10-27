package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {

	ip := flag.String("ip", "0.0.0.0:8080", "")
	//key := flag.String("key", "Quark", "")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hey!")
	})

	r.GET("/state", State)

	r.GET("/api", func(c *gin.Context) {
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			mt, msg, err := conn.ReadMessage()

			if err != nil {
				fmt.Println("Failed to read message: ", err)
				break
			}

			fmt.Printf("Received message: %s\n", string(msg))

			err = conn.WriteMessage(mt, msg)
			if err != nil {
				fmt.Println("Failed to write message: ", err)
				break
			}
		}
	})

	_ = r.Run(*ip)
}
