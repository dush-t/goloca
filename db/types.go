package db

// Store represents any database connection in the scope of this API
// only. I'm doing it like this to make sure that if we wish to use
// some other storage solution, we can.
type Store interface {
	Connect()
	Config() interface{}
	Insert(interface{}) error
}

// RideDataPoint represents a single data point that the client
// will send to this service
type RideDataPoint struct {
	ID                string  `json:"id,omitempty"`
	UserID            string  `json:"user_id,omitempty"`
	VehicleModelID    string  `json:"vehicle_model_id,omitempty"`
	PackageID         int     `json:"package_id,omitempty"`
	TravelTypeID      int     `json:"travel_type_id,omitempty"`
	FromAreaID        string  `json:"from_area_id,omitempty"`
	ToAreaID          string  `json:"to_area_id,omitempty"`
	FromCityID        string  `json:"from_city_id,omitempty"`
	ToCityID          string  `json:"to_city_id,omitempty"`
	FromDate          int64   `json:"from_date,omitempty"`
	ToDate            int64   `json:"to_date,omitempty"`
	OnlineBooking     bool    `json:"online_booking,omitempty"`
	MobileSiteBooking bool    `json:"mobile_site_booking,omitempty"`
	BookingCreated    int64   `json:"booking_created,omitempty"`
	FromLat           float64 `json:"from_lat,omitempty"`
	FromLong          float64 `json:"from_long,omitempty"`
	ToLat             float64 `json:"to_lat,omitempty"`
	ToLong            float64 `json:"to_long,omitempty"`
	CarCancellation   bool    `json:"Car_Cancellation,omitempty"`
}
