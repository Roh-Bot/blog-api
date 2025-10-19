package api

import (
	"github.com/labstack/echo/v4"
)

const (
	errInvalidUsername = "Invalid username"
)

// AuthLoginUser represents the login request payload
// swagger:model AuthLoginUser
type AuthLoginUserRequest struct {
	// Username for authentication
	// required: true
	// example: admin
	Username string `json:"Username" validate:"required"`
}

// AuthLoginResponse represents the successful login response
// swagger:model AuthLoginResponse
type AuthLoginUserResponse struct {
	// JWT access token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"AccessToken"`
}

// AuthLoginWrappedResponse wraps AuthLoginUserResponse inside standard Response
// swagger:model AuthLoginWrappedResponse
type AuthLoginWrappedResponse struct {
	// example: 1
	Status int `json:"Status"`

	// example: ""
	Error string `json:"Error"`

	// example: {"AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
	Data AuthLoginUserResponse `json:"Data"`
}

// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body AuthLoginUserRequest true "Login credentials"
// @Success 200 {object} AuthLoginWrappedResponse "Authentication successful"
// @Failure 400 {object} Response "Invalid request format"
// @Failure 401 {object} Response "Invalid credentials"
// @Failure 500 {object} Response "Internal server error"
// @Router /authentication/login [post]
func (s *Server) authLoginUser(ctx echo.Context) error {
	// Parse request body and get user details from user service
	user := new(AuthLoginUserRequest)
	if err := ctx.Bind(user); err != nil {
		return s.badRequest(ctx, err, err.Error())
	}
	if err := s.Validator.Struct(user); err != nil {
		return s.badRequest(ctx, err, validationToErrorMessage(err))
	}

	if isValid := s.App.Auth.IsValid(user.Username); !isValid {
		return s.unauthorized(ctx, nil, errInvalidUsername)
	}

	token, err := s.App.Auth.GenerateToken(user.Username)
	if err != nil {
		return s.internalServerError(ctx, err, err.Error())
	}

	return s.writeResponse(ctx, AuthLoginUserResponse{AccessToken: token})
}
