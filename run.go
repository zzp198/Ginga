package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

//go:embed frontend/*
var frontendFS embed.FS

//go:embed resource/*
var resourceFS embed.FS

func main() {
	name, _ := os.Executable()
	fmt.Println(name)

	ip := flag.String("ip", "0.0.0.0:8888", "")
	flag.Parse()

	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "daemon" {
		_ = exec.Command("pkill", "Crond").Run()

		cmd := exec.Command(os.Args[0], os.Args[2:]...)
		// 碰到的第一个坑,父进程结束时,会向子进程发送HUP,TERM指令,导致子进程会跟随父进程一块结束.
		// SysProcAttr.Setpgid设置为true,使子进程的进程组ID与其父进程不同.(KILL强杀也可以)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		} else {
			log.Println(fmt.Sprintf("%s [PID] %d running...", os.Args[0], cmd.Process.Pid))
		}
		return
	}

	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	frontend, err := template.New("").Delims("{%", "%}").ParseFS(frontendFS, "frontend/*")
	if err != nil {
		log.Fatal(err.Error())
	}
	r.SetHTMLTemplate(frontend)

	resource, err := fs.Sub(resourceFS, "resource")
	if err != nil {
		log.Fatal(err.Error())
	}
	r.StaticFS("/resource", http.FS(resource))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/m3u8", func(c *gin.Context) {
		c.HTML(http.StatusOK, "m3u8.html", nil)
	})

	srv := &http.Server{Addr: *ip, Handler: r}
	srv.RegisterOnShutdown(func() {
		log.Println(fmt.Sprintf("Server is shutting down"))
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println(fmt.Sprintf("Server closed under request"))
			} else {
				log.Println(err)
			}
		}
	}()

	down := make(chan os.Signal, 1)
	signal.Notify(down, syscall.SIGINT, syscall.SIGTERM)
	<-down

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	log.Println("Server has stopped gracefully.")
}
