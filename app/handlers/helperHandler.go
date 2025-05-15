package handlers

import (
	"camping/app/views/layouts"
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
	fmt.Println(isHTMX)
	if isHTMX {
		log.Println("HTMX request")
		return kit.Render(content)
	}
	log.Println("Not HTMX request")
	return kit.Render(layouts.App(content))
}
