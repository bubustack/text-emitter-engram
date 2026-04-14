package config

import "strings"

// Config captures text-emitter settings for Engram.spec.with.
type Config struct {
	Message string `json:"message" mapstructure:"message"`
	Enabled *bool  `json:"enabled" mapstructure:"enabled"`
	DelayMs int    `json:"delayMs" mapstructure:"delayMs"`
}

// Normalize applies defaults for text-emitter behavior.
func Normalize(cfg Config) Config {
	if strings.TrimSpace(cfg.Message) == "" {
		cfg.Message = "Hello! How can I help you today?"
	}
	if cfg.Enabled == nil {
		val := true
		cfg.Enabled = &val
	}
	if cfg.DelayMs < 0 {
		cfg.DelayMs = 0
	}
	return cfg
}
