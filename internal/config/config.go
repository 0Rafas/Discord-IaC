package config

type ObservabilityConfig struct {
	Prometheus bool `yaml:"prometheus"`
	Tracing    bool `yaml:"tracing"`
}

type ConcurrencyConfig struct {
	WorkerPoolSize int `yaml:"worker_pool_size"`
}

type DeploymentConfig struct {
	Kubernetes    bool `yaml:"kubernetes"`
	GithubActions bool `yaml:"github_actions"`
}

// BotConfig represents the bot.yaml structure
type BotConfig struct {
	Name     string   `yaml:"name"`
	Prefix   string   `yaml:"prefix"`
	Intents  []string `yaml:"intents"`
	Sharding bool     `yaml:"sharding"`
	Database string   `yaml:"database"`
	Cache    string   `yaml:"cache"`
	Features []string `yaml:"features"`

	// Enterprise Features
	Observability ObservabilityConfig `yaml:"observability"`
	Concurrency   ConcurrencyConfig   `yaml:"concurrency"`
	Deployment    DeploymentConfig    `yaml:"deployment"`
}

// DefaultBotConfig returns a sample configuration
func DefaultBotConfig() *BotConfig {
	return &BotConfig{
		Name:     "TitanBot",
		Prefix:   "!",
		Intents:  []string{"Guilds", "GuildMessages", "MessageContent"},
		Sharding: true,
		Database: "postgres",
		Cache:    "redis",
		Features: []string{"slash_commands", "raw_websocket_events"},
		Observability: ObservabilityConfig{
			Prometheus: true,
			Tracing:    true,
		},
		Concurrency: ConcurrencyConfig{
			WorkerPoolSize: 100,
		},
		Deployment: DeploymentConfig{
			Kubernetes:    true,
			GithubActions: true,
		},
	}
}
