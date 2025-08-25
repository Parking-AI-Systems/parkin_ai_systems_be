package user

import (
	"github.com/gogf/gf/v2/frame/g"
)

type RefreshTokenReq struct {
	g.Meta       `path:"/user/refresh" tags:"User" method:"post" summary:"Refresh Token"`
	RefreshToken string `json:"refresh_token" v:"required#Refresh token is required"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
