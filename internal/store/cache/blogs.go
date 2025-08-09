package cache

import (
	"context"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/database"
	"github.com/redis/go-redis/v9"
	"time"
)

type Blogs struct {
	db     *database.Cache
	config *config.AtomicConfig
}

func (b *Blogs) GetPosts(ctx context.Context, key string) (isExpired bool, result string, err error) {
	result, err = b.db.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return true, result, nil
		}
		return
	}
	return
}

func (b *Blogs) SetPost(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := b.db.SetNX(ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}
