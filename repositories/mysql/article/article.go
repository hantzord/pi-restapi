package article

import (
	"capstone/constants"
	"capstone/entities"
	articleEntities "capstone/entities/article"
	doctorEntities "capstone/entities/doctor"
	"time"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		db: db,
	}
}

func (repository *ArticleRepo) CreateArticle(article *articleEntities.Article, userId int) (*articleEntities.Article, error) {
	articleDB := Article{
		Title:       article.Title,
		Content:     article.Content,
		ImageUrl:    article.ImageUrl,
		DoctorID:    uint(userId), // Assuming userId is the ID of the doctor creating the article
		//ViewCount:   0,            // Initialize view count to 0
		ReadingTime: article.ReadingTime,
		Date:        time.Now(),
	}

	err := repository.db.Create(&articleDB).Error
	if err != nil {
		return nil, err
	}

	// Retrieve the created article with doctor details
	err = repository.db.Where("id = ?", articleDB.ID).Preload("Doctor").First(&articleDB).Error
	if err != nil {
		return nil, constants.ErrDataNotFound
	}

	articleResp := articleEntities.Article{
		ID:          articleDB.ID,
		Title:       articleDB.Title,
		Content:     articleDB.Content,
		Date:        articleDB.Date,
		ImageUrl:    articleDB.ImageUrl,
		//ViewCount:   articleDB.ViewCount,
		ReadingTime: articleDB.ReadingTime,
		DoctorID:    articleDB.DoctorID,
		Doctor: doctorEntities.Doctor{
			ID:   articleDB.Doctor.ID,
			Name: articleDB.Doctor.Name,
		},
	}

	return &articleResp, nil
}

func (repository *ArticleRepo) GetAllArticle(metadata entities.Metadata, userId int, search string) ([]articleEntities.Article, error) {
	var articlesDb []Article

	// Pagination
	err := repository.db.Limit(metadata.Limit).
		Offset((metadata.Page - 1) * metadata.Limit).
		Preload("Doctor"). // Preload doctor data
		Find(&articlesDb).Error
	if err != nil {
		return nil, err
	}

	// Search
	if search != "" {
		err = repository.db.Where("title LIKE ?", "%"+search+"%").Preload("Doctor").Find(&articlesDb).Error
		if err != nil {
			return nil, err
		}
	}

	articleLikes := make([]ArticleLikes, len(articlesDb))
	var counter int64
	var isLiked []bool

	// Check if articles are liked by the user
	for i := 0; i < len(articlesDb); i++ {
		articleLikes[i].UserId = uint(userId)
		articleLikes[i].ArticleID = articlesDb[i].ID
		err = repository.db.Model(&articleLikes[i]).Where("user_id = ? AND article_id = ?", articleLikes[i].UserId, articleLikes[i].ArticleID).Count(&counter).Error

		if err != nil {
			return nil, err
		}

		if counter > 0 {
			isLiked = append(isLiked, true)
		} else {
			isLiked = append(isLiked, false)
		}

		counter = 0
	}

	articlesEnt := make([]articleEntities.Article, len(articlesDb))
	for i := 0; i < len(articlesDb); i++ {
		articlesEnt[i] = articleEntities.Article{
			ID:          articlesDb[i].ID,
			Title:       articlesDb[i].Title,
			Content:     articlesDb[i].Content,
			Date:        articlesDb[i].Date,
			ImageUrl:    articlesDb[i].ImageUrl,
			//ViewCount:   articlesDb[i].ViewCount,
			ReadingTime: articlesDb[i].ReadingTime,
			DoctorID:    articlesDb[i].DoctorID,
			Doctor: doctorEntities.Doctor{
				ID:   articlesDb[i].Doctor.ID,
				Name: articlesDb[i].Doctor.Name,
			},
			IsLiked: isLiked[i],
		}
	}

	return articlesEnt, nil
}

func (repository *ArticleRepo) GetArticleById(articleId int, userId int) (articleEntities.Article, error) {
	var articleDb Article
	err := repository.db.Where("id = ?", articleId).Preload("Doctor").First(&articleDb).Error
	if err != nil {
		return articleEntities.Article{}, constants.ErrDataNotFound
	}

	var articleLikes ArticleLikes
	var isLiked bool
	var counter int64

	err = repository.db.Model(&articleLikes).Where("user_id = ? AND article_id = ?", userId, articleId).Count(&counter).Error

	if err != nil {
		return articleEntities.Article{}, constants.ErrServer
	}

	if counter > 0 {
		isLiked = true
	} else {
		isLiked = false
	}

	articleResp := articleEntities.Article{
		ID:          articleDb.ID,
		Title:       articleDb.Title,
		Content:     articleDb.Content,
		Date:        articleDb.Date,
		ImageUrl:    articleDb.ImageUrl,
		//ViewCount:   articleDb.ViewCount,
		ReadingTime: articleDb.ReadingTime,
		DoctorID:    articleDb.DoctorID,
		Doctor: doctorEntities.Doctor{
			ID:   articleDb.Doctor.ID,
			Name: articleDb.Doctor.Name,
		},
		IsLiked: isLiked,
	}

	err = repository.db.Model(&ArticleViews{}).Create(&ArticleViews{ArticleID: uint(articleId), UserId: uint(userId)}).Error
	if err != nil {
		return articleEntities.Article{}, constants.ErrServer
	}
	
	return articleResp, nil
}

func (repository *ArticleRepo) GetLikedArticle(metadata entities.Metadata, userId int) ([]articleEntities.Article, error) {
	var articleLikesDb []ArticleLikes
	err := repository.db.Limit(metadata.Limit).Offset((metadata.Page-1)*metadata.Limit).Where("user_id = ?", userId).Find(&articleLikesDb).Error
	if err != nil {
		return nil, constants.ErrDataNotFound
	}

	var likedArticleIDs []int
	for i := 0; i < len(articleLikesDb); i++ {
		likedArticleIDs = append(likedArticleIDs, int(articleLikesDb[i].ArticleID))
	}

	var articlesDb []Article
	err = repository.db.Where("id IN ?", likedArticleIDs).Preload("Doctor").Find(&articlesDb).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	articlesEnt := make([]articleEntities.Article, len(articlesDb))
	for i := 0; i < len(articlesDb); i++ {
		articlesEnt[i] = articleEntities.Article{
			ID:          articlesDb[i].ID,
			Title:       articlesDb[i].Title,
			Content:     articlesDb[i].Content,
			Date:        articlesDb[i].Date,
			ImageUrl:    articlesDb[i].ImageUrl,
			//ViewCount:   articlesDb[i].ViewCount,
			ReadingTime: articlesDb[i].ReadingTime,
			DoctorID:    articlesDb[i].DoctorID,
			Doctor: doctorEntities.Doctor{
				ID:   articlesDb[i].Doctor.ID,
				Name: articlesDb[i].Doctor.Name,
			},
			IsLiked: true,
		}
	}

	return articlesEnt, nil
}

func (repository *ArticleRepo) LikeArticle(articleId int, userId int) error {
	var articleLikes ArticleLikes

	err := repository.db.Where("user_id = ? AND article_id = ?", userId, articleId).First(&articleLikes).Error
	if err == nil {
		return constants.ErrAlreadyLiked
	}

	var article Article
	err = repository.db.Where("id = ?", articleId).First(&article).Error
	if err != nil {
		return constants.ErrServer
	}

	articleLikes = ArticleLikes{
		UserId:    uint(userId),
		ArticleID: uint(articleId),
	}

	err = repository.db.Create(&articleLikes).Error
	if err != nil {
		return constants.ErrServer
	}

	return nil
}

func (repository *ArticleRepo) UnlikeArticle(articleId int, userId int) error {
	var articleLikes ArticleLikes
	err := repository.db.Where("user_id = ? AND article_id = ?", userId, articleId).Delete(&articleLikes).Error
	if err != nil {
		return constants.ErrServer
	}
	return nil
}

func (repository *ArticleRepo) GetArticleByIdForDoctor(articleId int) (articleEntities.Article, error) {
	var articleDb Article
	err := repository.db.Where("id = ?", articleId).First(&articleDb).Error
	if err != nil {
		return articleEntities.Article{}, constants.ErrDataNotFound
	}

	return articleEntities.Article{
		ID:          articleDb.ID,
		Title:       articleDb.Title,
		Content:     articleDb.Content,
		Date:        articleDb.Date,
		ImageUrl:    articleDb.ImageUrl,
		//ViewCount:   articleDb.ViewCount,
		ReadingTime: articleDb.ReadingTime,
		DoctorID:    articleDb.DoctorID,
	}, nil
}

func (repository *ArticleRepo) GetAllArticleByDoctorId(metadata entities.MetadataFull, doctorId int) ([]articleEntities.Article, error) {
	var articleDb []Article

	query := repository.db.Where("doctor_id = ?", doctorId).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Order(metadata.Sort + " " + metadata.Order)

	if metadata.Search != "" {
		query = query.Where("title LIKE ?", "%"+metadata.Search+"%")
	}

	err := query.Find(&articleDb).Error
	if err != nil {
		return []articleEntities.Article{}, constants.ErrServer
	}

	articlesEnt := make([]articleEntities.Article, len(articleDb))
	for i := 0; i < len(articleDb); i++ {
		articlesEnt[i] = articleEntities.Article{
			ID:          articleDb[i].ID,
			Title:       articleDb[i].Title,
			Content:     articleDb[i].Content,
			Date:        articleDb[i].Date,
			ImageUrl:    articleDb[i].ImageUrl,
			//ViewCount:   articleDb[i].ViewCount,
			ReadingTime: articleDb[i].ReadingTime,
			DoctorID:    articleDb[i].DoctorID,
		}
	}

	return articlesEnt, nil
}

func (repository *ArticleRepo) CountArticleByDoctorId(doctorId int) (int, error) {
	var counter int64
	err := repository.db.Model(&Article{}).Where("doctor_id = ?", doctorId).Count(&counter).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(counter), nil
}

func (repository *ArticleRepo) CountArticleLikesByDoctorId(doctorId int) (int, error) {
	var counter int64
	err := repository.db.Table("article_likes").
		Joins("JOIN articles ON article_likes.article_id = articles.id").
		Where("articles.doctor_id = ?", doctorId).
		Count(&counter).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(counter), nil
}

func (repository *ArticleRepo) CountArticleViewByDoctorId(doctorId int) (int, error) {
	var articleDB []Article
	err := repository.db.Where("doctor_id = ?", doctorId).Find(&articleDB).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	var articleDBIDs []int
	for i := 0; i < len(articleDB); i++ {
		articleDBIDs = append(articleDBIDs, int(articleDB[i].ID))
	}

	var counter int64
	err = repository.db.Model(&ArticleViews{}).Where("article_id IN ?", articleDBIDs).Count(&counter).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(counter), nil
}

func (repository *ArticleRepo) CountArticleViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error) {
	var articleDB []Article
	err := repository.db.Where("doctor_id = ?", doctorId).Find(&articleDB).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	var articleDBIDs []int
	for i := 0; i < len(articleDB); i++ {
		articleDBIDs = append(articleDBIDs, int(articleDB[i].ID))
	}

	if len(articleDBIDs) == 0 {
		return nil, constants.ErrDataNotFound
	}

	var results []struct {
        Month int
        Views int
    }

	query := repository.db.Model(&ArticleViews{}).Select("MONTH(created_at) as month, COUNT(*) as views").
        Where("article_id IN ?", articleDBIDs).
        Where("created_at BETWEEN ? AND ?", startMonth+"-01", endMonth+"-31").
        Where("deleted_at IS NULL").
        Group("month").
        Order("month")

    err = query.Scan(&results).Error
    if err != nil {
        return nil, constants.ErrServer
    }

    viewsByMonth := make(map[int]int)
    for _, result := range results {
        viewsByMonth[result.Month] = result.Views
    }

    return viewsByMonth, nil
}

func (repository *ArticleRepo) EditArticle(article articleEntities.Article) (articleEntities.Article, error) {
	var articleDb Article

	err := repository.db.Where("id = ?", article.ID).First(&articleDb).Error
	if err != nil {
		return articleEntities.Article{}, constants.ErrDataNotFound
	}

	articleDb.Title = article.Title
	articleDb.Content = article.Content
	articleDb.Date = time.Now()

	if article.ImageUrl != "" {
		articleDb.ImageUrl = article.ImageUrl
	}

	err = repository.db.Save(&articleDb).Error
	if err != nil {
		return articleEntities.Article{}, constants.ErrServer
	}

	return articleEntities.Article{
		ID:          articleDb.ID,
		Title:       articleDb.Title,
		Content:     articleDb.Content,
		Date:        articleDb.Date,
		ImageUrl:    articleDb.ImageUrl,
		//ViewCount:   articleDb.ViewCount,
		ReadingTime: articleDb.ReadingTime,
		DoctorID:    articleDb.DoctorID,
	}, nil
}

func (repository *ArticleRepo) DeleteArticle(articleId int) error {
	return repository.db.Where("id = ?", articleId).Delete(&Article{}).Error
}

// func (repository *ArticleRepo) IncrementViewCount(articleId int) error {
// 	return repository.db.Model(&Article{}).Where("id = ?", articleId).Update("view_count", gorm.Expr("view_count + 1")).Error
// }
