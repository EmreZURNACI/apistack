package actor

import (
	"context"
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
	repository Repository
}

func NewUpdateActorHandler(repository Repository) *UpdateActorHandler {
	return &UpdateActorHandler{
		repository: repository,
	}
}

func (h *UpdateActorHandler) Handle(ctx context.Context, req *UpdateActorRequest) (*UpdateActorResponse, error) {

	if err := h.repository.UpdateActor(ctx, req.ID, req.FirstName, req.LastName); err != nil {
		return nil, err
	}

	return &UpdateActorResponse{
		ID: req.ID,
	}, nil
}
