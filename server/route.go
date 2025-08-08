package server

import (
	"github.com/EmreZURNACI/apistack.git/app"
	"github.com/EmreZURNACI/apistack.git/infra/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var tracer = otel.Tracer("stackapi")

func Route() {

	server := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
		Concurrency:  1024 * 1024,
	})

	validate := validator.New()

	server.Use(otelfiber.Middleware())

	handler, err := postgresql.GetPostgresHandler(tracer)
	if err != nil {
		zap.L().Error("Error getting postgres handler", zap.Error(err))
		return
	}

	server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	server.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	server.Get("/v1/actors", func(c *fiber.Ctx) error {

		type input struct {
			Search  string `json:"search"`
			Limit   int    `json:"limit"`
			Offset  int    `json:"offset"`
			OrderBy bool   `json:"order_by"`
		}

		var i input

		if err := c.QueryParser(&i); err != nil {
			zap.L().Error("Error parsing query", zap.Error(err))
			return c.JSON(err.Error())
		}

		ctx, span := tracer.Start(c.UserContext(), "Actors")
		defer span.End()

		ActorsHandler := app.NewGetActorsHandler(handler, tracer)
		res, err := ActorsHandler.Handle(ctx, &app.GetActorsRequest{
			Search:  i.Search,
			Limit:   i.Limit,
			Offset:  i.Offset,
			OrderBy: i.OrderBy,
		})

		if err != nil {
			zap.L().Error("Error getting actors", zap.Error(err))
			return c.JSON(err.Error())
		}

		return c.JSON(res.Actors)
	})

	server.Get("/v1/actor/:id", func(c *fiber.Ctx) error {

		var id = c.Params("id")

		type input struct {
			ID string `json:"id" validate:"required,numeric"`
		}

		i := input{ID: id}
		if err := validate.Struct(&i); err != nil {
			zap.L().Error("Error getting actor id", zap.Error(err))
			return c.JSON(err.Error())
		}

		ctx, span := tracer.Start(c.UserContext(), "Actor")
		defer span.End()

		ActorHandler := app.NewGetActorHandler(handler, tracer)
		res, err := ActorHandler.Handle(ctx, &app.GetActorRequest{
			ActorID: i.ID,
		})

		if err != nil {
			zap.L().Error("Error getting actor", zap.Error(err))
			return c.JSON(err.Error())
		}

		return c.JSON(res)
	})

	server.Post("/v1/actors", func(c *fiber.Ctx) error {

		type input struct {
			FirstName string `json:"FirstName" validate:"required"`
			LastName  string `json:"LastName" validate:"required"`
		}

		var i input
		if err := c.BodyParser(&i); err != nil {
			zap.L().Error("Error parsing actor", zap.Error(err))
			return c.JSON(err.Error())

		}

		if err := validate.Struct(&i); err != nil {
			zap.L().Error("Error validating", zap.Error(err))
			return c.JSON(err.Error())
		}

		ctx, span := tracer.Start(c.UserContext(), "CreateActor")
		defer span.End()

		CreateActorHandler := app.NewCreateActorHandler(handler, tracer)

		res, err := CreateActorHandler.Handle(ctx, &app.CreateActorRequest{
			FirstName: i.FirstName,
			LastName:  i.LastName,
		})

		if err != nil {
			zap.L().Error("Error creating actor", zap.Error(err))
			return c.JSON(err.Error())

		}

		return c.JSON(res)
	})

	server.Put("/v1/actor/:id", func(c *fiber.Ctx) error {

		var id string = c.Params("id")

		type input struct {
			ID        string `json:"id" validate:"required,numeric"`
			FirstName string `json:"FirstName" validate:"required"`
			LastName  string `json:"LastName" validate:"required"`
		}

		i := input{ID: id}
		if err := c.BodyParser(&i); err != nil {
			zap.L().Error("Error parsing actor", zap.Error(err))
			return c.JSON(err.Error())
		}

		if err := validate.Struct(&i); err != nil {
			zap.L().Error("Error validating", zap.Error(err))
			return c.JSON(err.Error())
		}

		ctx, span := tracer.Start(c.UserContext(), "UpdateActor")
		defer span.End()

		GetActorHandler := app.NewGetActorHandler(handler, tracer)
		_, err := GetActorHandler.Handle(ctx, &app.GetActorRequest{
			ActorID: i.ID,
		})

		if err != nil {
			zap.L().Error("Error getting actor", zap.Error(err))
			return c.JSON(err.Error())
		}

		UpdateActorHandler := app.NewUpdateActorHandler(handler, tracer)
		res, err := UpdateActorHandler.Handle(ctx, &app.UpdateActorRequest{
			ID:        i.ID,
			FirstName: i.FirstName,
			LastName:  i.LastName,
		})

		if err != nil {
			zap.L().Error("Error updating actor", zap.Error(err))
			return c.JSON(err.Error())
		}

		return c.JSON(res)

	})

	server.Delete("/v1/actor/:id", func(c *fiber.Ctx) error {

		var id = c.Params("id")

		type input struct {
			ID string `json:"id" validate:"required,numeric"`
		}

		i := input{ID: id}
		if err := validate.Struct(&i); err != nil {
			zap.L().Error("Error getting actor id", zap.Error(err))
			return c.JSON(err.Error())
		}

		ctx, span := tracer.Start(c.UserContext(), "DeleteActor")
		defer span.End()

		DeleteActorHandler := app.NewDeleteActorHandler(handler, tracer)
		res, err := DeleteActorHandler.Handle(ctx, &app.DeleteActorRequest{
			ID: i.ID,
		})
		if err != nil {
			zap.L().Error("Error deleting actor", zap.Error(err))
			return c.JSON(err.Error())
		}

		return c.JSON(res)
	})

	zap.L().Info("server started...", zap.String("port", os.Getenv("SERVER_PORT")))
	if err := server.Listen(os.Getenv("SERVER_PORT")); err != nil {
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
