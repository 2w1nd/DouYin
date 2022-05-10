module github.com/DouYin/service

go 1.16

replace github.com/DouYin/common => ../common

require (
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/spf13/viper v1.9.0
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.16
	github.com/DouYin/common v0.0.0
)


