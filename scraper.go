package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rui-han/rss-aggregator/internal/database"
)

func startScarping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scarping on %v go routines every %s durations", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// waiting for values to be received from the ticker.C channel
	for ; ; <-ticker.C {
		// for every interval (timeBetweenRequest)
		// grab the next batch of feeds to fetch
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), // global context
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}
		// fetch each feed individually at the same time
		waitGroup := &sync.WaitGroup{}
		for _, feed := range feeds {
			waitGroup.Add(1)                   // add 1 to the waitGroup for every feed
			go scrapeFeed(db, waitGroup, feed) // spawn go routines based on the number of feeds
		}
		waitGroup.Wait()
	}
}

func scrapeFeed(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		// if description is empty, set it to null in database
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		// parse the published at date
		pubAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post: ", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
