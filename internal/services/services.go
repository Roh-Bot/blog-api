package services

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	store2 "github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/pkg/logger"
)

type Service struct {
	Blog IBlog
	Auth IAuth
}

type IBlog interface {
	GetPosts(ctx context.Context, getPosts *GetPostsDto) (*store2.Post, error)
	AddPost(ctx context.Context, addPost *AddPostDto) error
	UpdatePost(ctx context.Context, addPost *UpdatePostDto) error
	DeletePost(ctx context.Context, addPost *DeletePostDto) error
}

type IAuth interface {
	GenerateToken(userId string) (string, error)
	IsValid(username string) bool
	ValidateToken(token string) (bool, error)
}

func NewService(store store2.Store, logger logger.Logger, config *config.AtomicConfig, auth *auth.Authentication) *Service {
	return &Service{
		Blog: &BlogService{
			logger: logger,
			store:  store,
			config: config,
		},
		Auth: &AuthService{
			config: config,
			logger: logger,
			store:  store,
			auth:   auth,
		},
	}
}
