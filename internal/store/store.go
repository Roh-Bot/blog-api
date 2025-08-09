package store

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/database"
	"github.com/Roh-Bot/blog-api/internal/entity"
)

type Store struct {
	Blogs IBlogs
}

type IBlogs interface {
	AddPost(context.Context, *AddPostQueryParam) error
	GetPosts(ctx context.Context, postParam *GetPostsQueryParams) ([]entity.Post, error)
	UpdatePost(context.Context, *UpdatePostQueryParams) error
	DeletePost(context.Context, *DeletePostQueryParams) error
}

func NewStorage(db *database.Database, config *config.AtomicConfig) Store {
	return Store{
		Blogs: &BlogStore{db: db, config: config},
	}
}
