package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) Register(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error) {
	res, err = service.User().SignUp(ctx, req)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError)
	}
	return
}
