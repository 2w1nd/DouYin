package core

import (
	"fmt"
	"github.com/DouYin/global"
	"github.com/DouYin/initialize"
	"log"
)

type server interface {
	ListenAndServer() error
}

func RunWindowsServer() {
	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, Router)
	log.Println(address)
	fmt.Println(`
		DouYin项目成功启动
		====================================================================
	`)
	log.Println(s.ListenAndServe().Error())
}
