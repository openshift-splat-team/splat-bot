package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Callback func(evt *slackevents.MessageEvent , args []string) ([]slack.MsgOption, error)

// Attributes define when and how to handle a message
type Attributes struct {
	// Regex when matched, the Callback is invoked.
	Regex          string
	compiledRegex  regexp.Regexp
	// The number of arguments a command must have. var args are not supported.
	RequiredArgs   int
	// Callback function called when the attributes are met
	Callback       Callback
	// Rank: Future - in a situation where multiple regexes match, this allows a priority to be assigned.
	Rank           int64
	// RequireMention when true, @splat-bot must be used to invoke the command.
	RequireMention 		bool
	// HelpMarkdown is markdown that is contributed with the bot shows help.
	HelpMarkdown       string
}

var attributes = []Attributes{}

func Initialize() {
	attributes = append(attributes, CreateAttributes)
	attributes = append(attributes, HelpAttributes)
	attributes = append(attributes, UnsizedAttributes)

	for idx, attribute := range attributes {
		attributes[idx].compiledRegex = *regexp.MustCompile(attribute.Regex)
	}
}

func tokenize(msgText string) []string{
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

func Handler(client *socketmode.Client, evt slackevents.EventsAPIEvent) error {
	switch evt.Type {
	case "message":
	case "event_callback":
	default:
		return nil
	}

	msg := &slackevents.MessageEvent {}
	switch ev := evt.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		fmt.Println("Received an AppMentionEvent")
		appMentionEvent := evt.InnerEvent.Data.(*slackevents.AppMentionEvent)
		msg = &	slackevents.MessageEvent {
			Channel: appMentionEvent.Channel,
			User:    appMentionEvent.User,
			Text:    appMentionEvent.Text,
			TimeStamp:      appMentionEvent.EventTimeStamp,
			ThreadTimeStamp: appMentionEvent.ThreadTimeStamp,
		}
	case *slackevents.MessageEvent:
		fmt.Println("Received a MessageEvent")
		msg = evt.InnerEvent.Data.(*slackevents.MessageEvent)
	default:
		return fmt.Errorf("received an unknown event type: %T", ev)
	}

	if len(msg.BotID) > 0 {
		// throw away bot messages
		return nil
	}

	//isAppMention := slackevents.EventsAPIType(reflect.TypeOf(evt.InnerEvent.Data).String()) == slackevents.AppMention

	for _, attribute := range attributes {
		if attribute.RequireMention && !ContainsBotMention(msg.Text) {
			fmt.Printf("requires mention: %s\n", msg.Text)
			continue
		}

		if attribute.compiledRegex.Match([]byte(msg.Text)) {
			var response []slack.MsgOption
			var err error
			args := tokenize(msg.Text)
			if attribute.RequireMention {
				args = args[1:]
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
				response, err = attribute.Callback(msg, args)
				if err != nil {
					fmt.Printf("failed processing message: %v", err)
				}
			}
			if len(response) > 0 {
				if len(GetThreadUrl(msg)) > 0{
					response = append(response, slack.MsgOptionTS(msg.ThreadTimeStamp))
				}
				_, _, err = client.PostMessage(msg.Channel, response...)
				if err != nil {
					fmt.Printf("failed responding to message: %v", err)
				}
			}
		}
	}

	return nil
}
