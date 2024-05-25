package conf

type LLMConfig struct {
	OpenAI   OpenAICfg   `json:"openai"`
	Moonshot MoonshotCfg `json:"moonshot"`
}

type OpenAICfg struct {
	APIKey     string `toml:"api_key" json:"api_key"`
	TimeoutSec int    `toml:"timeout_sec" json:"timeout_sec"`
}

type MoonshotCfg struct {
	APIKey     string `toml:"api_key" json:"api_key"`
	Model      string `toml:"model" json:"model"`
	TimeoutSec int    `toml:"timeout_sec" json:"timeout_sec"`
	//使用什么采样温度，介于 0 和 1 之间。较高的值（如 0.7）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定性 , Default 0.3
	Temperature float32 `toml:"temperature" json:"temperature"`
}
