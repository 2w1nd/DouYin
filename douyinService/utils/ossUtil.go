package utils

import (
	"fmt"
	"github.com/DouYin/service/global"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"golang.org/x/net/context"
	"mime/multipart"
)

func UploadVideo(key string, data *multipart.FileHeader) {
	//打开文件
	file, err := data.Open()
	defer file.Close()
	if err != nil {

	}
	var byteContainer []byte
	byteContainer = make([]byte, data.Size)
	file.Read(byteContainer)
	dataLen := int64(len(byteContainer))

	//上传凭证
	putPolicy := storage.PutPolicy{
		Scope: global.CONFIG.OSS.Bucket,
	}
	mac := qbox.NewMac(global.CONFIG.OSS.AccessKey, global.CONFIG.OSS.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(&global.CONFIG.OSS.Cfg)
	ret := storage.PutRet{}
	//拓展参数，目前为空
	putExtra := storage.PutExtra{Params: map[string]string{}}

	key = "video/" + key
	err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(byteContainer), dataLen, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
}
