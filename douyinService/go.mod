module github.com/DouYin/service

go 1.16

replace github.com/DouYin/common => ../common

require (
	github.com/DouYin/common v0.0.0
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.8.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/martian v2.1.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/qiniu/go-sdk/v7 v7.12.1
	github.com/spf13/viper v1.9.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.16
)
