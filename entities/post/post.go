package post

import (
	"capstone/entities"
	"capstone/entities/user"
	"mime/multipart"
)

type Post struct {
	ID       uint
	ForumId  uint
	UserId   uint
	Content  string
	ImageUrl string
	User     user.User
	Comments []PostComment
	NumberOfComments int
	IsLiked          bool
}

type PostComment struct {
	ID       uint
	Content  string
	PostID   uint
	UserID   uint
	User     user.User
	CreatedAt string
}

type RepositoryInterface interface {
	GetAllPostsByForumId(forumId uint, metadata entities.Metadata, userId uint) ([]Post, error)
	GetPostById(postId uint, userId uint) (Post, error)
	SendPost(post Post) (Post, error)
	LikePost(postId uint, userId uint) error
	UnlikePost(postId uint, userId uint) error
	SendComment(comment PostComment) (PostComment, error)
	GetAllCommentByPostId(postId uint, metadata entities.Metadata) ([]PostComment, error)
}

type UseCaseInterface interface {
	GetAllPostsByForumId(forumId uint, metadata entities.Metadata, userId uint) ([]Post, error)
	GetPostById(postId uint, userId uint) (Post, error)
	SendPost(post Post, file *multipart.FileHeader) (Post, error)
	LikePost(postId uint, userId uint) error
	UnlikePost(postId uint, userId uint) error
	SendComment(comment PostComment) (PostComment, error)
	GetAllCommentByPostId(postId uint, metadata entities.Metadata) ([]PostComment, error)
}