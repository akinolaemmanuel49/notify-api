package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/handlers"
	"github.com/akinolaemmanuel49/notify-api/middlewares"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func handleRequests(notificationHandler *handlers.NotificationHandler, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	// Define HTTP router
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()

	handleNotificationRequests(apiRouter, notificationHandler)
	handleUserRequests(apiRouter, userHandler)
	handleAuthRequest(apiRouter, authHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleNotificationRequests(apiRouter *mux.Router, notificationHandler *handlers.NotificationHandler) {
	// Notification Routes
	apiRouter.HandleFunc("/notifications/healthCheck", notificationHandler.NotificationHealthCheck).Methods("GET")
	apiRouter.HandleFunc("/notifications", middlewares.JWTAuthMiddleware(notificationHandler.CreateNotification)).Methods("POST")
	apiRouter.HandleFunc("/notifications", notificationHandler.GetAllNotifications).Methods("GET")
	apiRouter.HandleFunc("/notifications/{id}", notificationHandler.GetNotificationByID).Methods("GET")
	apiRouter.HandleFunc("/notifications/{id}", middlewares.JWTAuthMiddleware(notificationHandler.UpdateNotificationByID)).Methods("PUT")
	apiRouter.HandleFunc("/notifications/{id}", middlewares.JWTAuthMiddleware(notificationHandler.DeleteNotificationByID)).Methods("DELETE")
}

func handleUserRequests(apiRouter *mux.Router, userHandler *handlers.UserHandler) {
	// Users
	apiRouter.HandleFunc("/users/healthCheck", userHandler.UserHealthCheck).Methods("GET")
	apiRouter.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", middlewares.JWTAuthMiddleware(userHandler.UpdateUserByID)).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", middlewares.JWTAuthMiddleware(userHandler.DeleteUserByID)).Methods("DELETE")
}

func handleAuthRequest(apiRouter *mux.Router, authHandler *handlers.AuthHandler) {
	// Auth
	apiRouter.HandleFunc("/auth/token", authHandler.GenerateToken).Methods("POST")
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
	cfg.ReadFile("dev-config.yml") // For use in development
	// cfg.ReadFile("config.yml")
	cfg.ReadEnv()

	DATABASE_URI := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", cfg.Database.USER, cfg.Database.PASS, cfg.Database.NAME)
	db, err := sql.Open("postgres", DATABASE_URI)
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}
	defer db.Close()

	// Initialize repositories
	notificationRepository := repositories.NewNotificationRepository(db)
	userRepository := repositories.NewUserRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// Initialize services
	notificationService := services.NewNotificationService(notificationRepository)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(authRepository)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Start the server
	log.Println("Server started on port 8080")
	handleRequests(notificationHandler, userHandler, authHandler)
}
