package knowledge

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack/slackevents"
	"gopkg.in/yaml.v3"
)

const (
	simpleLogicYaml = `name: test
markdown: "test-prompt"
urls: ["test-url"]
on:
  type: "and"
  terms:
  - type: "or"
    tokens:
    - "virtx"
    - "virty"
  - type: "or"
    tokens:
    - "arm"
    - "x86"    
`
)

type testCase struct {
	name             string
	yamlSpec         string
	simulatedMessage []string
	expectedMatch    bool
}

var (
	tests = []testCase{
		{
			name:             "match 'and' with 2 'or' terms - matches",
			yamlSpec:         simpleLogicYaml,
			simulatedMessage: []string{"i have a problem with virtx on arm"},
			expectedMatch:    true,
		},
		{
			name:             "match 'and' with 2 'or' terms - mismatches",
			yamlSpec:         simpleLogicYaml,
			simulatedMessage: []string{"i have a problem with vitx on arm"},
			expectedMatch:    false,
		},
	}
)

func TestStripPunctuation(t *testing.T) {
	stripped := util.StripPunctuation("\"\"install-config?\"")
	if stripped != "install-config" {
		t.Fatalf("beginning and trailing puncuation are not removed: %s", stripped)
	}
}

func TestYamlLogic(t *testing.T) {
	var asset data.KnowledgeAsset
	err := yaml.Unmarshal([]byte(simpleLogicYaml), &asset)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, msg := range test.simulatedMessage {
				if IsStringMatch(asset, msg) != test.expectedMatch {
					matchExpected := "expected to match"
					if !test.expectedMatch {
						matchExpected = "expected not to match"
					}
					t.Fatalf("test: %s\n\n%s\n\n, got: %v", test.name, matchExpected, test.simulatedMessage)
				}
			}
		})
	}
}

func TestModelLoading(t *testing.T) {
	ctx := context.TODO()
	assets := knowledgeAssets
	slackClient = &util.StubInterface{}
	for _, asset := range assets {
		t.Run(asset.Name, func(t *testing.T) {
			channelName := "test"
			if len(asset.ShouldMatch) == 0 {
				t.Fatalf("unable to test knowledge prompt: expected at least one should_match")
				return
			}
			if len(asset.ShouldntMatch) == 0 {
				t.Fatalf("unable to test knowledge prompt: expected at least one shouldnt_match")
				return
			}

			if asset.ChannelContext != nil {
				if len(asset.ChannelContext.Channels) > 0 {
					channelName = asset.ChannelContext.Channels[0]
				}
			}
			for _, should := range asset.ShouldMatch {
				tokens := util.NormalizeTokensToSlice(strings.Split(should, " "))

				msgEvent := &slackevents.MessageEvent{
					Text:    should,
					Channel: channelName,
				}
				responses, err := defaultKnowledgeHandler(ctx, tokens, msgEvent)
				if err != nil || len(responses) == 0 {
					dump := DumpMatchTree(asset.On, nil, nil)
					t.Fatalf("expected to match: %s\nOn: %s", should, strings.Join(dump, "\n"))
					return
				}
				if !asset.WatchThreads {
					response, err := defaultKnowledgeHandler(ctx, tokens, &slackevents.MessageEvent{
						ThreadTimeStamp: time.Now().String(),
					})
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
						return
					}
					if len(response) > 0 {
						t.Fatalf("expected no response when not watching threads")
						return
					}
				}
			}
			for _, shouldnt := range asset.ShouldntMatch {
				tokens := util.NormalizeTokensToSlice(strings.Split(shouldnt, " "))
				msgEvent := &slackevents.MessageEvent{
					Text:    shouldnt,
					Channel: "random",
				}
				_, err := defaultKnowledgeHandler(ctx, tokens, msgEvent)
				if err != nil {
					dump := DumpMatchTree(asset.On, nil, nil)
					t.Fatalf("On: %s", strings.Join(dump, "\n"))
					t.Fatalf("expected not to match: %s\nOn: %s", shouldnt, strings.Join(dump, "\n"))

					return
				}
			}
		})
	}
}
