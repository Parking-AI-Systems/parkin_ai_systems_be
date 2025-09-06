package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserUpdateProfile(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error) {
	// Map API request to entity request
	input := &entity.UserUpdateProfileReq{
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
	updateRes, err := service.User().UserUpdateProfile(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserUpdateProfileRes{
		Success: updateRes.Success,
		Message: updateRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
