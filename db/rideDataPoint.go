package db

import (
	"errors"
	"math"
)

// Validate verifies that the data recieved from the client
// is clean
func (r RideDataPoint) Validate() error {
	var err error
	err = nil

	if r.PackageID < 1 || r.PackageID > 7 {
		err = errors.New("package_id must be between 1 and 7")
		return err
	}

	if r.TravelTypeID < 1 || r.TravelTypeID > 3 {
		err = errors.New("travel_type_id must be between 1 and 3")
		return err
	}

	if math.Abs(r.FromLat) > 90 || math.Abs(r.ToLat) > 90 {
		err = errors.New("latitude must be between -90 and +90")
		return err
	}

	if r.ToLong < 0 || r.ToLong > 180 || r.FromLong < 0 || r.FromLong > 180 {
		err = errors.New("longitude must be between 0 and 180")
		return err
	}

	return err
}
