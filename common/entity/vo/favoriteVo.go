package vo

import (
	"github.com/DouYin/common/entity/response"
)

type FavoriteListVo struct {
	//匿名嵌入，可以直接访问叶子属性而不需要给出完整路径
	response.Response
	VideoList []VideoVo `json:"video_list"`
}
