package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/openshift/must-gather-clean/pkg/obfuscator"
	"github.com/openshift/must-gather-clean/pkg/schema"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var (
	obfuscators = []obfuscator.ReportingObfuscator{}
)

func init() {
	tracker := obfuscator.NewSimpleTracker()
	newObfuscator, err := obfuscator.NewIPObfuscator(schema.ObfuscateReplacementTypeConsistent, tracker)
	if err != nil {
		// if we can't create obfuscators we need to crash out asap
		log.Panicf("unable to create ip obfuscator: %v", err)
	}
	obfuscators = append(obfuscators, newObfuscator)
	newObfuscator, err = obfuscator.NewMacAddressObfuscator(schema.ObfuscateReplacementTypeConsistent, tracker)
	if err != nil {
		// if we can't create obfuscators we need to crash out asap
		log.Panicf("unable to create mac obfuscator: %v", err)
	}
	obfuscators = append(obfuscators, newObfuscator)
	newObfuscator, err = obfuscator.NewRegexObfuscator(`^(?!:\/\/)(?=.{1,255}$)((.{1,63}\.){1,127}(?![0-9]*$)[a-z0-9-]+\.?)$`, tracker)
	if err != nil {
		// if we can't create obfuscators we need to crash out asap
		log.Panicf("unable to create mac obfuscator: %v", err)
	}
	obfuscators = append(obfuscators, newObfuscator)
}

func GetClient() (*socketmode.Client, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return nil, errors.New("SLACK_APP_TOKEN must be set")

	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("SLACK_BOT_TOKEN must be set")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	err := util.BindEnvVars()
	if err != nil {
		return nil, fmt.Errorf("unable to bind env vars: %v", err)
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return client, nil
}

func AnonymizeMessages(msgs []slack.Message) []slack.Message {
	opMap := map[string]string{}

	for idx, msg := range msgs {
		text := msg.Text
		user := msg.Username
		if _, exists := opMap[user]; !exists {
			if len(opMap) == 0 {
				opMap[user] = "op"
			} else {
				opMap[user] = fmt.Sprintf("contributor_%d", len(opMap)+1)
			}
		}
		msgs[idx].Username = opMap[user]
		msgs[idx].User = opMap[user]
		for _, obfuscator := range obfuscators {
			text = obfuscator.Contents(text)
		}
		// replace any inline mentions of the user
		for orgName, obfuscatedName := range opMap {
			text = strings.ReplaceAll(text, fmt.Sprintf("<@%s>", orgName), fmt.Sprintf("<@%s>", obfuscatedName))
		}
		msgs[idx].Text = text
	}
	return msgs
}
