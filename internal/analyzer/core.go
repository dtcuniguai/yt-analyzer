package analyzer

import (
	"fmt"
	"os"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB
var influx api.WriteAPIBlocking
var influxClient influxdb2.Client

func GetDB() (*sqlx.DB, error) {

	if db == nil {
		database := os.Getenv("DATABASE")
		dsn := os.Getenv("DATABASE_CONNECTION")

		d, err := sqlx.Open(database, dsn)
		if err != nil {
			return nil, err
		}

		err = d.Ping()
		if err != nil {
			return nil, err
		}

		db = d
	}

	return db, nil
}

func GetLogger() (influxdb2.Client, api.WriteAPIBlocking, error) {

	if influx == nil {

		//init config
		host := fmt.Sprintf("%v:%v", os.Getenv("INFLUX_HOST"), os.Getenv("INFLUX_PORT"))
		token := os.Getenv("INFLUX_TOKEN")
		org := os.Getenv("INFLUX_ORG")
		bucket := os.Getenv("INFLUX_BUCKET")

		client := influxdb2.NewClient(host, token)
		writeAPI := client.WriteAPIBlocking(org, bucket)

		influx = writeAPI
		influxClient = client
	}

	return influxClient, influx, nil

}
