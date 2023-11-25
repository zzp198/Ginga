package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ip := flag.String("ip", "0.0.0.0:8080", "")
	flag.Parse()

	//if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "daemon" {
	//	cmd := exec.Command(os.Args[0], append(os.Args[2:], "CronTag")...)
	//	// 碰到的第一个坑,父进程结束时,会向子进程发送HUP,TERM指令,导致子进程会跟随父进程一块结束.
	//	// SysProcAttr.Setpgid设置为true,使子进程的进程组ID与其父进程不同.(KILL强杀也可以)
	//	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	//	if err := cmd.Start(); err != nil {
	//		slog.Error(err.Error())
	//	} else {
	//		slog.Info(fmt.Sprintf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid))
	//	}
	//	return
	//}

	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("m3u8", func(c *gin.Context) {
		c.HTML(http.StatusOK, "m3u8.html", nil)
	})

	srv := &http.Server{Addr: *ip, Handler: r}
	srv.RegisterOnShutdown(func() {
		slog.Info(fmt.Sprintf("Server is shutting down"))
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Warn(fmt.Sprintf("Server closed under request"))
			} else {
				slog.Error(err.Error())
			}
		}
	}()

	down := make(chan os.Signal, 1)
	signal.Notify(down, syscall.SIGINT, syscall.SIGTERM)
	<-down

	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error(err.Error())
	}

	slog.Info(fmt.Sprintf("Server has stopped gracefully."))
}
