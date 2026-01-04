package main

import (
	"log"
	"net/http"

	"golang-http-patch/database"
	"golang-http-patch/handlers"
	"golang-http-patch/validation"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	database.InitDB()

	// Initialize validator
	validation.InitValidator()

	// Create router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.PatchUser).Methods("PATCH")

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
