// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorites is the golang structure for table favorites.
type Favorites struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	UserId    int64       `json:"userId"    orm:"user_id"    description:""`
	LotId     int64       `json:"lotId"     orm:"lot_id"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
type FavoritesInput struct {
	Id        int64       `json:"id"`
	UserId    int64       `json:"userId"`
	LotId     int64       `json:"lotId"`
	CreatedAt *gtime.Time `json:"createdAt"`
}
type FavoritesOutput struct {
	Id        int64       `json:"id"`
	UserId    int64       `json:"userId"`
	LotId     int64       `json:"lotId"`
	CreatedAt *gtime.Time `json:"createdAt"`
}
type FavoriteDelRes struct {
	Success bool `json:"success"`
	LotId   int64 `json:"lotId"`
	UserId  int64 `json:"userId"`
}
type FavoriteListReq struct {
	Page     int `json:"page"     v:"min:1#Page must be at least 1"    dc:"1"`
	PageSize int `json:"pageSize" v:"min:1|max:100#PageSize must be between 1 and 100" dc:"10"`
}
type FavouriteListRes struct {
	Favourites []FavouriteInfo `json:"favourites"`
	Page 	 int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Total 	 int             `json:"total"`
}
type FavouriteInfo struct {
	Id        int64  `json:"id"`
	LotId     int64  `json:"lotId"`
	LotName   string `json:"lotName"`
	Address   string `json:"address"`
	CreatedAt string `json:"createdAt"`
}
type FavouriteStatusReq struct {
	LotId int64 `json:"lotId" v:"required"`
}
type FavouriteStatusRes struct {
	IsFavourite bool `json:"isFavourite"`
}