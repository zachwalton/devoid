package config

const (
	ReasonerOllama Reasoner = "ollama"
)

type (
	Reasoner string

	Config struct {
		Prompt      string `mapstructure:"prompt"`
		ProjectPath string `mapstructure:"project-path"`
		LLM         LLM    `mapstructure:"llm"`
	}

	LLM struct {
		Type        Reasoner `mapstructure:"type"`
		Model       string   `mapstructure:"model"`
		Temperature float64  `mapstructure:"temperature"`
	}
)
