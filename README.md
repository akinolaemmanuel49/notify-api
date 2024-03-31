# Notify API Documentation

The Notify API is a web API built using the Go programming language with the `net/http` library. It allows users to perform CRUD operations on the notifications resource and provides endpoints for managing user resources. The API utilizes JWT for authentication, features a rate limiter, and is served over HTTPS using NGINX. Additionally, a Makefile is provided for ease of use in managing the application.

## Endpoints

### Notifications Resource

- `GET /notifications`: Retrieve all notifications.
- `GET /notifications/me`: Retrieve all notifications by current user.
- `GET /notifications/{id}`: Retrieve a specific notification by ID.
- `POST /notifications`: Create a new notification.
- `PUT /notifications/{id}`: Update an existing notification.
- `DELETE /notifications/{id}`: Delete a notification by ID.

### Users Resource

- `GET /users`: Retrieve all users (restricted to authenticated users).
- `GET /users/{id}`: Retrieve a specific user by ID.
- `POST /users`: Create a new user.
- `PUT /users/{id}`: Update an existing user (restricted to the owner).
- `DELETE /users/{id}`: Delete a user by ID (restricted to the owner).

## Authentication

- The API utilizes JSON Web Tokens (JWT) for authentication.
- Users must include a valid JWT token in the Authorization header for protected endpoints.

## Rate Limiting

- The API implements a rate limiter to prevent abuse and ensure fair usage.
- Each user is limited to a certain number of requests per time interval (e.g., 100 requests per hour).

## NGINX Setup

- NGINX is used to set up an HTTPS server for the Notify API.
- Configuration includes SSL certificate and key paths, as well as proxy pass directives to forward requests to the Go application.

## Makefile

The Makefile provides convenience commands for building, running, and managing the Notify API.

- `make`: Builds and runs the Notify API.
- `make build`: Builds the Go application.
- `make run`: Starts the Go application.
- `make stop`: Stops the Go application.
- `make start-nginx`: Starts NGINX.
- `make stop-nginx`: Stops NGINX.
- `make clean`: Cleans up generated files.

## Usage

1. Clone the Notify API repository.
```bash
git clone https://github.com/akinolaemmanuel49/notify-api.git
```
2. Set up NGINX with SSL certificates and keys.
3. Customize NGINX configuration to match your environment.
4. Run `make` to build and start the Notify API.
5. Access the API endpoints using a tool like cURL or a web browser, including a valid JWT token in the Authorization header for protected endpoints.

## Conclusion

The Notify API provides a secure and efficient way for users to manage notifications and user resources. With JWT authentication, rate limiting, and NGINX setup, it ensures reliability and security while offering ease of use and convenience through the provided Makefile.