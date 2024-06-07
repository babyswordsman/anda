package conf

type SearchConfig struct {
	Serper *SerperCfg `yaml:"serper" json:"serper"`
}

type SerperCfg struct {
	APIKey     string
	TimeoutSec int
}
