module github.com/DouYin/service

go 1.16

replace github.com/DouYin/common => ../common

require (
	github.com/DouYin/common v0.0.0
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/martian v2.1.0+incompatible
	github.com/spf13/viper v1.9.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.16
)
