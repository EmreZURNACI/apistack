package healthcheck

import (
	"github.com/EmreZURNACI/apistack/app/healthcheck"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (h *HealthCheckController) HealthCheck(c *fiber.Ctx) error {
	healthCheckhandler := healthcheck.NewHealthCheckHandler()
	res, err := healthCheckhandler.Handle(c.UserContext(), &healthcheck.HealthCheckRequest{})
	if err != nil {
		return c.JSON(err.Error())
	}
	return c.JSON(res.Message)
}
