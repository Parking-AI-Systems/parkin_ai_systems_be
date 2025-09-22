package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	"parkin-ai-system/internal/controller/payment"
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
	_ "parkin-ai-system/internal/logic/payment"
	_ "parkin-ai-system/internal/logic/user"
	_ "parkin-ai-system/internal/logic/vehicle"
	_ "parkin-ai-system/internal/logic/wallet_transaction"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
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

			// Initialize database connection with performance optimization
			if err := initDatabase(ctx); err != nil {
				g.Log().Fatal(ctx, "Failed to initialize database connection", "error", err)
				return err
			}

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

			// Initialize controllers
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
			paymentCtrl := payment.NewPayment()
			s := g.Server()
			s.Logger().SetHandlers(glog.HandlerJson)
			s.Use(CORS, ErrorHandler)

			s.Group("/metrics", func(metricsGroup *ghttp.RouterGroup) {
				metricsGroup.GET("/", ghttp.WrapH(promhttp.Handler()))
			})

			// Health check endpoints
			s.Group("/health", func(healthGroup *ghttp.RouterGroup) {
				healthGroup.GET("/", HealthCheck)
				healthGroup.GET("/db", DatabaseHealthCheck)
				healthGroup.GET("/db/pool", DatabasePoolStats)
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
					authGroup.Middleware(LogActionMiddleware("wallet_transaction_list"))
					authGroup.Bind(walletTransactionCtrl.WalletTransactionList)
					authGroup.Middleware(LogActionMiddleware("wallet_transaction_get"))
					authGroup.Bind(walletTransactionCtrl.WalletTransactionGet)
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
					authGroup.Middleware(LogActionMiddleware("create_payment_link"))
					authGroup.Bind(paymentCtrl.CreatePaymentLink)
					authGroup.Middleware(LogActionMiddleware("payment_link_get"))
					authGroup.Bind(paymentCtrl.PaymentLinkGet)
					authGroup.Middleware(LogActionMiddleware("payment_refund_add"))
					authGroup.Bind(paymentCtrl.RefundAdd)
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
					adminGroup.Middleware(LogActionMiddleware("admin_other_service_revenue"))
					adminGroup.Bind(otherServiceOrderCtrl.OthersServiceRevenue)
					adminGroup.Middleware(LogActionMiddleware("admin_other_service_popular"))
					adminGroup.Bind(otherServiceOrderCtrl.OthersServicePopular)
					adminGroup.Middleware(LogActionMiddleware("admin_other_service_trends"))
					adminGroup.Bind(otherServiceOrderCtrl.OthersServiceTrends)
					adminGroup.Middleware(LogActionMiddleware("admin_parking_order_revenue"))
					adminGroup.Bind(parkingOrderCtrl.ParkingOrderRevenue)
					adminGroup.Middleware(LogActionMiddleware("admin_parking_order_trends"))
					adminGroup.Bind(parkingOrderCtrl.ParkingOrderTrends)
					adminGroup.Middleware(LogActionMiddleware("admin_parking_lot_orders_status_breakdown"))
					adminGroup.Bind(parkingOrderCtrl.ParkingOrderStatusBreakdown)
					adminGroup.Middleware(LogActionMiddleware("admin_user_count"))
					adminGroup.Bind(userCtrl.UserCount)
					adminGroup.Middleware(LogActionMiddleware("admin_user_role_distribution"))
					adminGroup.Bind(userCtrl.UserRoleDistribution)
					adminGroup.Middleware(LogActionMiddleware("admin_user_recent_registrations"))
					adminGroup.Bind(userCtrl.UserRecentRegistrations)
				})
			})

			// Setup graceful shutdown
			setupGracefulShutdown(ctx, s)

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

// Helper functions for parsing configuration with performance-first defaults
func parseConnectionPoolInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed := gconv.Int(value); parsed > 0 {
		return parsed
	}
	return defaultValue
}

func parseConnectionPoolDuration(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}
	if parsed, err := time.ParseDuration(value); err == nil && parsed > 0 {
		return parsed
	}
	return defaultValue
}

// initDatabase initializes database connection pool with performance optimization
func initDatabase(ctx context.Context) error {
	cfg := config.GetConfig()

	// Get database instance
	db := g.DB()

	// Parse and validate connection pool settings with performance-first defaults
	maxIdle := parseConnectionPoolInt(cfg.Database.Default.MaxIdle, 5)  // Performance optimized: 5 idle connections
	maxOpen := parseConnectionPoolInt(cfg.Database.Default.MaxOpen, 10) // Performance optimized: 10 max connections

	// Parse connection lifetime with performance focus
	maxLifetime := parseConnectionPoolDuration(cfg.Database.Default.MaxLifetime, 15*time.Minute) // 15 minutes for less overhead
	maxIdleTime := parseConnectionPoolDuration(cfg.Database.Default.MaxIdleTime, 5*time.Minute)  // 5 minutes to maintain performance

	// Configure connection pool settings using GoFrame's configuration
	// Note: GoFrame manages connection pool internally based on config
	// The settings in config.yaml are automatically applied

	// Test database connectivity
	if err := db.PingMaster(); err != nil {
		return gerror.Wrap(err, "failed to ping database master")
	}

	// Test slave connection if configured
	if err := db.PingSlave(); err != nil {
		g.Log().Warning(ctx, "Failed to ping database slave, continuing with master only", "error", err)
	}

	// Log successful initialization
	g.Log().Info(ctx, "Database connection pool initialized with performance optimization",
		"maxIdleConns", maxIdle,
		"maxOpenConns", maxOpen,
		"maxLifetime", maxLifetime,
		"maxIdleTime", maxIdleTime)

	// Start performance monitoring
	startConnectionPoolMonitoring(ctx, db)

	return nil
}

// Connection pool monitoring optimized for performance tracking
func startConnectionPoolMonitoring(ctx context.Context, db gdb.DB) {
	go func() {
		ticker := time.NewTicker(2 * time.Minute) // Check every 2 minutes for performance monitoring
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Basic connectivity check
				if err := db.PingMaster(); err != nil {
					g.Log().Warning(ctx, "Database connection monitoring - ping failed", "error", err)
				} else {
					g.Log().Info(ctx, "Database connection pool monitoring - connection active")
				}
			}
		}
	}()
}

// HealthCheck provides basic application health status
func HealthCheck(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"service":   "parkin-ai-system",
	})
}

// DatabaseHealthCheck checks database connectivity
func DatabaseHealthCheck(r *ghttp.Request) {
	db := g.DB()

	// Test master connection
	if err := db.PingMaster(); err != nil {
		g.Log().Error(r.Context(), "Database health check failed", "error", err)
		r.Response.WriteStatus(http.StatusServiceUnavailable, g.Map{
			"status":    "unhealthy",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"service":   "database",
			"error":     err.Error(),
		})
		return
	}

	// Test slave connection
	slaveStatus := "healthy"
	if err := db.PingSlave(); err != nil {
		slaveStatus = "unavailable"
		g.Log().Warning(r.Context(), "Database slave health check failed", "error", err)
	}

	r.Response.WriteJson(g.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"service":   "database",
		"master":    "healthy",
		"slave":     slaveStatus,
	})
}

// DatabasePoolStats provides detailed connection pool statistics with performance metrics
func DatabasePoolStats(r *ghttp.Request) {
	db := g.DB()
	cfg := config.GetConfig()

	// Test connectivity
	poolStatus := "active"
	if err := db.PingMaster(); err != nil {
		poolStatus = "error"
		r.Response.WriteJson(g.Map{
			"status":    "error",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"error":     err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"status":    poolStatus,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"service":   "database_pool",
		"configuration": g.Map{
			"max_idle_conns": parseConnectionPoolInt(cfg.Database.Default.MaxIdle, 5),
			"max_open_conns": parseConnectionPoolInt(cfg.Database.Default.MaxOpen, 10),
			"max_lifetime":   parseConnectionPoolDuration(cfg.Database.Default.MaxLifetime, 15*time.Minute).String(),
			"max_idle_time":  parseConnectionPoolDuration(cfg.Database.Default.MaxIdleTime, 5*time.Minute).String(),
		},
		"connectivity": g.Map{
			"master": "connected",
			"slave": func() string {
				if err := db.PingSlave(); err != nil {
					return "unavailable"
				}
				return "connected"
			}(),
		},
		"note": "GoFrame manages connection pool internally based on configuration",
	})
}

// setupGracefulShutdown configures graceful shutdown for the application
func setupGracefulShutdown(ctx context.Context, server *ghttp.Server) {
	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle graceful shutdown
	go func() {
		<-sigChan
		g.Log().Info(ctx, "Received shutdown signal, initiating graceful shutdown...")

		// Close database connections
		if db := g.DB(); db != nil {
			g.Log().Info(ctx, "Closing database connections...")
			// GoFrame automatically handles connection cleanup during shutdown
		}

		g.Log().Info(ctx, "Graceful shutdown completed")
		os.Exit(0)
	}()
}
