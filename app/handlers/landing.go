package handlers

import (
	"camping/app/views/landing"

	"github.com/anthdm/superkit/kit"
)

// HandleLandingIndex renders the landing page index view using the default layout
func HandleLandingIndex(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.Index())
}

func HandleLandingAbout(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.About())

}
