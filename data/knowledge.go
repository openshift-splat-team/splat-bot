package data

import "github.com/expr-lang/expr/vm"

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

	// WatchThreads when true, the bot will apply this knowledge in a thread.
	// By default, the bot only watches channel level messages to see if it can
	// help.  This is intended to prevent the bot from posting multiple-times in a thread.
	// Additionally, Knowledge responses are only intended to provide an initial
	// touchpoint for a user to get more information.  If the user needs more
	// information, they can ask for it or we'll eventually check the channel.
	// This is a way to prevent the bot from being overly verbose aand spamming a thread.
	WatchThreads bool `yaml:"respond_in_threads"`

	// channels messages arriving on these channels will automatically have platform tokens
	// satisfied.
	ChannelContext *ChannelContext `yaml:"channel_context"`

	// ShouldMatch is a list of strings that should match
	ShouldMatch []string `yaml:"should_match"`

	// ShouldntMatch is a list of strings that shouldnt match
	ShouldntMatch []string `yaml:"shouldnt_match"`

	// RequireInChannel the attribute will only be recognized in a given channel(s).
	RequireInChannel []string `yaml:"must_be_in_channels"`
}

type ChannelContext struct {
	// contextPath is the path context to satisfy
	ContextPath string `yaml:"context_path"`

	// channels messages arriving on these channels will automatically have platform tokens
	Channels []string `yaml:"channels"`
}

type TokenMatch struct {
	Type         string       `yaml:"type"`
	Tokens       []string     `yaml:"tokens"`
	Terms        []TokenMatch `yaml:"terms"`
	CompiledExpr *vm.Program
	Expr         string `yaml:"expr"`
	Satisfied    bool
}
