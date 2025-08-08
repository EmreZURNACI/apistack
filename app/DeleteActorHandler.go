package app

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type DeleteActorRequest struct {
	ID string `json:"id"`
}
type DeleteActorResponse struct {
}

type DeleteActorHandler struct {
	Repository Repository
	Tracer     trace.Tracer
}

func NewDeleteActorHandler(repository Repository, trace trace.Tracer) *DeleteActorHandler {
	return &DeleteActorHandler{
		Repository: repository,
		Tracer:     trace,
	}
}

func (h *DeleteActorHandler) Handle(ctx context.Context, req *DeleteActorRequest) (*DeleteActorResponse, error) {

	ctx, span := h.Tracer.Start(ctx, "DeleteActorHandle")
	defer span.End()

	err := h.Repository.DeleteActor(ctx, req.ID)
	if err != nil {
		return &DeleteActorResponse{}, err
	}
	return &DeleteActorResponse{}, nil

}
