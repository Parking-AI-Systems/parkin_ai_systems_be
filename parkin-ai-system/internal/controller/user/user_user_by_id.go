package user

import (
	"context"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/middleware"
	"parkin-ai-system/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error) {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("Request context not found")
	}

	if !middleware.CheckResourceOwnership(r, strconv.FormatInt(req.Id, 10)) {
		return nil, gerror.New("Forbidden: You can only access your own profile")
	}

	res, err = service.User().UserById(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
