package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack/slackevents"
)

type AttributesTestCase struct {
	name       string
	attributes data.Attributes
}

var (
	defaultEvent = slackevents.EventsAPIEvent{
		Data: &slackevents.MessageEvent{},
	}
)

func buildAppMentionEvent(text, user, channel string, inThread bool) slackevents.EventsAPIEvent {
	timestamp := time.Now().String()
	threadedTimestamp := ""
	if inThread {
		threadedTimestamp = timestamp
	}

	return slackevents.EventsAPIEvent{
		Type: "message",
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Type: "message",
			Data: &slackevents.AppMentionEvent{
				Text:            text,
				Channel:         channel,
				User:            user,
				TimeStamp:       timestamp,
				ThreadTimeStamp: threadedTimestamp,
			},
		},
	}
}
func buildEvent(text, user, channel string, inThread bool) slackevents.EventsAPIEvent {
	timestamp := time.Now().String()
	threadedTimestamp := ""
	if inThread {
		threadedTimestamp = timestamp
	}

	return slackevents.EventsAPIEvent{
		Type: "message",
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Type: "message",
			Data: &slackevents.MessageEvent{
				Text:            text,
				Channel:         channel,
				User:            user,
				TimeStamp:       timestamp,
				ThreadTimeStamp: threadedTimestamp,
			},
		},
	}
}

const (
	SPLAT_BOT_USER_ID   = "testbot"
	SLACK_ALLOWED_USERS = "alloweduser1"
	POST_MESSAGE        = "failed responding to message: PostMessage"
	POST_EPHEMERAL      = "failed responding to message: PostEphemeral"
)

func checkRequireMention(tokens []string, client util.SlackClientInterface, attribute data.Attributes) error {
	ctx := context.TODO()
	msg := strings.Join(tokens, " ")
	if !attribute.RequireMention {
		return nil
	}
	if err := Handler(ctx, client, buildAppMentionEvent(msg, "test", "testchannel", false)); err != nil {
		return fmt.Errorf("expected no response when not mentioning bot: %v", err)
	}
	msg = fmt.Sprintf("<@%s> %s", SPLAT_BOT_USER_ID, msg)
	if err := Handler(ctx, client, buildAppMentionEvent(msg, "test", "testchannel", false)); err != nil {
		if attribute.ResponseIsEphemeral {
			if err.Error() != POST_EPHEMERAL {
				return fmt.Errorf("expected ephemeral response when mentioning bot")
			}
		} else if err.Error() != POST_MESSAGE {
			return fmt.Errorf("expected response when mentioning bot")
		}
	}
	return nil
}

func checkNotRequireMention(tokens []string, client util.SlackClientInterface, attribute data.Attributes) error {
	ctx := context.TODO()
	msg := strings.Join(tokens, " ")
	if attribute.RequireMention {
		return nil
	}
	if err := Handler(ctx, client, buildEvent(msg, "test", "testchannel", false)); err != nil {
		return fmt.Errorf("expected response: %v", err)
	}
	return nil
}

func TestHandler(t *testing.T) {
	mockClient := &util.StubInterface{}
	os.Setenv("SPLAT_BOT_USER_ID", SPLAT_BOT_USER_ID)
	os.Setenv("SLACK_ALLOWED_USERS", SLACK_ALLOWED_USERS)
	for _, attribute := range attributes {
		if err := checkRequireMention(attribute.Commands, mockClient, attribute); err != nil {
			t.Errorf("test failed for %v: %v", attribute.Commands, err)
		}
		if err := checkNotRequireMention(attribute.Commands, mockClient, attribute); err != nil {
			t.Errorf("test failed for %v: %v", attribute.Commands, err)
		}
	}
}

func TestCommands(t *testing.T) {
	attributes := getAttributes()

	var tests = []AttributesTestCase{}
	for _, attribute := range attributes {
		tests = append(tests, AttributesTestCase{name: strings.Join(attribute.Commands, " "), attributes: attribute})
	}

	t.Run("test commands", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if len(test.attributes.ShouldMatch) == 0 {
					t.Errorf("ShouldMatch is empty")
					return
				}
				if len(test.attributes.ShouldntMatch) == 0 {
					t.Errorf("ShouldntMatch is empty")
					return
				}
				for _, shouldMatch := range test.attributes.ShouldMatch {
					tokens := strings.Split(shouldMatch, " ")
					if !checkForCommand(tokens, test.attributes, "testchannel") {
						t.Errorf("Should have matched %s", shouldMatch)
					}
				}
				for _, shouldntMatch := range test.attributes.ShouldntMatch {
					tokens := strings.Split(shouldntMatch, " ")
					if checkForCommand(tokens, test.attributes, "testchannel") {
						t.Errorf("Shouldnt have matched %s", shouldntMatch)
					}
				}
			})
		}
	})
}
