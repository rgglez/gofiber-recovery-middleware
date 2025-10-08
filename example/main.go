/*
Copyright 2025 Rodolfo González González

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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