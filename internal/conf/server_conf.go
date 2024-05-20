package conf

import "strings"

type Origins []string

func (o *Origins) UnmarshalText(text []byte) error {
	*o = Origins(strings.Split(string(text), ","))
	return nil
}

type Headers []string

func (h *Headers) UnmarshalText(text []byte) error {
	*h = Headers(strings.Split(string(text), ","))
	return nil
}

type CORS struct {
	AllowOrigins   Origins  `env:"SERVER_ALLOWED_ORIGINS, required"`
	AllowedHeaders []string `env:"SERVER_ALLOWED_HEADERS, required"`
}

type server struct {
	ADDR string `env:"SERVER_ADDR, required"`
	SECRET string `env:"SERVER_SECRET, required"`
	CORS CORS
}
