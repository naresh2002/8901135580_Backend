package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/naresh2002/8901135580_Backend/db"
	"github.com/naresh2002/8901135580_Backend/handlers"
)

func main() {

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	logger := log.New(os.Stdout, "API: ", log.LstdFlags)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(logger, database)
	fileHandler := handlers.NewFileHandler(logger, database)

	router := mux.NewRouter()

	// GET Subrouter

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signup", userHandler.Signup)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/login", userHandler.Login)

	postRouter.HandleFunc("/upload", fileHandler.UploadFile).Methods("POST")
	getRouter.HandleFunc("/file", fileHandler.ServeFile).Methods("GET")

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Printf("Server starting on port 8000\n")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
