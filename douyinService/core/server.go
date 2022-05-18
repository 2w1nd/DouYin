package core

import (
	"fmt"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/initialize"
	"log"
)

type server interface {
	ListenAndServer() error
}

func RunWindowsServer() {
	initialize.Redis()
	initialize.Snowflake()
	Router := initialize.Routers() // 初始化路由

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, Router)
	log.Println("端口号" + address)
	fmt.Println(`
		DouYin项目成功启动
		====================================================================
	`)
	log.Println(s.ListenAndServe().Error())
}
