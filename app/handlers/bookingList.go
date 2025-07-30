package handlers

import (
	"explorer/app/views/landing"

	"github.com/anthdm/superkit/kit"
)

func HandelBooklist(kit *kit.Kit) error {

	return RenderWithLayout(kit, landing.PhotoView())
}
