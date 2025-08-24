package user

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UserLoginReq struct {
	g.Meta   `path:"/user/login" tags:"User" method:"post" summary:"User login"`
	Email    string `json:"email" v:"required|email#Email is required|Invalid email format"`
	Password string `json:"password" v:"required#Password is required"`
}

type UserLoginRes struct {
	AccessToken string `json:"access_token"`
}
