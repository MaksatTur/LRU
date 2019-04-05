package utils

import (
	"log"
	"strconv"

	"github.com/vaughan0/go-ini"
)

//Config contains conf data
type Config struct {
	ServerURL string
	Capacity  int
}

//NewConfig return Config
func NewConfig() *Config {
	file, err := ini.LoadFile("settings.ini")
	if err != nil {
		log.Fatal(err)
	}

	serverURL, ok := file.Get("server", "url")
	if !ok {
		log.Fatal("server url is nil")
	}

	capacity, ok := file.Get("others", "capacity")

	if !ok {
		log.Fatal("capacity is nil")
	}

	capacityInt, err := strconv.Atoi(capacity)

	if err != nil {
		log.Fatal("capacity is not int")
	}

	return &Config{
		ServerURL: serverURL,
		Capacity:  capacityInt,
	}
}
