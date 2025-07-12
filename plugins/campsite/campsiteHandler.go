package campsite

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// handlers/campsites.go
func HandleCampsiteCreate(kit *kit.Kit) error {
	session := kit.GetSession("user-session")
	busIDs := kit.Request.Form["bus_ids"]
	// Parse the form
	if err := kit.Request.ParseForm(); err != nil {
		session.AddFlash("Invalid form data", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Get form values
	camp := types.CampSite{
		Name:        kit.Request.FormValue("name"),
		Description: kit.Request.FormValue("description"),
		ImageURL:    kit.Request.FormValue("image_url"),
		Location:    kit.Request.FormValue("location"),
	}

	var err error
	camp.AvailableFrom, err = parseFormDate(kit.Request.FormValue("available_from"))
	if err != nil {
		session.AddFlash("Invalid From Date", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	camp.AvailableTo, err = parseFormDate(kit.Request.FormValue("available_to"))
	if err != nil {
		session.AddFlash("Invalid To Date", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Save to DB
	if err := db.Get().Create(&camp).Error; err != nil {
		log.Println("Failed to create campsite:", err)
		session.AddFlash("Failed to create campsite", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	session.AddFlash("Campsite created successfully", "success")
	session.Save(kit.Request, kit.Response)
	return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
}

type CampsiteFormValues struct {
	ID          uint   `form:"id"`
	Name        string `form:"name"`
	Description string `form:"description"`
	ImageURL    string `form:"image_url"`
}

func HandleCampsiteNewForm(kit *kit.Kit) error {
	var bus = []types.BuseType{}
	return kit.Render(NewCampsiteForm(bus))
}

// HandleCampsiteEditForm handles the HTTP request to display the edit form for a specific campsite.
// It retrieves the campsite by ID from the URL parameter, finds the campsite in the database,
// and renders the edit form with the campsite's current details. If the campsite ID is invalid
// or the campsite is not found, it adds a flash message and redirects to the area attraction page.
func HandleCampsiteEditForm(kit *kit.Kit) error {
	session := kit.GetSession("user-session")
	idParam := chi.URLParam(kit.Request, "ID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		session.AddFlash("Invalid campsite ID", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	var camp types.CampSite
	camp, err = FindCampByID(id)
	if err != nil {
		session.AddFlash("campsite not found", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	return kit.Render(EditCampsiteForm(camp))
}

func FindCampByID(id int) (types.CampSite, error) {
	var camp types.CampSite
	err := db.Get().First(&camp, id).Error
	return camp, err

}

// HandleCampsiteUpdate handles the HTTP request to update an existing campsite.
// It retrieves the campsite by ID, updates its fields with form values, and saves the changes.
// On success, it adds a success flash message and redirects to the area attraction page.
// On failure, it adds a failure flash message and returns an error response.
func HandleCampsiteUpdate(kit *kit.Kit) error {
	session := kit.GetSession("user-session")
	idParam := chi.URLParam(kit.Request, "ID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		session.AddFlash("Invalid campsite ID", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Parse form values
	if err := kit.Request.ParseForm(); err != nil {
		session.AddFlash("Failed to parse form", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")

	}

	// Get form values
	name := kit.Request.FormValue("name")
	description := kit.Request.FormValue("description")
	imageURL := kit.Request.FormValue("image_url")
	location := kit.Request.FormValue("location")
	fromDate := kit.Request.FormValue("available_from")
	toDate := kit.Request.FormValue("available_to")

	// Find the existing campsite
	var camp types.CampSite
	if err := db.Get().First(&camp, id).Error; err != nil {
		session.AddFlash("Campsite not found", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Update fields
	camp.Name = name
	camp.Description = description
	camp.ImageURL = imageURL
	camp.Location = location
	camp.AvailableFrom, err = parseFormDate(fromDate)
	if err != nil {
		session.AddFlash("Invalid From Date ", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Text(http.StatusBadRequest, "Invalid from date")
	}
	camp.AvailableTo, err = parseFormDate(toDate)
	if err != nil {
		session.AddFlash("Invalid To Date ", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Save the updated campsite
	if err := db.Get().Save(&camp).Error; err != nil {
		session.AddFlash("Failed to updated campsite ", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	session.AddFlash("Campsite updated successfully", "success")
	session.Save(kit.Request, kit.Response)

	// Redirect after success (optional)
	return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
}

func parseFormDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return &time.Time{}, nil // Return zero time for empty strings
	}

	// Common date formats you might encounter from forms
	dateLayouts := []string{
		"2006-01-02",          // HTML date input format (YYYY-MM-DD)
		"2006-01-02T15:04",    // HTML datetime-local input format
		"2006-01-02T15:04:05", // ISO format
		"01/02/2006",          // US format (MM/DD/YYYY)
		"02/01/2006",          // European format (DD/MM/YYYY)
		"2006-01-02 15:04:05", // DateTime format
		"January 2, 2006",     // Long format
		"Jan 2, 2006",         // Short format
	}

	for _, layout := range dateLayouts {
		if parsedTime, err := time.Parse(layout, dateStr); err == nil {
			return &parsedTime, nil
		}
	}

	return &time.Time{}, fmt.Errorf("unable to parse date '%s' - supported formats: YYYY-MM-DD, MM/DD/YYYY, DD/MM/YYYY", dateStr)
}

// HandleCampsiteDelete handles the deletion of a campsite by its ID.
// It retrieves the campsite from the database, deletes it, and redirects to the AreaAttraction page.
// If the campsite is not found or cannot be deleted, it adds a flash message and redirects.
func HandleCampsiteDelete(kit *kit.Kit) error {
	session := kit.GetSession("user-session")

	idParam := chi.URLParam(kit.Request, "ID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("Invalid campsite ID:", err)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Optional: Check if the record exists first (helps give clearer errors)
	var camp types.CampSite
	if err := db.Get().First(&camp, id).Error; err != nil {
		log.Println("Campsite not found:", err)
		session.AddFlash("Campsite not found", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	// Delete the record
	if err := db.Get().Delete(&camp).Error; err != nil {
		log.Println("Failed to delete campsite:", err)
		session.AddFlash("Campsite not deleted", "fail")
		session.Save(kit.Request, kit.Response)
		return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
	}

	session.AddFlash("Campsite deleted successfully", "success")
	session.Save(kit.Request, kit.Response)

	return kit.Redirect(http.StatusSeeOther, "/AreaAttraction")
}
