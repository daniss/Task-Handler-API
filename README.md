# Task Management API

## Overview

This is a RESTful API built with Go using the Gin framework and GORM for database interactions. The API allows users to manage tasks with features for user authentication and role-based access control.

## Features

- **User Registration**: Allows users to register and securely store passwords using hashing.
- **User Login**: Users can log in and receive a JWT token for authentication.
- **Task Management**: Users can create, retrieve, update, and delete tasks.
- **Role-Based Access Control**: Restricts certain actions based on user roles (e.g., Admin, User).

## Technologies Used

- **Go**: Programming language for building the API.
- **Gin**: HTTP web framework for Go.
- **GORM**: Object-Relational Mapping (ORM) library for Go.
- **PostgreSQL/MySQL**: Database for storing user and task data (specify your choice).

## Installation

### Prerequisites

- Go (1.17 or later)
- Database (PostgreSQL or MySQL)
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/task-management-api.git
cd task-management-api
```

### Set Up Environment Variables

Create a `.env` file in the root directory and add the following environment variables:

```plaintext
DATABASE=your_database_url
SECRETKEY=your_secret_key
```

### Install Dependencies

Make sure you have Go modules enabled and run:

```bash
go mod tidy
```

### Database Setup

Run the necessary migrations to create the database tables. Make sure to have GORM set up correctly for your database.

```go
// Example of running migrations
func setupDatabase() *gorm.DB {
    db, err := gorm.Open("your_database_driver", "your_database_url")
    if err != nil {
        log.Fatal(err)
    }
    
    db.AutoMigrate(&User{}, &Task{}) // Adjust as needed
    return db
}
```

## Usage

### Run the API

You can run the API using the following command:

```bash
go run main.go
```

The API will be accessible at `http://localhost:8080`.

### API Endpoints

#### User Registration

- **POST** `/register`
  - **Request Body**: `{ "username": "example", "password": "examplePassword" }`
  - **Response**: `{ "message": "User registered successfully" }`

#### User Login

- **POST** `/login`
  - **Request Body**: `{ "username": "example", "password": "examplePassword" }`
  - **Response**: `{ "token": "your_jwt_token" }`

#### Create Task

- **POST** `/createTask`
  - **Request Body**: `{ "title": "Task Title", "description": "Task Description", "completed": false }`
  - **Response**: `{ "ID": 1, "Title": "Task Title", "Description": "Task Description", "Completed": false }`

#### Get All Tasks

- **GET** `/tasks`
  - **Response**: `[{ "ID": 1, "Title": "Task Title", "Description": "Task Description", "Completed": false }]`

#### Update Task

- **PUT** `/updateTask/:id`
  - **Request Body**: `{ "title": "Updated Title", "description": "Updated Description", "completed": true }`
  - **Response**: `{ "message": "Task updated successfully" }`

#### Delete Task

- **DELETE** `/deleteTask/:id`
  - **Response**: `{ "message": "Task deleted successfully" }`

## Middleware

- **JWT Authentication**: Protects routes by requiring a valid JWT token in the `Authorization` header for certain endpoints.

## Testing

You can use tools like [Postman](https://www.postman.com/) to test the API endpoints.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Feel free to modify any sections or add additional details that are relevant to your specific implementation! Let me know if you need further adjustments or additions.
