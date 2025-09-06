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
type FavoriteAddReq struct {
	LotId int64 `json:"lotId"`
}

type FavoriteAddRes struct {
	Id int64 `json:"id"`
}

type FavoriteListReq struct {
	LotName  string `json:"lotName"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type FavoriteItem struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	LotId      int64  `json:"lot_id"`
	LotName    string `json:"lot_name"`
	LotAddress string `json:"lot_address"`
	CreatedAt  string `json:"created_at"`
}

type FavoriteListRes struct {
	List  []FavoriteItem `json:"list"`
	Total int            `json:"total"`
}

type FavoriteDeleteReq struct {
	Id int64 `json:"id"`
}

type FavoriteDeleteRes struct {
	Message string `json:"message"`
}