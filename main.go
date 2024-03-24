package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/handlers"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
)

func handleRequests(notificationHandler *handlers.NotificationHandler) {
	// Define HTTP router
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/notifications/healthCheck", notificationHandler.NotificationHealthCheck).Methods("GET")
	apiRouter.HandleFunc("/notifications", notificationHandler.CreateNotification).Methods("POST")
	apiRouter.HandleFunc("/notifications", notificationHandler.GetAllNotifications).Methods("GET")
	apiRouter.HandleFunc("/notifications/{id}", notificationHandler.GetNotificationByID).Methods("GET")
	apiRouter.HandleFunc("/notifications/{id}", notificationHandler.UpdateNotificationByID).Methods("PUT")
	apiRouter.HandleFunc("/notifications/{id}", notificationHandler.DeleteNotificationByID).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	var cfg config.Config
	cfg.ReadFile("config.yml")

	db, err := sql.Open("sqlite", cfg.Database.URI)
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
