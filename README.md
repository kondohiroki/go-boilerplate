# go-boilerplate :rocket:
[![Go Report Card](https://goreportcard.com/badge/github.com/kondohiroki/go-boilerplate)](https://goreportcard.com/report/github.com/kondohiroki/go-boilerplate)
![License](https://img.shields.io/github/license/kondohiroki/go-boilerplate)

This boilerplate is intended to be used as a starting point for a go application. It is not intended to be used as a but it is can be.

## Supported Features
- [x] Configuration with YAML
- [x] Logging with Zap Logger
- [x] CLI with Cobra
- [x] Scheduler with Cron
- [x] PostgreSQL
- [x] Docker
- [x] Discord

## Use Cases
- [x] As a Web Server
  - [x] HTTP API
  - [ ] gRPC API
  - [ ] GraphQL API
- [x] As a CLI Application
- [x] As a Scheduler for Cron Jobs

## Getting Started
### Prerequisites
-  Go 1.20
-  Docker

### Installation
1. Clone the repo
   ```sh
   git clone https://github.com/kondohiroki/go-boilerplate.git
    ```
2. Install Go dependencies
    ```sh
    go mod download
    ```
3. Install Docker dependencies
    ```sh
    docker compose build -pull
    ```
4. Run the application
    ```sh
    docker compose up
    ```
## Project Structure
```
.
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   ├── notify.go
│   ├── root.go
│   └── schedule.go
├── config
│   ├── config.example.yaml
│   ├── config.go
│   └── config.yaml
├── deploy
├── doc
│   └── images
│       ├── diagram.png
│       └── homepage.png
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── context.go
│   ├── discord
│   │   └── discord.go
│   ├── jobs
│   │   ├── discord.go
│   │   └── get_sku.go
│   ├── logger
│   │   └── zap_logger.go
│   ├── pgsql
│   │   ├── migrations
│   │   └── postgresql.go
│   └── scheduler
│       └── scheduler.go
├── main.go
└── pkg
    ├── color
    │   └── color_code.go
    ├── request
    │   ├── grpc.go
    │   └── http.go
    └── versioning
        └── versioning.go
```

## How to Use
### Configuration
- `config/config.yaml` (ignored by git)
  - Default configuration file
- `cmd/root.go`
  - `config/config.yaml` is loaded by default
  - You can specify the configuration file with the `--config` flag
- `internal/app/context.go`
  - You can access appContext from anywhere in the application
- `internal/logger/zap_logger.go`
  - You can see the log settings in the `NewZapLogger` function
- `jobs/`
  - You can add your own jobs here
- `scheduler/scheduler.go`
  - You can schedule your jobs here
  - You can configure the cron expression in `config/config.yaml`

## License
Distributed under the MIT License. See `LICENSE` for more information.
