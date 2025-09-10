package actor

import (
	"context"

	"github.com/EmreZURNACI/apistack/app"
	"github.com/EmreZURNACI/apistack/domain"

	"go.opentelemetry.io/otel/trace"
)

type GetActorsRequest struct {
	Search  string `json:"search"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	OrderBy bool   `json:"order_by"`
}
type GetActorsResponse struct {
	Actors []domain.Actor `json:"actors"`
}

type GetActorsHandler struct {
	Repository app.Repository
	Tracer     trace.Tracer
}

func NewGetActorsHandler(repository app.Repository, tracer trace.Tracer) *GetActorsHandler {
	return &GetActorsHandler{
		Repository: repository,
		Tracer:     tracer,
	}
}

func (h *GetActorsHandler) Handle(ctx context.Context, req *GetActorsRequest) (*GetActorsResponse, error) {

	ctx, span := h.Tracer.Start(ctx, "GetActorsHandle")
	defer span.End()

	actors, err := h.Repository.GetActors(ctx, req.Search, req.Offset, req.Limit, req.OrderBy)
	if err != nil {
		return nil, err
	}

	return &GetActorsResponse{
		Actors: actors,
	}, nil
}
