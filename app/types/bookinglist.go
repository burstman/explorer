package types

type BookingDetails struct {
	BookingID int
	CampID    int
	CampName  string
	User      User
	Guests    []Guest
}
