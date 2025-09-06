package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserRegister(ctx context.Context, req *user.UserRegisterReq) (res *user.UserRegisterRes, err error) {
	// Map API request to entity request
	input := &entity.UserRegisterReq{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Phone:     req.Phone,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		AvatarUrl: req.AvatarUrl,
	}

	// Call service
	registerRes, err := service.User().SignUp(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserRegisterRes{
		UserId:    registerRes.UserId,
		Username:  registerRes.Username,
		Email:     registerRes.Email,
		Phone:     registerRes.Phone,
		FullName:  registerRes.FullName,
		Gender:    registerRes.Gender,
		BirthDate: registerRes.BirthDate,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
