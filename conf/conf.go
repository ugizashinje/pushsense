package conf

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Db          DbConfig         `yaml:"db"`
	Collections map[string]Entry `yaml:"collections"`
}

type Entry struct {
	Sql string `yaml:"sql"`
}
type DbConfig struct {
	Url string `yaml:"url"`
}

func init() {
	var config Config
	dir := "./"

	yamlFile, err := os.ReadFile(dir + "conf.yaml")
	if err != nil {
		log.Fatal("yamlFile.Get err  # ", err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(1)
	}
	fmt.Println(config.Db.Url)
	DB = config.Db
	Collections = config.Collections
}

var DB DbConfig
var Collections map[string]Entry
