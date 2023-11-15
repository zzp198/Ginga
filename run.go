package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os/exec"
)

type CmdObject struct {
	Cmd    exec.Cmd
	Output bytes.Buffer
}

func main() {

	ip := flag.String("ip", "0.0.0.0:8080", "")

	r := gin.Default()

	// Get -> api/:Param/?arg=Query -> PostForm

	api := r.Group("api")

	tasks := make(map[string]*CmdObject)

	api.GET("cron", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	api.GET("cron/add", func(c *gin.Context) {
		id := c.Query("id")
		arg := c.Query("arg")

		tasks[id] = &CmdObject{
			Cmd: *exec.Command("/bin/bash", "-c", arg),
		}

		tasks[id].Cmd.Stdout = &tasks[id].Output
		tasks[id].Cmd.Stderr = &tasks[id].Output
		reader := bufio.NewReader(&tasks[id].Output)
		_ = tasks[id].Cmd.Start()

		go func() {
			for {
				line, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				tasks[id].Output.WriteString(line)
			}
		}()
		c.String(http.StatusOK, "ok")
	})
	api.GET("cron/detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, tasks[id].Output.String())
	})
	api.GET("cron/kill/:id", func(c *gin.Context) {
		id := c.Param("id")
		tasks[id].Cmd.Process.Kill()
		c.String(http.StatusOK, "ok")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hey!")
	})

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
