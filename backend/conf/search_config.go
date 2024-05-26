package conf

type SearchConfig struct {
	Serpers []*SerperCfg `yaml:"serpers" json:"serpers"`
}

type SerperCfg struct {
	APIKey     string
	TimeoutSec int
}
