package main

import (
	"io"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/matryer/is"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestCleaner(t *testing.T) {
	is := is.New(t)

	tweetCreatedAt := time.Now()
	createdAtFormatted := tweetCreatedAt.Format(time.RubyDate)

	twitterAPIMock := &TwitterAPIMock{
		GetSelfFunc: func(v url.Values) (anaconda.User, error) {
			return anaconda.User{Id: 1}, nil
		},
		UnRetweetFunc: func(id int64, trimUser bool) (anaconda.Tweet, error) {
			return anaconda.Tweet{}, nil
		},
		UnfavoriteFunc: func(id int64) (anaconda.Tweet, error) {
			return anaconda.Tweet{}, nil
		},
		DeleteTweetFunc: func(id int64, trimUser bool) (anaconda.Tweet, error) {
			return anaconda.Tweet{}, nil
		},
	}

	twitterAPIMock.GetUserTimelineFunc = func(v url.Values) ([]anaconda.Tweet, error) {
		if len(twitterAPIMock.GetUserTimelineCalls()) == 2 {
			return []anaconda.Tweet{
				{
					Id:        123,
					CreatedAt: createdAtFormatted,
					Favorited: true,
				}}, nil
		}

		return []anaconda.Tweet{
			{
				Id:        123,
				CreatedAt: createdAtFormatted,
				Favorited: true,
			},
			{
				Id:        1234,
				CreatedAt: createdAtFormatted,
				Retweeted: true,
			},
			{
				Id:        1234,
				CreatedAt: createdAtFormatted,
				Retweeted: true,
				Favorited: true,
			},
			{
				Id:        1,
				User:      anaconda.User{Id: 1},
				CreatedAt: createdAtFormatted,
				Favorited: true,
			},
			{
				Id:        1,
				User:      anaconda.User{Id: 1},
				CreatedAt: createdAtFormatted,
				Retweeted: true,
			},
			{
				Id:        1,
				CreatedAt: createdAtFormatted,
			},
		}, nil
	}

	twitterAPIMock.GetFavoritesFunc = func(v url.Values) ([]anaconda.Tweet, error) {
		if len(twitterAPIMock.GetFavoritesCalls()) == 2 {
			return []anaconda.Tweet{
				{
					Id:        123,
					CreatedAt: createdAtFormatted,
					Favorited: true,
				}}, nil
		}

		return []anaconda.Tweet{
			{
				Id:        123,
				CreatedAt: createdAtFormatted,
				Favorited: true,
			},
			{
				Id:        1,
				User:      anaconda.User{Id: 1},
				CreatedAt: createdAtFormatted,
				Favorited: true,
			},
		}, nil
	}

	c := NewCleaner(&Twitter{twitterAPIMock}, time.Second*1, time.Second*1, false)
	c.Init()

	go func() {
		is.NoErr(c.Start())
	}()

	time.Sleep(time.Millisecond * 1100)
	c.Stop()

	is.Equal(1, len(twitterAPIMock.GetSelfCalls()))
	is.Equal(2, len(twitterAPIMock.GetUserTimelineCalls()))
	is.Equal(2, len(twitterAPIMock.GetFavoritesCalls()))

	is.Equal(5, len(twitterAPIMock.UnfavoriteCalls()))
	is.Equal(3, len(twitterAPIMock.UnRetweetCalls()))
	is.Equal(2, len(twitterAPIMock.DeleteTweetCalls()))

}
