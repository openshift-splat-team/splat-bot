package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	RSS_FEED_PROMPT = `Can you give me a summary of the following text: %s`
)

var (
	providers = map[string][]string{
		"vsphere": {
			"https://feeds.feedburner.com/vmwareblogsfeed",
		},

		"aws": {
			"https://aws.amazon.com/blogs/aws/feed/",
		},

		"azure": {
			"https://azurecomcdn.azureedge.net/en-us/updates/feed/",
		},

		"gcp": {
			"https://cloud.google.com/feeds/gcp-release-notes.xml",
		},
	}
)

func getFeedSummary(lastNDays int, provider string, additionalContext ...string) (string, error) {
	//ctx := context.TODO()
	now := time.Now()
	fiveDaysAgo := now.AddDate(0, 0, lastNDays*-1)

	var urls []string

	if _urls, exists := providers[strings.ToLower(provider)]; !exists {
		return "", fmt.Errorf("%s is not a supported provider", provider)
	} else {
		urls = _urls
	}

	var response strings.Builder

	for _, url := range urls {
		feedItems, err := util.ParseFeed(url)
		if err != nil {
			return "", fmt.Errorf("unable to parse feed: %v", err)
		}
		inDateRange := util.GetItemsBetweenDates(feedItems, fiveDaysAgo, now)

		for _, item := range inDateRange {
			// TO-DO: once we get a performant way to get the summaries, we can turn this on
			/*summary, err := util.GenerateResponse(ctx, item.Description)
			if err != nil {
				return "", fmt.Errorf("failure while getting response from LLM: %v", err)
			}*/
			//response.WriteString(fmt.Sprintf("Title: %s\nPublished: %s\nLink: %s\nDescription: %s\n\n", item.Title, item.Published, item.Link, summary))
			response.WriteString(fmt.Sprintf("%s: <%s|%s>\n\n", item.Published, item.Link, item.Title))
		}

	}

	return response.String(), nil
}

var ProviderSummaryAttributes = data.Attributes{
	Commands:       []string{"provider-summary"},
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		provider := strings.ToLower(args[1])

		if _, exists := providers[provider]; !exists {
			return nil, fmt.Errorf("%s is not a supported provider", provider)
		}

		summary, err := getFeedSummary(5, provider)
		if err != nil {
			return nil, fmt.Errorf("unable to get feed summary: %v", err)
		}

		return StringToBlock(summary, false), nil
	},
	RequiredArgs: 2,
	HelpMarkdown: "summarize RSS feeds for various providers: `provider-summary [aws|vsphere|gcp|azure]`",
}
