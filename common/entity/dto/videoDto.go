package dto

import "github.com/DouYin/common/entity/response"

type VideoDto struct {
	Id            uint64  `json:"id,omitempty"`
	Author        UserDto `json:"author,omitempty"`
	PlayURL       string  `json:"play_url,omitempty"`
	CoverURL      string  `json:"cover_url,omitempty"`
	Title         string  `json:"title,omitempty"`
	FavoriteCount uint32  `json:"favorite_count"`
	CommentCount  uint32  `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
}
type VideoListDto struct {
	response.Response
	VideoList []VideoDto `json:"video_list,omitempty"`
}
