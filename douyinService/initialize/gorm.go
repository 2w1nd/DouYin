package initialize

import (
	model2 "github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

// MysqlTables
// @Description: 注册数据库表
// @param: db
func MysqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model2.Comment{},
		model2.User{},
		model2.Video{},
		model2.Follow{},
		model2.Favorite{},
	)
	if err != nil {
		log.Println("register table failed")
		os.Exit(0)
	}
	log.Println("register table success")
}

// GormMysql
// @Description: 初始化MySql数据库
// @return: *gorm.DB
func GormMysql() *gorm.DB {
	m := global.CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// gormConfig
// @Description: 数据库配置
// @return: *gorm.Config
func gormConfig() *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	return config
}
