package api

import (
	"errors"
	"github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/gofiber/fiber/v2"
)

type AddPostModel struct {
	Title       *string   `json:"Title" validate:"required,alphanum,lte=100"`
	Description *string   `json:"Description" validate:"required,alphanum"`
	Slug        *string   `json:"Slug" validate:"required"`
	Content     *string   `json:"Content" validate:"required"`
	AuthorName  *string   `json:"AuthorName" validate:"required,alpha"`
	Tags        *[]string `json:"Tags"`
	IsPublished *bool     `json:"IsPublished" validate:"required"`
}

type GetPostsModel struct {
	Id int `params:"Id"`
}

type UpdatePostModel struct {
	Id          int       `params:"Id" validate:"required"`
	Title       *string   `json:"Title" validate:"alphanum,lte=100"`
	Description *string   `json:"Description" validate:"alphanum"`
	Slug        *string   `json:"Slug" `
	Content     *string   `json:"Content" `
	AuthorName  *string   `json:"AuthorName" validate:"alpha"`
	Tags        *[]string `json:"Tags"`
	IsPublished *bool     `json:"IsPublished" `
}

type DeletePostsModel struct {
	Id int `params:"Id"`
}

func (s *Server) postAdd(ctx *fiber.Ctx) error {
	post := AddPostModel{}
	if err := ctx.BodyParser(&post); err != nil {
		return s.badRequest(ctx, err, err.Error())
	}

	if err := s.Validator.Struct(post); err != nil {
		return s.badRequest(ctx, err, validationToErrorMessage(err))
	}

	postDto := &services.AddPostDto{
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
		Content:     post.Content,
		AuthorName:  post.AuthorName,
		Tags:        post.Tags,
		IsPublished: post.IsPublished,
	}
	if err := s.Services.Blog.AddPost(ctx.UserContext(), postDto); err != nil {
		if errors.Is(err, store.ErrSlugAlreadyExists) {
			return s.conflict(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}

	return s.writeResponse(ctx, "Blog created successfully")
}

func (s *Server) postsGet(ctx *fiber.Ctx) error {
	post := GetPostsModel{}
	if err := ctx.ParamsParser(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	posts, err := s.Services.Blog.GetPosts(ctx.UserContext(), &services.GetPostsDto{
		Id: post.Id,
	})
	if err != nil {
		if errors.Is(err, store.ErrPostDoesNotExist) {
			return s.notFound(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}
	return s.writeResponse(ctx, posts)
}

func (s *Server) postUpdate(ctx *fiber.Ctx) error {
	post := UpdatePostModel{}

	if err := ctx.ParamsParser(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	if err := ctx.BodyParser(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	if err := s.Services.Blog.UpdatePost(ctx.UserContext(), &services.UpdatePostDto{
		Id:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
		Content:     post.Content,
		AuthorName:  post.AuthorName,
		Tags:        post.Tags,
		IsPublished: post.IsPublished,
	}); err != nil {
		switch {
		case errors.Is(err, store.ErrPostDoesNotExist):
			return s.notFound(ctx, err.Error())
		case errors.Is(err, store.ErrSlugAlreadyExists):
			return s.conflict(ctx, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}
	return s.writeResponse(ctx, nil)
}

func (s *Server) postDelete(ctx *fiber.Ctx) error {
	post := DeletePostsModel{}

	if err := ctx.ParamsParser(&post); err != nil {
		return s.badRequest(ctx, err, stringEmpty)
	}

	if err := s.Services.Blog.DeletePost(ctx.UserContext(), &services.DeletePostDto{
		Id: post.Id,
	}); err != nil {
		if errors.Is(err, store.ErrPostDoesNotExist) {
			return s.internalServerError(ctx, err, err.Error())
		}
		return s.internalServerError(ctx, err, stringEmpty)
	}
	return s.writeResponse(ctx, nil)
}
