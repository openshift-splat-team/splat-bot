package main

import (
	"fmt"
	"log"
	"os"

	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	slackutil "github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack/slackevents"
)

func main() {
	log.SetOutput(os.Stdout)

	client, err := slackutil.GetClient()
	if err != nil {
		fmt.Printf("unable to get slack client: %v", err)
		os.Exit(1)
	}

	err = commands.Initialize(client)
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
				err = commands.Handler(client, eventsAPIEvent)
				if err != nil {
					log.Printf("error encountered while processing event: %v", err)
				}
			case socketmode.EventTypeInteractive:

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
