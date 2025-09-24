package actor

import (
	"context"
)

type DeleteActorRequest struct {
	ID string `json:"id"`
}
type DeleteActorResponse struct {
	Message string `json:"message"`
}

type DeleteActorHandler struct {
	repository Repository
}

func NewDeleteActorHandler(repository Repository) *DeleteActorHandler {
	return &DeleteActorHandler{
		repository: repository,
	}
}

func (h *DeleteActorHandler) Handle(ctx context.Context, req *DeleteActorRequest) (*DeleteActorResponse, error) {

	err := h.repository.DeleteActor(ctx, req.ID)
	if err != nil {
		return &DeleteActorResponse{}, err
	}
	return &DeleteActorResponse{
		Message: "Aktör silindi",
	}, nil

}
