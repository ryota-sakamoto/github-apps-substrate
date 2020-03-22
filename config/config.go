package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	GitHub struct {
		AppID      int64  `required:"true" envconfig:"GITHUB_APP_ID"`
		Secret     string `required:"true" envconfig:"GITHUB_WEBHOOK_SECRET"`
		PrivateKey string `required:"true" envconfig:"GITHUB_PRIVATE_KEY"`
	}
	DB struct {
		Host     string `required:"true" envconfig:"DB_HOST"`
		User     string `required:"true" envconfig:"DB_USER"`
		Password string `required:"true" envconfig:"DB_PASSWORD"`
	}
}

func GetConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}
