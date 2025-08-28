package vehicles

import (
	"parkin-ai-system/internal/service"
)

type sVehicles struct{}

func init() {
	service.RegisterVehicles(&sVehicles{})
}
