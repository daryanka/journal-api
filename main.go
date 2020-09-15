package main

import (
	"api/clients"
	"api/routes"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Initialize .env
	err := godotenv.Load()
	if err != nil {
		panic("error getting .env file details" + err.Error())
	}

	// Initialize DB connection
	clients.InitializeOrm()
	// Setup DB logging
	f, err := os.OpenFile("db_query.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("unable to create db query log")
	}
	clients.ClientOrm.EnableQueryLogger(f)

	// Initialize routes
	routes.StartRouting()
}
