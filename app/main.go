package main

import (
	"fmt"
	"log"
	"net/http"

	"hotelm/db"
	"hotelm/routes"
)

func main() {
	db.Connect()
	defer db.Close()

	// Setup routes
	routes.SetupRoutes()

	// Start the server
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
