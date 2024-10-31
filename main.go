// main.go
// This file initializes the HTTP server for the receipt processing service.
// It sets up the routes and starts the server on port 8080.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/saurabhag23/receipt-processor/internal/handlers"
)

func main() {
	// Initialize a logger to output server logs to the console.
	// The logs are prefixed with "receipt-processor: " and include timestamps.
	logger := log.New(os.Stdout, "receipt-processor: ", log.LstdFlags)

	// Create a new router using Gorilla Mux for handling HTTP routes.
	r := mux.NewRouter()

	// Define the HTTP route for processing receipts.
	// This route listens for POST requests at /receipts/process and calls the ProcessReceipt handler.
	r.HandleFunc("/receipts/process", handlers.ProcessReceipt).Methods("POST")

	// Define the HTTP route for retrieving points for a specific receipt by ID.
	// This route listens for GET requests at /receipts/{id}/points and calls the GetPoints handler.
	r.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")

	// Start the HTTP server on port 8080 with the configured routes.
	// If the server encounters a fatal error, log it and exit.
	logger.Println("Server starting on port 8080...")
	logger.Fatal(http.ListenAndServe(":8080", r))
}
