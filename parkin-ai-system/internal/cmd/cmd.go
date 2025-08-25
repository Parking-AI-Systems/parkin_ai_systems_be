package cmd

import (
	"context"
	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/controller/user"
	"parkin-ai-system/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().SetHandlers(glog.HandlerJson)
			glog.SetHandlers(glog.HandlerJson)

			config.InitConfig(ctx)

			userCtrl := user.NewUser()

			s := g.Server()

			s.Logger().SetHandlers(glog.HandlerJson)

			// Register global middlewares
			s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)

			s.Group("/backend/parkin/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// Public routes (no auth required)
				group.POST("/user/register", userCtrl.Register)
				group.POST("/user/login", userCtrl.UserLogin)
				group.POST("/user/refresh", userCtrl.RefreshToken)

				// Protected routes (auth required)
				group.Group("/", func(authGroup *ghttp.RouterGroup) {
					authGroup.Middleware(middleware.Auth)
					authGroup.POST("/user/logout", userCtrl.UserLogout)
					authGroup.GET("/user/profile", userCtrl.UserProfile)
				})
			})
			s.Run()
			return nil
		},
	}
)
