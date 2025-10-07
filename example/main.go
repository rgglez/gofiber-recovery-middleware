package main

import (
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"
    gofiberrecovery "github.com/rgglez/gofiber-recovery-middleware/gofiberrecovery"
)

func main() {
    app := fiber.New()

    // Global logger
    logger := zerolog.New(os.Stdout).With().
        Str("service", "api").
        Timestamp().
        Logger()

    // Recovery middleware
    app.Use(gofiberrecovery.New(gofiberrecovery.Config{
        Logger:           &logger,
        TimeKeyFunc:      generateTimekey,
        EnableStackTrace: true,
    }))

    // Test path
    app.Get("/test-panic", func(c *fiber.Ctx) error {
        panic("¡This is a test!")
    })

    app.Listen(":3000")
}

func generateTimekey() string {
    return "2025-10-07T10:30:00Z"
}