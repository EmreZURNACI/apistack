package app

import (
	"context"
	"github.com/EmreZURNACI/apistack.git/domain"
	"go.opentelemetry.io/otel/trace"
)

type GetActorRequest struct {
	ActorID string `json:"actor_id"`
}

type GetActorResponse struct {
	Actor domain.Actor `json:"actor"`
}
type GetActorHandler struct {
	Repository Repository
	Tracer     trace.Tracer
}

func NewGetActorHandler(repository Repository, tracer trace.Tracer) *GetActorHandler {
	return &GetActorHandler{
		Repository: repository,
		Tracer:     tracer,
	}
}

func (h *GetActorHandler) Handle(ctx context.Context, req *GetActorRequest) (*GetActorResponse, error) {

	ctx, span := h.Tracer.Start(ctx, "GetActorHandle")
	defer span.End()

	actor, err := h.Repository.GetActor(ctx, req.ActorID)
	if err != nil {
		return nil, err
	}
	return &GetActorResponse{
		Actor: *actor,
	}, nil
}
