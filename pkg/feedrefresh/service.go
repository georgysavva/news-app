package feedrefresh

import (
	"context"

	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/article"
)

type Service interface {
	RefreshFeed(ctx context.Context) error
}

type Repository interface {
	ReplaceArticles(ctx context.Context, articles []*article.Article) error
}

type Downloader interface {
	DownloadArticles(ctx context.Context, providerURL string) ([]*article.Article, error)
}

type serviceImpl struct {
	repository   Repository
	downloader   Downloader
	providerURLs []string
}

func NewService(repository Repository, downloader Downloader, providerURLs []string) Service {
	return serviceImpl{
		repository:   repository,
		downloader:   downloader,
		providerURLs: providerURLs,
	}
}

func (s serviceImpl) RefreshFeed(ctx context.Context) error {
	var freshArticles []*article.Article
	// We could optimize this part and download for each provider in parallel.
	for _, pURL := range s.providerURLs {
		articles, err := s.downloader.DownloadArticles(ctx, pURL)
		if err != nil {
			return errors.Wrap(err, "can not download fresh articles")
		}
		freshArticles = append(freshArticles, articles...)
	}

	err := s.repository.ReplaceArticles(ctx, freshArticles)

	return errors.Wrap(err, "can not replacing articles with fresh ones in repository")
}
