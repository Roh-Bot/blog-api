package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Roh-Bot/blog-api/internal/application"
	"github.com/Roh-Bot/blog-api/internal/entity"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/validator"
	"github.com/labstack/echo/v4"
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
	GetPostsData    []entity.Post
	UpdatePostError error
	DeletePostError error
}

func (m *MockBlogsService) AddPost(ctx context.Context, addPost *application.AddPostDto) error {
	return m.AddPostError
}

func (m *MockBlogsService) GetPosts(ctx context.Context, getPosts *application.GetPostsDto) ([]entity.Post, error) {
	return m.GetPostsData, m.GetPostsError
}

func (m *MockBlogsService) UpdatePost(ctx context.Context, addPost *application.UpdatePostDto) error {
	return m.UpdatePostError
}

func (m *MockBlogsService) DeletePost(ctx context.Context, addPost *application.DeletePostDto) error {
	return m.DeletePostError
}

func setupBlogsTestServer(blogs application.IBlogUseCase) *echo.Echo {
	mockService := application.App{
		Blog: blogs,
		Auth: &MockAuthService{ShouldValidate: true},
	}

	validatorV10 := validator.NewValidator()

	server := &Server{
		Router:    echo.New(),
		Logger:    &MockLogger{},
		Validator: validatorV10,
		App:       mockService,
	}
	server.Router.Use(server.validateAuth)
	server.Router.GET(blogUrl, server.postsGet)
	server.Router.GET(blogUrlWithId, server.postsGet)
	server.Router.POST(blogUrl, server.postAdd)
	server.Router.PATCH(blogUrlWithId, server.postUpdate)
	server.Router.DELETE(blogUrlWithId, server.postDelete)
	return server.Router
}

func buildRequest(method, url string, body any) *http.Request {
	var buf bytes.Buffer
	switch b := body.(type) {
	case string:
		buf = *bytes.NewBuffer([]byte(b))
	default:
		err := json.NewEncoder(&buf).Encode(b)
		if err != nil {
			return nil
		}
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestPostAdd(t *testing.T) {
	type testCase struct {
		name           string
		requestBody    any
		mockService    *MockBlogsService
		expectedStatus int
	}

	sampleText := "sample"
	sampleBool := true
	validBody := AddPostRequest{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	tests := []testCase{
		{
			name:           "Success",
			requestBody:    validBody,
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid-json",
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation Error",
			requestBody: AddPostRequest{
				// missing Title
				Description: &sampleText,
				Slug:        &sampleText,
				Content:     &sampleText,
				AuthorName:  &sampleText,
				Tags:        &[]string{},
				IsPublished: &sampleBool,
			},
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Slug Already Exists",
			requestBody:    validBody,
			mockService:    &MockBlogsService{AddPostError: store.ErrSlugAlreadyExists},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Internal Server Error",
			requestBody:    validBody,
			mockService:    &MockBlogsService{AddPostError: errors.New("fail")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := setupBlogsTestServer(tc.mockService)

			req := buildRequest(http.MethodPost, blogUrl, tc.requestBody)

			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestPostGet(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockService    *MockBlogsService
		expectedStatus int
		expectData     bool
	}

	sampleId := 1
	sampleTitle := "Gaming"
	samplePosts := []entity.Post{{Title: sampleTitle}}

	tests := []testCase{
		{
			name:           "WithPostId_Success",
			url:            fmt.Sprintf("%s/%d", blogUrl, sampleId),
			mockService:    &MockBlogsService{GetPostsData: samplePosts},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:           "WithoutPostId_Success",
			url:            blogUrl,
			mockService:    &MockBlogsService{GetPostsData: samplePosts},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:           "BadParams",
			url:            fmt.Sprintf("%s/%s", blogUrl, "s"), // non-int
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PostDoesNotExist",
			url:            fmt.Sprintf("%s/%d", blogUrl, sampleId),
			mockService:    &MockBlogsService{GetPostsError: store.ErrPostDoesNotExist},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "InternalError_WithoutId",
			url:            blogUrl,
			mockService:    &MockBlogsService{GetPostsError: errors.New("internal error")},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "InternalError_WithId",
			url:            fmt.Sprintf("%s/%d", blogUrl, sampleId),
			mockService:    &MockBlogsService{GetPostsError: errors.New("internal error")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := setupBlogsTestServer(tc.mockService)

			req := buildRequest(http.MethodGet, tc.url, nil)

			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectData && rec.Code == http.StatusOK {
				var data GetPostsWrapperResponse
				err := json.NewDecoder(rec.Body).Decode(&data)
				assert.NoError(t, err)
				assert.NotEmpty(t, data.Data)
			}
		})
	}
}

func TestPostUpdate(t *testing.T) {
	type testCase struct {
		name           string
		postId         interface{} // can be int or string for URL param test
		body           any
		mockService    *MockBlogsService
		expectedStatus int
	}

	sampleText := "sample"
	sampleBool := true
	validBody := UpdatePostRequestBody{
		Title:       &sampleText,
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  &sampleText,
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	validationErrorBody := UpdatePostRequestBody{
		Description: &sampleText,
		Slug:        &sampleText,
		Content:     &sampleText,
		AuthorName:  func(s string) *string { return &s }("123"),
		Tags:        &[]string{},
		IsPublished: &sampleBool,
	}

	tests := []testCase{
		{
			name:           "Success",
			postId:         1,
			body:           validBody,
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON Request",
			postId:         1,
			body:           "invalid-json-string",
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Validation Error",
			postId:         1,
			body:           validationErrorBody,
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Bad Param (non-int ID)",
			postId:         "s",
			body:           nil,
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Post Does Not Exist",
			postId:         1,
			body:           validBody,
			mockService:    &MockBlogsService{UpdatePostError: store.ErrPostDoesNotExist},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Slug Already Exists",
			postId:         1,
			body:           validBody,
			mockService:    &MockBlogsService{UpdatePostError: store.ErrSlugAlreadyExists},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Internal Server Error",
			postId:         1,
			body:           validBody,
			mockService:    &MockBlogsService{UpdatePostError: errors.New("internal")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := setupBlogsTestServer(tc.mockService)

			url := fmt.Sprintf("%s/%v", blogUrl, tc.postId)

			req := buildRequest(http.MethodPatch, url, tc.body)

			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestPostDelete(t *testing.T) {
	type testCase struct {
		name           string
		postId         interface{}
		mockService    *MockBlogsService
		expectedStatus int
	}

	tests := []testCase{
		{
			name:           "Success",
			postId:         1,
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Bad Param (non-integer ID)",
			postId:         "s",
			mockService:    &MockBlogsService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Post Does Not Exist",
			postId:         1,
			mockService:    &MockBlogsService{DeletePostError: store.ErrPostDoesNotExist},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Internal Server Error",
			postId:         1,
			mockService:    &MockBlogsService{DeletePostError: errors.New("internal server error")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := setupBlogsTestServer(tc.mockService)

			url := fmt.Sprintf("%s/%v", blogUrl, tc.postId)
			req := buildRequest(http.MethodDelete, url, nil)

			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
