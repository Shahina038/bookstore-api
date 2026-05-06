# Bookstore API

A RESTful Book Store API built with Go and PostgreSQL.

## Tech Stack

- Go 1.24.4
- PostgreSQL 16 (Docker)
- JWT authentication
- bcrypt password hashing
- database/sql + lib/pq (no ORM)

## Project Structure

```
bookstore/
├── cmd/api/
│   ├── main.go        # Entry point
│   └── routes.go      # URL routing
├── internal/
│   ├── config/        # Environment config
│   ├── db/            # Database connection
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # JWT middleware
│   ├── model/         # Structs
│   └── repository/   # SQL queries
├── db_schema.sql
└── .env
```

## API Endpoints

| Method | URL | Description | Auth |
|--------|-----|-------------|------|
| POST | `/register` | Create account | No |
| POST | `/login` | Login, get JWT | No |
| DELETE | `/account` | Soft delete account | Yes |
| GET | `/books` | List all books | No |
| GET | `/books/{id}` | Get a book | No |
| POST | `/books` | Add a book | Yes |
| PUT | `/books/{id}` | Update a book | Yes |
| DELETE | `/books/{id}` | Delete a book | Yes |
| GET | `/cart` | View cart | Yes |
| POST | `/cart/items` | Add to cart | Yes |
| PUT | `/cart/items/{id}` | Update quantity | Yes |
| DELETE | `/cart/items/{id}` | Remove from cart | Yes |
| POST | `/orders` | Checkout | Yes |
| GET | `/orders` | Order history | Yes |
| GET | `/orders/{id}` | Order details | Yes |

## Setup

### 1. Start PostgreSQL

```bash
docker run --name bookstore-db -e POSTGRES_USER=bookstore_user -e POSTGRES_PASSWORD=bookstore_pass -e POSTGRES_DB=bookstore -p 5433:5432 -d postgres:16
```

### 2. Run schema

```bash
docker exec -i bookstore-db psql -U bookstore_user -d bookstore < db_schema.sql
```

### 3. Create .env

```
DB_HOST=localhost
DB_PORT=5433
DB_USER=bookstore_user
DB_PASSWORD=bookstore_pass
DB_NAME=bookstore
JWT_SECRET=supersecretkey123
SERVER_PORT=8080
```

### 4. Run

```bash
go run cmd/api/main.go
```

## Design Decisions

- No ORM — raw SQL for full control
- Repository interfaces for unit testability
- Soft delete for users — data is never lost
- UUID primary keys on all tables
- DB transaction on checkout — all or nothing
