package inmemory

import (
	"context"
	"sort"
	"sync"

	"github.com/georgysavva/news-app/pkg/article"
)

type Storage struct {
	sync.RWMutex
	articlesIndex map[string]*article.Article
	categoriesSet map[string]struct{}
	providersSet  map[string]struct{}
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) ReplaceArticles(_ context.Context, articles []*article.Article) error {
	articlesIndex := map[string]*article.Article{}
	categoriesSet := map[string]struct{}{}
	providersSet := map[string]struct{}{}
	for _, a := range articles {
		articlesIndex[a.ID] = a
		if a.Category != nil {
			categoriesSet[*a.Category] = struct{}{}
		}
		providersSet[a.Provider] = struct{}{}
	}
	s.Lock()
	defer s.Unlock()
	s.articlesIndex = articlesIndex
	s.categoriesSet = categoriesSet
	s.providersSet = providersSet
	return nil
}

func (s *Storage) GetArticles(_ context.Context, categories, providers []string) ([]*article.Article, error) {
	articles := s.retrieveArticles(categories, providers)
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PubDate.After(articles[j].PubDate)
	})
	return articles, nil
}

func (s *Storage) GetArticle(_ context.Context, articleID string) (*article.Article, error) {
	s.RLock()
	defer s.RUnlock()
	a := s.articlesIndex[articleID]
	return a, nil
}

func (s *Storage) GetCategories(_ context.Context) ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	var categories []string
	for category := range s.categoriesSet {
		categories = append(categories, category)
	}
	return categories, nil
}

func (s *Storage) GetProviders(_ context.Context) ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	var providers []string
	for provider := range s.providersSet {
		providers = append(providers, provider)
	}
	return providers, nil
}

func (s *Storage) retrieveArticles(categories, providers []string) []*article.Article {
	s.RLock()
	defer s.RUnlock()
	var articles []*article.Article
	for _, a := range s.articlesIndex {
		if a.Category != nil {
			if len(categories) > 0 && !stringInSlice(*a.Category, categories) {
				continue
			}
		}
		if len(providers) > 0 && !stringInSlice(a.Provider, providers) {
			continue
		}
		articles = append(articles, a)
	}
	return articles
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
