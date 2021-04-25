package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	version string
	commit  string
	date    string
)

const (
	accessTokenName       = "access-token"
	accessTokenSecretName = "access-token-secret"
	consumerKeyName       = "consumer-key"
	consumerSecretName    = "consumer-secret"
	tweetAgeName          = "tweet-age"
	dryRunName            = "dry-run"
	checkIntervalName     = "check-interval"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	log.SetOutput(stdout)

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		accessToken       = flags.String(accessTokenName, "", "Access token")
		accessTokenSecret = flags.String(accessTokenSecretName, "", "Access token secret")
		consumerKey       = flags.String(consumerKeyName, "", "Consumer key")
		consumerSecret    = flags.String(consumerSecretName, "", "Consumer secret")
		tweetAge          = flags.Duration(tweetAgeName, time.Hour*4380, "Tweets older than this duration will be deleted")
		dryRun            = flags.Bool(dryRunName, false, "Check that something can be deleted, no real deletion is made")
		checkInterval     = flags.Duration(checkIntervalName, time.Hour*24, "Cleanup interval")
		showVersion       = flags.Bool("v", false, "Show version")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)
		return nil
	}

	if *accessToken == "" {
		return fmt.Errorf("%s can't be empty", accessTokenName)
	}
	if *accessTokenSecret == "" {
		return fmt.Errorf("%s can't be empty", accessTokenSecretName)
	}
	if *consumerKey == "" {
		return fmt.Errorf("%s can't be empty", consumerKeyName)
	}
	if *consumerSecret == "" {
		return fmt.Errorf("%s can't be empty", consumerSecretName)
	}

	tw := &Twitter{
		api: anaconda.NewTwitterApiWithCredentials(*accessToken, *accessTokenSecret, *consumerKey, *consumerSecret),
	}

	c := NewCleaner(tw, *tweetAge, *checkInterval, *dryRun)

	if err := c.Init(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	log.Println("successfully started")
	if *dryRun {
		log.Println("running in \"dry run\" mode")
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		log.Println("shutdown")
		c.Stop()
	}()

	return c.Start()
}
