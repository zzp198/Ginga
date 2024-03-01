package main

import (
	"fmt"
	"image"
	"os"
)

//import (
//	"context"
//	"errors"
//	"flag"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//	"os"
//	"os/exec"
//	"os/signal"
//	"strings"
//	"syscall"
//	"time"
//)
//
//var Cmds = make(map[string]*exec.Cmd)
//
//func main() {
//	ip := flag.String("ip", "0.0.0.0:5200", "")
//	flag.Parse()
//
//	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "daemon" {
//
//		arr, _ := exec.Command("lsof", "-t", "Ginga.lock").Output()
//		if len(arr) > 0 {
//			fmt.Println(fmt.Sprintf("检测到已有Ginga程序运行, PID %s", string(arr)))
//			_ = exec.Command("kill", string(arr)).Run()
//		}
//
//		// 碰到的第一个坑,父进程结束时,会向子进程发送HUP,TERM指令,导致子进程会跟随父进程一块结束.
//		// SysProcAttr.Setpgid设置为true,使子进程的进程组ID与其父进程不同.(KILL强杀也可以)
//		cmd := exec.Command(os.Args[0], os.Args[2:]...)
//		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//
//		if err := cmd.Start(); err != nil {
//			log.Fatal(err)
//		} else {
//			log.Println(fmt.Sprintf("%s [PID] %d running...", os.Args[0], cmd.Process.Pid))
//		}
//
//		return
//	}
//
//	//LOCK_SH 共享锁,多个进程可以使用同一把锁,常用作读共享锁.
//	//LOCK_EX 排他锁,同时只允许一个进程使用,常被用作写锁.
//	//LOCK_UN 释放锁.
//	//        如果文件被其他进程锁住,进程会被阻塞直到锁释放.
//	//LOCK_NB 如果文件被其他进程锁住,会返回错误 EWOULDBLOCK
//	lock, e := os.Create("Ginga.lock")
//	if e != nil {
//		log.Fatalln(e)
//	}
//	defer lock.Close()
//
//	HysteriaInit()
//
//	//gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//	r.Use(gin.Recovery(), gin.Logger())
//
//	r.GET("/hy2", func(c *gin.Context) {
//		isrun := true
//
//		if Cmds["Hysteria"] == nil || Cmds["Hysteria"].ProcessState != nil && Cmds["Hysteria"].ProcessState.Exited() || Cmds["Hysteria"].Process == nil {
//			isrun = false
//		}
//
//		config := HysteriaGetConfig()
//		logcont := HysteriaLog()
//
//		c.String(http.StatusOK, fmt.Sprintf("%t\n%s\n%s", isrun, config, logcont))
//	})
//
//	r.GET("/hy2/start", func(c *gin.Context) {
//		e := HysteriaStart()
//		if e != nil {
//			c.String(http.StatusOK, e.Error())
//			return
//		}
//		c.String(http.StatusOK, "ok")
//	})
//
//	r.GET("/hy2/stop", func(c *gin.Context) {
//		HysteriaStop()
//		c.String(http.StatusOK, "ok")
//	})
//
//	r.GET("/chunked", func(c *gin.Context) {
//
//		c.Header("Content-Type", "text/html")
//		c.Writer.WriteHeader(http.StatusOK)
//
//		c.Writer.Write([]byte(`<html><body>`))
//		c.Writer.Flush()
//
//		for i := 0; i < 10; i++ {
//			c.Writer.Write([]byte(fmt.Sprintf(`<h3>%d<h3>`, i)))
//			c.Writer.Flush()
//			time.Sleep(1 * time.Second)
//		}
//
//		c.Writer.Write([]byte(`</body></html>`))
//		c.Writer.Flush()
//	})
//
//	srv := &http.Server{Addr: *ip, Handler: r}
//	srv.RegisterOnShutdown(func() {
//		log.Println(fmt.Sprintf("Server is shutting down"))
//	})
//
//	go func() {
//		if err := srv.ListenAndServe(); err != nil {
//			if errors.Is(err, http.ErrServerClosed) {
//				log.Println(fmt.Sprintf("Server closed under request"))
//			} else {
//				log.Println(err)
//			}
//		}
//	}()
//
//	down := make(chan os.Signal, 1)
//	signal.Notify(down, syscall.SIGINT, syscall.SIGTERM)
//	<-down
//
//	if err := srv.Shutdown(context.Background()); err != nil {
//		log.Println(err)
//	}
//	log.Println("Server has stopped gracefully.")
//}
//
//func HysteriaInit() {
//	path := "data/Hysteria/hysteria-linux-amd64"
//
//	os.Chmod(path, 0777)
//
//	Cmds["Hysteria"] = exec.Command(path, "server", "-c", "data/Hysteria/config.yaml")
//
//	lf, _ := os.Create("data/Hysteria/Hysteria.log")
//	Cmds["Hysteria"].ExtraFiles = []*os.File{lf}
//
//	Cmds["Hysteria"].Stdout = Cmds["Hysteria"].ExtraFiles[0]
//	Cmds["Hysteria"].Stderr = Cmds["Hysteria"].ExtraFiles[0]
//}
//
//func HysteriaGetConfig() string {
//	data, _ := os.ReadFile("data/Hysteria/config.yaml")
//	return string(data)
//}
//
//func HysteriaSetConfig(config string) {
//	_ = os.WriteFile("data/Hysteria/config.yaml", []byte(config), 0644)
//}
//
//func HysteriaLog() string {
//	data, _ := os.ReadFile("data/Hysteria/Hysteria.log")
//	return string(data)
//}
//
//func HysteriaStart() error {
//	return Cmds["Hysteria"].Start()
//}
//
//func HysteriaStop() {
//	Cmds["Hysteria"].Process.Release()
//	Cmds["Hysteria"].Wait()
//}
//
//func IsProcessRunning(c *exec.Cmd) bool {
//	if c.Process != nil && c.ProcessState == nil {
//		return true
//	}
//	return false
//}
//
//func GetResult() {
//	//	https://api.streamtape.com
//	//ftp.streamtape.com
//}

func main() {
	// 打开文件夹
	dir, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	// 遍历文件夹中的文件
	files, err := dir.ReadDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 获取文件名
		name := file.Name()

		// 打开文件
		f, err := os.Open(fmt.Sprintf("input/%s", name))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// 读取文件内容
		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 显示图片
		switch img.(type) {
		case *image.RGBA:
			fmt.Println("这是 RGBA 图片")
		case *image.NRGBA:
			fmt.Println("这是 NRGBA 图片")
		case *image.Gray:
			fmt.Println("这是 Gray 图片")
		case *image.Gray16:
			fmt.Println("这是 Gray16 图片")
		case *image.CMYK:
			fmt.Println("这是 CMYK 图片")
		case *image.YCbCr:
			fmt.Println("这是 YCbCr 图片")
		default:
			fmt.Println("无法识别的图片格式")
		}
	}
}
