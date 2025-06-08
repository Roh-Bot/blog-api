package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	requestId = "request_id"
)

func (s *Server) validateAuth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if len(authHeader) < len("Bearer ") || authHeader[:7] != "Bearer " {
		return s.unauthorized(ctx, nil, stringEmpty)
	}
	token := authHeader[7:]
	isValid, err := s.Services.Auth.ValidateToken(token)
	if err != nil || !isValid {
		s.Logger.ErrorlnWithRequestId(ctx.UserContext(), err)
		if errors.Is(err, auth.ErrTokenExpired) {
			return s.unauthorized(ctx, err, auth.ErrTokenExpired.Error())
		}
		return s.unauthorized(ctx, err, stringEmpty)
	}

	return ctx.Next()
}

func (s *Server) requestLogger(ctx *fiber.Ctx) error {
	//Log the request
	ctxUser := context.WithValue(context.Background(), requestId, uuid.New().String())
	ctx.SetUserContext(ctxUser)

	headers := ctx.GetReqHeaders()
	body := string(ctx.Request().Body())

	requestLog := map[string]any{
		"Headers": headers,
		"Body":    body,
	}
	requestLogJson, err := json.Marshal(requestLog)
	if err != nil {
		return s.internalServerError(ctx, err, stringEmpty)
	}
	s.Logger.InfolnWithRequestId(ctx.UserContext(), string(requestLogJson))
	return ctx.Next()
}

func (s *Server) responseLogger(ctx *fiber.Ctx) error {
	err := ctx.Next()
	if err != nil {
		s.Logger.ErrorlnWithRequestId(ctx.UserContext(), err.Error())
	}

	body := ctx.Response().Body()
	// Log the response
	s.Logger.InfolnWithRequestId(ctx.UserContext(), string(body))

	return err
}
