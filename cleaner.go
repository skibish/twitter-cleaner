package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type Cleaner struct {
	twitter       *Twitter
	tweetAge      time.Duration
	checkInterval time.Duration
	dryRun        bool
	userID        int64
	ticker        *time.Ticker
	shutdown      chan bool
}

func NewCleaner(twitter *Twitter, tweetAge, checkInterval time.Duration, dryRun bool) *Cleaner {
	return &Cleaner{
		twitter:       twitter,
		tweetAge:      tweetAge,
		checkInterval: checkInterval,
		dryRun:        dryRun,
		ticker:        time.NewTicker(checkInterval),
		shutdown:      make(chan bool),
	}
}

func (c *Cleaner) Init() error {
	userID, err := c.twitter.Self()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}
	c.userID = userID

	return nil
}

func (c *Cleaner) Start() error {
	var wg sync.WaitGroup

	for {
		select {
		case <-c.ticker.C:
			wg.Add(2)
			if err := c.cleanTimeline(); err != nil {
				return fmt.Errorf("failed to clean timeline: %w", err)
			}
			wg.Done()

			if err := c.cleanFavorites(); err != nil {
				return fmt.Errorf("failed to clean favorites: %w", err)
			}
			wg.Done()
		case <-c.shutdown:
			wg.Wait()
			return nil
		}
	}
}

func (c *Cleaner) Stop() {
	c.ticker.Stop()
	c.shutdown <- true
}

func (c *Cleaner) cleanTimeline() error {
	var oldestTweetID int64
	for {
		tweets, err := c.twitter.GetUserTimeline(oldestTweetID)
		if err != nil {
			return fmt.Errorf("failed to get tweets from user timeline: %w", err)
		}

		log.Printf("scanned through %d timeline tweets", len(tweets))

		if len(tweets) <= 1 {
			break
		}

		// get oldest tweet ID from the response
		// and set it as a limit for the next call
		oldestTweetID = tweets[len(tweets)-1].Id

		for _, tweet := range tweets {
			if err := c.remove(tweet); err != nil {
				return fmt.Errorf("failed to remove tweet: %w", err)
			}
		}
	}

	return nil
}

func (c *Cleaner) cleanFavorites() error {
	var oldestTweetID int64
	for {
		tweets, err := c.twitter.GetUserFavorites(oldestTweetID)
		if err != nil {
			return fmt.Errorf("failed to get tweets from favorites: %w", err)
		}

		log.Printf("scanned through %d favorite tweets", len(tweets))

		if len(tweets) <= 1 {
			break
		}

		// get oldest tweet ID from the response
		// and set it as a limit for the next call
		oldestTweetID = tweets[len(tweets)-1].Id

		for _, tweet := range tweets {
			if err := c.remove(tweet); err != nil {
				return fmt.Errorf("failed to remove tweet: %w", err)
			}
		}
	}

	return nil
}

func (c *Cleaner) remove(tweet anaconda.Tweet) error {
	createdAt, err := tweet.CreatedAtTime()
	if err != nil {
		return fmt.Errorf("failed to get createdAt time of a tweet: %w", err)
	}

	// if duration is less than retention period specified,
	// skip it
	if time.Since(createdAt) < c.tweetAge {
		return nil
	}

	if tweet.Favorited {
		log.Printf("UNFAVORITING\t%d", tweet.Id)

		if !c.dryRun {
			if err := c.twitter.UnFavorite(tweet.Id); err != nil {

				// 144: No status found with that ID
				if !strings.Contains(err.Error(), "144") {
					return fmt.Errorf("failed to unfavorite the tweet %d: %w", tweet.Id, err)
				}
			}
		}
	}

	if tweet.Retweeted {
		log.Printf("UNRETWEETING\t%d", tweet.Id)

		if !c.dryRun {
			if err := c.twitter.UnRetweet(tweet.Id); err != nil {
				return fmt.Errorf("failed to unretweet the tweet %d: %w", tweet.Id, err)
			}
		}

		// retweeted tweet can't be deleted
		return nil
	}

	// if tweet is not a users tweet,
	// return earlier
	if tweet.User.Id != c.userID {
		return nil
	}

	log.Printf("DELETING\t%d", tweet.Id)
	if !c.dryRun {
		if err := c.twitter.Delete(tweet.Id); err != nil {
			return fmt.Errorf("failed to delete the tweet %d: %w", tweet.Id, err)
		}
	}

	return nil
}
