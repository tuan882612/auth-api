package databases

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewPostgres(pgUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	log.Info().Msg("postgres connected...")
	return pool, nil
}
