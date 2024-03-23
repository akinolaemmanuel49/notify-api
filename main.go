package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/handlers"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
)

func handleRequests(notificationHandler *handlers.NotificationHandler) {
	// Define HTTP router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/notifications", notificationHandler.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications", notificationHandler.GetAllNotifications).Methods("GET")
	router.HandleFunc("/notifications/{id}", notificationHandler.GetNotificationByID).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	db, err := sql.Open("sqlite", "notifications.db")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Migrate notification table
	notificationRepository := repositories.NewNotificationRepository(db)
	err = notificationRepository.MigrateNotificationTable()
	if err != nil {
		log.Fatal("Error migrating notification table:", err)
	}

	// Initialize services
	notificationService := services.NewNotificationService(notificationRepository)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Start the server
	log.Println("Server started on port 8080")
	handleRequests(notificationHandler)
}
