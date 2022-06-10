package initialize

import (
	"github.com/DouYin/service/global"
	"github.com/qiniu/go-sdk/v7/storage"
)

func InitOSS() {
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	global.CONFIG.OSS.Cfg = cfg
}
