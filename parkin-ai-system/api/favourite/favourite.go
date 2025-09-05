package favourite

import (
	"github.com/gogf/gf/v2/frame/g"
)

type FavouriteAddReq struct {
	g.Meta `path:"/favourites" method:"post" tags:"Favourite" summary:"Thêm bãi xe vào yêu thích" security:"BearerAuth"`
	LotId  int64 `json:"lot_id" v:"required"`
}

type FavouriteAddRes struct {
	Success bool `json:"success"`
}

type FavouriteDeleteReq struct {
	g.Meta `path:"/favourites/{lot_id}" method:"delete" tags:"Favourite" summary:"Xoá bãi xe khỏi yêu thích" security:"BearerAuth"`
	LotId  int64 `json:"lot_id" v:"required"`
}

type FavouriteDeleteRes struct {
	Success bool `json:"success"`
}

type FavouriteListReq struct {
	g.Meta   `path:"/favourites" method:"get" tags:"Favourite" summary:"Danh sách bãi xe yêu thích" security:"BearerAuth"`
	Page     int `json:"page"     v:"min:1#Page must be at least 1"    dc:"1"`
	PageSize int `json:"page_size" v:"min:1|max:100#PageSize must be between 1 and 100" dc:"10"`
}

type FavouriteListRes struct {
	Favourites []FavouriteInfo `json:"favourites"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	Total      int             `json:"total"`
}

type FavouriteInfo struct {
	Id        int64  `json:"id"`
	LotId     int64  `json:"lotId"`
	LotName   string `json:"lotName"`
	Address   string `json:"address"`
	CreatedAt string `json:"createdAt"`
}

type FavouriteStatusReq struct {
	g.Meta `path:"/favourites/{lot_id}/status" method:"get" tags:"Favourite" summary:"Kiểm tra trạng thái yêu thích" security:"BearerAuth"`
	LotId  int64 `json:"lot_id" v:"required"`
}

type FavouriteStatusRes struct {
	IsFavourite bool `json:"is_favourite"`
}
