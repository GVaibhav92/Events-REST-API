## Go Events REST API

A production-style REST API built in Go for managing events and user registrations. Designed to demonstrate clean backend architecture, JWT authentication, validation, relational modeling, and secure middleware design.

---

## 📌 Overview

This project implements a layered backend architecture using Go and Gin, supporting:

- User authentication with JWT access tokens and refresh token rotation
- Secure password hashing (bcrypt)
- Role-based access control (RBAC)
- Event CRUD operations with ownership enforcement
- Many-to-many event registrations with duplicate prevention
- Structured request validation
- Pagination
- Middleware-driven request lifecycle with per-request timeout
- Context cancellation propagated through the full request stack
- Graceful server shutdown with active-request draining
- Clean separation of concerns

---

## 🏗 Architecture
```
Client
   ↓
Gin Router
   ↓
Middleware (RequestID / Timeout / Logger / Auth / RBAC)
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
| `middleware/` | Authentication, RBAC, timeout & logging |
| `utils/` | JWT, hashing, validation |

This separation ensures maintainability, testability, and scalability.

---

## ✨ Features

- 🔐 JWT-based stateless authentication (access token + refresh token rotation)
- 🔑 Secure password hashing with bcrypt
- 🛡 Protected routes via custom middleware stack (RequestID → Timeout → Logger → Auth)
- 🎭 Role-based access control — admin and user roles enforced at middleware and handler level
- 👤 Ownership enforcement — only the event creator or an admin can update or delete
- 📋 Full CRUD for events
- 🔁 Many-to-many event registrations with DB-level duplicate prevention (`UNIQUE` constraint)
- 📄 Pagination with `page` and `limit` query params; response includes `total` and `totalPages`
- 🧪 Structured request validation (`go-playground/validator`) with custom `future_date` rule
- ⏱ Per-request timeout middleware with configurable duration (default 30s)
- 🔗 Full context cancellation propagation — request context flows from handler → model → DB
- 🛑 Graceful shutdown — drains active requests (10s window), then closes DB connection
- ⚙️ Environment-based configuration via `.env`
- 🪵 Custom request logging middleware with colored output and per-request ID tracing
- 🧾 Unit tests using Go's standard `testing` package

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
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a3f8c2d1e4b7..."
}
```

Use the access token in the request header:
```
Authorization: <access_token>
```

**Refresh the access token:**
```
POST /auth/refresh
{ "refresh_token": "<refresh_token>" }
```

**Logout (invalidates refresh token):**
```
POST /auth/logout
{ "refresh_token": "<refresh_token>" }
```

---

## 📌 API Endpoints

| Method | Route | Description | Protected |
|---|---|---|---|
| `POST` | `/signup` | Create user | ❌ |
| `POST` | `/login` | Login — returns access + refresh token | ❌ |
| `POST` | `/auth/refresh` | Rotate refresh token | ❌ |
| `POST` | `/auth/logout` | Invalidate refresh token | ❌ |
| `GET` | `/events` | List events (paginated) | ❌ |
| `GET` | `/events/:id` | Get event by ID | ❌ |
| `POST` | `/events` | Create event | ✅ |
| `PUT` | `/events/:id` | Update event (owner or admin) | ✅ |
| `DELETE` | `/events/:id` | Delete event (owner or admin) | ✅ |
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
- JWT access tokens signed with HMAC SHA256; expiry enforced on every request
- Refresh token rotation — old token invalidated on every refresh
- Role-based access control enforced at middleware and handler level
- Ownership validation on event updates and deletes
- Foreign key constraints enabled in SQLite
- Duplicate registration prevention enforced at the database level (`UNIQUE` constraint)
- Per-request context cancellation prevents hanging DB operations on client disconnect

---

## 🧪 Testing
```bash
go test ./...
```

---

## 🔮 Future Improvements

- PostgreSQL migration
- Rate limiting middleware
- Docker support
- Redis caching
- CI/CD integration
- API documentation with Swagger

---

## 🎯 Purpose

Built to deeply understand:

- Backend system architecture in Go
- Middleware lifecycle and context propagation
- Stateless authentication with token rotation
- Role-based authorization patterns
- Database modeling, constraints, and relationships
- Secure and production-correct API design
- Graceful shutdown and resource cleanup
- Testable, modular backend structure

---

## 📜 License

MIT License
