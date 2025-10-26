package app

import (
	"encoding/json"
	"explorer/app/handlers"
	"explorer/app/views/errors"
	"explorer/plugins/auth"
	"explorer/plugins/booking"
	buses "explorer/plugins/busesConfig"
	"explorer/plugins/campsite"
	"explorer/plugins/paymentservices"
	"explorer/plugins/services"
	"explorer/plugins/status"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

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
		//app.Get("/about", kit.Handler(handlers.HandleLandingAbout))
		//app.Get("/help", kit.Handler(handlers.HandleHelp))
		app.Get("/photo+view", kit.Handler(handlers.HandlePhotoView))
		app.Get("/AreaAttraction", kit.Handler(handlers.HandleCampSites))
		app.Get("/book-new/{campID}", kit.Handler(handlers.HandleBookNew))
		app.Get("/campsites/description/{campID}", kit.Handler(handlers.CampDescription))
	})

	// Authenticated routes
	//
	// Routes that "must" have an authenticated user or else they
	// will be redirected to the configured redirectURL, set in the²
	// AuthenticationConfig.
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, true)) // strict set to true

		app.Get("/admin/campsites/new", kit.Handler(campsite.HandleCampsiteNewForm))
		app.Get("/admin/campsites/edit/{ID}", kit.Handler(campsite.HandleCampsiteEditForm))
		app.Get("/admin/buses", kit.Handler(buses.HandleModal))
		app.Get("/admin/services", kit.Handler(services.HandleServices))
		app.Get("/user/bookings/status", kit.Handler(status.BookingHandler))
		app.Get("/admin/booking/list", kit.Handler(handlers.HandelBooklist))
		app.Get("/admin/bookings/{id}/edit", kit.Handler(handlers.EditBooking))
		app.Get("/admin/bookings/{id}/showDetail", kit.Handler(handlers.BookingShowDetail))
		app.Get("/admin/bookings/{user_id}/new", kit.Handler(handlers.BookingAdmin))
		app.Get("/admin/bookings/search", kit.Handler(handlers.BookingSearch))
		app.Get("/admin/bookings/print", kit.Handler(handlers.PrintBookings))
		app.Get("/admin/carousel", kit.Handler(handlers.Carousel))

		app.Post("/admin/buses/create", kit.Handler(buses.HandleCreate))
		app.Post("/admin/campsites/create", kit.Handler(campsite.HandleCampsiteCreate))
		app.Post("/admin/campsites/edit/{ID}", kit.Handler(campsite.HandleCampsiteUpdate))
		app.Post("/admin/services/create", kit.Handler(services.HandleServiceCreate))
		app.Post("/user/bookings", kit.Handler(booking.HandelCreateBooking))
		app.Post("/admin/bookings/{userID}/create", kit.Handler(handlers.AdminBookingAdd))
		app.Post("/admin/bookings/{Bookid}/edit", kit.Handler(handlers.EditPostBooking))
		app.Post("/admin/carousel/create", kit.Handler(handlers.CarouselImageCreate))
		app.Post("/user/konnect/payment/Init", kit.Handler(func(kit *kit.Kit) error {
			service := &paymentservices.KonnectService{
				APIKey:     os.Getenv("KONNECT_API_KEY"),
				BaseURL:    os.Getenv("KONNECT_API_BASE_URL"),
				WebhookURL: os.Getenv("WEBHOOK_URL"),
				SuccessURL: os.Getenv("SUCCESS_URL"),
				FailURL:    os.Getenv("FAIL_URL"),
				WalletID:   os.Getenv("RECEIVER_WALLET_ID"),
			}
			return handlers.KonnectInitPayment(kit, service)
		}))
		app.Post("/payments/webhook", kit.Handler(func(kit *kit.Kit) error {
			body, err := io.ReadAll(kit.Request.Body)
			if err != nil {
				return err
			}

			log.Println("📩 Konnect WEBHOOK received:", string(body))

			// just echo back JSON so you see it in browser / logs
			kit.Response.Header().Set("Content-Type", "application/json")
			kit.Response.WriteHeader(http.StatusOK)
			kit.Response.Write(body)
			return nil
		}))

		app.Get("/payments/success", kit.Handler(func(kit *kit.Kit) error {
			q := kit.Request.URL.Query()
			log.Println("✅ Konnect SUCCESS redirect:", q)

			// send back JSON with all query params
			kit.Response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(kit.Response).Encode(q)
			return nil
		}))

		app.Get("/payments/fail", kit.Handler(func(kit *kit.Kit) error {
			q := kit.Request.URL.Query()
			log.Println("❌ Konnect FAIL redirect:", q)

			kit.Response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(kit.Response).Encode(q)
			return nil
		}))

		// Deletion routes

		app.Post("/admin/campsites/delete/{ID}", kit.Handler(campsite.HandleCampsiteDelete))
		app.Post("/admin/buses/{id}/delete", kit.Handler(buses.HandleDelete))
		app.Post("/admin/services/{id}/delete", kit.Handler(services.HandleServiceDelete))

		app.Delete("/admin/bookings/list/{bookID}", kit.Handler(handlers.HandelDeleteBookingList))
		app.Delete("/admin/carousel/{id}/delete", kit.Handler(handlers.CaroucelImageDelete))
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
