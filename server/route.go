package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EmreZURNACI/apistack/cache/redis"
	"github.com/EmreZURNACI/apistack/controller/actor"
	"github.com/EmreZURNACI/apistack/controller/healthcheck"
	"github.com/EmreZURNACI/apistack/infra/postgresql"
	"github.com/spf13/viper"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("stackapi")

func Route() {

	server := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
		Concurrency:  1024 * 1024,
	})

	handler, err := postgresql.GetPostgresHandler(tracer)
	if err != nil {
		zap.L().Error("Error getting postgres handler", zap.Error(err))
		return
	}

	cacher, err := redis.Connection()
	if err != nil {
		zap.L().Error("Error getting postgres handler", zap.Error(err))
		return
	}

	actorController := actor.NewActorController(handler, cacher)
	healthcheckController := healthcheck.NewHealthCheckController()

	server.Use(otelfiber.Middleware())

	v1 := server.Group("/v1/actors")

	server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	server.Get("/healthcheck", healthcheckController.HealthCheck)
	v1.Get("/", actorController.GetActors)
	v1.Get("/:id", actorController.GetActor)
	v1.Post("/", actorController.CreateActor)
	v1.Put("/:id", actorController.UpdateActor)
	v1.Delete("/:id", actorController.DeleteActor)

	zap.L().Info("server started...", zap.Int("port", viper.GetInt("server.port")))
	if err := server.Listen(":" + viper.GetString("server.port")); err != nil {
		zap.L().Fatal("server stopped", zap.Error(err))
	}

	GracefulShutdown(server)
}

func GracefulShutdown(app *fiber.App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	zap.L().Sugar().Info("Shutting down server")

	if err := app.Shutdown(); err != nil {
		zap.L().Sugar().Error("Shutting down server. : %s", zap.Error(err))
	}
	zap.L().Sugar().Info("Server gracefully stopped")

}
