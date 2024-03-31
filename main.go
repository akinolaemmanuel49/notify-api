package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/handlers"
	"github.com/akinolaemmanuel49/notify-api/middlewares"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/akinolaemmanuel49/notify-api/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func handleRequests(notificationHandler *handlers.NotificationHandler, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	// Define HTTP router
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middlewares.RateLimitMiddleware)

	handleNotificationRequests(apiRouter, notificationHandler)
	handleUserRequests(apiRouter, userHandler)
	handleAuthRequest(apiRouter, authHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Println("Server started on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func handleNotificationRequests(apiRouter *mux.Router, notificationHandler *handlers.NotificationHandler) {
	// Notification Routes
	apiRouter.HandleFunc("/notifications/healthCheck", notificationHandler.NotificationHealthCheck).Methods("GET")
	apiRouter.HandleFunc("/notifications", middlewares.JWTAuthMiddleware(notificationHandler.CreateNotification)).Methods("POST")
	apiRouter.HandleFunc("/notifications/me", middlewares.JWTAuthMiddleware(notificationHandler.GetOwnNotifications)).Methods("GET")
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

func main() {
	utils.LoadEnv()

	var cfg config.Config
	cfg.ReadFile("dev-config.yml") // For use in development
	cfg.ReadEnv()

	DATABASE_URI := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	db, err := sql.Open("postgres", DATABASE_URI)
	if err != nil {
		log.Println("Error connecting to the database:", err)
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

	// Handle requests
	handleRequests(notificationHandler, userHandler, authHandler)
}
