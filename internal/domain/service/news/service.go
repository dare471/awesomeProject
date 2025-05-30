package news

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/news"
)

type NewsService struct {
	newsRepo news.Repository
}

func NewNewsService() *NewsService {
	return &NewsService{
		newsRepo: news.NewsRepository(database.GetDB()),
	}
}

func (s *NewsService) GetNewsByID(id uint) (news.News, error) {
	return s.newsRepo.FindByID(id)
}

func (s *NewsService) CreateNews(news *news.News) error {
	return s.newsRepo.Create(news)
}

func (s *NewsService) GetAllNews() ([]news.News, error) {
	return s.newsRepo.FindAll()
}
