### Notify Web API Documentation

#### Introduction
The Notify Web API provides functionality for managing users and notifications. It employs RESTful principles for interaction and utilizes JSON Web Tokens (JWT) for authentication. This documentation outlines the endpoints, their functionalities, and any access restrictions.

#### Base URL
```
https://api.notify.com
```

#### Authentication
The Notify API uses JWT for authentication. To authenticate, include the generated token in the Authorization header with the format: `Bearer <token>`.

#### Data Structures

##### AuthCredentials
```go
type AuthCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

##### Notification
```go
type Notification struct {
    ID          int64    `json:"id"`
    Title       string   `json:"title"`
    Message     string   `json:"message"`
    Priority    Priority `json:"priority"`
    PublisherID int64    `json:"publisher_id"`
    CreatedAt   string   `json:"created_at"`
    UpdatedAt   string   `json:"updated_at"`
}
```

##### NotificationInput
```go
type NotificationInput struct {
    Title    string   `json:"title"`
    Message  string   `json:"message"`
    Priority Priority `json:"priority"`
}
```

##### NotificationResponse
```go
type NotificationResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
```

##### User
```go
type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
```

##### UserProfile
```go
type UserProfile struct {
    ID        int64  `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

##### UserInput
```go
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
```

##### UserInputWithPassword
```go
type UserInputWithPassword struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}
```

##### UserResponse
```go
type UserResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
```

#### Endpoints

##### 1. Authentication

###### Generate Token
- **Endpoint:** `/auth/token`
- **Method:** POST
- **Description:** Generates a JWT token for authentication.
- **Request Body:** AuthCredentials
- **Access:** Unprotected
- **Sample Response:**
    ```json
    {
	"token": <token>
    }
    ```

##### 2. User Management

###### Create User
- **Endpoint:** `/users`
- **Method:** POST
- **Description:** Creates a new user.
- **Request Body:** UserInputWithPassword
- **Access:** Unprotected
- **Sample Response:**
    ```json
    {
	"code": 201,
	"message": "User was successfully created"
    }
    ```

###### Get User By ID
- **Endpoint:** `/users/{userId}`
- **Method:** GET
- **Description:** Retrieves user information by ID.
- **Access:** Unprotected
- **Sample Response:**
    ```json
    {
        "code": 200,
        "data": {
            "id": 2,
            "first_name": "John",
            "last_name": "Doe",
            "email": "johndoe@mail.com",
            "created_at": "2024-03-26T20:43:55+01:00",
            "updated_at": "2024-03-26T20:43:55+01:00"
        },
        "message": "User with ID: 2 was successfully retrieved"
    }
    ```

###### Update User
- **Endpoint:** `/users/{userId}`
- **Method:** PUT
- **Description:** Updates user information.
- **Request Body:** UserProfile
- **Access:** Protected (only the user can update their own information)

###### Delete User
- **Endpoint:** `/users/{userId}`
- **Method:** DELETE
- **Description:** Deletes a user.
- **Access:** Protected (only the user can delete their own account)

###### Get All Users
- **Endpoint:** `/users`
- **Method:** GET
- **Description:** Retrieves information for all users.
- **Access:** Unprotected

##### 3. Notification Management

###### Create Notification
- **Endpoint:** `/notifications`
- **Method:** POST
- **Description:** Creates a new notification.
- **Request Body:** NotificationInput
- **Access:** Protected
- **Sample Response:**
    ```json
    {
	"code": 201,
	"message": "Notification was successfully created"
    }
    ```


###### Get Notification By ID
- **Endpoint:** `/notifications/{notificationId}`
- **Method:** GET
- **Description:** Retrieves a notification by ID.
- **Access:** Unprotected
- **Sample Response:**
    ```json
    {
        "code": 200,
        "data": {
            "id": 2,
            "title": "New Notification Title 2",
            "message": "This is a sample notification message 2.",
            "priority": 2,
            "publisher_id": 1,
            "created_at": "2024-03-26T10:00:00+01:00",
            "updated_at": "2024-03-29T16:56:21+01:00"
        },
        "message": "Notification with ID: 2 was successfully retrieved"
    }
    ```

###### Update Notification
- **Endpoint:** `/notifications/{notificationId}`
- **Method:** PUT
- **Description:** Updates a notification.
- **Request Body:** NotificationInput
- **Access:** Protected (only the publisher can update their own notification)
- **Request Headers:**
    ```http
    Authorization: Bearer <token>
    Content-Type: application/json
    ```
- **Sample Response:**
    ```json
    {
        "code": 200,
        "message": "Notification with ID: 2 was successfully updated"
    }
    ```

###### Delete Notification
- **Endpoint:** `/notifications/{notificationId}`
- **Method:** DELETE
- **Description:** Deletes a notification.
- **Access:** Protected (only the publisher can delete their own notification)
- **Sample Response:**
    ```json
    {
        "code": 200,
        "message": "Notification with ID: 2 was successfully deleted"
    }
    ```

###### Get All Notifications
- **Endpoint:** `/notifications`
- **Method:** GET
- **Description:** Retrieves all notifications.
- **Access:** Unprotected
- **Query Parameters:**
  - `page` (optional): Specifies the page number for pagination. Default is 1.
  - `pageSize` (optional): Specifies the number of notifications per page. Default is 10.
- **Sample Response:**
    ```json
    {
	"code": 200,
	"data": [
		{
			"id": 4,
			"title": "New Notification Title 4",
			"message": "This is a sample notification message 4.",
			"priority": 1,
			"publisher_id": 1,
			"created_at": "2024-03-26T10:00:00+01:00",
			"updated_at": "2024-03-26T10:00:00+01:00"
		},
		{
			"id": 3,
			"title": "New Notification Title 3",
			"message": "This is a sample notification message 3.",
			"priority": 1,
			"publisher_id": 1,
			"created_at": "2024-03-26T10:00:00+01:00",
			"updated_at": "2024-03-26T10:00:00+01:00"
		},
		{
			"id": 6,
			"title": "John's Post.",
			"message": "This post is by John Doe.",
			"priority": 1,
			"publisher_id": 2,
			"created_at": "2024-03-27T12:38:24+01:00",
			"updated_at": "2024-03-27T12:39:10+01:00"
		}
	],
	"message": "Notifications successfully retrieved."
    }
    ```

#### Error Handling
- The API follows standard HTTP status codes for error handling.
- Detailed error messages are provided in the response body for better understanding of issues.

#### Rate Limiting
- The API implements rate limiting to prevent abuse. Exceeding the rate limit will result in HTTP 429 Too Many Requests status code.

#### Sample Request Headers
```http
Authorization: Bearer <token>
Content-Type: application/json
```

<!-- #### Sample Response
```json
{
  "message": "Notification created successfully.",
  "notificationId": "123456",
  "timestamp": "2024-03-27T12:00:00Z"
}
``` -->

#### Conclusion
This documentation provides an overview of the Notify Web API endpoints, their functionalities, and access restrictions. Developers can use this API to manage users and notifications securely.