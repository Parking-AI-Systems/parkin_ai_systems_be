package favourite

import "github.com/gogf/gf/v2/frame/g"

type FavoriteAddReq struct {
	g.Meta `path:"/favorites" tags:"Favorite" method:"POST" summary:"Add a parking lot to favorites" description:"Adds a parking lot to the authenticated user's favorite list." middleware:"middleware.Auth"`
	LotId  int64 `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
}

type FavoriteAddRes struct {
	Id int64 `json:"id"`
}

type FavoriteListReq struct {
	g.Meta   `path:"/favorites" tags:"Favorite" method:"GET" summary:"List favorite parking lots" description:"Retrieves a paginated list of the authenticated user's favorite parking lots or all favorites for admins." middleware:"middleware.Auth"`
	LotName  string `json:"lot_name"`
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
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
	g.Meta `path:"/favorites/:id" tags:"Favorite" method:"DELETE" summary:"Delete a favorite parking lot" description:"Removes a parking lot from the authenticated user's favorite list. Admins can delete any favorite." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Favorite ID is required|Favorite ID must be positive"`
}

type FavoriteDeleteRes struct {
	Message string `json:"message"`
}
