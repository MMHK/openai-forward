package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	IssuerURL      string   `json:"issuer_url,omitempty"`
	ClientID       string   `json:"client_id,omitempty"`
	ClientSecret   string   `json:"client_secret,omitempty"`
	RedirectURL    string   `json:"redirect_url,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	Debug          bool     `json:"debug,omitempty"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
}

func LoadOIDCConfigFromEnv() *OIDCConfig {
	conf := &OIDCConfig{
		IssuerURL:      os.Getenv("OIDC_ISSUER_URL"),
		ClientID:       os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret:   os.Getenv("OIDC_CLIENT_SECRET"),
		RedirectURL:    os.Getenv("OIDC_REDIRECT_URL"),
		Debug:          os.Getenv("OIDC_DEBUG") == "true",
		AllowedDomains: strings.Split(os.Getenv("OIDC_ALLOWED_DOMAINS"), ","),
	}

	scopes := os.Getenv("OIDC_SCOPES")
	if scopes != "" {
		conf.Scopes = strings.Split(scopes, ",")
	} else {
		conf.Scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	// 清理空的域名
	cleanedDomains := []string{}
	for _, domain := range conf.AllowedDomains {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			cleanedDomains = append(cleanedDomains, domain)
		}
	}
	conf.AllowedDomains = cleanedDomains

	return conf
}

type OIDCService struct {
	OIDCConfig
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

func NewOIDCService(conf *OIDCConfig) (*OIDCService, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, conf.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: conf.ClientID,
	}

	oauth2Config := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	return &OIDCService{
		OIDCConfig:   *conf,
		provider:     provider,
		verifier:     provider.Verifier(oidcConfig),
		oauth2Config: oauth2Config,
	}, nil
}

type OIDCToken struct {
	OAuth2Token *oauth2.Token
	IDToken     *oidc.IDToken
	RawIDToken  string
}

type UserInfo struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
	Subject  string `json:"sub"`
}

func (s *OIDCService) ValidateEmailDomain(email string) bool {
	if len(s.AllowedDomains) == 0 {
		return true
	}

	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return false
	}

	domain := emailParts[1]
	for _, allowedDomain := range s.AllowedDomains {
		if strings.EqualFold(domain, allowedDomain) {
			return true
		}
	}

	return false
}

func (s *OIDCService) Exchange(ctx context.Context, code string) (*OIDCToken, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	oauth2Token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// 获取 ID Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in token response")
	}

	// 验证 ID Token
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// 检查邮箱域
	userInfo, err := s.GetUserInfo(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	if !s.ValidateEmailDomain(userInfo.Email) {
		return nil, fmt.Errorf("email domain not allowed: %s", userInfo.Email)
	}

	return &OIDCToken{
		OAuth2Token: oauth2Token,
		IDToken:     idToken,
		RawIDToken:  rawIDToken,
	}, nil
}

func (s *OIDCService) GetUserInfo(ctx context.Context, idToken *oidc.IDToken) (*UserInfo, error) {
	var claims json.RawMessage
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to get ID token claims: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(claims, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

func (s *OIDCService) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return s.oauth2Config.AuthCodeURL(state, opts...)
}
