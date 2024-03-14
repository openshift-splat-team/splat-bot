package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func GetClient() (*socketmode.Client, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return nil, errors.New("SLACK_APP_TOKEN must be set")

	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("SLACK_BOT_TOKEN must be set")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	err := util.BindEnvVars()
	if err != nil {
		return nil, fmt.Errorf("unable to bind env vars: %v", err)
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return client, nil
}
