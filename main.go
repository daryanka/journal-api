package main

import (
	"api/clients"
	"api/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize .env
	err := godotenv.Load()
	if err != nil {
		panic("error getting .env file details" + err.Error())
	}

	// Initialize DB connection
	clients.InitializeOrm()

	// Initialize routes
	routes.StartRouting()
}
