package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Callback func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error)

// Attributes define when and how to handle a message
type Attributes struct {
	// Commands when matched, the Callback is invoked.
	Commands []string
	// The number of arguments a command must have. var args are not supported.
	RequiredArgs int
	// Callback function called when the attributes are met
	Callback Callback
	// Rank: Future - in a situation where multiple regexes match, this allows a priority to be assigned.
	Rank int64
	// RequireMention when true, @splat-bot must be used to invoke the command.
	RequireMention bool
	// HelpMarkdown is markdown that is contributed with the bot shows help.
	HelpMarkdown string
	// RespondInDM responds in a DM to the user.
	RespondInDM bool
	// MustBeInThread the attribute will only be recognized in a thread.
	MustBeInThread bool
	// AllowNonSplatUsers by default, only members of @splat-team can interact with the bot
	AllowNonSplatUsers bool
}

var (
	attributes   = []Attributes{}
	allowedUsers = map[string]bool{}
)

// AddCommand adds a handler to the list of handlers
func AddCommand(attribute Attributes) {
	log.Printf("adding command: %v", attribute.Commands)
	attributes = append(attributes, attribute)
}

func Initialize(client *socketmode.Client) error {
	attributes = append(attributes, CreateSummaryAttributes)
	attributes = append(attributes, CreateAttributes)
	attributes = append(attributes, SummarizeAttributes)
	attributes = append(attributes, HelpAttributes)
	attributes = append(attributes, UnsizedAttributes)
	attributes = append(attributes, ProwAttributes)
	attributes = append(attributes, ProwGraphAttributes)

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
func tokenize(msgText string) []string {
	var tokens []string
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

		args := tokenize(msg.Text)
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

func checkForCommand(args []string, attribute Attributes) bool {
	match := true
	for index, command := range attribute.Commands {
		if command != args[index] {
			match = false
			break
		}
	}
	return match
}
