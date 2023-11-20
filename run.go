package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
)

// Get -> api/:[Param]/?arg=[Query] -> [PostForm]
func main() {
	ip := flag.String("ip", "0.0.0.0:8080", "")
	if len(os.Args) > 1 {
		if strings.ToUpper(os.Args[1]) == "DEBUG" {
			fmt.Println(debug.ReadBuildInfo())
			return
		}
		if strings.ToUpper(os.Args[1]) == "START" {

		}
		if strings.ToUpper(os.Args[1]) == "STOP" {

		}
	}

	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("api")

	tasks := make(map[string]*exec.Cmd)

	api.GET("cron", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	api.GET("cron/add", func(c *gin.Context) {
		id := c.Query("id")
		arg := c.Query("arg")

		logf, _ := os.Create(fmt.Sprintf("log/%s.log", id))
		tasks[id] = exec.Command("/bin/bash", "-c", arg)
		tasks[id].Stdout = logf
		tasks[id].Stderr = logf

		go func() {
			_ = tasks[id].Run()
		}()

		c.String(http.StatusOK, "ok")
	})

	api.GET("cron/detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		data, _ := os.ReadFile(fmt.Sprintf("log/%s.log", id))
		c.String(http.StatusOK, string(data))
	})

	api.GET("cron/kill/:id", func(c *gin.Context) {
		id := c.Param("id")
		_ = tasks[id].Process.Kill()
		_ = tasks[id].Wait()
		c.String(http.StatusOK, "ok")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hey!")
	})

	//var wsUpgrader = websocket.Upgrader{
	//	ReadBufferSize:  1024,
	//	WriteBufferSize: 1024,
	//	CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	},
	//}

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

	server := &http.Server{
		Addr:    *ip,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("L&S error: %s\n", err)
		}
	}()

	down := make(chan os.Signal, 1)
	signal.Notify(down, syscall.SIGINT, syscall.SIGTERM)
	<-down

	fmt.Println("Server is shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}

	fmt.Println("Server has stopped gracefully.")
}
