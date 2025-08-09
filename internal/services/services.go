package services

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/entity"
	store2 "github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/store/cache"
	"github.com/Roh-Bot/blog-api/pkg/logger"
)

type Service struct {
	Blog IBlog
	Auth IAuth
}

type IBlog interface {
	GetPosts(ctx context.Context, getPosts *GetPostsDto) ([]entity.Post, error)
	AddPost(ctx context.Context, addPost *AddPostDto) error
	UpdatePost(ctx context.Context, addPost *UpdatePostDto) error
	DeletePost(ctx context.Context, addPost *DeletePostDto) error
}

type IAuth interface {
	GenerateToken(username string) (string, error)
	IsValid(username string) bool
	ValidateToken(token string) (bool, error)
}

func NewService(config *config.AtomicConfig, auth auth.Authentication, store store2.Store, cache cache.Cache, logger logger.Logger) Service {
	return Service{
		Blog: &BlogService{
			logger: logger,
			config: config,
			cache:  cache,
			store:  store,
		},
		Auth: &AuthService{
			config: config,
			logger: logger,
			cache:  cache,
			auth:   auth,
		},
	}
}
