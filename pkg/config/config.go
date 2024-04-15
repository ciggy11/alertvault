package config

import (
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/ciggy11/alertvault/pkg/db"
)

type Config struct {
	HTTPListenAddress string `yaml:"http_listen_address"`
	LogLevel          string `yaml:"log_level"`
	Backend           string `default:"redis" yaml:"backend"`
	Tenant            struct {
		InLabel      bool   `yaml:"in_label"`
		InAnnotation bool   `yaml:"in_annotation"`
		Label        string `yaml:"label"`
		Annotation   string `yaml:"annotation"`
		UniqueName   string `yaml:"unique_name"`
		Header       string `yaml:"header"`
	}
	VaultDB db.Config `yaml:"vaultdb"`
}

func LoadConfig(file string) (*Config, error) {
	cfg := &Config{}
	if file != "" {
		y, err := os.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config file")
		}
		if err := yaml.UnmarshalStrict(y, cfg); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal config file")
		}
	}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "Unable to parse env vars")
	}

	// Set default values
	if cfg.HTTPListenAddress == "" {
		cfg.HTTPListenAddress = "0.0.0.0:8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.Backend == "" {
		cfg.Backend = "redis"
	}
	if cfg.Tenant.Label == "" {
		cfg.Tenant.Label = "tenantID"
	}
	if cfg.Tenant.Annotation == "" {
		cfg.Tenant.Annotation = "tenantID"
	}
	if cfg.Tenant.UniqueName == "" {
		cfg.Tenant.UniqueName = "fingerprint"
	}
	if cfg.Tenant.Header == "" {
		cfg.Tenant.Header = "X-Scope-OrgID"
	}
	return cfg, nil
}
