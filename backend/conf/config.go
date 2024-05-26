package conf

type Config struct {
	ServerAddr   string        `yaml:"server_addr"`
	LogLevel     string        `yaml:"log_level"`
	LLMConfig    *LLMConfig    `yaml:"llm" json:"llm"`
	SearchConfig *SearchConfig `yaml:"search" json:"search"`
}
