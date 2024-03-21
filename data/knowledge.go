package data

// Knowledge defines a peice of knowledge that the bot can respond with
type Knowledge struct {
	Attributes

	// MarkdownPrompt message that is returned when the prompt matches
	MarkdownPrompt string `yaml:"markdown"`

	// URLS urls to be appended to a response. if MarkdownPrompt isn't defined, URLS will be
	// attached to a reasonable default message.
	URLS []string `yaml:"urls"`

	// when true, the message is sent to an LLM to construct an answer.
	InvokeLLM bool `yaml:"invoke_llm"`
}

type KnowledgeAsset struct {
	// Name of the knowledge asset
	Name string `yaml:"name"`

	// MarkdownPrompt message that is returned when the prompt matches
	MarkdownPrompt string `yaml:"markdown"`

	// URLS urls to be appended to a response. if MarkdownPrompt isn't defined, URLS will be
	// attached to a reasonable default message.
	URLS []string `yaml:"urls"`

	// when true, the message is sent to an LLM to construct an answer.
	InvokeLLM bool `yaml:"invoke_llm"`

	// When the prompt is matched
	On TokenMatch `yaml:"on"`

	// ShouldMatch is a list of strings that should match
	ShouldMatch []string `yaml:"should_match"`

	// ShouldntMatch is a list of strings that shouldnt match
	ShouldntMatch []string `yaml:"shouldnt_match"`
}

type TokenMatch struct {
	Type   string       `yaml:"type"`
	Tokens []string     `yaml:"tokens"`
	Terms  []TokenMatch `yaml:"terms"`
}
