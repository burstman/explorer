package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"

	"github.com/anthdm/superkit/kit"
)

// HandleLandingIndex renders the landing page index view using the default layout
func HandleLandingIndex(kit *kit.Kit) error {
	session := kit.GetSession("user-session")
	successFlashes := session.Flashes("success")
	failFlashes := session.Flashes("fail")
	session.Save(kit.Request, kit.Response)

	successMessages := make([]string, len(successFlashes))
	for i, flash := range successFlashes {
		successMessages[i] = flash.(string)
	}

	failMessages := make([]string, len(failFlashes))
	for i, flash := range failFlashes {
		failMessages[i] = flash.(string)
	}

	return RenderWithLayout(kit, landing.Index(successMessages, failMessages))
}

func HandleLandingAbout(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.About())

}

func HandleHelp(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.Help())
}

func HandlePhotoView(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.PhotoView())
}

// HandleCampSites retrieves camp sites and renders the camp sites page with user role, site data, and optional flash messages
func HandleCampSites(kit *kit.Kit) error {
	var data []types.CampSite
	db.Get().Order("title asc").Find(&data)

	user, ok := kit.Auth().(types.AuthUser)
	role := "user"
	if ok {
		role = user.GetRole()
	}

	session := kit.GetSession("user-session")
	successFlashes := session.Flashes("success")
	failFlashes := session.Flashes("fail")
	session.Save(kit.Request, kit.Response)

	var flashType, flashMsg string
	if len(successFlashes) > 0 {
		flashType = "success"
		flashMsg = successFlashes[0].(string)
	} else if len(failFlashes) > 0 {
		flashType = "fail"
		flashMsg = failFlashes[0].(string)
	}

	return RenderWithLayout(kit, landing.CampSites(role, data, flashType, flashMsg))
}

func HandleBookNew(kit *kit.Kit) error {
	return RenderWithLayout(kit, landing.NewBooking())
}
