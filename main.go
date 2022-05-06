package main

import (
	"github.com/DouYin/core"
	"github.com/DouYin/global"
	"github.com/DouYin/initialize"
)

func main() {

	global.VP = core.Viper()      // 初始化viper
	global.DB = initialize.Gorm() // gorm连接数据库
	if global.DB != nil {
		initialize.MysqlTables(global.DB) // 初始化表
	}
	core.RunWindowsServer()
}
