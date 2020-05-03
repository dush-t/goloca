package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dush-t/goloca/db/query"

	_ "github.com/lib/pq" // Postgres database driver
)

// PostgresConn contains info about the active database connection
type PostgresConn struct {
	db *sql.DB
}

// Connect will establish a connection with Postgres. It will retry at an interval
// of 10 seconds upto 10 times if it fails to connect
func (p *PostgresConn) Connect() {
	maxConnectAttempts, _ := strconv.Atoi(os.Getenv("MAX_DB_CONNECTION_ATTEMPTS"))
	reconnectInterval := 10 * time.Second
	attempts := 0
	for {
		err := p.connect()
		if err != nil {
			if attempts < maxConnectAttempts {
				log.Println("Error connecting to database, retrying")
				attempts++
				time.Sleep(reconnectInterval)
			} else {
				log.Fatal("Unable to connect to database:", err)
				break
			}
		} else {
			break
		}
	}
}

func (p *PostgresConn) connect() error {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", psqlInfo)
	(*p).db = db
	err = db.Ping()
	if err != nil {
		return err
	}
	log.Println("Successfully connected to Postgres database")
	log.Println("Creating datapoint table if it does not already exist")
	_, tableErr := (*p).db.Exec(query.CreateTable)
	if tableErr != nil {
		return tableErr
	}

	return nil
}

// Config returns data about the database connection
func (p *PostgresConn) Config() interface{} {
	return *p
}

// Insert is used to add rows to the database with the given endpoint
func (p *PostgresConn) Insert(data interface{}) error {
	db := (*p).db
	d, ok := data.(RideDataPoint)
	if !ok {
		err := errors.New("Invalid data")
		return err
	}

	// d := RideDataPoint(data.(RideDataPoint))

	result, err := db.Exec(query.InsertDataPoint,
		d.ID,
		d.UserID,
		d.VehicleModelID,
		d.PackageID,
		d.TravelTypeID,
		d.FromAreaID,
		d.ToAreaID,
		d.FromCityID,
		d.ToCityID,
		d.FromDate,
		d.ToDate,
		d.OnlineBooking,
		d.MobileSiteBooking,
		d.BookingCreated,
		d.FromLat,
		d.FromLong,
		d.ToLat,
		d.ToLong,
		d.CarCancellation,
	)

	if err != nil {
		return err
	}

	log.Println("Record added to database: ", result)

	return nil
}
