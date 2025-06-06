package api

import (
	"github.com/gofiber/fiber/v2"
)

const (
	errBadRequest          = "Invalid request received"
	errInternalServerError = "Internal Server Error"
	errUnauthorized        = "Unauthorized"
	stringEmpty            = ""
	success                = "Success"
)

type response struct {
	Status int
	Error  string
	Data   any
}

func (s *Server) writeResponse(ctx *fiber.Ctx, data any) error {
	if data == nil {
		return ctx.Status(fiber.StatusOK).JSON(response{
			Status: 1,
			Data:   success,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response{
		Status: 1,
		Data:   data,
	})
}

func (s *Server) writeErrorResponse(ctx *fiber.Ctx, statusCode int, error string) error {
	return ctx.Status(statusCode).JSON(response{
		Status: -1,
		Error:  error,
	})
}

func (s *Server) internalServerError(ctx *fiber.Ctx, err error, errorMessage string) error {
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), err.Error())
	if errorMessage == "" {
		errorMessage = errInternalServerError
	}
	return s.writeErrorResponse(ctx, fiber.StatusInternalServerError, errorMessage)
}

func (s *Server) badRequest(ctx *fiber.Ctx, err error, errorMessage string) error {
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), err.Error())
	if errorMessage == "" {
		errorMessage = errBadRequest
	}
	return s.writeErrorResponse(ctx, fiber.StatusBadRequest, errBadRequest)
}

func (s *Server) unauthorized(ctx *fiber.Ctx, err error, errorMessage string) error {
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), err.Error())
	if errorMessage == "" {
		errorMessage = errUnauthorized
	}
	return s.writeErrorResponse(ctx, fiber.StatusUnauthorized, errUnauthorized)
}

func (s *Server) notFound(ctx *fiber.Ctx, errorMessage string) error {
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), errorMessage)
	return s.writeErrorResponse(ctx, fiber.StatusNotFound, errorMessage)
}

func (s *Server) conflict(ctx *fiber.Ctx, errorMessage string) error {
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), errorMessage)
	return s.writeErrorResponse(ctx, fiber.StatusConflict, errorMessage)
}

//func (s *Server) handleServiceError(c echo.Context, err error) error {
//	// Map service errors to appropriate HTTP status codes
//	switch err {
//	case services.ErrInvalidInput:
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
//	case services.ErrUnauthorized:
//		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
//	case services.ErrForbidden:
//		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
//	case services.ErrNotFound:
//		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
//	case services.ErrInternalServer:
//		fallthrough // Fallback for unknown errors
//	default:
//		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
//	}
//}
