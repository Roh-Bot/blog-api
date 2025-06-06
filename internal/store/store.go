package store

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Blogs IBlogs
}

type IBlogs interface {
	AddPost(context.Context, *AddPostQueryParam) error
	GetPosts(ctx context.Context, postParam *GetPostsQueryParams) ([]Post, error)
	UpdatePost(context.Context, *UpdatePostQueryParams) error
	DeletePost(context.Context, *DeletePostQueryParams) error
}

func NewStorage(db *pgxpool.Pool, config *config.AtomicConfig) Store {
	return Store{
		Blogs: &BlogStore{db: db, config: config},
	}
}
