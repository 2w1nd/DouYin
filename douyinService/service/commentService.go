package service

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type CommentService struct {
}

// AddCommentDemo
// @Description: 测试栗子
// @receiver: e
// @param: c
// @return: err
func (cs *CommentService) AddCommentDemo(c model.Comment) (err error) {
	err = global.DB.Create(&c).Error
	return err
}
