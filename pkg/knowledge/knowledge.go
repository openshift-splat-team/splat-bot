package knowledge

import "github.com/openshift-splat-team/splat-bot/pkg/commands"


const (
	DEFAULT_URL_PROMPT = `This may be a topic that I can help with. Check out these URLs:`
	DEFAULT_LLM_PROMPT = `Can you provide a short response that attempts to answer this question: `
)

type Knowledge struct {
	commands.Attributes

	// MarkdownPrompt message that is returned when the prompt matches
	MarkdownPrompt string

	// URLS urls to be appended to a response. if MarkdownPrompt isn't defined, URLS will be
	// attached to a reasonable default message.
	URLS []string

	// when true, the message is sent to an LLM to construct an answer.
	InvokeLLM bool
}