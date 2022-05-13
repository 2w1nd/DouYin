package service

import (
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/service/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type FeedService struct {
}

func (fs *FeedService) Feed(c *gin.Context) []vo.VideoVo {
	db := global.DB

	var videoVo []vo.VideoVo
	//db.Preload("Author").
	//	Raw("SELECT \n\tv.video_id AS video_id, \n\tv.author_id AS author_id, \n\tu.username AS `name`, \n\tu.follow_count AS follow_count, \n\tu.follower_count AS follower_count, \n\tf.is_deleted AS is_follow, \n\tv.path AS play_url, \n\tv.cover_path AS cover_url, \n\tv.favorite_count AS favorite_count, \n\tv.comment_count AS comment_count,\n\t1 AS is_favorite\nFROM \n\tdouyin_video v,\n\tdouyin_user u,\n\tdouyin_follow f\nWHERE\n\tv.author_id = u.user_id AND f.followed_user_id = u.user_id").
	//	Scan(&videoVo)

	db.Debug().Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Debug().
			Table("douyin_user u").
			Select("\tu.user_id AS user_id,\n\tu.username AS `name`,\n\tu.follow_count AS follow_count,\n\tu.follower_count AS follower_count,\n\t0 AS is_follow")
	}).
		Table("douyin_video v, douyin_user u").
		Where("v.author_id = u.user_id").
		Select("\tv.video_id AS video_id,\n\tu.user_id AS user_id,\n\tu.username AS `name`,\n\tu.follow_count AS follow_count,\n\tu.follower_count AS follower_count,\n\t0 AS is_follow,\n\tv.path AS play_url,\n\tv.cover_path AS cover_url,\n\tv.favorite_count AS favorite_count,\n\tv.comment_count AS comment_count,\n\t0 AS is_favorite").
		Find(&videoVo)

	//var count int64
	//rows, _ := db.Select("\tv.video_id AS video_id,\n\tv.author_id AS author_id,\n\tu.username AS `name`,\n\tu.follow_count AS follow_count,\n\tu.follower_count AS follower_count,\n\t0 AS is_follow,\n\tv.path AS play_url,\n\tv.cover_path AS cover_url,\n\tv.favorite_count AS favorite_count,\n\tv.comment_count AS comment_count,\n\t0 AS is_favorite").
	//	Table("douyin_video v").
	//	Joins("left join douyin_user u on v.author_id = u.user_id").Rows()
	//
	//log.Println(count)
	//
	//
	//for rows.Next() {
	//	var video vo.VideoVo
	//	_ = db.ScanRows(rows, video)
	//	videoVo = append(videoVo, video)
	//}

	log.Println(videoVo)
	return videoVo
}
