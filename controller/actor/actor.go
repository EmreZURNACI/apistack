package actor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EmreZURNACI/apistack/app/actor"
	"github.com/EmreZURNACI/apistack/cache/redis"
	"github.com/EmreZURNACI/apistack/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("stackapi")

var validate = validator.New()

func (h *ActorController) GetActors(c *fiber.Ctx) error {

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

	key := fmt.Sprintf("actors:search=%s", i.Search)

	actors, err := h.cache.Get(ctx, key)

	if err == nil {
		var actorList []domain.Actor
		if err := json.Unmarshal(actors, &actorList); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"actors": actorList,
		})
	}

	ActorsHandler := actor.NewGetActorsHandler(h.db)
	res, err := ActorsHandler.Handle(ctx, &actor.GetActorsRequest{
		Search:  i.Search,
		Limit:   i.Limit,
		Offset:  i.Offset,
		OrderBy: i.OrderBy,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	bs, err := json.Marshal(res.Actors)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = h.cache.Set(ctx, redis.Message{
		Key:      []byte(key),
		Value:    bs,
		Duration: time.Minute * 3,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"actors": res.Actors,
	})
}
func (h *ActorController) GetActor(c *fiber.Ctx) error {
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

	ActorHandler := actor.NewGetActorHandler(h.db)
	res, err := ActorHandler.Handle(ctx, &actor.GetActorRequest{
		ActorID: i.ID,
	})

	if err != nil {
		zap.L().Error("Error getting actor", zap.Error(err))
		return c.JSON(err.Error())
	}

	return c.JSON(res)
}
func (h *ActorController) CreateActor(c *fiber.Ctx) error {
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

	CreateActorHandler := actor.NewCreateActorHandler(h.db)

	res, err := CreateActorHandler.Handle(ctx, &actor.CreateActorRequest{
		FirstName: i.FirstName,
		LastName:  i.LastName,
	})

	if err != nil {
		zap.L().Error("Error creating actor", zap.Error(err))
		return c.JSON(err.Error())

	}

	return c.JSON(res)
}
func (h *ActorController) UpdateActor(c *fiber.Ctx) error {
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

	GetActorHandler := actor.NewGetActorHandler(h.db)
	_, err := GetActorHandler.Handle(ctx, &actor.GetActorRequest{
		ActorID: i.ID,
	})

	if err != nil {
		zap.L().Error("Error getting actor", zap.Error(err))
		return c.JSON(err.Error())
	}

	UpdateActorHandler := actor.NewUpdateActorHandler(h.db)
	res, err := UpdateActorHandler.Handle(ctx, &actor.UpdateActorRequest{
		ID:        i.ID,
		FirstName: i.FirstName,
		LastName:  i.LastName,
	})

	if err != nil {
		zap.L().Error("Error updating actor", zap.Error(err))
		return c.JSON(err.Error())
	}

	return c.JSON(res)
}
func (h *ActorController) DeleteActor(c *fiber.Ctx) error {
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

	DeleteActorHandler := actor.NewDeleteActorHandler(h.db)
	res, err := DeleteActorHandler.Handle(ctx, &actor.DeleteActorRequest{
		ID: i.ID,
	})
	if err != nil {
		zap.L().Error("Error deleting actor", zap.Error(err))
		return c.JSON(err.Error())
	}

	return c.JSON(res)
}
