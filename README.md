# Go Events REST API

> A production-style REST API built in Go for managing events and user registrations.
> Designed to demonstrate clean backend architecture, JWT authentication, validation, and testability.

---

## ✨ Features

- 🔐 JWT-based authentication
- 🛡 Password hashing with bcrypt
- 🔑 Protected routes using custom middleware
- 📋 Full CRUD operations for events
- 👤 Ownership checks — only event creators can update/delete
- 🗓 Event registration & cancellation
- 📄 Pagination for scalable listing
- 🧪 Structured validation using go-playground/validator
- ⚙️ Environment-based configuration
- 🧾 Unit tests with Go testing package
- 🪵 Custom logging middleware (response times + status codes)

---

## 🛠 Tech Stack

- **Go**
- **Gin** HTTP framework
- **SQLite** (modernc.org/sqlite, pure Go)
- **JWT** authentication
- **bcrypt** password hashing
- **go-playground/validator**
- **godotenv** for env variables

---

## 🧱 Project Structure

```
go-events-api/
├── api-test/        # Tests for routes and handlers
├── config/          # Environment setup
├── db/              # Database connection and setup
├── middleware/      # Auth & logging middleware
├── models/          # Database models & queries
├── routes/          # API route handlers
├── utils/           # Utils (JWT, validation, hashing)
├── .gitignore
├── go.mod
└── main.go
```

This structure cleanly separates concerns for maintainability and clarity.

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/GVaibhav92/Events-REST-API.git
cd Events-REST-API
```

### 2️⃣ Install Dependencies

```bash
go mod tidy
```

### 3️⃣ Create a `.env` File

```
PORT=8080
DB_PATH=api.db
JWT_SECRET=your-secret-key
```

> ⚠️ Do **not** commit `.env` or any secret values.

### 4️⃣ Run the Server

```bash
go run main.go
```

The API will start at:

```
http://localhost:8080
```

---

## 🔐 Authentication

Protected routes require:

```
Authorization: Bearer <token>
```

Token is obtained via:

```
POST /login
```

---

## 📌 API Overview

| Method | Route | Description | Protected |
|--------|-------|-------------|-----------|
| POST | `/signup` | Create a new user | ❌ |
| POST | `/login` | Login a user | ❌ |
| GET | `/events` | List events (paginated) | ❌ |
| GET | `/events/:id` | Get a specific event | ❌ |
| POST | `/events` | Create event | ✅ |
| PUT | `/events/:id` | Update event | ✅ |
| DELETE | `/events/:id` | Delete event | ✅ |
| POST | `/events/:id/register` | Register for event | ✅ |
| DELETE | `/events/:id/register` | Cancel registration | ✅ |

---

## 📄 Pagination

Supports:

```
GET /events?page=1&limit=10
```

Response includes:

- total records
- current page
- total pages
- limit

---

## 🛡 Validation Rules

- **email** — required, valid email format
- **password** — required, 6–72 characters
- **event name** — required, 3–100 characters
- **event description** — required, 10–500 characters
- **event location** — required, 3–100 characters
- **event dateTime** — must be a future date

Validation errors return a structured JSON response.

---

## 📦 Error Format

**Standard Error:**
```json
{
  "message": "description of what went wrong"
}
```

**Validation Error:**
```json
{
  "message": "validation failed",
  "errors": [
    { "field": "name", "message": "must be at least 3 characters" }
  ]
}
```

---

## 🧪 Testing

Run all tests:

```bash
go test ./...
```

## 🔮 Future Improvements

- PostgreSQL migration
- Refresh token support
- OAuth login
- API rate limiting
- Containerization (Docker)
- gRPC microservices
- Redis caching

---

## 📌 Why This Project

Built to understand:

- Backend architecture in Go
- Middleware design
- Auth flows & secure routes
- Database modeling with ownership
- Testable and maintainable code

---
