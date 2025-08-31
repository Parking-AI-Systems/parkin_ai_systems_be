package cmd

import (
	"context"
	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/controller/parking_lot"
	"parkin-ai-system/internal/controller/parking_lot_review"
	"parkin-ai-system/internal/controller/other_service"
	"parkin-ai-system/internal/controller/favourite"
	"parkin-ai-system/internal/controller/parking_slot"
	"parkin-ai-system/internal/controller/user"
	"parkin-ai-system/internal/controller/vehicles"
	"parkin-ai-system/internal/middleware"

	_ "parkin-ai-system/internal/logic/parking_lot"
	_ "parkin-ai-system/internal/logic/parking_lot_review"
	_ "parkin-ai-system/internal/logic/other_service"
	_ "parkin-ai-system/internal/logic/favourite"
	_ "parkin-ai-system/internal/logic/parking_slot"

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
					   parkingLotReviewCtrl := parking_lot_review.NewParkingLotReview()
					   otherServiceCtrl := other_service.NewOtherService()
					   favouriteCtrl := favourite.NewFavourite()
					   parkingSlotCtrl := parking_slot.NewParkingSlot()

			s := g.Server()

			s.Logger().SetHandlers(glog.HandlerJson)
			s.Use(CORS, ghttp.MiddlewareHandlerResponse)

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
				       authGroup.GET("/parking-lots", parkingLotCtrl.ParkingLotList)
				       authGroup.PUT("/parking-lots/{id}", parkingLotCtrl.ParkingLotUpdate)
				       authGroup.DELETE("/parking-lots/{id}", parkingLotCtrl.ParkingLotDelete)

				       // Parking Lot Review CRUD
				       authGroup.POST("/parking-lot-reviews", parkingLotReviewCtrl.ParkingLotReviewAdd)
				       authGroup.GET("/parking-lot-reviews", parkingLotReviewCtrl.ParkingLotReviewList)
				       authGroup.GET("/parking-lot-reviews/{id}", parkingLotReviewCtrl.ParkingLotReviewDetail)
				       authGroup.PUT("/parking-lot-reviews/{id}", parkingLotReviewCtrl.ParkingLotReviewUpdate)
				       authGroup.DELETE("/parking-lot-reviews/{id}", parkingLotReviewCtrl.ParkingLotReviewDelete)

				       // Other Service CRUD
				       authGroup.POST("/other-services", otherServiceCtrl.OtherServiceAdd)
				       authGroup.GET("/other-services", otherServiceCtrl.OtherServiceList)
				       authGroup.GET("/other-services/{id}", otherServiceCtrl.OtherServiceDetail)
				       authGroup.PUT("/other-services/{id}", otherServiceCtrl.OtherServiceUpdate)
				       authGroup.DELETE("/other-services/{id}", otherServiceCtrl.OtherServiceDelete)

				       // Favourite CRUD
				       authGroup.POST("/favourites", favouriteCtrl.FavouriteAdd)
				       authGroup.GET("/favourites", favouriteCtrl.FavouriteList)
				       authGroup.GET("/favourites/{lot_id}/status", favouriteCtrl.FavouriteStatus)
				       authGroup.DELETE("/favourites/{lot_id}", favouriteCtrl.FavouriteDelete)

									   // Parking Slot CRUD
									   authGroup.POST("/parking-slots", parkingSlotCtrl.ParkingSlotAdd)
									   authGroup.GET("/parking-slots", parkingSlotCtrl.ParkingSlotList)
									   authGroup.PUT("/parking-slots/{id}", parkingSlotCtrl.ParkingSlotUpdate)
									   authGroup.DELETE("/parking-slots/{id}", parkingSlotCtrl.ParkingSlotDelete)
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
