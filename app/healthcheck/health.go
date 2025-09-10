package healthcheck

import "context"

type HealthCheckRequest struct {
}
type HealthCheckResponse struct {
	Message string `json:"message"`
}
type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Handle(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Message: "OK"}, nil
}
