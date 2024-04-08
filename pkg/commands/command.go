package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	attributeMu  sync.Mutex
	attributes   = []data.Attributes{}
	allowedUsers = map[string]bool{}
)

// AddCommand adds a handler to the list of handlers. Matching of the message can be overriden
// by providing a MessageOfInterest function.
func AddCommand(attribute data.Attributes, handler ...data.MessageOfInterest) {
	attributeMu.Lock()
	defer attributeMu.Unlock()
	log.Printf("adding command: %v", attribute.Commands)
	if len(handler) > 0 {
		attribute.MessageOfInterest = handler[0]
	} else {
		attribute.MessageOfInterest = checkForCommand
	}
	attributes = append(attributes, attribute)
}

func getAttributes() []data.Attributes {
	attributeMu.Lock()
	defer attributeMu.Unlock()

	newAttributes := make([]data.Attributes, len(attributes))

	copy(newAttributes, attributes)
	return newAttributes
}

func init() {
	AddCommand(CreateSummaryAttributes)
	AddCommand(CreateAttributes)
	AddCommand(SummarizeAttributes)
	AddCommand(HelpAttributes)
	AddCommand(UnsizedAttributes)
	AddCommand(ProwAttributes)
	AddCommand(ProwGraphAttributes)
	AddCommand(ProviderSummaryAttributes)
	AddCommand(CreateJiraWithThreadAttributes)
}

func Initialize() error {
	// TODO:  Global allowed users means we cannot make some actions available to some users while others not.  This could
	//        be beefed up in the future to be allowed users per command from config provided by a yaml file or something of
	//        that nature.
	allowed := os.Getenv("SLACK_ALLOWED_USERS")
	if len(allowed) == 0 {
		log.Printf("Disabling user enforcement.  Please configure SLACK_ALLOWED_USERS if you wish to enforce allowed users on certain commands.")
	} else {
		allowedUsersIDs := strings.Split(allowed, ",")
		for _, user := range allowedUsersIDs {
			allowedUsers[user] = true
		}
	}
	return nil
}

func isAllowedUser(evt *slackevents.MessageEvent) error {
	fmt.Printf("User size: %d\n", len(allowedUsers))
	if _, found := allowedUsers[evt.User]; !found && len(allowedUsers) > 0 {
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

func getDMChannelID(client util.SlackClientInterface, evt *slackevents.MessageEvent) (string, error) {
	user := evt.User
	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user},
	})
	if err != nil {
		return "", fmt.Errorf("failed to open conversation: %v", err)
	}

	return channel.Latest.Channel, nil
}

func Handler(ctx context.Context, client util.SlackClientInterface, evt slackevents.EventsAPIEvent) error {
	isAppMentionEvent := false

	switch evt.Type {
	case "message":
	case "event_callback":
	default:
		return nil
	}

	msg := &slackevents.MessageEvent{}
	switch ev := evt.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		isAppMentionEvent = true
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

	if len(msg.BotID) > 0 && IsSPLATBotID(msg.BotID) {
		log.Printf("throwing away message from bot: %s", msg.BotID)
		// throw away bot messages
		return nil
	}

	for _, attribute := range getAttributes() {
		log.Printf("Checking command: %v", attribute.Commands)

		// For app mention logic, there are two scenarios: 1.) Channel Msg.  2.) Direct Message
		// For Channel messages, we want the event to be an AppMention if attribute.RequireMention.
		// For Direct messages, we will want event to be Message, Channel = "im", and ContainsBotMention
		// Note, for AppMessage, InnerEvent is AppMessageEvent, for Message, its MessageEvent.
		if attribute.RequireMention {
			if isAppMentionEvent && !ContainsBotMention(msg.Text) {
				continue
			} else if !isAppMentionEvent {
				ieData := evt.InnerEvent.Data.(*slackevents.MessageEvent)
				channelType := ieData.ChannelType
				if channelType == slack.TYPE_CHANNEL || (channelType == slack.TYPE_IM && !ContainsBotMention(msg.Text)) {
					continue
				}
			}
		}

		if len(attribute.RequireInChannel) > 0 {
			allowedInChannel := false
			for _, channel := range attribute.RequireInChannel {
				if allowedInChannel = channel == msg.Channel; allowedInChannel {
					break
				}
			}
			if !allowedInChannel {
				continue
			}
		}

		args := tokenize(msg.Text, !attribute.DontGlobQuotes)
		if ContainsBotMention(msg.Text) {
			args = args[1:]
		}

		if checkForCommand(args, attribute, msg.Channel) {
			log.Printf("Found command: %v", attribute.Commands)
			// Now that we found command, make sure it can be used by current user.
			if !attribute.AllowNonSplatUsers {
				err := isAllowedUser(msg)
				if err != nil {
					return fmt.Errorf("user not allowed: %v", err)
				}
			}

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
					log.Printf("failed processing message: %v", err)
				}
			}
			if len(response) > 0 {
				log.Printf("responding to message: %v", response)
				if attribute.RespondInDM {
					channelID, err := getDMChannelID(client, msg)
					if err != nil {
						fmt.Printf("failed getting channel ID: %v", err)
					}
					msg.Channel = channelID
				} else if !attribute.RespondInChannel {
					response = append(response, slack.MsgOptionTS(msg.TimeStamp))
				} else if len(GetThreadUrl(msg)) > 0 {
					response = append(response, slack.MsgOptionTS(msg.ThreadTimeStamp))
				}

				log.Printf("responding to message in channel: %s", msg.Channel)
				if attribute.ResponseIsEphemeral {
					_, err = client.PostEphemeral(msg.Channel, msg.User, response...)
				} else {
					_, _, err = client.PostMessage(msg.Channel, response...)
				}
				if err != nil {
					return fmt.Errorf("failed responding to message: %v", err)
				}
				return nil
			}
			log.Printf("next\n")
		}
	}

	return nil
}

func checkForCommand(args []string, attribute data.Attributes, channel string) bool {
	match := true
	for index, command := range attribute.Commands {
		if command != args[index] {
			match = false
			break
		}
	}
	return match
}
