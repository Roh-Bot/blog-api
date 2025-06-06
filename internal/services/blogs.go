package services

import (
	"context"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/config"
	store2 "github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/pkg/logger"
)

type BlogService struct {
	logger logger.Logger
	store  store2.Store
	config *config.AtomicConfig
}

type GetPostsDto struct {
	Id int
}

type AddPostDto struct {
	Title       *string
	Description *string
	Slug        *string
	Content     *string
	AuthorName  *string
	Tags        *[]string
	IsPublished *bool
}

type UpdatePostDto struct {
	Id          int
	Title       *string
	Description *string
	Slug        *string
	Content     *string
	AuthorName  *string
	Tags        *[]string
	IsPublished *bool
}

type DeletePostDto struct {
	Id int
}

var (
	errInvalidInput = errors.New("invalid input")
)

func (b *BlogService) GetPosts(ctx context.Context, post *GetPostsDto) ([]store2.Post, error) {
	return b.store.Blogs.GetPosts(ctx, &store2.GetPostsQueryParams{
		Id: post.Id,
	})
}

func (b *BlogService) AddPost(ctx context.Context, post *AddPostDto) error {
	if err := b.store.Blogs.AddPost(ctx, &store2.AddPostQueryParam{
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
		Content:     post.Content,
		AuthorName:  post.AuthorName,
		Tags:        post.Tags,
		IsPublished: post.IsPublished,
	}); err != nil {
		return err
	}
	return nil
}

func (b *BlogService) UpdatePost(ctx context.Context, post *UpdatePostDto) error {
	return b.store.Blogs.UpdatePost(ctx, &store2.UpdatePostQueryParams{
		Id:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
		Content:     post.Content,
		AuthorName:  post.AuthorName,
		Tags:        post.Tags,
		IsPublished: post.IsPublished,
	})
}

func (b *BlogService) DeletePost(ctx context.Context, post *DeletePostDto) error {
	return b.store.Blogs.DeletePost(ctx, &store2.DeletePostQueryParams{
		Id: post.Id,
	})
}

func (b *BlogService) ErrInvalidInput() error {
	return errInvalidInput
}
