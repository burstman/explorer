package handlers

import (
	"explorer/app/types"
	"explorer/app/views/layouts"
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/anthdm/superkit/kit"
)

// RenderWithLayout renders the given content, either directly or wrapped in the application layout.
// If the request is an HTMX request, it renders the content directly. Otherwise, it renders
// the content within the standard application layout.
func RenderWithLayout(kit *kit.Kit, content templ.Component) error {
	isHTMX := kit.Request.Header.Get("HX-Request") == "true"

	if isHTMX {
		return kit.Render(content)
	}
	isLoggedIn := kit.Auth().Check()
	var role string
	var userId uint
	if isLoggedIn {
		// Get the authenticated user and extract the role
		if user, ok := kit.Auth().(types.AuthUser); ok {
			role = user.GetRole()
			userId = user.GetUserID()

		}
	}
	log.Println("role in render with layout", role)
	log.Println("userid in render with layout", userId)

	fmt.Println("is loggedin", isLoggedIn)

	return kit.Render(layouts.App(content, role, isLoggedIn, userId))
}
