package global

import (
	"github.com/DouYin/service/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	REDIS  *redis.Client
	CONFIG config.Server
	VP     *viper.Viper
)
