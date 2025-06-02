package news

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/news"
	"context"
	"fmt"
	"sync"
	"time"
)

type NewsService struct {
	newsRepo news.Repository
}

// NewsWithDetails содержит новость с дополнительными данными
type NewsWithDetails struct {
	news.News
	CommentsCount int       `json:"comments_count"`
	LikesCount    int       `json:"likes_count"`
	LastUpdated   time.Time `json:"last_updated"`
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

// GetAllNewsWithDetails возвращает все новости с дополнительными данными
func (s *NewsService) GetAllNewsWithDetails(ctx context.Context) ([]NewsWithDetails, error) {
	// Получаем базовый список новостей
	newsList, err := s.newsRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make([]NewsWithDetails, len(newsList))
	errors := make(chan error, len(newsList))

	// Для каждой новости запускаем горутину для получения дополнительных данных
	for i, n := range newsList {
		wg.Add(1)
		go func(index int, newsItem news.News) {
			defer wg.Done()

			// Создаем каналы для получения данных
			commentsChan := make(chan int, 1)
			likesChan := make(chan int, 1)
			updateChan := make(chan time.Time, 1)

			// Запускаем горутины для получения разных типов данных
			go func() {
				defer close(commentsChan)
				// Здесь можно добавить реальную логику получения количества комментариев
				commentsChan <- 0 // Заглушка
			}()

			go func() {
				defer close(likesChan)
				// Здесь можно добавить реальную логику получения количества лайков
				likesChan <- 0 // Заглушка
			}()

			go func() {
				defer close(updateChan)
				// Получаем время последнего обновления
				updateChan <- newsItem.UpdatedAt
			}()

			// Используем select для получения данных с таймаутом
			select {
			case comments := <-commentsChan:
				select {
				case likes := <-likesChan:
					select {
					case lastUpdate := <-updateChan:
						results[index] = NewsWithDetails{
							News:          newsItem,
							CommentsCount: comments,
							LikesCount:    likes,
							LastUpdated:   lastUpdate,
						}
					case <-time.After(500 * time.Millisecond):
						errors <- fmt.Errorf("timeout getting update time for news %d", newsItem.ID)
					}
				case <-time.After(500 * time.Millisecond):
					errors <- fmt.Errorf("timeout getting likes for news %d", newsItem.ID)
				}
			case <-time.After(500 * time.Millisecond):
				errors <- fmt.Errorf("timeout getting comments for news %d", newsItem.ID)
			case <-ctx.Done():
				errors <- ctx.Err()
			}
		}(i, n)
	}

	// Запускаем горутину для закрытия канала ошибок после завершения всех операций
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Собираем ошибки
	var errs []error
	for err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Если есть ошибки, возвращаем их
	if len(errs) > 0 {
		return results, fmt.Errorf("errors occurred while fetching news details: %v", errs)
	}

	return results, nil
}
