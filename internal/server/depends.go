package server

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"authapi/internal/auth"
	"authapi/internal/auth/token"
	"authapi/internal/conf"
	"authapi/internal/databases"
	"authapi/internal/user"
)

type Depends struct {
	Config   *conf.Configuration
	Postgres *pgxpool.Pool
	Prov     *token.Provider
	AuthHdlr *auth.Handler
}

func NewDependencies() (*Depends, error) {
	cfg, err := conf.New()
	if err != nil {
		return nil, err
	}

	pg, err := databases.NewPostgres(cfg.Databases.PGURL)
	if err != nil {
		return nil, err
	}

	authRepo := user.NewRepository(pg)
	prov := token.NewProvider(cfg.Server.SECRET)
	authSvc := auth.NewSFAService(authRepo, prov)
	authHandler := auth.NewHandler(authSvc)

	return &Depends{
		Config:   cfg,
		Postgres: pg,
		Prov:     prov,
		AuthHdlr: authHandler,
	}, nil
}

func (d *Depends) Close() error {
	d.Postgres.Close()
	return nil
}
