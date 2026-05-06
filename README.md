# go-crud

[![Go Version](https://img.shields.io/badge/Go-1.26-blue)](https://go.dev/doc/devel/release#go1.26.2)
[![License](https://img.shields.io/badge/License-MIT-green)](#license)

A Contacts CRUD API built with Go using **Clean Architecture**, demonstrating layered separation of concerns, dependency injection, and repository pattern implementation.

---

## Project Overview

This project was built as a learning exercise to understand Clean Architecture in Go, inspired by [evrone/go-clean-template](https://github.com/evrone/go-clean-template). It implements a RESTful API for managing contact information with full CRUD operations.

---

## Features

- **RESTful API** for Contact management (Create, Read, Update, Delete)
- **Clean Architecture** with clear layer separation
- **Swappable Repository Implementations** — switch between in-memory and PostgreSQL storage without changing business logic
- **Middleware Support** — logging and panic recovery
- **Proper Error Handling** — structured error responses
- **Query Filtering** — filter contacts by name, email, phone, and relationship type with pagination

---

## Tech Stack

| Component     | Technology                                         |
| ------------- | -------------------------------------------------- |
| Language      | Go 1.26                                            |
| Web Framework | Gin                                                |
| Database      | PostgreSQL (with in-memory option for development) |
| Architecture  | Clean Architecture                                 |

---

## Project Structure

The project follows Clean Architecture principles with clear separation between layers:

```
go-crud/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── entity/
│   │   └── contact.go           # Domain entities (business models)
│   ├── usecase/
│   │   ├── contact/             # Business logic layer
│   │   └── contracts.go         # UseCase interfaces
│   ├── repo/
│   │   ├── contracts.go         # Repository interfaces (ports)
│   │   ├── inmem/              # In-memory implementation
│   │   └── persistent/        # PostgreSQL implementation
│   └── controller/
│       └── restapi/
│           ├── v1/             # REST API handlers
│           └── middleware/      # HTTP middleware
└── pkg/
    ├── logger/                  # Logging utilities
    ├── postgres/               # PostgreSQL connection
    └── httpserver/             # HTTP server utilities
```

### Layer Responsibilities

| Layer        | Responsibility                                                                 |
| ------------ | ------------------------------------------------------------------------------ |
| `entity`     | Domain models and business rules                                               |
| `usecase`    | Application business logic (orchestrating repository calls)                    |
| `repo`       | Data access abstraction — defines interfaces that implementations must satisfy |
| `controller` | HTTP request handling, validation, and response formatting                     |

The key principle: **each layer only knows about the layer immediately below it**. UseCases depend on repository interfaces (contracts), not concrete implementations. This enables swapping data sources without touching business logic.

---

## API Endpoints

| Method   | Endpoint               | Description                                           |
| -------- | ---------------------- | ----------------------------------------------------- |
| `GET`    | `/api/v1/contacts`     | List all contacts (supports filtering and pagination) |
| `GET`    | `/api/v1/contacts/:id` | Get a contact by ID                                   |
| `POST`   | `/api/v1/contacts`     | Create a new contact                                  |
| `PUT`    | `/api/v1/contacts/:id` | Update an existing contact                            |
| `DELETE` | `/api/v1/contacts/:id` | Delete a contact                                      |

### Query Parameters (List Endpoint)

| Parameter      | Type   | Description                                                       |
| -------------- | ------ | ----------------------------------------------------------------- |
| `first_name`   | string | Filter by first name (partial match)                              |
| `last_name`    | string | Filter by last name (partial match)                               |
| `email`        | string | Filter by email (partial match)                                   |
| `phone_number` | string | Filter by phone number (partial match)                            |
| `relationship` | string | Filter by relationship (`Friend`, `Family`, `Colleague`, `Other`) |
| `limit`        | uint64 | Number of results to return (default: 10)                         |
| `offset`       | uint64 | Number of results to skip (default: 0)                            |

---

## Quick Start

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- (Optional) PostgreSQL for persistent storage

### Running the Application

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-crud.git
cd go-crud
```

2. Run the application:

```bash
go run cmd/app/main.go
```

The server starts on `http://localhost:8080` by default.

### Running Tests

```bash
go test ./...
```

---

## Architecture Highlights

### Swappable Repository Pattern

One of the key features implemented in this project is the ability to swap between different repository implementations without modifying business logic.

**How it works:**

1. The `usecase` layer depends on repository **interfaces** (contracts), not concrete types
2. At application startup, the concrete repository is injected via dependency injection
3. Switching implementations only requires changing which repository is passed to the use case

```go
// internal/repo/contracts.go
type ContactRepository interface {
    Create(ctx context.Context, contact entity.Contact) (entity.Contact, error)
    GetByID(ctx context.Context, id string) (entity.Contact, error)
    List(ctx context.Context, filter ContactFilter) ([]entity.Contact, int64, error)
    Update(ctx context.Context, id string, update ContactUpdate) (entity.Contact, error)
    Delete(ctx context.Context, id string) error
}
```

Currently, two implementations are provided:

- **`inmem`**: In-memory storage using Go maps — useful for development and testing
- **`persistent`**: PostgreSQL implementation for production use

### What I Learned

Implementing the swappable repository pattern in this project taught me:

1. **Dependency Inversion Principle** — depending on abstractions rather than concrete implementations gives flexibility
2. **Interface Segregation** — defining minimal interfaces makes implementations simpler and more focused
3. **Testability** — with interfaces, mocking the repository for unit tests becomes straightforward
4. **Layer Separation** — clean architecture makes it easy to reason about where changes need to be made
5. **Configuration Management** — handling different environments (dev vs prod) requires thoughtful setup

This experience reinforced why Clean Architecture matters: changing the data layer doesn't require touching the business logic, making the codebase easier to maintain and extend.

---

## Configuration

### Environment Variables

| Variable      | Description                               | Default     |
| ------------- | ----------------------------------------- | ----------- |
| `APP_ENV`     | Application environment (`dev` or `prod`) | `dev`       |
| `PORT`        | Server port                               | `8080`      |
| `DB_HOST`     | PostgreSQL host                           | `localhost` |
| `DB_PORT`     | PostgreSQL port                           | `5432`      |
| `DB_USER`     | PostgreSQL username                       | `postgres`  |
| `DB_PASSWORD` | PostgreSQL password                       | `postgres`  |
| `DB_NAME`     | PostgreSQL database name                  | `go_crud`   |

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Acknowledgments

- [evrone/go-clean-template](https://github.com/evrone/go-clean-template) — inspiration for Clean Architecture implementation in Go
