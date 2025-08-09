package services

import (
	"context"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/entity"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockBlogsStore struct {
	AddPostError    error
	GetPostsError   error
	GetPostsData    []entity.Post
	UpdatePostError error
	DeletePostError error
}

func (m *MockBlogsStore) AddPost(ctx context.Context, param *store.AddPostQueryParam) error {
	return m.AddPostError
}

func (m *MockBlogsStore) GetPosts(ctx context.Context, postParam *store.GetPostsQueryParams) ([]entity.Post, error) {
	return m.GetPostsData, m.GetPostsError
}

func (m *MockBlogsStore) UpdatePost(ctx context.Context, params *store.UpdatePostQueryParams) error {
	return m.UpdatePostError

}

func (m *MockBlogsStore) DeletePost(ctx context.Context, params *store.DeletePostQueryParams) error {
	return m.DeletePostError
}

func newMockBlogService(blogsStore store.IBlogs) *BlogService {
	mockStore := store.Store{Blogs: blogsStore}
	return &BlogService{store: mockStore}
}

func TestGetPosts_EmptyDataSuccess(t *testing.T) {
	mockStore := &MockBlogsStore{}
	service := newMockBlogService(mockStore)

	id := 1
	_, err := service.GetPosts(context.Background(), &GetPostsDto{Id: &id})
	assert.NoError(t, err)
}

func TestGetPosts_WithDataSuccess(t *testing.T) {
	mockStore := &MockBlogsStore{GetPostsData: []entity.Post{{Id: 1}}}
	service := newMockBlogService(mockStore)

	posts, err := service.GetPosts(context.Background(), &GetPostsDto{Id: utils.Ptr(1)})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, 1, posts[0].Id)
}

func TestGetPosts_Error(t *testing.T) {
	mockStore := &MockBlogsStore{GetPostsError: errors.New("an error occurred")}
	service := newMockBlogService(mockStore)

	_, err := service.GetPosts(context.Background(), &GetPostsDto{Id: utils.Ptr(1)})
	assert.Error(t, err)
}

func TestAddPosts_Success(t *testing.T) {
	mockStore := &MockBlogsStore{}
	service := newMockBlogService(mockStore)

	err := service.AddPost(context.Background(), &AddPostDto{
		Title:       utils.Ptr("Title"),
		Description: utils.Ptr("Desc"),
		Slug:        utils.Ptr("slug"),
		Content:     utils.Ptr("content"),
		AuthorName:  utils.Ptr("author"),
		Tags:        &[]string{"go", "test"},
		IsPublished: utils.Ptr(true),
	})

	assert.NoError(t, err)
}

func TestAddPosts_Error(t *testing.T) {
	mockStore := &MockBlogsStore{AddPostError: errors.New("insert failed")}
	service := newMockBlogService(mockStore)

	err := service.AddPost(context.Background(), &AddPostDto{
		Title: utils.Ptr("Title"),
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "insert failed")
}

func TestUpdatePosts_Success(t *testing.T) {
	mockStore := &MockBlogsStore{}
	service := newMockBlogService(mockStore)

	err := service.UpdatePost(context.Background(), &UpdatePostDto{
		Id:    1,
		Title: utils.Ptr("Updated"),
	})

	assert.NoError(t, err)
}

func TestUpdatePosts_Error(t *testing.T) {
	mockStore := &MockBlogsStore{UpdatePostError: errors.New("update failed")}
	service := newMockBlogService(mockStore)

	err := service.UpdatePost(context.Background(), &UpdatePostDto{
		Id: 1,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
}

func TestDeletePosts_Success(t *testing.T) {
	mockStore := &MockBlogsStore{}
	service := newMockBlogService(mockStore)

	err := service.DeletePost(context.Background(), &DeletePostDto{Id: 1})

	assert.NoError(t, err)
}

func TestDeletePosts_Error(t *testing.T) {
	mockStore := &MockBlogsStore{DeletePostError: errors.New("delete failed")}
	service := newMockBlogService(mockStore)

	err := service.DeletePost(context.Background(), &DeletePostDto{Id: 1})

	assert.Error(t, err)
	assert.EqualError(t, err, "delete failed")
}
