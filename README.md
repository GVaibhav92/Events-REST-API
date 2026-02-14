# Go Events REST API

> A production-style REST API built in Go for managing events and user registrations.  
> Designed to demonstrate backend architecture, authentication, validation, and clean project structure using Go.

---

## ✨ Features

- 🔐 JWT-based authentication
- 🔑 Password hashing using bcrypt
- 🛡 Protected routes via custom middleware
- 📋 Full CRUD operations for events
- 👤 Ownership checks (only creators can update/delete events)
- 📝 Event registration & cancellation
- 📄 Pagination support for scalable event listing
- ✅ Structured request validation
- ⚙️ Environment-based configuration (no hardcoded secrets)
- 🧪 Unit testing with Go’s built-in testing package
- 🧾 Custom logging middleware (status codes + response times)

---

## 🛠 Tech Stack

- **Go**
- **Gin** (HTTP framework)
- **SQLite** (modernc.org/sqlite — pure Go, no CGO required)
- **golang-jwt/jwt**
- **bcrypt**
- **go-playground/validator**
- **godotenv**

---

## 🧱 Project Structure

```
go-events-api/
├── config/         # Environment variable loading
├── db/             # Database connection & schema setup
├── middleware/     # Authentication & logging middleware
├── models/         # Database models & query logic
├── routes/         # Route handlers
├── utils/          # Shared utilities (JWT, hashing, validation)
├── .gitignore
├── go.mod
└── main.go
```

### Architecture Overview

- **Routes** handle HTTP layer and request/response lifecycle  
- **Models** handle database operations  
- **Middleware** handles authentication and request logging  
- **Utils** provides reusable helpers (JWT, hashing, validation)  
- **Config** loads environment variables  

The project follows a clean separation of concerns to keep business logic independent from routing and middleware layers.

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/go-events-api.git
cd go-events-api
```

### 2️⃣ Install Dependencies

```bash
go mod tidy
```

### 3️⃣ Create `.env` File

Create a `.env` file in the root directory:

```
PORT=8080
DB_PATH=api.db
JWT_SECRET=your-secret-key
```

> ⚠️ Never commit your `.env` file.

### 4️⃣ Run the Server

```bash
go run main.go
```

Server starts at:

```
http://localhost:8080
```

---

## 🔐 Authentication

Protected routes require a JWT token in the header:

```
Authorization: Bearer <token>
```

You can obtain a token via:

```
POST /login
```

---

## 📌 API Overview

| Method | Route | Description | Protected |
|--------|-------|------------|-----------|
| POST | /signup | Create user | ❌ |
| POST | /login | Login user | ❌ |
| GET | /events | List events (paginated) | ❌ |
| GET | /events/:id | Get single event | ❌ |
| POST | /events | Create event | ✅ |
| PUT | /events/:id | Update event | ✅ |
| DELETE | /events/:id | Delete event | ✅ |
| POST | /events/:id/register | Register for event | ✅ |
| DELETE | /events/:id/register | Cancel registration | ✅ |

---

## 📄 Pagination

Event listing supports:

```
GET /events?page=1&limit=10
```

Response includes:

- total records  
- current page  
- total pages  
- limit per page  

---

## 🛡 Validation Rules

- Email must be valid format
- Password: 6–72 characters
- Event name: 3–100 characters
- Event description: 10–500 characters
- Event location: 3–100 characters
- Event dateTime must be a future date

Validation errors return structured responses.

---

## 🧾 Error Format

Standard error:

```json
{
  "message": "description of what went wrong"
}
```

Validation error:

```json
{
  "message": "validation failed",
  "errors": [
    { "field": "name", "message": "must be at least 3 characters" }
  ]
}
```

---

## 🧪 Running Tests

```bash
go test ./...
```

---

## 🔮 Future Improvements

- PostgreSQL migration
- Refresh token implementation
- OAuth integration
- API rate limiting
- Docker containerization
- gRPC microservice split
- Redis caching layer

---

## 📌 Why This Project?

This project was built to deeply understand:

- Backend architecture in Go
- Middleware design
- Authentication flows
- Database modeling & ownership constraints
- Clean separation of concerns
- Writing maintainable and testable code

---
