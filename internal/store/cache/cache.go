package cache

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/database"
	"time"
)

type Cache struct {
	Blogs IBlogs
}

type IBlogs interface {
	GetPosts(ctx context.Context, key string) (bool, string, error)
	SetPost(ctx context.Context, key string, value any, expiration time.Duration) error
}

func NewCache(dbCache *database.Cache, config *config.AtomicConfig) Cache {
	return Cache{
		Blogs: &Blogs{
			db:     dbCache,
			config: config,
		},
	}
}
