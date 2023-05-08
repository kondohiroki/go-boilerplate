# go-boilerplate :rocket:
[![Go Report Card](https://goreportcard.com/badge/github.com/kondohiroki/go-boilerplate)](https://goreportcard.com/report/github.com/kondohiroki/go-boilerplate)
[![Go with Confidence](https://github.com/kondohiroki/go-boilerplate/actions/workflows/go_with_confidence.yml/badge.svg)](https://github.com/kondohiroki/go-boilerplate/actions/workflows/go_with_confidence.yml)

This boilerplate is intended to be used as a starting point for a go application. It is not intended to be used as a but it is can be.

<p align="center">
<img src="https://user-images.githubusercontent.com/49369000/236752939-05e510db-a5ae-42ad-b1aa-da1c0222418b.png"  width="600" />
</p>

## Getting Started
### Prerequisites
-  Go 1.20
-  Docker
-  sonar-scanner - for coverage test in local
   ```sh
   brew install sonar-scanner
   ```

### Installation
1. Clone the repo
   ```sh
   git clone https://github.com/kondohiroki/go-boilerplate.git
    ```
2. Install Go dependencies
    ```sh
    go mod download
    ```
3. Copy the default configuration file
    ```sh
    cp config/config.example.yaml config/config.yaml
    ```
4. Start the database
    ```sh
    docker compose up -d
    ```
5. Migrate Database
    ```sh
    go run main.go migrate
    ```
6. Run the application
    ```sh
    # Run normally
    go run main.go serve-api

    # Run with hot reload
    air serve-api
    ```
7. Testing (optional)
    ```sh
    # Run unit-test
    make unit-test

    # Run api-test
    make api-test

    # Create sonar scret
    touch .sonar.secret
    echo "your-sonar-token" > .sonar.secret

    # Add secret to .sonar.secret
    # Get from sonar web
    ```
 
 ## Standard and Styles Guide

 ### Coding Standard

 1. For those `const`, use capitalized SNAKE_CASE for public constant. For private, constant name should led by _ (underscore).

    **Good Example**

    ```go
    // public
    const BAD_REQUEST int = 400

    // private
    const _UNAUTHORIZED int = 401
    ```

    **Bad Example**

    ```go
    const BadRequest   int = 400
    const unauthorized int = 401
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


## Supported Features
- [x] Configuration with YAML
- [x] Logging with Zap Logger
- [x] CLI with Cobra
- [x] Scheduler with Cron
- [x] PostgreSQL
- [x] Redis Cache
- [x] Docker
- [x] Fiber Router 
- [x] Add Redis Reliable Queue

## Use Cases
- [x] As a Web Server
  - [x] HTTP API
  - [ ] gRPC API
- [x] As a CLI Application
- [x] As a Scheduler for Cron Jobs

## Roadmap
- [ ] Add gRPC API
- [ ] Document the code

