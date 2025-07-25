package buses

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"

	"log"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// HandleModal renders the bus configuration modal, fetching all buses and displaying any success or failure messages from the session
func HandleModal(kit *kit.Kit) error {

	var buses []types.BuseType
	if err := db.Get().Find(&buses).Error; err != nil {
		return err
	}

	return kit.Render(BusConfigModal(buses))
}

// HandleCreate processes a form submission to create a new bus configuration.
// It validates the form data, converts the capacity to an integer, creates a new bus record,
// and manages session flash messages to provide user feedback about the creation process.
func HandleCreate(kit *kit.Kit) error {

	// Parse form
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	name := kit.Request.FormValue("name")

	capacityStr := kit.Request.FormValue("capacity")

	if name == "" {
		return fmt.Errorf("bus name is required")
	}

	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		return err
	}

	if capacity <= 0 {
		return fmt.Errorf("capacity must be greater than zero")
	}

	bus := types.BuseType{Name: name, Capacity: capacity}
	if err := db.Get().Create(&bus).Error; err != nil {
		return err
	}

	var buses []types.BuseType
	if err := db.Get().Find(&buses).Error; err != nil {
		return err
	}

	return kit.Render(BusConfigModal(buses))
}

func HandleDelete(kit *kit.Kit) error {
	log.Println("Request path:", kit.Request.URL.Path)

	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	var bus types.BuseType

	if err := db.Get().First(&bus, id).Error; err != nil {
		log.Println("bus does not exist in database", err)
		return err
	}

	if err := db.Get().Delete(&bus).Error; err != nil {
		log.Println("Failed to delete bus:", err)

		return err
	}

	var buses []types.BuseType
	if err := db.Get().Find(&buses).Error; err != nil {
		return err
	}

	// Re-render modal with flash and updated bus list
	return kit.Render(BusConfigModal(buses))
}
