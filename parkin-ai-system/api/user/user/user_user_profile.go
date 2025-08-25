package user

import "github.com/gogf/gf/v2/frame/g"

type UserProfileReq struct {
	g.Meta `path:"/user/profile" method:"get" tags:"User" summary:"Get user profile"`
}

type UserProfileRes struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}
