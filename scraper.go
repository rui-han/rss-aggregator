package main

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Println("Found post: ", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
