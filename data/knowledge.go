package data

// Knowledge defines a peice of knowledge that the bot can respond with
type Knowledge struct {
	Attributes

	// MarkdownPrompt message that is returned when the prompt matches
	MarkdownPrompt string

	// URLS urls to be appended to a response. if MarkdownPrompt isn't defined, URLS will be
	// attached to a reasonable default message.
	URLS []string

	// when true, the message is sent to an LLM to construct an answer.
	InvokeLLM bool
}
