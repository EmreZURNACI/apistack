package app

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type UpdateActorRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateActorResponse struct {
	ID string `json:"id"`
}

type UpdateActorHandler struct {
	Repository Repository
	Tracer     trace.Tracer
}

func NewUpdateActorHandler(repository Repository, tracer trace.Tracer) *UpdateActorHandler {
	return &UpdateActorHandler{
		Repository: repository,
		Tracer:     tracer,
	}
}

func (h *UpdateActorHandler) Handle(ctx context.Context, req *UpdateActorRequest) (*UpdateActorResponse, error) {

	ctx, span := h.Tracer.Start(ctx, "UpdateActor")
	defer span.End()

	if err := h.Repository.UpdateActor(ctx, req.ID, req.FirstName, req.LastName); err != nil {
		return nil, err
	}

	return &UpdateActorResponse{
		ID: req.ID,
	}, nil
}
