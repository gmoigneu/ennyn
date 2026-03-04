package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gmoigneu/ennyn/internal/server"
	"gopkg.in/yaml.v3"
)

// Service represents a single service entry in the config file.
type Service struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Command string `yaml:"command"`
}

// SplitCommand splits the command string into executable and arguments.
func (s Service) SplitCommand() (string, []string) {
	parts := strings.Fields(s.Command)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// Config represents the ennyn.yml project configuration.
type Config struct {
	Services []Service `yaml:"services"`
}

// Load reads and parses an ennyn.yml file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	seen := make(map[string]bool)
	for i, svc := range cfg.Services {
		if !server.ValidateHost(svc.Host) {
			return fmt.Errorf("service %d: invalid host %q: must be lowercase alphanumeric and hyphens, cannot start or end with a hyphen", i+1, svc.Host)
		}
		if svc.Port < 1 || svc.Port > 65535 {
			return fmt.Errorf("service %d (%s): port %d out of range: must be 1–65535", i+1, svc.Host, svc.Port)
		}
		if strings.TrimSpace(svc.Command) == "" {
			return fmt.Errorf("service %d (%s): command is required", i+1, svc.Host)
		}
		if seen[svc.Host] {
			return fmt.Errorf("duplicate host %q", svc.Host)
		}
		seen[svc.Host] = true
	}
	return nil
}
