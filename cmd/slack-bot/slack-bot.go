package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"

	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	_ "github.com/openshift-splat-team/splat-bot/pkg/controllers"
	_ "github.com/openshift-splat-team/splat-bot/pkg/knowledge"
	slackutil "github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack/slackevents"
)

func main() {
	ctx := context.TODO()
	log.SetOutput(os.Stdout)

	client, err := slackutil.GetClient()
	if err != nil {
		fmt.Printf("unable to get slack client: %v", err)
		os.Exit(1)
	}

	err = commands.Initialize()
	if err != nil {
		fmt.Printf("unable to get users in group")
		os.Exit(1)
	}

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)
				err = commands.Handler(ctx, client, eventsAPIEvent)
				if err != nil {
					log.Printf("error encountered while processing event: %v", err)
				}
			case socketmode.EventTypeInteractive:
				fmt.Printf("GOT INTERACTIVE EVENT: %v\n", evt)
				client.Ack(*evt.Request)

				data := evt.Data.(slack.InteractionCallback)

				// This outputs the event data for debugging
				buffer := bytes.NewBuffer([]byte{})
				if err := json.NewEncoder(buffer).Encode(data); err != nil {
					fmt.Printf("Error: %v", err)
				} else {
					fmt.Print(buffer.String())
				}

				// For now, only close if text of action
				if data.ActionCallback.BlockActions[0].Text.Text == "Close" {
					_, _, err = client.PostMessage(data.Channel.ID, slack.MsgOptionDeleteOriginal(data.ResponseURL))
					if err != nil {
						log.Printf("Error occurred handling interative event: %v", err)
					}
				}
			case socketmode.EventTypeSlashCommand:
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	err = client.Run()
	if err != nil {
		log.Printf("error encountered while running client: %v", err)
	}
}
