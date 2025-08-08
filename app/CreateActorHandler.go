package app

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type CreateActorRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateActorResponse struct {
	ID int64 `json:"id"`
}
type CreateActorHandler struct {
	Repository Repository
	Tracer     trace.Tracer
}

func NewCreateActorHandler(repository Repository, trace trace.Tracer) *CreateActorHandler {
	return &CreateActorHandler{
		Repository: repository,
		Tracer:     trace,
	}
}

func (h *CreateActorHandler) Handle(ctx context.Context, req *CreateActorRequest) (*CreateActorResponse, error) {

	ctx, span := h.Tracer.Start(ctx, "CreateActorHandle")
	defer span.End()

	id, err := h.Repository.CreateActor(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}
	return &CreateActorResponse{
		ID: id,
	}, nil
}
