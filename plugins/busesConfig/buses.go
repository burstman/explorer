package buses

import (
	"explorer/app/db"
	"explorer/app/types"
	"log"
	"net/http"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// HandleModal renders the bus configuration modal, fetching all buses and displaying any success or failure messages from the session
func HandleModal(kit *kit.Kit) error {
	session := kit.GetSession("user-session")

	successFlashes := session.Flashes("success")
	failFlashes := session.Flashes("fail")
	session.Save(kit.Request, kit.Response)

	var buses []types.BuseType
	if err := db.Get().Find(&buses).Error; err != nil {
		failFlashes = append(failFlashes, "Failed to fetch buses")
	}

	return kit.Render(BusConfigModal(buses, flashStrings(successFlashes), flashStrings(failFlashes)))
}

func flashStrings(flashes []any) []string {
	msgs := make([]string, len(flashes))
	for i, f := range flashes {
		msgs[i] = f.(string)
	}
	return msgs
}

// HandleCreate processes a form submission to create a new bus configuration.
// It validates the form data, converts the capacity to an integer, creates a new bus record,
// and manages session flash messages to provide user feedback about the creation process.
func HandleCreate(kit *kit.Kit) error {
	session := kit.GetSession("user-session")

	// Parse form
	if err := kit.Request.ParseForm(); err != nil {
		session.AddFlash("Invalid form data", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	name := kit.Request.FormValue("name")
	if len(name) == 0 {
		session.AddFlash("Bus name is required", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	capacityStr := kit.Request.FormValue("capacity")
	if len(capacityStr) == 0 {
		session.AddFlash("Bus capacity is required", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		session.AddFlash("Capacity must be a number ", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	bus := types.BuseType{Name: name, Capacity: capacity}
	if err := db.Get().Create(&bus).Error; err != nil {
		session.AddFlash("Failed to create bus", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	session.AddFlash("Bus created successfully", "success")
	session.Save(kit.Request, kit.Response)
	return renderModalWithFlashes(kit)
}

// renderModalWithFlashes retrieves session flash messages for success and failure,
// saves the session, fetches all bus configurations from the database,
// and renders the bus configuration modal with the buses and flash messages.
func renderModalWithFlashes(kit *kit.Kit) error {
	session := kit.GetSession("user-session")

	successFlashes := session.Flashes("success")
	failFlashes := session.Flashes("fail")
	session.Save(kit.Request, kit.Response)

	var buses []types.BuseType
	db.Get().Find(&buses)

	return kit.Render(BusConfigModal(buses, toStringSlice(successFlashes), toStringSlice(failFlashes)))
}

func toStringSlice(flashes []interface{}) []string {
	out := make([]string, 0, len(flashes))
	for _, f := range flashes {
		if s, ok := f.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func HandleDelete(kit *kit.Kit) error {
	log.Println("Request path:", kit.Request.URL.Path)
	session := kit.GetSession("user-session")
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.Text(http.StatusBadRequest, "Invalid bus ID")
	}

	if err := db.Get().Delete(&types.BuseType{}, id).Error; err != nil {
		log.Println("Failed to delete bus:", err)

		session.AddFlash("Failed to delete bus", "fail")
		session.Save(kit.Request, kit.Response)
		return renderModalWithFlashes(kit)
	}

	// Success flash
	session.AddFlash("Bus deleted successfully", "success")
	session.Save(kit.Request, kit.Response)

	// Re-render modal with flash and updated bus list
	return renderModalWithFlashes(kit)
}
