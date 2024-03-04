package commands

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Callback func(eventsAPIEvent slackevents.EventsAPIEvent) ([]slack.MsgOption, error)

// Attributes define when and how a
type Attributes struct {
	Regex          string
	compiledRegex  regexp.Regexp
	RequiredArgs   int64
	Channels       []string
	Callback       Callback
	Rank           int64
	RequireMention bool
}

var attributes = []Attributes{}

func Initialize() {
	attributes = append(attributes, CreateAttributes)
	attributes = append(attributes, HelpAttributes)

	for idx, attribute := range attributes {
		attributes[idx].compiledRegex = *regexp.MustCompile(attribute.Regex)
	}
}

func Handler(client *socketmode.Client, evt slackevents.EventsAPIEvent) error {
	switch evt.Type {
	case "message":
	case "event_callback":
	default:
		return nil
	}

	msg := evt.InnerEvent.Data.(*slackevents.MessageEvent)
	if len(msg.BotID) > 0 {
		// throw away bot messages
		return nil
	}

	isAppMention := slackevents.EventsAPIType(reflect.TypeOf(evt.InnerEvent.Data).String()) == slackevents.AppMention

	for _, attribute := range attributes {
		if attribute.RequireMention && !isAppMention {
			fmt.Printf("requires mention\n")
			continue
		}

		spew.Dump(attribute.compiledRegex)
		if attribute.compiledRegex.Match([]byte(msg.Text)) {
			response, err := attribute.Callback(evt)
			if err != nil {
				fmt.Printf("failed processing message: %v", err)
			}
			if len(response) > 0 {
				_, _, err = client.PostMessage(msg.Channel, response...)
				if err != nil {
					fmt.Printf("failed responding to message: %v", err)
				}
			}
		}
	}

	return nil
}
