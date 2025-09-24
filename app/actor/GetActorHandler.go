package actor

import (
	"context"

	"github.com/EmreZURNACI/apistack/domain"
)

type GetActorRequest struct {
	ActorID string `json:"actor_id"`
}

type GetActorResponse struct {
	Actor domain.Actor `json:"actor"`
}
type GetActorHandler struct {
	repository Repository
}

func NewGetActorHandler(repository Repository) *GetActorHandler {
	return &GetActorHandler{
		repository: repository,
	}
}

func (h *GetActorHandler) Handle(ctx context.Context, req *GetActorRequest) (*GetActorResponse, error) {

	actor, err := h.repository.GetActor(ctx, req.ActorID)
	if err != nil {
		return nil, err
	}
	return &GetActorResponse{
		Actor: *actor,
	}, nil
}
