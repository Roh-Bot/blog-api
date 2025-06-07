package api

import (
	"github.com/gofiber/fiber/v2"
)

const (
	errInvalidUsername = "Invalid username"
)

type AuthLoginUser struct {
	Username string `json:"username" validate:"required"`
}

func (s *Server) authLoginUser(ctx *fiber.Ctx) error {
	// Parse request body and get user details from user service
	user := new(AuthLoginUser)
	if err := ctx.BodyParser(user); err != nil {
		return s.badRequest(ctx, err, err.Error())
	}
	if err := s.Validator.Struct(user); err != nil {
		return s.badRequest(ctx, err, validationToErrorMessage(err))
	}

	if isValid := s.Services.Auth.IsValid(user.Username); !isValid {
		return s.writeResponse(ctx, errInvalidUsername)
	}

	token, err := s.Services.Auth.GenerateToken(user.Username)
	if err != nil {
		return s.internalServerError(ctx, err, err.Error())
	}

	response := map[string]any{
		"access_token": token,
	}
	return s.writeResponse(ctx, response)
}
