package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	RSS_FEED_PROMPT = `Can you give me a summary of the following text: %s`
)

var (
	providers = map[string][]string{
		"vsphere": {
			"https://feeds.feedburner.com/vmwarekbfeed",
		},

		"aws": {
			"https://aws.amazon.com/blogs/aws/feed/",
		},

		"azure": {
			"https://azurecomcdn.azureedge.net/en-us/updates/feed/",
		},

		"gcp": {
			"https://blog.google/rss/",
		},
	}
)

func getFeedSummary(lastNDays int, provider string, additionalContext ...string) ([]string, error) {
	ctx := context.TODO()
	now := time.Now()
	fiveDaysAgo := now.AddDate(0, 0, lastNDays*-1)
	blocks := []string{}
	var urls []string

	if _urls, exists := providers[strings.ToLower(provider)]; !exists {
		return nil, fmt.Errorf("%s is not a supported provider", provider)
	} else {
		urls = _urls
	}

	for _, url := range urls {
		feedItems, err := util.ParseFeed(url)
		if err != nil {
			return nil, fmt.Errorf("unable to parse feed: %v", err)
		}
		inDateRange := util.GetItemsBetweenDates(feedItems, fiveDaysAgo, now)
		for _, item := range inDateRange {
			summary, err := util.GenerateResponse(ctx, item.Description)
			if err != nil {
				return nil, fmt.Errorf("failure while getting response from LLM: %v", err)
			}
			blocks = append(blocks, fmt.Sprintf("<%s|%s>\n*Published*: %s\n*Generated Article Summary:* %s\n\n", item.Link, item.Title, item.Published, summary))
		}
	}

	return blocks, nil
}

var ProviderSummaryAttributes = data.Attributes{
	Commands:       []string{"provider-summary"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		provider := strings.ToLower(args[1])

		if _, exists := providers[provider]; !exists {
			return nil, fmt.Errorf("%s is not a supported provider", provider)
		}

		summary, err := getFeedSummary(5, provider)
		if err != nil {
			log.Printf("unable to get feed summary: %v", err)
			summary = append(summary, fmt.Sprintf("unavailable: %v", err))
		}

		return StringsToBlockUnfurl(summary, false, false), nil
	},
	RequiredArgs: 2,
	HelpMarkdown: "summarize RSS feeds for various providers: `provider-summary [aws|vsphere|gcp|azure]`",
	ShouldMatch: []string{
		"provider-summary aws",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
