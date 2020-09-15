package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/feed"
	"github.com/georgysavva/news-app/pkg/feedrefresh"
	"github.com/georgysavva/news-app/pkg/feedrefresh/download"
	"github.com/georgysavva/news-app/pkg/inmemory"
)

// This can come from configs/environment to be more flexible.
var providerURLs = []string{
	"http://feeds.bbci.co.uk/news/uk/rss.xml",
	"http://feeds.bbci.co.uk/news/technology/rss.xml",
	"http://feeds.skynews.com/feeds/rss/uk.xml",
	"http://feeds.skynews.com/feeds/rss/technology.xml",
}

func main() {
	storage := &inmemory.Storage{}
	feedService := feed.NewService(storage)
	articlesDownloader := download.NewDownloader()
	feedRefreshService := feedrefresh.NewService(storage, articlesDownloader, providerURLs)
	feedHandler := feed.MakeHTTPHandler(feedService)
	feedrefreshHandler := feedrefresh.MakeHTTPHandler(feedRefreshService)

	http.Handle("/feed/", http.StripPrefix("/feed", feedHandler))
	http.Handle("/feedrefresh/", http.StripPrefix("/feedrefresh", feedrefreshHandler))

	// We could add graceful shutdown here.
	log.Println("Starting http server on 8080 port")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(errors.Wrap(err, "http server failed"))
	}
}
