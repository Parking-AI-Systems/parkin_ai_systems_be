package user

import "github.com/gogf/gf/v2/frame/g"

type UserUpdateProfileReq struct {
	g.Meta    `path:"/user/update-profile" method:"post" tags:"User" summary:"Update user profile"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type UserUpdateProfileRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
