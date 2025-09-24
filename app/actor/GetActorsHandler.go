package actor

import (
	"context"

	"github.com/EmreZURNACI/apistack/domain"
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
	repository Repository
}

func NewGetActorsHandler(repository Repository) *GetActorsHandler {
	return &GetActorsHandler{
		repository: repository,
	}
}

func (h *GetActorsHandler) Handle(ctx context.Context, req *GetActorsRequest) (*GetActorsResponse, error) {

	actors, err := h.repository.GetActors(ctx, req.Search, req.Offset, req.Limit, req.OrderBy)
	if err != nil {
		return nil, err
	}

	return &GetActorsResponse{
		Actors: actors,
	}, nil
}
