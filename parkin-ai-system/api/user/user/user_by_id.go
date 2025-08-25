package user

import "github.com/gogf/gf/v2/frame/g"

type UserByIdReq struct {
	g.Meta `path:"/users/:id" method:"get" tags:"User" summary:"Get user by ID"`
	Id     int64 `json:"id" in:"path" name:"id"`
}

type UserByIdRes struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}
