package actor

import (
	"github.com/EmreZURNACI/apistack/cache/redis"
	"github.com/EmreZURNACI/apistack/infra/postgresql"
)

type ActorController struct {
	cache *redis.Handler
	db    *postgresql.PostgresHandler
}

func NewActorController(db *postgresql.PostgresHandler, cache *redis.Handler) *ActorController {
	return &ActorController{
		cache: cache,
		db:    db,
	}
}
