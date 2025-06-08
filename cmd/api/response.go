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

// Response represents a standard response model
// swagger:model Response
type Response struct {
	// Status
	// example: 1 / -1
	Status int `json:"Status"`

	// Error message
	// example: Something went wrong
	Error string `json:"Error"`

	// Data actual response data
	// example: Success
	Data any `json:"Data"`
}

func (s *Server) writeResponse(ctx *fiber.Ctx, data any) error {
	if data == nil {
		return ctx.Status(fiber.StatusOK).JSON(Response{
			Status: 1,
			Data:   success,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Status: 1,
		Data:   data,
	})
}

func (s *Server) writeResponseWithStatusCode(ctx *fiber.Ctx, statusCode int, data any) error {
	if data == nil {
		return ctx.Status(statusCode).JSON(Response{
			Status: 1,
			Data:   success,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Status: 1,
		Data:   data,
	})
}

func (s *Server) writeErrorResponse(ctx *fiber.Ctx, statusCode int, error string) error {
	return ctx.Status(statusCode).JSON(Response{
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
	if err != nil {
		errorMessage = err.Error()
	}
	s.Logger.ErrorlnWithRequestId(ctx.UserContext(), errorMessage)
	if errorMessage == "" {
		errorMessage = errUnauthorized
	}
	return s.writeErrorResponse(ctx, fiber.StatusUnauthorized, errorMessage)
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
