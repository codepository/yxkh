package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/codepository/yxkh/config"
	"github.com/codepository/yxkh/conmgr"
	"github.com/codepository/yxkh/controller"
	"github.com/codepository/yxkh/model"
	"github.com/codepository/yxkh/router"
)

var conf = *config.Config

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(10)
	if err := goMain(); err != nil {
		os.Exit(1)
	}
}
func goMain() error {
	// 启动数据库连接
	model.StartDB()
	defer func() {
		model.StopDB()
	}()
	// 启动连接管理器
	conmgr.New()
	defer func() {
		conmgr.Conmgr.Stop()
	}()
	// 启动远程接口调用客户端
	model.StartRemoteClient()
	defer func() {
		model.StopRemoteClient()
	}()
	// 启动函数路由
	controller.SetRouters()
	// http 路由
	mux := router.Mux
	// 启动服务
	readTimeout, err := strconv.Atoi(conf.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := strconv.Atoi(conf.WriteTimeout)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(writeTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	// 监听关闭请求和关闭信号（Ctrl+C）
	interrupt := interruptListener(server)
	log.Printf("the application start up at port%s\n", server.Addr)
	if conf.TLSOpen == "true" {
		err = server.ListenAndServeTLS(conf.TLSCrt, conf.TLSKey)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Printf("Server err: %v", err)
		return err
	}
	<-interrupt
	return nil
}
