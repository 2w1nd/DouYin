package vo

import (
	"github.com/DouYin/common/entity/response"
)

type VideoListVo struct {
	response.Response
	NextTime  int64     `json:"next_time"`
	VideoList []VideoVo `json:"video_list,omitempty"`
}

type VideoData struct {
	NextTime  int64     `json:"next_time"`
	VideoList []VideoVo `json:"video_list,omitempty"`
}

type VideoVo struct {
	VideoID       uint64   `json:"id,omitempty"`
	Author        AuthorVo `json:"author"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount uint32   `json:"favorite_count,omitempty"`
	CommentCount  uint32   `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
	Title   	  string   `json:"title,omitempty"`
}

type AuthorVo struct {
	UserID        uint64 `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint32 `json:"follow_count,omitempty"`
	FollowerCount uint32 `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
