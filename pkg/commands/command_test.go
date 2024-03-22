package commands

import (
	"strings"
	"testing"

	"github.com/openshift-splat-team/splat-bot/data"
)

type AttributesTestCase struct {
	name       string
	attributes data.Attributes
}

func TestCommands(t *testing.T) {
	attributes := getAttributes()

	var tests = []AttributesTestCase{}
	for _, attribute := range attributes {
		tests = append(tests, AttributesTestCase{name: strings.Join(attribute.Commands, " "), attributes: attribute})
	}

	t.Run("test commands", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if len(test.attributes.ShouldMatch) == 0 {
					t.Errorf("ShouldMatch is empty")
					return
				}
				if len(test.attributes.ShouldntMatch) == 0 {
					t.Errorf("ShouldntMatch is empty")
					return
				}
				for _, shouldMatch := range test.attributes.ShouldMatch {
					tokens := strings.Split(shouldMatch, " ")
					if !checkForCommand(tokens, test.attributes) {
						t.Errorf("Should have matched %s", shouldMatch)
					}
				}
				for _, shouldntMatch := range test.attributes.ShouldntMatch {
					tokens := strings.Split(shouldntMatch, " ")
					if checkForCommand(tokens, test.attributes) {
						t.Errorf("Shouldnt have matched %s", shouldntMatch)
					}
				}
			})
		}
	})
}
