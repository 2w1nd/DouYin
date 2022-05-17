package config

import (
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuOSS struct {
	Cfg       storage.Config
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	AccessKey string `mapstructure:"accessKey" json:"accessKey" yaml:"accessKey"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
}
