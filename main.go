package main

import (
	"lru/api"
	"lru/utils"
)

func main() {
	config := utils.NewConfig()
	service := api.NewService(config.ServerURL, config.Capacity)
	service.StartService()
}
