package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/typesense/typesense-go/typesense/api"
)

type Config struct {
	Db          DbConfig         `json:"db"`
	Collections map[string]Entry `json:"collections"`
	Typsense    TypesenseConfig  `json:"typesense"`
}

type Entry struct {
	Sql        string                              `json:"sql"`
	TableName  string                              `json:"tableName"`
	Processors map[string]string                   `json:"processors"`
	Schema     api.CreateCollectionJSONRequestBody `json:"schema"`
}
type DbConfig struct {
	Url string `json:"url"`
}

type TypesenseConfig struct {
	Url                       string `json:"url"`
	ApiKey                    string `json:"apiKey"`
	ConnectionTimeout         int    `json:"connectionTimeout"`
	CircuitBreakerMaxRequest  int    `json:"circuitBreakerMaxRequest"`
	CircuitBreakerMaxInterval int    `json:"circuitBreakerMaxInterval"`
	CircuitBreakerMaxTimeout  int    `json:"circuitBreakerMaxTimeout"`
}

func init() {
	var config Config
	dir := "./"

	jsonFile, err := os.ReadFile(dir + "conf.json")
	if err != nil {
		log.Fatal("jsonConfig.Get err  # ", err.Error())
	}
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(1)
	}
	fmt.Println(config.Db.Url)
	DB = config.Db
	Collections = config.Collections
	Typesense = config.Typsense
}

var DB DbConfig
var Collections map[string]Entry
var Typesense TypesenseConfig
