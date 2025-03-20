# Movie App

A RESTful API for managing movies, built with Go and Gin.

## Features

- CRUD operations for movies
- User authentication with JWT
- PostgreSQL database
- Docker support
- Swagger documentation

## Environment Setup

1. Copy the example environment file:

   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your configuration:
   - Database credentials
   - Server port
   - JWT secret (make sure to use a strong secret in production)

## Installation

### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/movie_app.git
   cd movie_app
   ```

2. Install dependencies:

   ```bash
   make deps
   ```

3. Run the application:
   ```bash
   make run
   ```

## API Documentation

The API documentation is available at `/swagger/index.html` when the server is running.

### Authentication

The API uses JWT for authentication. To access protected endpoints:

1. Register a new user:

   ```bash
   curl -X POST http://localhost:8080/api/v1/users/register \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "password123"}'
   ```

2. Login to get a token:

   ```bash
   curl -X POST http://localhost:8080/api/v1/users/login \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "password123"}'
   ```

3. Use the token in subsequent requests:
   ```bash
   curl -X GET http://localhost:8080/api/v1/movies \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```