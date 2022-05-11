package vo

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/model"
)

type VideoListVo struct {
	response.Response
	NextTime  int64   `json:"next_time,omitempty"`
	VideoList []VideoVo `json:"video_list,omitempty"`
}

type VideoVo struct {
	Id            int64  `json:"id,omitempty"`
	Author        model.User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}
