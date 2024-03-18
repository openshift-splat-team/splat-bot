package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	RSS_FEED_PROMPT = `Can you give me a short summary of the attached RSS feed without mentioning article metadata: %s`
)

var (
	providers = map[string][]string{
		"vsphere": {
			"https://feeds.feedburner.com/vmwarekbfeed",
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
	ctx := context.TODO()
	now := time.Now()
	fiveDaysAgo := now.AddDate(0, 0, lastNDays*-1)

	var urls []string

	if _urls, exists := providers[strings.ToLower(provider)]; !exists {
		return "", fmt.Errorf("%s is not a supported provider", provider)
	} else {
		urls = _urls
	}

	var response string
	for _, url := range urls {
		feedItems, err := util.ParseFeed(url)
		if err != nil {
			return "", fmt.Errorf("unable to parse feed: %v", err)
		}
		inDateRange := util.GetItemsBetweenDates(feedItems, fiveDaysAgo, now)
		feedJSON, err := json.MarshalIndent(inDateRange, "", "  ")
		if err != nil {
			log.Fatalf("Error converting feed to JSON: %v", err)
		}
		prompt := fmt.Sprintf(RSS_FEED_PROMPT, feedJSON)
		response, err = util.GenerateResponse(ctx, prompt)
		if err != nil {
			return "", fmt.Errorf("failure while getting response from LLM: %v", err)
		}
	}

	return response, nil
}

var ProviderSummaryAttributes = commands.Attributes{
	Regex:          `provider-summary`,
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {

		return nil, nil
	},
	RequiredArgs: 2,
	HelpMarkdown: `summarizes RSS feeds for various providers. supported providers are: AWS, Azure, vSphere, GCP`,
}
