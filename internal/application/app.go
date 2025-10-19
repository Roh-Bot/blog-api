package application

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/entity"
	store2 "github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/store/cache"
	"github.com/Roh-Bot/blog-api/pkg/logger"
)

type App struct {
	Blog IBlogUseCase
	Auth IAuthUseCase
}

type IBlogUseCase interface {
	GetPosts(ctx context.Context, getPosts *GetPostsDto) ([]entity.Post, error)
	AddPost(ctx context.Context, addPost *AddPostDto) error
	UpdatePost(ctx context.Context, addPost *UpdatePostDto) error
	DeletePost(ctx context.Context, addPost *DeletePostDto) error
}

type IAuthUseCase interface {
	GenerateToken(username string) (string, error)
	IsValid(username string) bool
	ValidateToken(token string) (bool, error)
}

func NewService(config *config.AtomicConfig, auth auth.Authentication, store store2.Store, cache cache.Cache, logger logger.Logger) App {
	return App{
		Blog: &BlogUseCase{
			logger: logger,
			config: config,
			cache:  cache,
			store:  store,
		},
		Auth: &AuthUseCase{
			config: config,
			logger: logger,
			cache:  cache,
			auth:   auth,
		},
	}
}
