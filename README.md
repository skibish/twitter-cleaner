# Twitter Cleaner

![Test](https://github.com/skibish/twitter-cleaner/workflows/run%20tests/badge.svg)
![Release](https://github.com/skibish/twitter-cleaner/workflows/release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/skibish/twitter-cleaner)](https://goreportcard.com/report/github.com/skibish/twitter-cleaner)

Clean your Twitter.

Once in 24 hours all tweets that are older than half a year (by default) will be deleted, un-retweeted and un-favorited from your feed.

## Motivation

I observed that tweets liveability is a few days maximum.
Why then I need to store them all in my timeline?
For me Twitter is a place to interact with others and not to use it as a diary for archeologists.

If you would like to do regular cleanups too, this tool is for you.

## Install

Download binary from [releases page](https://github.com/skibish/twitter-cleaner/releases).

If you have Go:

```sh
go get github.com/skibish/twitter-cleaner
```

## Example

```sh
$ ./twitter-retention \
  -access-token aaa \
  -access-token-secret ttt \
  -consumer-key kkk \
  -consumer-secret xxx
2021/04/25 13:19:53 successfully started
2021/04/25 13:20:52 scanned through 44 timeline tweets
2021/04/25 13:20:52 DELETING    XXXXXXXXXXXXXX240
2021/04/25 13:20:52 UNRETWEETING        XXXXXXXXXXXXXX624
2021/04/25 13:20:53 scanned through 0 timeline tweets
2021/04/25 13:20:53 scanned through 199 favorite tweets
2021/04/25 13:20:54 scanned through 105 favorite tweets
2021/04/25 13:20:54 UNFAVORITING        XXXXXXXXXXXXXX416
2021/04/25 13:20:54 UNFAVORITING        XXXXXXXXXXXXXX183
2021/04/25 13:20:54 UNFAVORITING        XXXXXXXXXXXXXX100
2021/04/25 13:20:54 UNFAVORITING        XXXXXXXXXXXXXX225
2021/04/25 13:20:55 UNFAVORITING        XXXXXXXXXXXXXX508
2021/04/25 13:20:55 UNFAVORITING        XXXXXXXXXXXXXX317
^C
2021/04/25 13:21:20 shutdown
```

## Usage

```sh
./twitter-cleaner:
  -access-token string
        Access token
  -access-token-secret string
        Access token secret
  -check-interval duration
        Cleanup interval (default 24h0m0s)
  -consumer-key string
        Consumer key
  -consumer-secret string
        Consumer secret
  -dry-run
        Check that something can be deleted, no real deletion is made
  -tweet-age duration
        Tweets older than this duration will be deleted (default 4380h0m0s)
  -v    Show version
```

You can [find out here](https://developer.twitter.com/en/docs/basics/authentication/guides/access-tokens) how to create all needed tokens.

## Development

```sh
go get
```
