package handler

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
	"v1/internal/infrastructure/logger"
)

type Server struct {
	app *fiber.App
	s3  *s3.Client
	l   *zap.SugaredLogger
}

func (s *Server) setupRoutes() {
	s.app.Post("/upload", s.handleUpload)
	s.app.Get("/download", s.handleDownload)
	s.app.Get("/list", s.handleList)
}

func New(s3 *s3.Client) *Server {
	l := logger.GetLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			l.Error("Unhandled error",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		l.Info("HTTP request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)))
		return err
	})

	server := &Server{
		app: app,
		s3:  s3,
		l:   logger.GetLogger(),
	}

	server.setupRoutes()

	return server
}
