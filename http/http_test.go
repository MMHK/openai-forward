package http

import (
	"testing"

	_ "openai-forward/test"
)

func TestServer_Start(t *testing.T) {
	cfg := &HTTPConfig{}
	cfg.MarginWithENV()

	server := NewServer(cfg)

	if err := server.Start(); err != nil {
		t.Errorf("Failed to start HTTP server: %v", err)
	}
}
