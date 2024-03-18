package data

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Callback func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error)

type MessageOfInterest func(args []string, attribute Attributes) bool

// Attributes define when and how to handle a message
type Attributes struct {
	// Commands when matched, the Callback is invoked.
	Commands []string
	// MessageOfInterest supercedes Commands. If MessageOfInterest is set, Commands is ignored.  This is useful for more complex matching.
	MessageOfInterest MessageOfInterest
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
	// This command will not be included in the help message.
	ExcludeFromHelp bool
	// DontGlobQuotes when true, quotes are not globbed.  This is useful for knowledge commands that need discrete tokens.
	DontGlobQuotes bool
	// RespondInChannel responds in the channel to the user. If false, responds in a thread.
	RespondInChannel bool
	// ResponseIsEphemeral specifies if the response should be ephemeral.
	ResponseIsEphemeral bool
}
