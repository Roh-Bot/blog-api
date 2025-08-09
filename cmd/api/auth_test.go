package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/internal/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	authUrl = "/api/authentication/login"
)

type MockAuthService struct {
	ShouldValidate bool
	Token          string
	TokenErr       error
}

func (m *MockAuthService) IsValid(username string) bool {
	return m.ShouldValidate
}

func (m *MockAuthService) GenerateToken(username string) (string, error) {
	return m.Token, m.TokenErr
}

func (m *MockAuthService) ValidateToken(token string) (bool, error) {
	return true, nil
}

func setupAuthTestServer(auth services.IAuth) *fiber.App {
	mockService := services.Service{
		Auth: auth,
	}

	validatorV10 := validator.NewValidator()

	server := &Server{
		Services:  mockService,
		Validator: validatorV10,
		Logger:    &MockLogger{},
		Router:    fiber.New(),
	}

	server.Router.Post(authUrl, server.authLoginUser)
	return server.Router
}

func TestAuthLoginUser_Success(t *testing.T) {
	mockAuth := &MockAuthService{
		ShouldValidate: true,
		Token:          "mock-token",
	}

	app := setupAuthTestServer(mockAuth)

	body := AuthLoginUserRequest{Username: "devadiga.rohit"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, authUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthLoginUser_InvalidUsername(t *testing.T) {
	mockAuth := &MockAuthService{
		ShouldValidate: false,
	}

	app := setupAuthTestServer(mockAuth)

	body := AuthLoginUserRequest{Username: "invalid"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, authUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthLoginUser_GenerateTokenError(t *testing.T) {
	mockAuth := &MockAuthService{
		ShouldValidate: false,
		TokenErr:       errors.New("token generation failed"),
	}

	app := setupAuthTestServer(mockAuth)

	body := AuthLoginUserRequest{Username: "invalid"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, authUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthLoginUser_InvalidPayload(t *testing.T) {
	mockAuth := &MockAuthService{}
	app := setupAuthTestServer(mockAuth)

	req := httptest.NewRequest(http.MethodPost, authUrl, bytes.NewReader([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthLoginUser_ValidationFailure(t *testing.T) {
	mockAuth := &MockAuthService{}
	app := setupAuthTestServer(mockAuth)

	body := map[string]string{
		"Username": "",
	}
	payload, err := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, authUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
