package conf

type databases struct {
	PGURL string `env:"DATABASE_PGURL, required"`
}
