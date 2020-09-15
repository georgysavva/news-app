package feed

import (
	"context"

	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/article"
)

type Service interface {
	GetArticles(ctx context.Context, opts ...FeedOption) ([]*article.Article, error)
	GetArticle(ctx context.Context, articleID string) (*article.Article, error)
	GetProviders(ctx context.Context) ([]string, error)
	GetCategories(ctx context.Context) ([]string, error)
}

type FeedOption func(o *feedOptions)

type feedOptions struct {
	categories []string
	providers  []string
}

func WithCategories(categories []string) FeedOption {
	return func(o *feedOptions) {
		o.categories = categories
	}
}

func WithProvides(providers []string) FeedOption {
	return func(o *feedOptions) {
		o.providers = providers
	}
}

type Repository interface {
	GetArticles(ctx context.Context, categories, providers []string) ([]*article.Article, error)
	GetArticle(ctx context.Context, articleID string) (*article.Article, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetProviders(ctx context.Context) ([]string, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository: repository}
}

func (s *serviceImpl) GetArticles(ctx context.Context, opts ...FeedOption) ([]*article.Article, error) {
	fo := &feedOptions{}
	for _, o := range opts {
		o(fo)
	}
	articles, err := s.repository.GetArticles(ctx, fo.categories, fo.providers)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve articles list from the repository")
	}
	return articles, nil
}

func (s *serviceImpl) GetArticle(ctx context.Context, articleID string) (*article.Article, error) {
	a, err := s.repository.GetArticle(ctx, articleID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not retrieve article from the repository")
	}
	return a, nil
}

func (s *serviceImpl) GetProviders(ctx context.Context) ([]string, error) {
	providers, err := s.repository.GetProviders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve providers list from the repository")
	}
	return providers, nil
}

func (s *serviceImpl) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.repository.GetCategories(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve categories list from the repository")
	}
	return categories, nil
}
