package post

import (
	"capstone/constants"
	"capstone/entities"
	postEntities "capstone/entities/post"
	"capstone/utilities"
	"mime/multipart"
)

type PostUseCase struct {
	postRepository postEntities.RepositoryInterface
}

func NewPostUseCase(postRepository postEntities.RepositoryInterface) *PostUseCase {
	return &PostUseCase{
		postRepository: postRepository,
	}
}

func (postUseCase *PostUseCase) GetAllPostsByForumId(forumId uint, metadata entities.Metadata, userId uint) ([]postEntities.Post, error) {
	posts, err := postUseCase.postRepository.GetAllPostsByForumId(forumId, metadata, userId)
	if err != nil {
		return []postEntities.Post{}, err
	}
	return posts, nil
}

func (postUseCase *PostUseCase) GetPostById(postId uint, userId uint) (postEntities.Post, error) {
	post, err := postUseCase.postRepository.GetPostById(postId, userId)
	if err != nil {
		return postEntities.Post{}, err
	}
	return post, nil
}

func (postUseCase *PostUseCase) SendPost(post postEntities.Post, file *multipart.FileHeader) (postEntities.Post, error) {
	if post.ForumId == 0 || post.Content == "" {
		return postEntities.Post{}, constants.ErrEmptyInputPost
	}

	if file != nil {
		secureUrl, err := utilities.UploadImage(file)
		if err != nil {
			return postEntities.Post{}, constants.ErrUploadImage
		}
		post.ImageUrl = secureUrl
	}

	post, err := postUseCase.postRepository.SendPost(post)
	if err != nil {
		return postEntities.Post{}, err
	}
	return post, nil
}

func (postUseCase *PostUseCase) LikePost(postId uint, userId uint) error {
	if postId == 0 {
		return constants.ErrEmptyInputLike
	}

	err := postUseCase.postRepository.LikePost(postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (postUseCase *PostUseCase) UnlikePost(postId uint, userId uint) error {
	if postId == 0 {
		return constants.ErrEmptyInputLike
	}

	err := postUseCase.postRepository.UnlikePost(postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (postUseCase *PostUseCase) SendComment(comment postEntities.PostComment) (postEntities.PostComment, error) {
	if comment.PostID == 0 || comment.Content == "" {
		return postEntities.PostComment{}, constants.ErrEmptyInputComment
	}

	comment, err := postUseCase.postRepository.SendComment(comment)
	if err != nil {
		return postEntities.PostComment{}, err
	}
	return comment, nil
}

func (postUseCase *PostUseCase) GetAllCommentByPostId(postId uint, metadata entities.Metadata) ([]postEntities.PostComment, error) {
	comments, err := postUseCase.postRepository.GetAllCommentByPostId(postId, metadata)
	if err != nil {
		return []postEntities.PostComment{}, err
	}
	return comments, nil
}