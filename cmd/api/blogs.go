package api

import (
	"errors"
	"github.com/Roh-Bot/blog-api/internal/application"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// Request Models
type (

	// AddPostRequest represents the payload to add a new blog post.
	// swagger:model AddPostRequest
	//
	// This contains all required fields to create a blog post.
	AddPostRequest struct {
		// Post title
		// required: true
		// example: MyFirstPost
		Title *string `json:"Title" validate:"required,alphanum,lte=100"`

		// Post description
		// required: true
		// example: A sample post
		Description *string `json:"Description" validate:"required,alphanum"`

		// Unique slug for the post
		// required: true
		// example: my-first-post
		Slug *string `json:"Slug" validate:"required"`

		// Main content of the post
		// required: true
		// example: Hello World!
		Content *string `json:"Content" validate:"required"`

		// Author's name
		// required: true
		// example: Rohit
		AuthorName *string `json:"AuthorName" validate:"required,alpha"`

		// Tags associated with the post
		// example: ["go", "fiber"]
		Tags *[]string `json:"Tags"`

		// Post visibility status
		// required: true
		// example: true
		IsPublished *bool `json:"IsPublished" validate:"required"`
	}

	// GetPostsRequest represents the path parameter for fetching post(s).
	// swagger:model GetPostsRequest
	//
	// GetPostsRequest optionally takes an ID to fetch a single post.
	GetPostsRequest struct {
		// Optional: specific post ID
		// example: 1
		Id *int `params:"Id"`
	}

	// UpdatePostPathParams represents the path parameters for updating a post.
	// swagger:parameters updatePost
	UpdatePostPathParams struct {
		// Post ID to update
		// in: path
		// required: true
		// example: 1
		Id int `params:"Id" validate:"required"`
	}

	// UpdatePostRequestBody contains fields for updating a blog post.
	// swagger:model UpdatePostRequestBody
	//
	// All fields are optional.
	UpdatePostRequestBody struct {
		// Optional: updated title
		// example: Updated Post Title
		Title *string `json:"Title" validate:"alphanum,lte=100"`

		// Optional: updated description
		// example: Updated description
		Description *string `json:"Description" validate:"alphanum"`

		// Optional: updated slug
		// example: updated-post-title
		Slug *string `json:"Slug"`

		// Optional: updated content
		// example: Updated blog content here...
		Content *string `json:"Content"`

		// Optional: updated author name
		// example: RohitDev
		AuthorName *string `json:"AuthorName" validate:"alpha"`

		// Optional: updated tags
		// example: ["go", "api"]
		Tags *[]string `json:"Tags"`

		// Optional: updated visibility
		// example: false
		IsPublished *bool `json:"IsPublished"`
	}

	// DeletePostRequest contains the post ID to delete.
	// swagger:model DeletePostRequest
	//
	// DeletePostRequest uses the path parameter ID to delete a post.
	DeletePostRequest struct {
		// Post ID to be deleted
		// example: 2
		Id int `params:"Id"`
	}
)

// Response Wrappers for documentation
type (

	// GetPostsResponse represents the response returned for a blog post.
	// swagger:model GetPostsResponse
	//
	// GetPostsResponse contains complete blog post details.
	GetPostsResponse struct {
		// Post ID
		// example: 1
		Id int `json:"Id"`

		// Post title
		// example: My First Post
		Title string `json:"Title"`

		// Post description
		// example: This is a short description of the blog post.
		Description string `json:"Description"`

		// Post slug
		// example: my-first-post
		Slug string `json:"Slug"`

		// Post content
		// example: This is the main content of the post.
		Content string `json:"Content"`

		// Author's name
		// example: Rohit
		AuthorName string `json:"AuthorName"`

		// Tags associated with the post
		// example: ["go", "fiber"]
		Tags []string `json:"Tags"`

		// Published at datetime
		// example: 2006-01-02 15:05:05
		PublishedAt time.Time `json:"PublishedAt"`
	}

	// GetPostsWrapperResponse wraps GetPostsResponse inside standard Response
	// swagger:model GetPostsWrapperResponse
	GetPostsWrapperResponse struct {
		// example: 1
		Status int `json:"Status"`

		// example: ""
		Error string `json:"Error"`

		// example: {"Title":"Gaming..."}
		Data []GetPostsResponse `json:"Data"`
	}
)

// postAdd godoc
// @Summary Add a new blog post
// @Description Adds a new blog post to the system
// @Tags Blog
// @Accept json
// @Produce json
// @Param data body api.AddPostRequest true "Post details"
// @Success 200 {object} Response "Post deleted successfully"
// @Failure 400 {object} Response "Invalid request format"
// @Failure 401 {object} Response "Invalid credentials"
// @Failure 404 {object} Response "Post Not Found"
// @Failure 500 {object} Response "Internal server error"
// @Security ApiKeyAuth
// @Router /blog-post/{id} [post]
func (s *Server) postAdd(ctx echo.Context) error {
	post := AddPostRequest{}
	if err := ctx.Bind(&post); err != nil {
		return s.badRequest(ctx, err, err.Error())
	}

	if err := s.Validator.Struct(post); err != nil {
		return s.badRequest(ctx, err, validationToErrorMessage(err))
	}

	postDto := &application.AddPostDto{
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
		Content:     post.Content,
		AuthorName:  post.AuthorName,
		Tags:        post.Tags,
		IsPublished: post.IsPublished,
	}
	if err := s.App.Blog.AddPost(ctx.Request().Context(), postDto); err != nil {
		if errors.Is(err, store.ErrSlugAlreadyExists) {
			return s.conflict(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}

	return s.writeResponse(ctx, "Post created successfully")
}

// postsGet godoc
// @Summary Get blog post(s)
// @Description Get a single blog post by ID or all posts if ID is 0
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path int false "Post ID"
// @Success 200 {object} api.GetPostsWrapperResponse
// @Failure 400 {object} Response "Invalid request format"
// @Failure 401 {object} Response "Invalid credentials"
// @Failure 404 {object} Response "Post Not Found"
// @Failure 500 {object} Response "Internal server error"
// @Security ApiKeyAuth
// @Router /blog-post [get]
// @Router /blog-post/{id} [get]
func (s *Server) postsGet(ctx echo.Context) error {
	post := GetPostsRequest{}
	if err := ctx.Bind(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	posts, err := s.App.Blog.GetPosts(ctx.Request().Context(), &application.GetPostsDto{
		Id: post.Id,
	})

	if err != nil {
		if errors.Is(err, store.ErrPostDoesNotExist) {
			return s.notFound(ctx, err.Error())
		}
		return s.internalServerError(ctx, nil, stringEmpty)
	}
	return s.writeResponseWithStatusCode(ctx, http.StatusCreated, posts)
}

// postUpdate godoc
// @Summary Update an existing blog post
// @Description Updates blog post by ID
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param data body api.UpdatePostRequestBody true "Updated post data"
// @Success 200 {object} Response "Post updated successfully""
// @Failure 400 {object} Response "Invalid request format"
// @Failure 401 {object} Response "Invalid credentials"
// @Failure 500 {object} Response "Internal server error"
// @Security ApiKeyAuth
// @Router /blog-post/{id} [patch]
func (s *Server) postUpdate(ctx echo.Context) error {
	postParams := UpdatePostPathParams{}
	if err := ctx.Bind(&postParams); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	postBody := &UpdatePostRequestBody{}
	if err := ctx.Bind(postBody); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	if err := s.Validator.Struct(postBody); err != nil {
		return s.badRequest(ctx, err, validationToErrorMessage(err))
	}

	if err := s.App.Blog.UpdatePost(ctx.Request().Context(), &application.UpdatePostDto{
		Id:          postParams.Id,
		Title:       postBody.Title,
		Description: postBody.Description,
		Slug:        postBody.Slug,
		Content:     postBody.Content,
		AuthorName:  postBody.AuthorName,
		Tags:        postBody.Tags,
		IsPublished: postBody.IsPublished,
	}); err != nil {
		switch {
		case errors.Is(err, store.ErrPostDoesNotExist):
			return s.notFound(ctx, err.Error())
		case errors.Is(err, store.ErrSlugAlreadyExists):
			return s.conflict(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}
	return s.writeResponse(ctx, "Post updated successfully")
}

// postDelete godoc
// @Summary Delete a blog post
// @Description Deletes blog post by ID
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} Response "Post deleted successfully"
// @Failure 400 {object} Response "Invalid request format"
// @Failure 401 {object} Response "Invalid credentials"
// @Failure 404 {object} Response "Post Not Found"
// @Failure 500 {object} Response "Internal server error"
// @Security ApiKeyAuth
// @Router /blog-post/{id} [delete]
func (s *Server) postDelete(ctx echo.Context) error {
	post := DeletePostRequest{}

	if err := ctx.Bind(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	if err := s.App.Blog.DeletePost(ctx.Request().Context(), &application.DeletePostDto{
		Id: post.Id,
	}); err != nil {
		if errors.Is(err, store.ErrPostDoesNotExist) {
			return s.notFound(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}
	return s.writeResponse(ctx, "Post deleted successfully")
}
