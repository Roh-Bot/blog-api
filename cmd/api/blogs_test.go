package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	blogUrl       = "/api/blog-post"
	blogUrlWithId = "/api/blog-post/:id"
)

type MockBlogsService struct {
	AddPostError    error
	GetPostsError   error
	GetPostsData    []store.Post
	UpdatePostError error
	DeletePostError error
}

func (m *MockBlogsService) AddPost(ctx context.Context, addPost *services.AddPostDto) error {
	return m.AddPostError
}

func (m *MockBlogsService) GetPosts(ctx context.Context, getPosts *services.GetPostsDto) ([]store.Post, error) {
	return m.GetPostsData, m.GetPostsError
}

func (m *MockBlogsService) UpdatePost(ctx context.Context, addPost *services.UpdatePostDto) error {
	return m.UpdatePostError
}

func (m *MockBlogsService) DeletePost(ctx context.Context, addPost *services.DeletePostDto) error {
	return m.DeletePostError
}

func setupBlogsTestServer(blogs services.IBlog) *fiber.App {
	mockService := &services.Service{
		Blog: blogs,
	}

	validatorV10 := validator.NewValidator()

	server := &Server{
		Router:    fiber.New(),
		Logger:    &MockLogger{},
		Validator: validatorV10,
		Services:  mockService,
	}
	server.Router.Get(blogUrl, server.postsGet)
	server.Router.Get(blogUrlWithId, server.postsGet)
	server.Router.Post(blogUrl, server.postAdd)
	server.Router.Patch(blogUrlWithId, server.postUpdate)
	server.Router.Delete(blogUrlWithId, server.postDelete)
	return server.Router
}

func TestPostAdd_Success(t *testing.T) {
	mockService := &MockBlogsService{
		AddPostError: nil,
	}

	app := setupBlogsTestServer(mockService)

	sampleText := "sample"
	sampleBool := true
	body := AddPostRequest{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, blogUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestPostAdd_InvalidRequest(t *testing.T) {
	mockService := &MockBlogsService{
		AddPostError: nil,
	}

	app := setupBlogsTestServer(mockService)
	payload, _ := json.Marshal(`body`)

	req := httptest.NewRequest(http.MethodPost, blogUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostAdd_ValidationError(t *testing.T) {
	mockService := &MockBlogsService{
		AddPostError: nil,
	}

	app := setupBlogsTestServer(mockService)
	sampleText := "sample"
	sampleBool := true
	body := AddPostRequest{
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, blogUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostAdd_SlugAlreadyExists(t *testing.T) {
	mockService := &MockBlogsService{
		AddPostError: store.ErrSlugAlreadyExists,
	}

	app := setupBlogsTestServer(mockService)
	sampleText := "sample"
	sampleBool := true
	body := AddPostRequest{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, blogUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestPostAdd_InternalServerError(t *testing.T) {
	mockService := &MockBlogsService{
		AddPostError: errors.New("internal server error"),
	}

	app := setupBlogsTestServer(mockService)
	sampleText := "sample"
	sampleBool := true
	body := AddPostRequest{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, blogUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPostGet_WithPostIdSuccess(t *testing.T) {
	mockService := &MockBlogsService{
		GetPostsData: []store.Post{{Title: "Gaming"}},
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, err := app.Test(req)

	if !assert.NoError(t, err) {
		return
	}

	getPostsResponse := &GetPostsWrapperResponse{}

	err = json.NewDecoder(resp.Body).Decode(getPostsResponse)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, getPostsResponse.Data)
}

func TestPostGet_WithoutPostIdSuccess(t *testing.T) {
	mockService := &MockBlogsService{
		GetPostsData: []store.Post{{Title: "Gaming"}},
	}

	app := setupBlogsTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, blogUrl, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestPostGet_BadParamsRequest(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := "s"

	url := fmt.Sprintf("%s/%s", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostGet_PostDoesNotExistsError(t *testing.T) {
	mockService := &MockBlogsService{
		GetPostsError: store.ErrPostDoesNotExist,
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPostGet_WithoutIdInternalServerError(t *testing.T) {
	mockService := &MockBlogsService{
		GetPostsError: errors.New("internal server error"),
	}

	app := setupBlogsTestServer(mockService)

	req := httptest.NewRequest(http.MethodGet, blogUrl, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPostGet_WithIdInternalServerError(t *testing.T) {
	mockService := &MockBlogsService{
		GetPostsError: errors.New("internal server error"),
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPostUpdate_Success(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"
	sampleBool := true
	body := UpdatePostRequestBody{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestPostUpdate_InvalidRequest(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"

	payload, _ := json.Marshal(sampleText)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostUpdate_ValidationError(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"
	sampleAuthorName := "123"
	sampleBool := true
	body := UpdatePostRequestBody{
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleAuthorName,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostUpdate_BadParamsRequest(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := "s"

	url := fmt.Sprintf("%s/%s", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostUpdate_PostDoesNotExistError(t *testing.T) {
	mockService := &MockBlogsService{
		UpdatePostError: store.ErrPostDoesNotExist,
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"
	sampleBool := true
	body := UpdatePostRequestBody{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPostUpdate_SlugAlreadyExistsError(t *testing.T) {
	mockService := &MockBlogsService{
		UpdatePostError: store.ErrSlugAlreadyExists,
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"
	sampleBool := true
	body := UpdatePostRequestBody{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestPostUpdate_InternalServerError(t *testing.T) {
	mockService := &MockBlogsService{
		UpdatePostError: errors.New("internal server error"),
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	sampleText := "sample"
	sampleBool := true
	body := UpdatePostRequestBody{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPostDelete_Success(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestPostDelete_BadParamsRequest(t *testing.T) {
	mockService := &MockBlogsService{}

	app := setupBlogsTestServer(mockService)

	sampleId := "s"

	url := fmt.Sprintf("%s/%s", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPostDelete_PostDoesNotExistsError(t *testing.T) {
	mockService := &MockBlogsService{
		DeletePostError: store.ErrPostDoesNotExist,
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPostDelete_InternalServerError(t *testing.T) {
	mockService := &MockBlogsService{
		DeletePostError: errors.New("internal server error"),
	}

	app := setupBlogsTestServer(mockService)

	sampleId := 1

	url := fmt.Sprintf("%s/%d", blogUrl, sampleId)
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
