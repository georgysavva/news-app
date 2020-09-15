package download

import (
	"context"
	"encoding/base64"

	"github.com/mmcdole/gofeed"

	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/article"
)

type Downloader struct {
	fp *gofeed.Parser
}

func NewDownloader() *Downloader {
	return &Downloader{fp: gofeed.NewParser()}
}

func (d *Downloader) DownloadArticles(ctx context.Context, providerURL string) ([]*article.Article, error) {
	feed, err := d.fp.ParseURLWithContext(providerURL, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "gofeed could not download or parse the feed")
	}

	var category *string
	if len(feed.Categories) >= 1 {
		category = &feed.Categories[0]
	}

	var articles []*article.Article
	for _, item := range feed.Items {
		if item.GUID == "" && item.PublishedParsed == nil {
			continue
		}
		a := &article.Article{
			ID:          base64.URLEncoding.EncodeToString([]byte(item.GUID)),
			Title:       item.Title,
			Provider:    feed.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     *item.PublishedParsed,
			Category:    category,
		}
		thumbnails := item.Extensions["media"]["thumbnail"]
		if len(thumbnails) >= 1 {
			thumbnail := thumbnails[0].Attrs
			url, urlPresent := thumbnail["url"]
			width, widthPresent := thumbnail["width"]
			height, heightPresent := thumbnail["height"]
			if urlPresent && widthPresent && heightPresent {
				a.Thumbnail = &article.Image{
					URL:    url,
					Width:  width,
					Height: height,
				}
			}
		}
		articles = append(articles, a)
	}
	return articles, nil
}
