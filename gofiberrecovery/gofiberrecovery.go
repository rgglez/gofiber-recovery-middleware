package gofiberrecovery

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

//-----------------------------------------------------------------------------

// PanicInfo contains detailed information about the panic
type PanicInfo struct {
	File     string
	Line     int
	Function string
	Stack    []StackFrame
}

// StackFrame represents a stack trace frame
type StackFrame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

//-----------------------------------------------------------------------------

type Config struct {
	// Next defines a function to skip this middleware when it returns true.
	Next func(c *fiber.Ctx) bool

	// Logger is the zerolog instance for logging panics
	Logger *zerolog.Logger

	// TimeKeyFunc function to generate the time_key in logs
	TimeKeyFunc func() string

	// SkipFrames number of frames to skip in the stack trace (default: 4)
	SkipFrames int

	// IncludeRuntime includes runtime frames in the stack trace
	IncludeRuntime bool

	// EnableStackTrace enables logging of the complete stack trace
	EnableStackTrace bool

	// CustomResponse allows customizing the HTTP response in case of panic
	CustomResponse func(c *fiber.Ctx, panicInfo *PanicInfo, e any) error

	// OnPanic additional callback that executes when a panic occurs
	OnPanic func(c *fiber.Ctx, panicInfo *PanicInfo, e any)
}

//-----------------------------------------------------------------------------

// ConfigDefault is the default configuration
var ConfigDefault = Config{
	Next:             nil,
	Logger:           nil,
	TimeKeyFunc:      nil,
	SkipFrames:       4,
	IncludeRuntime:   false,
	EnableStackTrace: true,
	CustomResponse:   nil,
	OnPanic:          nil,
}

//-----------------------------------------------------------------------------

// New creates a new panic recovery middleware
func New(config ...Config) fiber.Handler {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	// Required configuration validation
	if cfg.Logger == nil {
		panic("gofiberrecovery: Logger is required")
	}

	// Default values
	if cfg.SkipFrames == 0 {
		cfg.SkipFrames = 4
	}
	if cfg.TimeKeyFunc == nil {
		cfg.TimeKeyFunc = func() string { return "" }
	}

	return func(c *fiber.Ctx) error {
		// Execute Next if defined
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Defer to catch panics
		defer func() {
			if r := recover(); r != nil {
				// Extract detailed panic information
				panicInfo := extractPanicInfo(cfg.SkipFrames, cfg.IncludeRuntime)

				// Build log event
				logEvent := cfg.Logger.Error().
					Str("tag", "panic").
					Str("method", c.Method()).
					Str("url", c.OriginalURL()).
					Str("path", c.Path()).
					Int("status", fiber.StatusInternalServerError).
					Str("remote_addr", c.IP()).
					Str("user_agent", c.Get("User-Agent")).
					Interface("panic_value", r).
					Str("panic_file", panicInfo.File).
					Int("panic_line", panicInfo.Line).
					Str("panic_function", panicInfo.Function).
					Str("panic_location", fmt.Sprintf("%s:%d", panicInfo.File, panicInfo.Line))

				// Add time_key if configured
				if cfg.TimeKeyFunc != nil {
					logEvent = logEvent.Str("time_key", cfg.TimeKeyFunc())
				}

				// Add complete stack trace if enabled
				if cfg.EnableStackTrace {
					logEvent = logEvent.Interface("stack_trace", panicInfo.Stack)
				}

				// Log the panic
				logEvent.Msgf("PANIC RECOVERED: %s:%d in %s - %v",
					panicInfo.File,
					panicInfo.Line,
					panicInfo.Function,
					r)

				// Execute custom callback if it exists
				if cfg.OnPanic != nil {
					cfg.OnPanic(c, panicInfo, r)
				}

				// Send response to client
				if cfg.CustomResponse != nil {
					cfg.CustomResponse(c, panicInfo, r)
				} else {
					c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": "An unexpected error occurred",
					})
				}
			}
		}()

		// Continue to next handler
		return c.Next()
	}
}

//-----------------------------------------------------------------------------

// extractPanicInfo extracts detailed information about the panic
func extractPanicInfo(skipFrames int, includeRuntime bool) *PanicInfo {
	var pcs [32]uintptr
	n := runtime.Callers(skipFrames, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	info := &PanicInfo{
		Stack: make([]StackFrame, 0),
	}
	firstRelevantFrame := true

	for {
		frame, more := frames.Next()

		// Determine if the frame is relevant
		isRelevant := includeRuntime || (!strings.Contains(frame.File, "runtime/") &&
			!strings.Contains(frame.File, "gofiber/fiber") &&
			!strings.Contains(frame.File, "panic.go"))

		if isRelevant {
			stackFrame := StackFrame{
				File:     frame.File,
				Line:     frame.Line,
				Function: frame.Function,
			}
			info.Stack = append(info.Stack, stackFrame)

			// Capture the first relevant frame as the panic origin
			if firstRelevantFrame && !strings.Contains(frame.File, "runtime/") {
				info.File = frame.File
				info.Line = frame.Line
				info.Function = frame.Function
				firstRelevantFrame = false
			}
		}

		if !more {
			break
		}
	}

	return info
}

//-----------------------------------------------------------------------------

// GetPanicLocation returns a formatted string with the panic location
func GetPanicLocation(info *PanicInfo) string {
	if info == nil {
		return "unknown location"
	}
	return fmt.Sprintf("%s:%d (%s)", info.File, info.Line, info.Function)
}