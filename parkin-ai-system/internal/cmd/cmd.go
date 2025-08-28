package cmd

import (
	"context"
	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/controller/parking_lot"
	"parkin-ai-system/internal/controller/user"
	"parkin-ai-system/internal/controller/vehicles"
	"parkin-ai-system/internal/middleware"

	_ "parkin-ai-system/internal/logic/parking_lot"

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
			vehiclesCtrl := &vehicles.ControllerVehicles{}
			parkingLotCtrl := parking_lot.NewParkingLot()

			s := g.Server()

			s.Logger().SetHandlers(glog.HandlerJson)
			s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)

			s.Group("/backend/parkin/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// Public routes (no auth required)
				group.POST("/user/register", userCtrl.Register)
				group.POST("/user/login", userCtrl.UserLogin)
				group.POST("/user/refresh", userCtrl.RefreshToken)
				group.GET("/parking-lots/{id}", parkingLotCtrl.ParkingLotDetail)


				// Protected routes (auth required)
				group.Group("/", func(authGroup *ghttp.RouterGroup) {
					authGroup.Middleware(middleware.Auth)
					// Vehicles CRUD
					authGroup.POST("/vehicles", vehiclesCtrl.Create)
					authGroup.GET("/vehicles", vehiclesCtrl.List)
					authGroup.GET("/vehicles/{id}", vehiclesCtrl.Get)
					authGroup.PUT("/vehicles/{id}", vehiclesCtrl.Update)
					authGroup.DELETE("/vehicles/{id}", vehiclesCtrl.Delete)

					authGroup.GET("/user/profile", userCtrl.UserProfile)
					authGroup.POST("/parking-lots", parkingLotCtrl.ParkingLotAdd)
				})

				group.Group("/", func(userGroup *ghttp.RouterGroup) {
					userGroup.Middleware(middleware.UserOrAdmin)
					userGroup.GET("/users/:id", userCtrl.UserById)
				})
				//admin
				group.Group("/admin", func(adminGroup *ghttp.RouterGroup) {
					adminGroup.Middleware(middleware.AdminOnly)
					adminGroup.GET("/users", userCtrl.GetAllUsers)
					adminGroup.DELETE("/users/:id", userCtrl.DeleteUser)
					adminGroup.PUT("/users/:id/role", userCtrl.UpdateUserRole)
				})
			})
			s.Run()
			return nil
		},
	}
)

func CORS(r *ghttp.Request) {
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(200)
		r.Exit()
		return
	}
	r.Middleware.Next()
}
