package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/handlers"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}
func main() {
	LoadEnv()

	var cfg config.Config
	cfg.ReadFile("config.yml")
	cfg.ReadEnv()

	db, err := sql.Open("postgres", cfg.Database.URI)
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}
	defer db.Close()

	// Initialize repositories
	notificationRepository := repositories.NewNotificationRepository(db)

	// Initialize services
	notificationService := services.NewNotificationService(notificationRepository)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Start the server
	log.Println("Server started on port 8080")
	handleRequests(notificationHandler)
}
