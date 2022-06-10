package cache

import (
	"context"
	"encoding/json"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/service/global"
	"log"
	"strconv"
)

type CommentCache struct {
}

// GetCommentCount 查评论数量
func (cc *CommentCache) GetCommentCount(videoId uint64) uint32 {
	CommentString := "videoComment:comment"
	var commentVos []vo.CommentVo
	var commentCount uint32
	data, _ := global.REDIS.Get(context.Background(), CommentString+strconv.FormatUint(videoId, 10)).Result()
	if data != "" {
		log.Println("评论从缓存中找")
		err := json.Unmarshal([]byte(data), &commentVos)
		if err != nil {
			log.Println(err)
		}
		commentCount = uint32(len(commentVos))
	} else { // 从数据库中找
		log.Println("评论从数据库中找")
		video := videoRepository.GetVideoByVideoId(videoId)
		commentCount = video.CommentCount
	}
	return commentCount
}
