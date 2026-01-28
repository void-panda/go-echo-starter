package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-echo-starter/pkg/logger"
)

// RequestIDHeader is the header name for request ID
const RequestIDHeader = "X-Request-ID"

// Setup configures all middlewares for the Echo instance
func Setup(e *echo.Echo, log *logger.Logger) {
	// Recovery middleware
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error().
				Err(err).
				Str("stack", string(stack)).
				Msg("Panic recovered")
			return nil
		},
	}))

	// Request ID middleware
	e.Use(RequestID())

	// Logger middleware
	e.Use(Logger(log))

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, RequestIDHeader},
	}))

	// Secure middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://validator.swagger.io",
	}))
}

// RequestID middleware adds a unique request ID to each request
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Request().Header.Set(RequestIDHeader, requestID)
			c.Response().Header().Set(RequestIDHeader, requestID)
			return next(c)
		}
	}
}

// Logger middleware logs each HTTP request
func Logger(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			log.Info().
				Str("request_id", req.Header.Get(RequestIDHeader)).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("remote_ip", c.RealIP()).
				Int("status", res.Status).
				Dur("latency", time.Since(start)).
				Str("user_agent", req.UserAgent()).
				Msg("HTTP Request")

			return err
		}
	}
}
