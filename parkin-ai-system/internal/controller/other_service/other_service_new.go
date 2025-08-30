package other_service

import (
	"parkin-ai-system/api/other_service"
)

type ControllerOtherService struct{}

func NewOtherService() other_service.IOtherService {
	return &ControllerOtherService{}
}
