package conf

import (
	"encoding/json"
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
	Mapping    map[string]string                   `json:"mapping"`
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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func init() {
	var config Config
	configFile := getenv("PUSHSENSE_CONFIG_FILE", "/app/config/conf.json")

	jsonFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("cant read  "+configFile, err.Error())
	}
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(1)
	}
	log.Println("database connected")
	DB = config.Db
	Collections = config.Collections
	Typesense = config.Typsense
}

var DB DbConfig
var Collections map[string]Entry
var Typesense TypesenseConfig
