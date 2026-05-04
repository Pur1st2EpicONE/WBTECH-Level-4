## L2.18

![calendar banner](assets/banner.png)

<h3 align="center">A simple event calendar service written in Go, featuring REST API, structured logging, graceful shutdown, modular architecture, and Swagger documentation.</h3> <br>

<br>

## Architecture

The system is built with clean layered architecture:

* App layer — lifecycle management, initialization, shutdown.

* HTTP layer — Gin handlers, routing, middleware.

* Service layer — business logic and validation.

* Repository layer — in-memory storage with support for CRUD operations.

The server supports graceful shutdown, request logging, UUID-based event IDs, and Swagger UI for API exploration.

<br>

## Cool features

### Fully interface-driven architecture

All core components — server, handler, service, repository, logger — interact via interfaces, enabling easy mocking, dependency injection, and strict layer isolation.

### Versioned API (v1 and beyond)

REST API is fully versioned (/api/v1/...) with clear separation of versions and zero shared state, allowing backward-compatible evolution.

### Request logging middleware with latency, request IDs & full context

Captures latency, request IDs, client info, query strings, protocol, and Gin errors for full observability.

### Extensive multi-layer validation

Validation occurs at every stage, from JSON parsing and semantic checks in handlers to business rules in service and consistency enforcement in storage.

### High-performance in-memory storage

Efficient, size-controlled repository with hierarchical userID → date → events mapping, auxiliary lookup maps for O(1) access, preallocated maps, zero-copy updates, and thread safety via RWMutex.

### Production-ready codebase with 100% test coverage

Handler, service, and repository layers are fully tested, covering all parsing, validation, business rules, repository operations, error handling, and update/no-update scenarios.

<br>

## Installation
⚠️ Prerequisite: Go 1.25.1

Optionally, edit [the config file](config.yaml) to customize your preferences, then build and run the project using the Makefile command:

```bash
make
```
The server will start, run in foreground mode, and gracefully shut down upon receiving a SIGINT.

After that, you can use:

```bash
make clean
```

to remove the application executable and the logs directory.

<br>

## Testing & Linting

Run tests and ensure code quality:

```bash
make test        # Unit tests
make lint        # Linting checks
```

<br>

## Documentation

API documentation with examples is available via Swagger UI at:

```bash
http://localhost:<port>/swagger/index.html
```

Replace \<port> with the HTTP port configured in [config.yaml](config.yaml)

You can explore all endpoints, view request/response schemas, and try out the API directly from the browser.