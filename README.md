# Hadia Parfums

[![CI](https://github.com/nurkenti/hadiaParfums/actions/workflows/ci.yml/badge.svg)](https://github.com/nurkenti/hadiaParfums/actions)

Telegram bot backend for a perfume store. Built with Go, PostgreSQL, sqlc, and Docker.

## Features

- **RBAC & FSM:** Separate user and admin flows with state handling.
- **Client Flow:** Product catalog, item view, order placement.
- **Admin Flow:** Catalog management (CRUD) and incoming order processing.
- **Clean UI:** Automatically deletes obsolete bot messages to keep chat history clean.
- **Database Init:** Auto-executes raw SQL migrations and indexes on container startup.

## Tech Stack

- **Language:** Go 1.22
- **Database Tools:** PostgreSQL, sqlc
- **Infra & Tooling:** Docker, Docker Compose, Makefile, GitHub Actions (`go vet` & build check)

## Project Structure

```text
.
├── cmd/
│   └── main.go          # Application entry point
├── internal/            # Core business logic, handlers, and generated sqlc code
├── migration/           # SQL schema migrations
├── Makefile             # Helper commands for build and execution
├── sqlc.yaml            # sqlc code generator configuration
└── docker-compose.yaml  # Environment orchestrator
```


## Local Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/nurkenti/hadiaParfums.git
   cd hadiaParfums
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   ```

3. **Run via Docker:**
   ```bash
   docker-compose up --build
   ```
