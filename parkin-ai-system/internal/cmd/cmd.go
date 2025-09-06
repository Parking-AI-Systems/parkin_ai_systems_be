package cmd

import (
	"context"
	"net/http"
	"sync"
	"time"

	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/controller/favourite"
	"parkin-ai-system/internal/controller/notification"
	"parkin-ai-system/internal/controller/other_service_orders"
	"parkin-ai-system/internal/controller/others_service"
	"parkin-ai-system/internal/controller/parking_lot"
	"parkin-ai-system/internal/controller/parking_lot_review"
	"parkin-ai-system/internal/controller/parking_order"
	"parkin-ai-system/internal/controller/parking_slot"
	"parkin-ai-system/internal/controller/user"
	"parkin-ai-system/internal/controller/vehicle"
	"parkin-ai-system/internal/controller/wallet_transaction"
	"parkin-ai-system/internal/middleware"

	_ "parkin-ai-system/internal/logic/favourite"
	_ "parkin-ai-system/internal/logic/notification"
	_ "parkin-ai-system/internal/logic/other_service"
	_ "parkin-ai-system/internal/logic/other_service_order"
	_ "parkin-ai-system/internal/logic/parking_lot"
	_ "parkin-ai-system/internal/logic/parking_lot_review"
	_ "parkin-ai-system/internal/logic/parking_order"
	_ "parkin-ai-system/internal/logic/parking_slot"
	_ "parkin-ai-system/internal/logic/user"
	_ "parkin-ai-system/internal/logic/vehicle"
	_ "parkin-ai-system/internal/logic/wallet_transaction"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	csrfTokens     = make(map[string]time.Time)
	csrfTokensLock sync.RWMutex
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "Start the HTTP server for the Parkin AI System",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			config.InitConfig(ctx)
			g.Log().SetHandlers(glog.HandlerJson)
			glog.SetHandlers(glog.HandlerJson)

			cfg := config.GetConfig()
			if cfg.Auth.SecretKey == "" {
				g.Log().Fatal(ctx, "JWT_SECRET is required in configuration")
			}
			if cfg.Auth.AccessTokenExpireMinute <= 0 {
				g.Log().Async().Warning(ctx, "Invalid AccessTokenExpireMinute; using default: 15 minutes")
				cfg.Auth.AccessTokenExpireMinute = 15
			}
			if cfg.Auth.RefreshTokenExpireMinute <= 0 {
				g.Log().Async().Warning(ctx, "Invalid RefreshTokenExpireMinute; using default: 7 days")
				cfg.Auth.RefreshTokenExpireMinute = 7 * 24 * 60
			}

			userCtrl := user.NewUser()
			vehiclesCtrl := vehicle.NewVehicle()
			parkingLotCtrl := parking_lot.NewParking_lot()
			parkingLotReviewCtrl := parking_lot_review.NewParking_lot_review()
			otherServiceCtrl := others_service.NewOthers_service()
			otherServiceOrderCtrl := other_service_orders.NewOther_service_orders()
			notificationCtrl := notification.NewNotification()
			walletTransactionCtrl := wallet_transaction.NewWallet_transaction()
			favouriteCtrl := favourite.NewFavourite()
			parkingSlotCtrl := parking_slot.NewParking_slot()
			parkingOrderCtrl := parking_order.NewParking_order()

			s := g.Server()
			s.Logger().SetHandlers(glog.HandlerJson)
			s.Use(CORS, ErrorHandler)

			s.Group("/metrics", func(metricsGroup *ghttp.RouterGroup) {
				metricsGroup.GET("/", ghttp.WrapH(promhttp.Handler()))
			})

			s.Group("/backend/parkin/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(CSRFMiddleware)

				// Public endpoints
				group.Middleware(LogActionMiddleware("user_register"))
				group.Bind(userCtrl.UserRegister)

				group.Middleware(LogActionMiddleware("user_login"))
				group.Bind(userCtrl.UserLogin)

				group.Middleware(LogActionMiddleware("user_refresh"))
				group.Bind(userCtrl.UserRefreshToken)

				// Public parking lot get
				group.Bind(parkingLotCtrl.ParkingLotGet)

				group.Group("/", func(authGroup *ghttp.RouterGroup) {
					authGroup.Middleware(middleware.Auth)
					authGroup.Middleware(LogActionMiddleware("user_logout"))
					authGroup.Bind(userCtrl.UserLogout)
					authGroup.Bind(userCtrl.UserProfile)
					authGroup.Middleware(LogActionMiddleware("user_profile_update"))
					authGroup.Bind(userCtrl.UserUpdateProfile)
					authGroup.Middleware(LogActionMiddleware("vehicles_create"))
					authGroup.Bind(vehiclesCtrl.VehicleAdd)
					authGroup.Bind(vehiclesCtrl.VehicleList)
					authGroup.Bind(vehiclesCtrl.VehicleGet)
					authGroup.Middleware(LogActionMiddleware("vehicles_update"))
					authGroup.Bind(vehiclesCtrl.VehicleUpdate)
					authGroup.Middleware(LogActionMiddleware("vehicles_delete"))
					authGroup.Bind(vehiclesCtrl.VehicleDelete)
					authGroup.Middleware(LogActionMiddleware("parking_lot_create"))
					authGroup.Bind(parkingLotCtrl.ParkingLotAdd)
					authGroup.Bind(parkingLotCtrl.ParkingLotList)
					authGroup.Middleware(LogActionMiddleware("parking_lot_update"))
					authGroup.Bind(parkingLotCtrl.ParkingLotUpdate)
					authGroup.Middleware(LogActionMiddleware("parking_lot_delete"))
					authGroup.Bind(parkingLotCtrl.ParkingLotDelete)
					authGroup.Middleware(LogActionMiddleware("parking_lot_image_delete"))
					authGroup.Bind(parkingLotCtrl.ParkingLotImageDelete)
					authGroup.Middleware(LogActionMiddleware("parking_lot_review_create"))
					authGroup.Bind(parkingLotReviewCtrl.ParkingLotReviewAdd)
					authGroup.Bind(parkingLotReviewCtrl.ParkingLotReviewList)
					authGroup.Bind(parkingLotReviewCtrl.ParkingLotReviewGet)
					authGroup.Middleware(LogActionMiddleware("parking_lot_review_update"))
					authGroup.Bind(parkingLotReviewCtrl.ParkingLotReviewUpdate)
					authGroup.Middleware(LogActionMiddleware("parking_lot_review_delete"))
					authGroup.Bind(parkingLotReviewCtrl.ParkingLotReviewDelete)
					authGroup.Middleware(LogActionMiddleware("other_service_create"))
					authGroup.Bind(otherServiceCtrl.OthersServiceAdd)
					authGroup.Bind(otherServiceCtrl.OthersServiceList)
					authGroup.Bind(otherServiceCtrl.OthersServiceGet)
					authGroup.Middleware(LogActionMiddleware("other_service_update"))
					authGroup.Bind(otherServiceCtrl.OthersServiceUpdate)
					authGroup.Middleware(LogActionMiddleware("other_service_delete"))
					authGroup.Bind(otherServiceCtrl.OthersServiceDelete)
					authGroup.Middleware(LogActionMiddleware("service_order_create"))
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderAdd)
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderList)
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderGet)
					authGroup.Middleware(LogActionMiddleware("service_order_update"))
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderUpdate)
					authGroup.Middleware(LogActionMiddleware("service_order_cancel"))
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderCancel)
					authGroup.Middleware(LogActionMiddleware("service_order_delete"))
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderDelete)
					authGroup.Middleware(LogActionMiddleware("service_order_payment"))
					authGroup.Bind(otherServiceOrderCtrl.OthersServiceOrderPayment)
					authGroup.Middleware(LogActionMiddleware("notification_mark_read"))
					authGroup.Bind(notificationCtrl.NotificationMarkRead)
					authGroup.Bind(notificationCtrl.NotificationList)
					authGroup.Bind(notificationCtrl.NotificationGet)
					authGroup.Middleware(LogActionMiddleware("notification_delete"))
					authGroup.Bind(notificationCtrl.NotificationDelete)
					authGroup.Middleware(LogActionMiddleware("wallet_transaction_create"))
					authGroup.Bind(walletTransactionCtrl.WalletTransactionAdd)
					authGroup.Bind(walletTransactionCtrl.WalletTransactionList)
					authGroup.Middleware(LogActionMiddleware("favourite_create"))
					authGroup.Bind(favouriteCtrl.FavoriteAdd)
					authGroup.Bind(favouriteCtrl.FavoriteList)
					authGroup.Middleware(LogActionMiddleware("favourite_delete"))
					authGroup.Bind(favouriteCtrl.FavoriteDelete)
					authGroup.Middleware(LogActionMiddleware("parking_slot_create"))
					authGroup.Bind(parkingSlotCtrl.ParkingSlotAdd)
					authGroup.Bind(parkingSlotCtrl.ParkingSlotList)
					authGroup.Bind(parkingSlotCtrl.ParkingSlotGet)
					authGroup.Middleware(LogActionMiddleware("parking_slot_update"))
					authGroup.Bind(parkingSlotCtrl.ParkingSlotUpdate)
					authGroup.Middleware(LogActionMiddleware("parking_slot_delete"))
					authGroup.Bind(parkingSlotCtrl.ParkingSlotDelete)
					authGroup.Middleware(LogActionMiddleware("parking_order_create"))
					authGroup.Bind(parkingOrderCtrl.ParkingOrderAdd)
					authGroup.Bind(parkingOrderCtrl.ParkingOrderList)
					authGroup.Bind(parkingOrderCtrl.ParkingOrderGet)
					authGroup.Middleware(LogActionMiddleware("parking_order_update"))
					authGroup.Bind(parkingOrderCtrl.ParkingOrderUpdate)
					authGroup.Middleware(LogActionMiddleware("parking_order_cancel"))
					authGroup.Bind(parkingOrderCtrl.ParkingOrderCancel)
					authGroup.Middleware(LogActionMiddleware("parking_order_delete"))
					authGroup.Bind(parkingOrderCtrl.ParkingOrderDelete)
					authGroup.Middleware(LogActionMiddleware("parking_order_payment"))
					authGroup.Bind(parkingOrderCtrl.ParkingOrderPayment)
				})

				group.Group("/", func(userGroup *ghttp.RouterGroup) {
					userGroup.Middleware(middleware.UserOrAdmin)
					userGroup.Bind(userCtrl.UserById)
				})

				group.Group("/admin", func(adminGroup *ghttp.RouterGroup) {
					adminGroup.Middleware(middleware.AdminOnly)
					adminGroup.Middleware(LogActionMiddleware("admin_users_list"))
					adminGroup.Bind(userCtrl.UserList)
					adminGroup.Middleware(LogActionMiddleware("admin_user_delete"))
					adminGroup.Bind(userCtrl.UserDelete)
					adminGroup.Middleware(LogActionMiddleware("admin_user_role_update"))
					adminGroup.Bind(userCtrl.UserUpdateRole)
					adminGroup.Middleware(LogActionMiddleware("admin_user_wallet_update"))
					adminGroup.Bind(userCtrl.UserUpdateWalletBalance)
				})
			})

			s.Run()
			return nil
		},
	}
)

func CORS(r *ghttp.Request) {
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-CSRF-Token,Refresh-Token")
	r.Response.Header().Set("Access-Control-Expose-Headers", "New-Access-Token,New-Refresh-Token,X-CSRF-Token")
	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(200)
		r.Exit()
		return
	}
	r.Middleware.Next()
}

func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		g.Log().Error(r.Context(), "Request error occurred", "error", err, "path", r.URL.Path, "method", r.Method)
		code := consts.CodeInternalError
		if gerror.HasCode(err, consts.CodeInvalidInput) {
			code = consts.CodeInvalidInput
		}
		r.Response.WriteStatus(http.StatusBadRequest, g.Map{
			"code":    code.HttpStatus(),
			"message": err.Error(),
			"data":    nil,
		})
	}
}

func CSRFMiddleware(r *ghttp.Request) {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
		csrfToken := GenerateCSRFToken()
		r.Response.Header().Set("X-CSRF-Token", csrfToken)
		if r.URL.Path == "/backend/parkin/v1/refresh-token" {
			token := r.Header.Get(middleware.HeaderCSRFToken)
			g.Log().Info(r.Context(), "CSRF token check", "token", token)
			if token == "" || !verifyCSRFToken(token) {
				g.Log().Warning(r.Context(), "Invalid or missing CSRF token", "path", r.URL.Path)
				r.Response.WriteStatus(http.StatusForbidden, g.Map{
					"code":    consts.CodeInvalidInput.HttpStatus(),
					"message": "Invalid or missing CSRF token",
					"data":    nil,
				})
				r.Exit()
				return
			}
		}
	}
	r.Middleware.Next()
}

func LogActionMiddleware(action string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		userID := middleware.GetUserIDFromCtx(r.Context())
		g.Log().Async().Info(r.Context(), "Action performed", "action", action, "user_id", userID, "path", r.URL.Path, "method", r.Method)
		r.Middleware.Next()
	}
}

func GenerateCSRFToken() string {
	token := guid.S()
	csrfTokensLock.Lock()
	defer csrfTokensLock.Unlock()
	csrfTokens[token] = time.Now().Add(time.Hour)
	return token
}

func verifyCSRFToken(token string) bool {
	csrfTokensLock.RLock()
	defer csrfTokensLock.RUnlock()
	expireTime, exists := csrfTokens[token]
	if !exists || time.Now().After(expireTime) {
		return false
	}
	return true
}
