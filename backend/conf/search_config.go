package conf

type SearchConfig struct {
	Serper []SerperCfg `json:"search_engine"`
}

type SerperCfg struct {
	APIKey     string
	TimeoutSec int
}
