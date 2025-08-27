package user

import "github.com/gogf/gf/v2/frame/g"

type RegisterReq struct {
	g.Meta    `path:"/user/register" method:"post" tags:"User" summary:"User Registration"`
	Username  string `json:"username" v:"required#Username is required"`
	Password  string `json:"password" v:"required#Password is required"`
	Email     string `json:"email" v:"required#Email is required|email#Invalid email format"`
	Phone     string `json:"phone" v:"regex:^\\d{10,15}$#Invalid phone number"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender" v:"in:male,female,other#Gender must be male, female, or other"`
	BirthDate string `json:"birth_date" v:"date#Invalid date format"`
}

type RegisterRes struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}
