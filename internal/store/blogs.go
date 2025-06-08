package store

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

const (
	stateP1 = "P0001"
	stateP2 = "P0002"
)

var (
	ErrSlugAlreadyExists = errors.New("Slug already exists. Please choose a different one")
	ErrPostDoesNotExist  = errors.New("Post does not exist.")
)

type BlogStore struct {
	db     *pgxpool.Pool
	config *config.AtomicConfig
}

// Post is a data model for users table
type Post struct {
	Id          int
	Title       string
	Description string
	Slug        string
	Content     string
	AuthorName  string
	Tags        []string
	PublishedAt time.Time
}

type GetPostsQueryParams struct {
	Id *int
}

type AddPostQueryParam struct {
	Title       *string
	Description *string
	Slug        *string
	Content     *string
	AuthorName  *string
	Tags        *[]string
	IsPublished *bool
}

type UpdatePostQueryParams struct {
	Id          int
	Title       *string
	Description *string
	Slug        *string
	Content     *string
	AuthorName  *string
	Tags        *[]string
	IsPublished *bool
}

type DeletePostQueryParams struct {
	Id int
}

func (u *BlogStore) AddPost(ctx context.Context, postParam *AddPostQueryParam) error {
	if _, err := u.db.Exec(ctx, `SELECT * FROM posts_insert($1, $2, $3, $4, $5, $6, $7)`,
		postParam.Title, postParam.Description, postParam.Slug, postParam.Content, postParam.AuthorName, postParam.Tags, postParam.IsPublished); err != nil {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) {
			return err
		}
		if pgError.Code == stateP1 {
			return ErrSlugAlreadyExists
		}
		return err
	}
	return nil
}

func (u *BlogStore) GetPosts(ctx context.Context, postParam *GetPostsQueryParams) ([]Post, error) {
	rows, err := u.db.Query(ctx,
		`SELECT * FROM posts_get($1)`, postParam.Id)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(
			&post.Id, &post.Title, &post.Description, &post.Slug, &post.Content,
			&post.AuthorName, &post.Tags, &post.PublishedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) && pgError.Code == stateP1 {
			return nil, ErrPostDoesNotExist
		}
		return nil, err
	}

	return posts, nil
}

func (u *BlogStore) UpdatePost(ctx context.Context, postParam *UpdatePostQueryParams) error {
	if _, err := u.db.Exec(ctx, `SELECT * FROM posts_update($1, $2, $3, $4, $5, $6, $7, $8)`,
		postParam.Id, postParam.Title, postParam.Description, postParam.Slug, postParam.Content,
		postParam.AuthorName, postParam.Tags, postParam.IsPublished); err != nil {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) {
			return err
		}
		switch {
		case pgError.Code == stateP1:
			return ErrPostDoesNotExist
		case pgError.Code == stateP2:
			return ErrSlugAlreadyExists
		}
		return err
	}
	return nil
}

func (u *BlogStore) DeletePost(ctx context.Context, postParam *DeletePostQueryParams) error {
	if _, err := u.db.Exec(ctx, `SELECT * FROM posts_delete($1)`, postParam.Id); err != nil {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) {
			return err
		}
		if pgError.Code == stateP1 {
			return ErrPostDoesNotExist
		}
		return err
	}
	return nil
}
