package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	attributes   = []data.Attributes{}
	allowedUsers = map[string]bool{}
)

// AddCommand adds a handler to the list of handlers. Matching of the message can be overriden
// by providing a MessageOfInterest function.
func AddCommand(attribute data.Attributes, handler ...data.MessageOfInterest) {
	log.Printf("adding command: %v", attribute.Commands)
	if len(handler) > 0 {
		attribute.MessageOfInterest = handler[0]
	} else {
		attribute.MessageOfInterest = checkForCommand
	}
	attributes = append(attributes, attribute)
}

func Initialize(client *socketmode.Client) error {
	AddCommand(CreateSummaryAttributes)
	AddCommand(CreateAttributes)
	AddCommand(SummarizeAttributes)
	AddCommand(HelpAttributes)
	AddCommand(UnsizedAttributes)
	AddCommand(ProwAttributes)
	AddCommand(ProwGraphAttributes)
	AddCommand(ProviderSummaryAttributes)

	allowed := os.Getenv("SLACK_ALLOWED_USERS")
	if len(allowed) == 0 {
		log.Printf("no allowed users specified with SLACK_ALLOWED_USERS. some commands may not work.")
	}
	allowedUsersIDs := strings.Split(allowed, ",")
	for _, user := range allowedUsersIDs {
		allowedUsers[user] = true
	}
	return nil
}

func isAllowedUser(evt *slackevents.MessageEvent) error {
	if _, found := allowedUsers[evt.User]; !found {
		return errors.New("user not allowed")
	}
	return nil
}
func tokenize(msgText string, glob bool) []string {
	var tokens []string
	if glob {
		re := regexp.MustCompile(`"([^"]*?)"|(\S+)`)
		matches := re.FindAllStringSubmatch(msgText, -1)

		for _, match := range matches {
			if match[1] != "" {
				// Remove leading and trailing quotation marks
				tokens = append(tokens, strings.Trim(match[1], "\""))
			} else {
				tokens = append(tokens, match[2])
			}
		}
		return tokens
	} else {
		return strings.Split(msgText, " ")
	}
}

func getDMChannelID(client *socketmode.Client, evt *slackevents.MessageEvent) (string, error) {
	user := evt.User
	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user},
	})
	if err != nil {
		return "", fmt.Errorf("failed to open conversation: %v", err)
	}

	return channel.Latest.Channel, nil
}

func Handler(ctx context.Context, client *socketmode.Client, evt slackevents.EventsAPIEvent) error {
	switch evt.Type {
	case "message":
	case "event_callback":
	default:
		return nil
	}

	msg := &slackevents.MessageEvent{}
	switch ev := evt.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		appMentionEvent := evt.InnerEvent.Data.(*slackevents.AppMentionEvent)
		msg = &slackevents.MessageEvent{
			Channel:         appMentionEvent.Channel,
			User:            appMentionEvent.User,
			Text:            appMentionEvent.Text,
			TimeStamp:       appMentionEvent.TimeStamp,
			ThreadTimeStamp: appMentionEvent.ThreadTimeStamp,
		}
	case *slackevents.MessageEvent:
		msg = evt.InnerEvent.Data.(*slackevents.MessageEvent)
	default:
		return fmt.Errorf("received an unknown event type: %T", ev)
	}

	if len(msg.BotID) > 0 {
		// throw away bot messages
		return nil
	}

	for _, attribute := range attributes {
		if attribute.RequireMention && !ContainsBotMention(msg.Text) {
			continue
		}

		if !attribute.AllowNonSplatUsers {
			err := isAllowedUser(msg)
			if err != nil {
				return fmt.Errorf("user not allowed: %v", err)
			}
		}

		args := tokenize(msg.Text, !attribute.DontGlobQuotes)
		if attribute.RequireMention {
			args = args[1:]
		}

		if checkForCommand(args, attribute) {
			var response []slack.MsgOption
			var err error
			inThread := len(GetThreadUrl(msg)) > 0
			if attribute.MustBeInThread && !inThread {
				continue
			}
			if len(args) < attribute.RequiredArgs {
				response = []slack.MsgOption{
					slack.MsgOptionText(fmt.Sprintf("command requires %d arguments.\n%s\n", attribute.RequiredArgs, attribute.HelpMarkdown), true),
				}
			} else if attribute.RequiredArgs > 0 && len(args) > attribute.RequiredArgs {
				response = []slack.MsgOption{
					slack.MsgOptionText(fmt.Sprintf("command requires %d arguments. if an argument is greater than one word, be sure to wrap that argument in quotes.\n%s\n", attribute.RequiredArgs, attribute.HelpMarkdown), true),
				}
			} else {
				response, err = attribute.Callback(ctx, client, msg, args)
				if err != nil {
					fmt.Printf("failed processing message: %v", err)
				}
			}
			if len(response) > 0 {
				if attribute.RespondInDM {
					channelID, err := getDMChannelID(client, msg)
					if err != nil {
						fmt.Printf("failed getting channel ID: %v", err)
					}
					msg.Channel = channelID
				} else if len(GetThreadUrl(msg)) > 0 {
					response = append(response, slack.MsgOptionTS(msg.ThreadTimeStamp))
				}
				_, _, err = client.PostMessage(msg.Channel, response...)
				if err != nil {
					fmt.Printf("failed responding to message: %v", err)
				}
				return nil
			}
		}
	}

	return nil
}

func checkForCommand(args []string, attribute data.Attributes) bool {
	match := true
	for index, command := range attribute.Commands {
		if command != args[index] {
			match = false
			break
		}
	}
	return match
}
