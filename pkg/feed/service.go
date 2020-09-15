package feed

import (
	"context"

	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/article"
)

type Service interface {
	GetArticles(ctx context.Context, opts ...Option) ([]*article.Article, error)
	GetArticle(ctx context.Context, articleID string) (*article.Article, error)
	GetProviders(ctx context.Context) ([]string, error)
	GetCategories(ctx context.Context) ([]string, error)
}

type Option func(o *options)

type options struct {
	categories []string
	providers  []string
}

func WithCategories(categories []string) Option {
	return func(o *options) {
		o.categories = categories
	}
}

func WithProvides(providers []string) Option {
	return func(o *options) {
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
	return serviceImpl{repository: repository}
}

func (s serviceImpl) GetArticles(ctx context.Context, opts ...Option) ([]*article.Article, error) {
	fo := &options{}
	for _, o := range opts {
		o(fo)
	}
	articles, err := s.repository.GetArticles(ctx, fo.categories, fo.providers)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve articles list from the repository")
	}
	return articles, nil
}

func (s serviceImpl) GetArticle(ctx context.Context, articleID string) (*article.Article, error) {
	a, err := s.repository.GetArticle(ctx, articleID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not retrieve article from the repository")
	}
	return a, nil
}

func (s serviceImpl) GetProviders(ctx context.Context) ([]string, error) {
	providers, err := s.repository.GetProviders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve providers list from the repository")
	}
	return providers, nil
}

func (s serviceImpl) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.repository.GetCategories(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can not retrieve categories list from the repository")
	}
	return categories, nil
}
