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
- [x] Fiber Router 

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

## How to Use
### Configuration
- `config/config.yaml` (ignored by git)
  - Default configuration file
- `cmd/root.go`
  - `config/config.yaml` is loaded by default
  - You can specify the configuration file with the `--config` flag
- `internal/app/<your-handler>/<xxx>.go`
  - Define your handler functions for your endpoint
- `internal/logger/zap_logger.go`
  - You can see the log settings in the `NewZapLogger` function
- `job/`
  - You can add your own jobs here
- `scheduler/scheduler.go`
  - You can schedule your jobs here
  - You can configure the cron expression in `config/config.yaml`

## License
Distributed under the MIT License. See `LICENSE` for more information.
