# API Documentation

### Back-End

This API is built using the **Fiber framework** in **Go** and uses **SQLite** as the database engine. It includes various routes for managing tasks, authentication, and more. The API is designed to be secure, performant, and highly available.

## Getting Started

### Prerequisites
- Go 1.16 or higher
- SQLite
- Redis

### Installation
1. Clone the repository:
```bash
git clone https://github.com/pageton/todo-list.git
```

2. Change to the project directory:
```bash
cd todo-list
```

3. Install dependencies:
```bash
go mod tidy
```

4. Set up environment variables:
```bash
cp .env.example .env
```

5. Edit `.env` file with your configuration

6. Run the application:
```bash
go run cmd/main.go
```

The server will start at `http://localhost:3000`

## Features

- **Compression & Caching**: Supports HTTP compression and caching with Redis to enhance performance.
- **Rate Limiting**: Implements rate limiting using Redis to prevent abuse.
- **Enhanced Security**: Implements comprehensive security measures including advanced security headers like `X-Content-Type-Options`, `X-Frame-Options`, and `Strict-Transport-Security`, robust encryption, and multi-layer protection mechanisms for maximum security.
- **Error Handling**: Custom error handler for internal server errors.
- **CORS**: Configured to allow requests from any origin (CORS support).
- **Timeouts**: Includes request timeouts to ensure high responsiveness.

## Framework

- **Fiber**: A fast and lightweight web framework for Go.

## Database Engine

- **SQLite**: A lightweight and self-contained SQL database engine used to store tasks.

---

## Routes

### Task Management

1. **POST /api/task/create**
   - Creates a new task.
   - Handler: `CreateTaskHandler`

2. **GET /api/tasks**
   - Retrieves all tasks.
   - Handler: `GetTasksHandler`

3. **GET /api/tasks/:limit**
   - Retrieves tasks with a limit query parameter.
   - Handler: `TaskLimiterHandler`

4. **GET /api/task/:task_id**
   - Retrieves a specific task by its ID.
   - Handler: `TaskByIdHandler`

5. **PUT /api/task/:task_id**
   - Updates an existing task by its ID.
   - Handler: `UpdateTaskHandler`

6. **DELETE /api/task/:task_id**
   - Deletes a task by its ID.
   - Handler: `DeleteTaskHandler`

### Authentication

1. **POST /api/auth/register**
   - Registers a new user.
   - Handler: `RegisterHandler`

2. **POST /api/auth/login**
   - Logs in an existing user.
   - Handler: `LoginHandler`

---

## Security

- **Encryption**: The API uses strong encryption mechanisms to handle sensitive data. This ensures that passwords and other sensitive information are protected both in transit and at rest.
- **SSL/TLS**: Communication between the client and server is encrypted using SSL/TLS to ensure data confidentiality.

---

## Front-End (Coming Soon)

A front-end application will be available soon to interact with this API.

---

### Notes

- The API currently supports rate limiting and caching, but more features will be added in the future.
- The API uses a JWT (JSON Web Token) for authentication. The JWT is generated using the `jwt` package and contains information about the user's ID and other relevant details.
- The API uses a Redis store for caching and rate limiting. The Redis store is configured to use a connection pool to improve performance.
- The API uses a SQLite database for storing tasks. The SQLite database is configured to use a connection pool to improve performance.
- The API uses environment variables for configuration. The environment variables are loaded using the `godotenv` package.

## Example Response

```json
{
  "status": "success",
  "data": {
    "task_id": 1,
    "task_name": "Finish the API documentation",
    "status": "pending"
  }
}
```

---

## Rate Limiting

This API has rate limiting in place to ensure fair usage. Each IP address is limited to 100 requests per minute.

There is duplicate information about Redis and rate limiting in the Notes section that has already been covered in Features and Rate Limiting sections.

---

## Contributing

We welcome contributions or bug/vulnerability reports through issues. Feel free to report any problems or suggest improvements to help make this API better.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
