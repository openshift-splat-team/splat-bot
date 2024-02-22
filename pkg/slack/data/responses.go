package data

type Response struct {
	Questions []string `yaml:"questions"`
	Response  string   `yaml:"response"`
	Channel   string   `yaml:"channel"`
}

type Responses struct {
	Responses []Response `yaml:"responses"`
	StopWords []string   `yaml:"stop-words"`
}
