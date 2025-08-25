package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func Register(s *ghttp.Server) {
	s.Use(CORS)

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(Auth)
	})
}

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
