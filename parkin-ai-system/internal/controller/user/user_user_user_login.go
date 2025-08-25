package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserLogin(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	res, err = service.User().Login(ctx, req)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	return
}
