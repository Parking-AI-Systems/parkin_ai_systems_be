package favourite

import (
	"parkin-ai-system/api/favourite"
)

type ControllerFavourite struct{}

func NewFavourite() favourite.IFavourite {
	return &ControllerFavourite{}
}
