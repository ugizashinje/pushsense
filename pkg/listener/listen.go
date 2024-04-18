package listener

import (
	"fmt"
	"time"

	"github.com/GoWebProd/uuid7"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ugizashinje/pushsense/conf"
	"github.com/ugizashinje/pushsense/pkg/sender"
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
		updated_at TIME,
		last_transaction TIME,
		collection VARCHAR(32)
	);
`

var getLastStatus = `
	SELECT * from pushsense_statuses ps WHERE ps.collection = $1 ORDER BY ps.updated_at LIMIT 1
`
var startLog = `
	INSERT INTO pushsense_statuses (updated_at, last_transaction, collection)
	VALUES (now(), '1970-01-01 00:00:00'::timestamp , $1)
`
var logPush = `
	UPDATE pushsense_statuses SET updated_at = now() , last_transaction = $2
	WHERE collection = $1

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
			_, err = db.Exec(startLog, name)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}
		_, err = sender.CreateCollection(name, colConfig.Schema)
		if err != nil {
			panic(err)
		}
		go listenCollection(name, colConfig)
	}

}

func listenCollection(collection string, colConfig conf.Entry) {
	latest := time.Time{}

	for {
		rows, err := db.Queryx(colConfig.Sql, latest)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}
		// colTypes, err := rows.ColumnTypes()
		// for _, s := range colTypes {
		// 	fmt.Println("cols type:", s.Name(), s.DatabaseTypeName())
		// }
		allUpdates := []map[string]any{}
		allDeletions := []string{}
		for rows.Next() {
			mapRow := make(map[string]any)
			err = rows.MapScan(mapRow)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}
			if mapRow["is_deleted"] == true {
				id := mapRow["id"]
				allDeletions = append(allDeletions, id.(string))
			} else {
				allUpdates = append(allUpdates, mapRow)
			}

			for k, v := range colConfig.Processors {
				_, ok := mapRow[k]
				if ok {
					mapRow[k] = processors[v](mapRow[k])
				}
			}
			updatedAt, ok := mapRow["updated_at"]
			if ok {
				curent, ok := updatedAt.(time.Time)
				if ok && !curent.Before(latest) {
					latest = curent
				}
			}
		}
		if len(allDeletions) > 0 {
			err = sender.Delete(collection, allDeletions)
			if err != nil {
				fmt.Println("Error deleting ", err.Error())
			}
		}
		if len(allUpdates) > 0 {
			err = sender.Send(collection, allUpdates)
			if err != nil {
				fmt.Println("Error sending ", err.Error())
			}
		}

		_, err = db.Exec(logPush, collection, latest)
		if err != nil {
			fmt.Println("ERROR", err.Error())
		}
		if len(allUpdates) < 100 {
			time.Sleep(time.Second * 10)
		}
	}
}
