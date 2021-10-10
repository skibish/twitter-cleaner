package main

import (
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

//go:generate moq -out twitterapi_moq_test.go . TwitterAPI
type TwitterAPI interface {
	GetSelf(v url.Values) (u anaconda.User, err error)
	GetUserTimeline(v url.Values) (timeline []anaconda.Tweet, err error)
	GetFavorites(v url.Values) (favorites []anaconda.Tweet, err error)
	DeleteTweet(id int64, trimUser bool) (tweet anaconda.Tweet, err error)
	Unfavorite(id int64) (rt anaconda.Tweet, err error)
	UnRetweet(id int64, trimUser bool) (rt anaconda.Tweet, err error)
}

type Twitter struct {
	api TwitterAPI
}

func (t *Twitter) Self() (int64, error) {
	user, err := t.api.GetSelf(nil)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (t *Twitter) GetUserTimeline(limitTweetID int64) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Set("include_rts", "1")
	v.Set("count", "200")
	v.Set("trim_user", "true")

	// if 0, will return latest 200 tweets.
	// otherwise, reteurn 200 before specified tweetID.
	if limitTweetID != 0 {
		v.Set("max_id", strconv.Itoa(int(limitTweetID)))
	}

	return t.api.GetUserTimeline(v)
}

func (t *Twitter) GetUserFavorites(limitTweetID int64) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Set("count", "200")
	v.Set("include_entities", "false")

	// if 0, will return latest 200 tweets.
	// otherwise, return 200 before specified tweetID.
	if limitTweetID != 0 {
		v.Set("max_id", strconv.Itoa(int(limitTweetID)))
	}

	return t.api.GetFavorites(v)
}

func (t *Twitter) Delete(tweetID int64) error {
	_, err := t.api.DeleteTweet(tweetID, true)
	return err
}

func (t *Twitter) UnFavorite(tweetID int64) error {
	_, err := t.api.Unfavorite(tweetID)
	return err
}

func (t *Twitter) UnRetweet(tweetID int64) error {
	_, err := t.api.UnRetweet(tweetID, true)
	return err
}
