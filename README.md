# gofiber-recovery-middleware

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-recovery-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-recovery-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-recovery-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-recovery-middleware/gofiberip)](https://goreportcard.com/report/github.com/rgglez/gofiber-recovery-middleware/gofiberip)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-recovery-middleware.svg)](https://github.com/rgglez/gofiber-recovery-middleware/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/gofiber-recovery-middleware?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/gofiber-recovery-middleware?style=social)

**gofiber-recovery-middleware** is a [gofiber](https://gofiber.io/) [middleware](https://docs.gofiber.io/category/-middleware/) intended to recover and log panics.

## Installation

```bash
go get github.com/rgglez/gofiber-recovery-middleware
```

## Usage

```go
import gofiberrecovery "github.com/rgglez/gofiber-recovery-middleware/gofiberrecovery"

// Initialize Fiber app and middleware
app := fiber.New()
app.Use(gofiberrecovery.New(gofiberrecovery.Config{/* your configuration here */}))
```

## Configuration

There are some configuration options available in the ```Config``` struct:

* ``Next`` defines a function to skip this middleware when it returns true.
* ``Logger`` it is the zerolog instance to record panics.
* ``TimeKeyFunc`` function to generate the time_key in the logs.
* ``SkipFrames`` number of frames to skip in the stack trace (default: 4).
* ``EnableStackTrace`` enables full stack trace logging.
* ``CustomResponse`` allows you to customize the HTTP response in case of panic.
* ``OnPanic`` additional callback that executes when a panic occurs.

## Example

An example is included in the [example](example/) directory. To execute it:

1. Enter the example directory.
1. Run the example:
   ```bash
   go run .
   ```
1. In your browser open

   [http://127.0.0.1:3000/test-panic](http://127.0.0.1:3000/test-panic)


## Dependencies

* [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)
* [github.com/rs/zerolog](https://github.com/rs/zerolog)

## License

Copyright (c) 2025 Rodolfo González González.

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.