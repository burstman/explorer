package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"fmt"
	"strconv"

	"log"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
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

type SeatInfo struct {
	CampsiteID int
	TotalSeats int
}

type BookingInfo struct {
	CampID      int
	BookedSeats int
}

// HandleCampSites retrieves camp sites and renders the camp sites page with user role, site data, and optional flash messages
func HandleCampSites(kit *kit.Kit) error {
	var camps []types.CampSite
	db.Get().Order("title asc").Find(&camps)

	user, _ := kit.Auth().(types.AuthUser)

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

	var buses []types.BuseType
	db.Get().Find(&buses)

	var seatData []SeatInfo

	err := db.Get().Table("campsite_buses").
		Select("campsite_buses.campsite_id, SUM(campsite_buses.quantity * bus_types.capacity) AS total_seats").
		Joins("JOIN bus_types ON bus_types.id = campsite_buses.bus_type_id").
		Group("campsite_buses.campsite_id").
		Scan(&seatData).Error
	if err != nil {
		log.Println("Failed to fetch camp site data:", err)
		return err
	}

	var bookings []BookingInfo

	err = db.Get().Raw(`
	SELECT 
		b.camp_id,
		COUNT(DISTINCT b.id) + COUNT(g.id) AS booked_seats
	FROM bookings b
	LEFT JOIN guests g ON g.booking_id = b.id AND g.deleted_at IS NULL
	WHERE b.deleted_at IS NULL
	GROUP BY b.camp_id
`).Scan(&bookings).Error

	if err != nil {
		log.Println("Failed to fetch booking data:", err)
		return err
	}
	//map campsite id to total seat
	totalSeatsMap := make(map[int]int)
	bookedSeatsMap := make(map[int]int)

	for _, row := range seatData {
		totalSeatsMap[row.CampsiteID] = row.TotalSeats
	}

	for _, s := range seatData {
		log.Printf("Seats for Camp ID %d: %d", s.CampsiteID, s.TotalSeats)
	}

	for _, b := range bookings {
		bookedSeatsMap[b.CampID] = b.BookedSeats
	}

	remainingSeatsMap := make(map[int]int)
	for campID, total := range totalSeatsMap {
		booked := bookedSeatsMap[campID]
		remaining := total - booked
		if remaining < 0 {
			remaining = 0 // safety
		}
		remainingSeatsMap[campID] = remaining
	}
	log.Printf("Booking data with guests: %+v", bookings)
	log.Printf("Remaining seats per camp: %+v", remainingSeatsMap)

	return RenderWithLayout(kit, landing.CampSites(user, camps, buses, remainingSeatsMap, flashType, flashMsg))
}

func HandleBookNew(kit *kit.Kit) error {
	auth := kit.Auth().(types.AuthUser)
	userID := auth.UserID

	campID := chi.URLParam(kit.Request, "campID")

	var camp types.CampSite
	var user types.User
	err := db.Get().Where("id = ?", campID).First(&camp).Error
	if err != nil {
		return err
	}

	err = db.Get().Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}
	var services []types.Service

	err = db.Get().Find(&services).Error
	if err != nil {
		return err
	}

	return RenderWithLayout(kit, landing.NewBooking(camp, user, services))
}

func CampDescription(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "campID")

	id, err := strconv.Atoi(idStr)
	//log.Println("camp id description:", id)
	if err != nil {
		return fmt.Errorf("error parsing camp id in CampDescriptionHandler: %v", err)
	}

	camp, err := GetCampByID(id)
	//log.Println("camp:", camp.Description)

	if err != nil {
		return fmt.Errorf("err retrieving data in CampDescriptionHandler %v", err)
	}
	return RenderWithLayout(kit, landing.Description(camp))

}

func GetCampByID(id int) (types.CampSite, error) {
	var camp types.CampSite
	err := db.Get().Where("id = ?", id).Find(&camp).Error
	if err != nil {
		return types.CampSite{}, err
	}
	return camp, nil
}
