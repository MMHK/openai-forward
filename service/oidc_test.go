package service

import (
	"context"
	"github.com/google/uuid"
	"openai-forward/logging"
	"openai-forward/test"
	"os"
	"testing"
)

func LoadTestConfig() *OIDCConfig {
	cfg := LoadOIDCConfigFromEnv()

	logging.Logger.Debugf("Loaded OIDC config: %+v", test.ToJSON(cfg))

	return cfg
}

func TestOIDCService_AuthCodeURL(t *testing.T) {
	cfg := LoadTestConfig()

	service, err := NewOIDCService(cfg)

	if err != nil {
		t.Fatalf("Failed to create OIDC service: %v", err)
	}

	logging.Logger.Debugf("Created OIDC service: %+v", test.ToJSON(service))

	redirectURL := service.AuthCodeURL(uuid.New().String())

	logging.Logger.Debugf("Redirect URL: %s", redirectURL)
}

func TestOIDCService_Exchange(t *testing.T) {
	cfg := LoadTestConfig()

	service, err := NewOIDCService(cfg)

	if err != nil {
		t.Fatalf("Failed to create OIDC service: %v", err)
	}

	logging.Logger.Debugf("Created OIDC service: %+v", test.ToJSON(service))

	code := os.Getenv("TEST_OIDC_CODE")

	token, err := service.Exchange(context.Background(), code)

	if err != nil {
		t.Fatalf("Failed to exchange code for token: %v", err)
	}

	logging.Logger.Debugf("Exchanged code for token: %+v", test.ToJSON(token))

}

func TestOIDCService_GetUserInfo(t *testing.T) {
	cfg := LoadTestConfig()

	service, err := NewOIDCService(cfg)
	if err != nil {
		t.Fatalf("Failed to create OIDC service: %v", err)
	}

	logging.Logger.Debugf("Created OIDC service: %+v", test.ToJSON(service))

	code := os.Getenv("TEST_OIDC_CODE")

	token, err := service.Exchange(context.Background(), code)
	if err != nil {
		t.Fatalf("Failed to exchange code for token: %v", err)
	}

	logging.Logger.Debugf("Exchanged code for token: %+v", test.ToJSON(token))

	userInfo, err := service.GetUserInfo(context.Background(), token.IDToken)
	if err != nil {
		t.Fatalf("Failed to get user info: %v", err)
	}

	logging.Logger.Debugf("Got user info: %+v", test.ToJSON(userInfo))
}
