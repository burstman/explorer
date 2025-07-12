package handlers

import (
	"explorer/app/types"
	"explorer/app/views/layouts"
	"fmt"

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
	
	if isLoggedIn {
		// Get the authenticated user and extract the role
		if user,ok := kit.Auth().(types.AuthUser); ok {
			role = user.GetRole()
		}
	}
	fmt.Println("role in render with layout", role)
	fmt.Println("is loggedin", isLoggedIn)
	

	return kit.Render(layouts.App(content, role, isLoggedIn))
}
