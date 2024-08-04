package configurator

import (
	"encoding/json"
)

type Config struct {
	Preset   string   `json:"preset"`
	RawRules any      `json:"rules"`
	Exclude  []string `json:"exclude"`
	NotName  []string `json:"notName"`
	NotPath  []string `json:"notPath"`
}

const (
	defaultPreset = "psr12"
)

func NewConfig(content string) (*Config, error) {
	if content == "" {
		return &Config{Preset: defaultPreset}, nil
	}

	return Parse(content)
}

func Parse(configContent string) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal([]byte(configContent), config)

	if err != nil {
		return nil, err
	}

	if config.Preset == "" {
		config.Preset = defaultPreset
	}

	return config, nil
}

func (c *Config) GetRules() []Rule {
	var rules []Rule

	if c.RawRules == nil {
		return rules
	}

	for k, v := range c.RawRules.(map[string]any) {
		if val, ok := v.(bool); ok {
			rules = append(rules, Rule{rule: k, status: val})
		}

		if val, ok := v.(map[string]any); ok {
			rule := Rule{rule: k, status: true, settings: map[string]bool{}}

			for k2, v2 := range val {
				if boolVal, ok := v2.(bool); ok {
					rule.settings[k2] = boolVal
				}
			}
			rules = append(rules, rule)
		}

	}

	return rules
}
