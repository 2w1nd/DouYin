package initialize

import (
	"github.com/DouYin/global"
	"github.com/DouYin/model"
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
		model.Comment{},
		model.User{},
		model.Video{},
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
