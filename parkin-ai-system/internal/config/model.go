package config

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var cfg *Config

func InitConfig(ctx context.Context) {
	data, err := g.Cfg().Data(ctx)
	if err != nil {
		panic(err)
	}
	err = gconv.Scan(data, &cfg)
	if err != nil {
		panic(err)
	}
}
func GetConfig() Config {
	if cfg == nil {
		panic("config not initialized")
	}
	return *cfg
}
