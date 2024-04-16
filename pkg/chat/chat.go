package chat

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tmc/langchaingo/llms"
)

// TO-DO: provide context for each thread in a DM with the bot
func HandleChatInteraction(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent) ([]slack.MsgOption, error) {
	url := util.GetThreadUrl(evt)
	promptContext := []llms.MessageContent{}
	if len(url) > 0 {
		more := true
		var nextCursor string
		var err error
		msgs := []slack.Message{}
		// AddToContext
		// TO-DO: build context from thread
		for more {
			msgs, more, nextCursor, err = client.GetConversationReplies(&slack.GetConversationRepliesParameters{
				ChannelID: evt.Channel,
				Timestamp: evt.ThreadTimeStamp,
				Cursor:    nextCursor,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to get conversation replies: %v", err)
			}
			for _, msg := range msgs {
				role := "generic"
				if len(msg.BotID) > 0 {
					role = "system"
				}
				promptContext = util.AddToContext(role, msg.Text, promptContext)
			}
		}
		spew.Dump(promptContext)
	}

	response, err := util.GenerateResponse(ctx, evt.Text, promptContext...)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	msgOptions := util.StringsToBlockUnfurl([]string{response}, false, false)
	if len(evt.ThreadTimeStamp) > 0 {
		msgOptions = append(msgOptions, slack.MsgOptionTS(evt.ThreadTimeStamp))
	} else {
		msgOptions = append(msgOptions, slack.MsgOptionTS(evt.TimeStamp))
	}

	return msgOptions, nil
}
