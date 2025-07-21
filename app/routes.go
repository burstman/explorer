package app

import (
	"explorer/app/handlers"
	"explorer/app/views/errors"
	"explorer/plugins/auth"
	buses "explorer/plugins/busesConfig"
	"explorer/plugins/campsite"
	"explorer/plugins/services"
	"log/slog"

	"github.com/anthdm/superkit/kit"
	"github.com/anthdm/superkit/kit/middleware"
	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Define your global middleware
func InitializeMiddleware(router *chi.Mux) {
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.WithRequest)
}

// Define your routes in here
func InitializeRoutes(router *chi.Mux) {
	// Authentication plugin
	//
	// By default the auth plugin is active, to disable the auth plugin
	// you will need to pass your own handler in the `AuthFunc`` field
	// of the `kit.AuthenticationConfig`.
	//  authConfig := kit.AuthenticationConfig{
	//      AuthFunc: YourAuthHandler,
	//      RedirectURL: "/login",
	//  }
	auth.InitializeRoutes(router)
	// authConfig := kit.AuthenticationConfig{
	// 	AuthFunc:    auth.AuthenticateUser,
	// 	RedirectURL: "/login",
	// }
	authConfig := kit.AuthenticationConfig{
		AuthFunc:    auth.AuthenticateUser,
		RedirectURL: "/login",
	}

	// Routes that "might" have an authenticated user
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, false)) // strict set to false

		// Routes
		app.Get("/", kit.Handler(handlers.HandleLandingIndex))
		app.Get("/about", kit.Handler(handlers.HandleLandingAbout))
		app.Get("/help", kit.Handler(handlers.HandleHelp))
		app.Get("/photo+view", kit.Handler(handlers.HandlePhotoView))
		app.Get("/AreaAttraction", kit.Handler(handlers.HandleCampSites))
		app.Get("/book-new/{campID}", kit.Handler(handlers.HandleBookNew))
	})

	// Authenticated routes
	//
	// Routes that "must" have an authenticated user or else they
	// will be redirected to the configured redirectURL, set in the²
	// AuthenticationConfig.
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, true)) // strict set to true
		app.Post("/admin/services/delete/{id}", kit.Handler(services.HandleDeleteService))
		app.Get("/admin/campsites/new", kit.Handler(campsite.HandleCampsiteNewForm))
		app.Get("/admin/campsites/edit/{ID}", kit.Handler(campsite.HandleCampsiteEditForm))
		app.Get("/admin/buses", kit.Handler(buses.HandleModal))
		app.Get("/admin/services", kit.Handler(services.HandleServices))

		app.Post("/admin/buses/create", kit.Handler(buses.HandleCreate))
		app.Post("/admin/campsites/create", kit.Handler(campsite.HandleCampsiteCreate))
		app.Post("/admin/campsites/edit/{ID}", kit.Handler(campsite.HandleCampsiteUpdate))
		app.Post("/admin/services/add", kit.Handler(services.HandleAddService))

		app.Post("/admin/campsites/delete/{ID}", kit.Handler(campsite.HandleCampsiteDelete))
		app.Post("/admin/buses/{id}/delete", kit.Handler(buses.HandleDelete))

		// app.Post("/admin/services/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		// 	id := chi.URLParam(r, "id")
		// 	fmt.Println("Captured ID is:", id)
		// 	w.Write([]byte("Captured ID = " + id))
		// })
		// Routes
		// app.Get("/path", kit.Handler(myHandler.HandleIndex))
	})
}

// NotFoundHandler that will be called when the requested path could
// not be found.
func NotFoundHandler(kit *kit.Kit) error {
	return kit.Render(errors.Error404())
}

// ErrorHandler that will be called on errors return from application handlers.
func ErrorHandler(kit *kit.Kit, err error) {
	slog.Error("internal server error", "err", err.Error(), "path", kit.Request.URL.Path)
	kit.Render(errors.Error500())
}
