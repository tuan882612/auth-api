package conf

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	Server    server
	Databases databases
}

// New creates a new Configuration instance
func New() (*Configuration, error) {
	cfg := &Configuration{}

	if err := envconfig.Process(context.Background(), cfg); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return cfg, nil
}
