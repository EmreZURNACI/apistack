package actor

import (
	"context"
)

type CreateActorRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateActorResponse struct {
	ID int64 `json:"id"`
}
type CreateActorHandler struct {
	repository Repository
}

func NewCreateActorHandler(repository Repository) *CreateActorHandler {
	return &CreateActorHandler{
		repository: repository,
	}
}

func (h *CreateActorHandler) Handle(ctx context.Context, req *CreateActorRequest) (*CreateActorResponse, error) {

	id, err := h.repository.CreateActor(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}
	return &CreateActorResponse{
		ID: id,
	}, nil
}
