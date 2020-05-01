package query

// InsertDataPoint is the Postgresql query for inserting a data point
var InsertDataPoint string = `
INSERT INTO datapoint (
id,
user_id, 
vehicle_model_id, 
package_id, 
travel_type_id,
from_area_id, 
to_area_id, 
from_city_id, 
to_city_id, 
from_date, 
to_date, 
online_booking, 
mobile_site_booking,
booking_created, 
from_lat, 
from_long, 
to_lat, 
to_long, 
Car_Cancellation
) VALUES (
$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
)
`
