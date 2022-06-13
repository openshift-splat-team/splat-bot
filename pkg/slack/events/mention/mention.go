package mention

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/openshift-eng/splat-sandbox/pkg/slack/events"
)

const (
	defaultResponse = "I'm just a bot, if you want to reach the SPLAT team you can message @splat-team, reach out in #forum-splat, or open a card on our <https://issues.redhat.com/secure/RapidBoard.jspa?projectKey=SPLAT&rapidView=5962|board>. "
)

type messagePoster interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

// Handler returns a handler that knows how to respond to
// new messages that mention the robot by showing users
// which interactive workflows they might be interested in,
// based on the phrasing that they used to mention the bot.
func Handler(client messagePoster) events.PartialHandler {
	return events.PartialHandlerFunc("mention", func(callback *slackevents.EventsAPIEvent, logger *logrus.Entry) (handled bool, err error) {
		if callback.Type != slackevents.CallbackEvent {
			return false, nil
		}
		event, ok := callback.InnerEvent.Data.(*slackevents.AppMentionEvent)
		if !ok {
			return false, nil
		}
		logger.Info("Handling app mention...")
		timestamp := event.TimeStamp
		if event.ThreadTimeStamp != "" {
			timestamp = event.ThreadTimeStamp
		}
		responseChannel, responseTimestamp, err := client.PostMessage(event.Channel, slack.MsgOptionBlocks(responseFor(event.Text)...), slack.MsgOptionTS(timestamp))
		if err != nil {
			logger.WithError(err).Warn("Failed to post response to app mention")
		} else {
			logger.Infof("Posted response to app mention in channel %s at %s", responseChannel, responseTimestamp)
		}
		return true, err
	})
}

func responseFor(message string) []slack.Block {

	var blocks []slack.Block

	if len(blocks) == 0 {
		blocks = append(blocks, &slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				//Text: "Sorry, I don't know how to help with that. Here are all the things I know how to do:",
				Text: defaultResponse,
			},
		})
	}
	return blocks
}
