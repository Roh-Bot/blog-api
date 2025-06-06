package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Health(c *fiber.Ctx) error {
	return c.Send(nil)
}
