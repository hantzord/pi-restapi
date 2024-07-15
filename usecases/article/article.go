package article

import (
	"capstone/constants"
	"capstone/entities"
	articleEntities "capstone/entities/article"
	"capstone/utilities"
	"mime/multipart"
)

type ArticleUseCase struct {
	articleRepository articleEntities.ArticleRepositoryInterface
}

func NewArticleUseCase(articleRepository articleEntities.ArticleRepositoryInterface) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepository: articleRepository,
	}
}

func (useCase *ArticleUseCase) CreateArticle(article *articleEntities.Article, userId int) (*articleEntities.Article, error) {
	if article.Title == "" || article.Content == "" {
		return nil, constants.ErrEmptyInputArticle
	}

	createdArticle, err := useCase.articleRepository.CreateArticle(article, userId)
	if err != nil {
		return nil, err
	}

	return createdArticle, nil
}

func (useCase *ArticleUseCase) GetAllArticle(metadata entities.Metadata, userId int, search string) ([]articleEntities.Article, error) {
	articles, err := useCase.articleRepository.GetAllArticle(metadata, userId, search)
	if err != nil {
		return []articleEntities.Article{}, err
	}
	return articles, nil
}

func (useCase *ArticleUseCase) GetArticleById(articleId int, userId int) (articleEntities.Article, error) {
	// err := useCase.articleRepository.IncrementViewCount(articleId)
	// if err != nil {
	// 	return articleEntities.Article{}, err
	// }

	article, err := useCase.articleRepository.GetArticleById(articleId, userId)
	if err != nil {
		return articleEntities.Article{}, err
	}
	return article, nil
}

func (useCase *ArticleUseCase) GetLikedArticle(metadata entities.Metadata, userId int) ([]articleEntities.Article, error) {
	articles, err := useCase.articleRepository.GetLikedArticle(metadata, userId)
	if err != nil {
		return []articleEntities.Article{}, err
	}
	return articles, nil
}

func (useCase *ArticleUseCase) LikeArticle(articleId int, userId int) error {
	err := useCase.articleRepository.LikeArticle(articleId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *ArticleUseCase) UnlikeArticle(articleId int, userId int) error {
	err := useCase.articleRepository.UnlikeArticle(articleId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *ArticleUseCase) GetArticleByIdForDoctor(articleId int) (articleEntities.Article, error) {
	articles, err := useCase.articleRepository.GetArticleByIdForDoctor(articleId)
	if err != nil {
		return articleEntities.Article{}, err
	}
	return articles, nil
}

func (useCase *ArticleUseCase) GetAllArticleByDoctorId(metadata entities.MetadataFull, doctorId int) ([]articleEntities.Article, error) {
	articles, err := useCase.articleRepository.GetAllArticleByDoctorId(metadata, doctorId)
	if err != nil {
		return []articleEntities.Article{}, err
	}
	return articles, nil
}

func (useCase *ArticleUseCase) CountArticleByDoctorId(doctorId int) (int, error) {
	count, err := useCase.articleRepository.CountArticleByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (useCase *ArticleUseCase) CountArticleLikesByDoctorId(doctorId int) (int, error) {
	count, err := useCase.articleRepository.CountArticleLikesByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (useCase *ArticleUseCase) CountArticleViewByDoctorId(doctorId int) (int, error) {
	count, err := useCase.articleRepository.CountArticleViewByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (useCase *ArticleUseCase) CountArticleViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error) {
	if startMonth == "" || endMonth == "" {
		return map[int]int{}, constants.ErrEmptyInputViewByMonth
	}

	count, err := useCase.articleRepository.CountArticleViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return map[int]int{}, err
	}
	return count, nil
}

func (useCase *ArticleUseCase) EditArticle(article articleEntities.Article, file *multipart.FileHeader) (articleEntities.Article, error) {
	if article.Title == "" || article.Content == "" {
		return articleEntities.Article{}, constants.ErrEmptyInputArticle
	}
	if file != nil {
		secureUrl, err := utilities.UploadImage(file)
		if err != nil {
			return articleEntities.Article{}, constants.ErrUploadImage
		}
		article.ImageUrl = secureUrl
	} else {
		article.ImageUrl = ""
	}

	article, err := useCase.articleRepository.EditArticle(article)
	if err != nil {
		return articleEntities.Article{}, err
	}
	return article, nil
}

func (useCase *ArticleUseCase) DeleteArticle(articleId int) error {
	err := useCase.articleRepository.DeleteArticle(articleId)
	if err != nil {
		return err
	}
	return nil
}
