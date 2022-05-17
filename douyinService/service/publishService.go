package service

import (
	"github.com/DouYin/service/context"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"sync"
)

type PublishService struct {
}

func (ps *PublishService) Publish(userContext context.UserContext, data *multipart.FileHeader) {
	task := sync.WaitGroup{}
	task.Add(1)
	go func() {
		utils.UploadVideo("123", data)
		task.Done()
	}()
	task.Wait()
}

func (ps *PublishService) PublishList(c *gin.Context) {

}
