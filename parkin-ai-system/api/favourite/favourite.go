package favourite

import (
	"github.com/gogf/gf/v2/frame/g"
)

type FavouriteAddReq struct {
	g.Meta   `path:"/favourites" method:"post" tags:"Favourite" summary:"Thêm bãi xe vào yêu thích" security:"BearerAuth"`
	LotId    int64  `json:"lot_id" v:"required"`
}

type FavouriteAddRes struct {
	Success bool `json:"success"`
}

type FavouriteDeleteReq struct {
	g.Meta   `path:"/favourites/{lot_id}" method:"delete" tags:"Favourite" summary:"Xoá bãi xe khỏi yêu thích" security:"BearerAuth"`
	LotId    int64  `json:"lot_id" v:"required"`
}

type FavouriteDeleteRes struct {
	Success bool `json:"success"`
}

type FavouriteListReq struct {
	g.Meta   `path:"/favourites" method:"get" tags:"Favourite" summary:"Danh sách bãi xe yêu thích" security:"BearerAuth"`
}

type FavouriteListRes struct {
	Favourites []FavouriteInfo `json:"favourites"`
}

type FavouriteInfo struct {
	Id        int64  `json:"id"`
	LotId     int64  `json:"lot_id"`
	CreatedAt string `json:"created_at"`
}

type FavouriteStatusReq struct {
	g.Meta   `path:"/favourites/{lot_id}/status" method:"get" tags:"Favourite" summary:"Kiểm tra trạng thái yêu thích" security:"BearerAuth"`
	LotId    int64  `json:"lot_id" v:"required"`
}

type FavouriteStatusRes struct {
	IsFavourite bool `json:"is_favourite"`
}
