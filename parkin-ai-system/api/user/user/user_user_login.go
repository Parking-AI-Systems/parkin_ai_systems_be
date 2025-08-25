package user

import (
	"parkin-ai-system/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type UserLoginReq struct {
	g.Meta `path:"/user/login" tags:"User" method:"post" summary:"User login"`
	model.SignInInput
}

type UserLoginRes struct {
	model.SignInOutput
}
