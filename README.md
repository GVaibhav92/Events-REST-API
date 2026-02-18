## Go Events REST API

A production-style REST API built in Go for managing events and user registrations. Designed to demonstrate clean backend architecture, JWT authentication, validation, relational modeling, and secure middleware design.

---

## 📌 Overview

This project implements a layered backend architecture using Go and Gin, supporting:

- User authentication with JWT
- Secure password hashing (bcrypt)
- Event CRUD operations
- Ownership-based authorization
- Many-to-many event registrations
- Structured validation
- Pagination
- Middleware-driven request lifecycle
- Clean separation of concerns

---

## 🏗 Architecture
```
Client
   ↓
Gin Router
   ↓
Middleware (Logger / Auth)
   ↓
Route Handlers (Controllers)
   ↓
Models (Business Logic + DB Queries)
   ↓
Database (SQLite)
```

### Layers

| Directory | Responsibility |
|---|---|
| `config/` | Environment configuration |
| `db/` | Database connection & pooling |
| `models/` | Business logic & SQL queries |
| `routes/` | HTTP handlers |
| `middleware/` | Authentication & logging |
| `utils/` | JWT, hashing, validation |

This separation ensures maintainability, testability, and scalability.

---

## ✨ Features

- 🔐 JWT-based stateless authentication
- 🔑 Secure password hashing with bcrypt
- 🛡 Protected routes via custom middleware
- 👤 Ownership enforcement (only creators can modify events)
- 📋 Full CRUD for events
- 🔁 Many-to-many event registrations
- 📄 Pagination support
- 🧪 Structured validation (`go-playground/validator`)
- ⚙️ Environment-based configuration (`.env`)
- 🪵 Custom request logging middleware
- 🧾 Unit tests using Go's testing package

---

## 🛠 Tech Stack

- **Go**
- **Gin** — HTTP framework
- **SQLite** — `modernc.org/sqlite` (pure Go driver)
- **JWT** — `github.com/golang-jwt/jwt/v5`
- **bcrypt** — `golang.org/x/crypto/bcrypt`
- **go-playground/validator**
- **godotenv**

---

## 📂 Project Structure
```
go-events-api/
├── api-test/        # Tests
├── config/          # Environment configuration
├── db/              # Database initialization & pooling
├── middleware/      # Auth & logging middleware
├── models/          # Data models & queries
├── routes/          # HTTP handlers
├── utils/           # JWT, hashing, validation
├── go.mod
└── main.go
```

---

## 🚀 Getting Started

### 1. Clone
```bash
git clone https://github.com/GVaibhav92/Events-REST-API.git
cd Events-REST-API
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Create `.env`
```env
PORT=8080
DB_PATH=api.db
JWT_SECRET=your-secret-key
```

> ⚠️ Never commit `.env` to version control.

### 4. Run Server
```bash
go run main.go
```

Server runs at: `http://localhost:8080`

---

## 🔐 Authentication

**Login:**
```
POST /login
```

Request body:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Use the token in the request header:
```
Authorization: <token>
```

> If updated to Bearer format: `Authorization: Bearer <token>`

---

## 📌 API Endpoints

| Method | Route | Description | Protected |
|---|---|---|---|
| `POST` | `/signup` | Create user | ❌ |
| `POST` | `/login` | Login user | ❌ |
| `GET` | `/events` | List events | ❌ |
| `GET` | `/events/:id` | Get event | ❌ |
| `POST` | `/events` | Create event | ✅ |
| `PUT` | `/events/:id` | Update event | ✅ |
| `DELETE` | `/events/:id` | Delete event | ✅ |
| `POST` | `/events/:id/register` | Register for event | ✅ |
| `DELETE` | `/events/:id/register` | Cancel registration | ✅ |

---

## 📄 Pagination
```
GET /events?page=1&limit=10
```

Response includes: `data`, `total`, `page`, `limit`, `totalPages`

---

## 🛡 Security Features

- Passwords hashed with bcrypt
- JWT signed using HMAC SHA256
- Token expiration enforced
- Ownership validation on updates/deletes
- Foreign key constraints enabled
- Duplicate registration prevention

---

## 🧪 Testing
```bash
go test ./...
```

---

## 🔮 Future Improvements

- PostgreSQL migration
- Refresh tokens
- Role-based access control
- Rate limiting middleware
- Docker support
- Redis caching
- CI/CD integration
- API documentation with Swagger

---

## 🎯 Purpose

Built to deeply understand:

- Backend system architecture in Go
- Middleware lifecycle
- Stateless authentication
- Database modeling & relationships
- Secure API design
- Testable, modular backend structure

---

## 📜 License

MIT License
