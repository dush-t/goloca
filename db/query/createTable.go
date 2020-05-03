package query

// CreateTable is the SQL query used to create the
// datapoints table
// var CreateTable string = `

// CREATE TABLE IF NOT EXISTS datapoint(
// 	id bigserial PRIMARY KEY,
// 	user_id VARCHAR(100) UNIQUE NOT NULL,
// 	vehicle_model_id VARCHAR(100) NOT NULL,
// 	package_id SMALLINT,
// 	travel_type_id SMALLINT,
// 	from_area_id VARCHAR(100),
// 	to_area_id VARCHAR(100) ,
// 	from_city_id VARCHAR(100),
// 	to_city_id VARCHAR(100),
// 	from_date BIGINT,
// 	to_date BIGINT,
// 	online_booking BOOLEAN,
// 	mobile_site_booking BOOLEAN,
// 	booking_created BIGINT,
// 	from_lat float8,
// 	from_long float8,
// 	to_lat float8,
// 	to_long float8,
// 	Car_Cancellation BOOLEAN
// )
// `

// CreateTable is the SQL query used to create the
// datapoints table
var CreateTable string = `

CREATE TABLE IF NOT EXISTS datapoint(
	id bigserial,
	user_id VARCHAR(100) NOT NULL, 
	vehicle_model_id VARCHAR(100) NOT NULL, 
	package_id SMALLINT, 
	travel_type_id SMALLINT,
	from_area_id VARCHAR(100), 
	to_area_id VARCHAR(100) ,
	from_city_id VARCHAR(100), 
	to_city_id VARCHAR(100), 
	from_date BIGINT, 
	to_date BIGINT, 
	online_booking BOOLEAN, 
	mobile_site_booking BOOLEAN,
	booking_created BIGINT, 
	from_lat float8, 
	from_long float8, 
	to_lat float8, 
	to_long float8, 
	Car_Cancellation BOOLEAN
)
` // This schema is used for benchmarking
