package conf

import (
	"testing"

	"github.com/joho/godotenv"
)

func Test_Config(t *testing.T) {
	if err := godotenv.Load("../../assets/test.env"); err != nil {
		t.Fatalf("godotenv.Load() failed: %v", err)
	}

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if cfg.Server.ADDR != "testAddr" {
		t.Errorf("cfg.Server.Address = %v, want %v", cfg.Server.ADDR, "testAddr")
	}

	if len(cfg.Server.CORS.AllowOrigins) != 2 {
		t.Errorf("len(cfg.Server.CORS.AllowOrigins) = %v, want %v", len(cfg.Server.CORS.AllowOrigins), 2)
	}

	if len(cfg.Server.CORS.AllowedHeaders) != 2 {
		t.Errorf("len(cfg.Server.CORS.AllowedHeaders) = %v, want %v", len(cfg.Server.CORS.AllowedHeaders), 2)
	}
}
