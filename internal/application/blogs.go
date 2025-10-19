package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/entity"
	store2 "github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/store/cache"
	"github.com/Roh-Bot/blog-api/pkg/logger"
	"strconv"
	"time"
)

type BlogUseCase struct {
	logger logger.Logger
	store  store2.Store
	cache  cache.Cache
	config *config.AtomicConfig
}

type GetPostsDto struct {
	Id *int
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

func (b *BlogUseCase) GetPosts(ctx context.Context, post *GetPostsDto) (posts []entity.Post, err error) {
	// setting cache key for the post
	cacheKey := fmt.Sprintf("blogapi:blogs:public:posts:")
	if post.Id == nil {
		cacheKey = cacheKey + "all"
	} else {
		cacheKey = cacheKey + strconv.Itoa(*post.Id)
	}

	// Getting the value of the key from cache
	isExpired, result, err := b.cache.Blogs.GetPosts(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	// Checking if the key has expired

	// If not return the result from the cache
	if !isExpired {
		if err := json.Unmarshal([]byte(result), &posts); err != nil {
			return nil, err
		}
		return posts, nil
	}

	posts, err = b.store.Blogs.GetPosts(ctx, &store2.GetPostsQueryParams{Id: post.Id})
	if err != nil {
		return nil, err
	}

	// generating json string for cache
	postsJson, err := json.Marshal(posts)
	if err != nil {
		return nil, err
	}

	if err := b.cache.Blogs.SetPost(ctx, cacheKey, string(postsJson), time.Second*300); err != nil {
		return nil, err
	}
	return posts, nil
}

func (b *BlogUseCase) AddPost(ctx context.Context, post *AddPostDto) error {
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

func (b *BlogUseCase) UpdatePost(ctx context.Context, post *UpdatePostDto) error {
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

func (b *BlogUseCase) DeletePost(ctx context.Context, post *DeletePostDto) error {
	return b.store.Blogs.DeletePost(ctx, &store2.DeletePostQueryParams{
		Id: post.Id,
	})
}
