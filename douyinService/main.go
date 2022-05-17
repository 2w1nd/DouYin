package main

import (
	"github.com/DouYin/service/core"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/initialize"
)

func main() {

	global.VP = core.Viper()      // 初始化viper
	global.DB = initialize.Gorm() // gorm连接数据库
	initialize.InitOSS()          //初始化OSS配置
	if global.DB != nil {
		//initialize.MysqlTables(global.DB) // 初始化表
	}
	core.RunWindowsServer()
}
