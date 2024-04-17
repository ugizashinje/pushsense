package listener

import (
	"fmt"
	"os"
	"time"

	"github.com/GoWebProd/uuid7"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ugizashinje/pushsense/conf"
)

type Status struct {
	ID              int64  `db:"id"`
	TableName       string `db:"table_name"`
	UpdatedAt       string `db:"updated_at"`
	LastTransaction string `db:"last_transaction"`
	Collection      string `db:"collection"`
}

var tables_statuses map[string]time.Time = make(map[string]time.Time)

var GenUUID *uuid7.Generator
var db *sqlx.DB
var dropAll = `
DROP TABLE IF EXISTS pushsense_statuses CASCADE;
DROP SEQUENCE IF EXISTS ps_sequence;
`
var createSchema string = `

	CREATE SEQUENCE IF NOT EXISTS ps_sequence START 1 CYCLE ;

	CREATE TABLE IF NOT EXISTS pushsense_statuses (
		id INTEGER PRIMARY KEY default nextval('ps_sequence'),
		table_name VARCHAR(32),
		updated_at TIME,
		last_transaction TIME,
		collection VARCHAR(32)
	);
`

var getLastStatus = `
	SELECT * from pushsense_statuses ps WHERE ps.table_name = $1 ORDER BY ps.updated_at LIMIT 1
`
var startLog = `
	INSERT INTO pushsense_statuses (table_name, updated_at, last_transaction, collection)
	VALUES ($1, $2, $3, $4)
`

func init() {
	db = sqlx.MustOpen("postgres", conf.DB.Url)
	createSchema = dropAll + createSchema
	db.MustExec(createSchema)
	GenUUID = uuid7.New()
}

func Start() {
	status := Status{}
	for name, colConfig := range conf.Collections {
		err := db.Get(&status, getLastStatus, name)
		if err.Error() == "sql: no rows in result set" {
			db.Exec(startLog)
		} else if err != nil {
			fmt.Println("Fatal:" + err.Error())
			os.Exit(1)
		}
		fmt.Println(name, colConfig)
		go listenCollection(name, colConfig)
	}

}

func listenCollection(collection string, colConfig conf.Entry) {

	for {
		rows, err := db.Queryx(colConfig.Sql)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}

		allUpdates := []map[string]any{}
		latest = time.Time{}
		for rows.Next() {
			mapRow := make(map[string]any)
			err = rows.MapScan(mapRow)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}
			allUpdates = append(allUpdates, mapRow)
		}
		fmt.Println("No updates", len(allUpdates))
		time.Sleep(time.Second)
	}
}
