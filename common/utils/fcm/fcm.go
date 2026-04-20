package fcm

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cth.release/common"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

const defaultScope = "https://www.googleapis.com/auth/firebase.messaging"

type Credentials struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

type Config struct {
	Credentials Credentials
	BaseURL     string
	Scope       string
	HTTPClient  *resty.Client
}

type Sender struct {
	credentials Credentials
	baseURL     string
	scope       string
	httpClient  *resty.Client
	now         func() time.Time
}

type Notification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Image string `json:"image,omitempty"`
}

type AndroidNotification struct {
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Color    string `json:"color,omitempty"`
	ImageURL string `json:"image,omitempty"`
}

type AndroidConfig struct {
	Priority     string               `json:"priority,omitempty"`
	TTL          string               `json:"ttl,omitempty"`
	CollapseKey  string               `json:"collapse_key,omitempty"`
	RestrictedTo string               `json:"restricted_package_name,omitempty"`
	Data         map[string]string    `json:"data,omitempty"`
	Notification *AndroidNotification `json:"notification,omitempty"`
}

type WebpushNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Image string `json:"image,omitempty"`
}

type WebpushConfig struct {
	Headers      map[string]string      `json:"headers,omitempty"`
	Data         map[string]string      `json:"data,omitempty"`
	Notification *WebpushNotification   `json:"notification,omitempty"`
	FCMOptions   map[string]interface{} `json:"fcm_options,omitempty"`
}

type APNSConfig struct {
	Headers    map[string]string      `json:"headers,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	FCMOptions map[string]interface{} `json:"fcm_options,omitempty"`
}

type Message struct {
	Token        string            `json:"token,omitempty"`
	Topic        string            `json:"topic,omitempty"`
	Condition    string            `json:"condition,omitempty"`
	Notification *Notification     `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Android      *AndroidConfig    `json:"android,omitempty"`
	APNS         *APNSConfig       `json:"apns,omitempty"`
	Webpush      *WebpushConfig    `json:"webpush,omitempty"`
}

type SendRequest struct {
	ValidateOnly bool    `json:"validate_only,omitempty"`
	Message      Message `json:"message"`
}

type SendResponse struct {
	Name string `json:"name"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewSender(cfg Config) (*Sender, error) {
	credentials := cfg.Credentials
	if credentials.ProjectID == "" || credentials.ClientEmail == "" || credentials.PrivateKey == "" {
		return nil, fmt.Errorf("fcm credentials are incomplete")
	}
	if credentials.TokenURI == "" {
		credentials.TokenURI = "https://oauth2.googleapis.com/token"
	}

	scope := cfg.Scope
	if scope == "" {
		scope = defaultScope
	}

	client := cfg.HTTPClient
	if client == nil {
		client = resty.New().SetTimeout(10 * time.Second)
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://fcm.googleapis.com"
	}

	return &Sender{
		credentials: credentials,
		baseURL:     baseURL,
		scope:       scope,
		httpClient:  client,
		now:         time.Now,
	}, nil
}

func NewSenderFromAppConfig(cfg *common.Config) (*Sender, error) {
	credentials, err := CredentialsFromAppConfig(cfg)
	if err != nil {
		return nil, err
	}

	return NewSender(Config{
		Credentials: credentials,
		BaseURL:     cfg.FCMBaseURL,
		Scope:       cfg.FCMScope,
	})
}

func CredentialsFromAppConfig(cfg *common.Config) (Credentials, error) {
	if strings.TrimSpace(cfg.FCMCredsJSON) != "" {
		return CredentialsFromJSON([]byte(cfg.FCMCredsJSON))
	}
	if strings.TrimSpace(cfg.FCMCredsFile) != "" {
		data, err := os.ReadFile(cfg.FCMCredsFile)
		if err != nil {
			return Credentials{}, err
		}
		return CredentialsFromJSON(data)
	}

	credentials := Credentials{
		ProjectID:   strings.TrimSpace(cfg.FCMProjectID),
		ClientEmail: strings.TrimSpace(cfg.FCMClientEmail),
		PrivateKey:  normalizePrivateKey(cfg.FCMPrivateKey),
		TokenURI:    strings.TrimSpace(cfg.FCMTokenURI),
	}
	if credentials.TokenURI == "" {
		credentials.TokenURI = "https://oauth2.googleapis.com/token"
	}
	if credentials.ProjectID == "" || credentials.ClientEmail == "" || credentials.PrivateKey == "" {
		return Credentials{}, fmt.Errorf("missing fcm configuration")
	}

	return credentials, nil
}

func CredentialsFromJSON(data []byte) (Credentials, error) {
	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return Credentials{}, err
	}

	credentials.ProjectID = strings.TrimSpace(credentials.ProjectID)
	credentials.ClientEmail = strings.TrimSpace(credentials.ClientEmail)
	credentials.PrivateKey = normalizePrivateKey(credentials.PrivateKey)
	credentials.TokenURI = strings.TrimSpace(credentials.TokenURI)
	if credentials.TokenURI == "" {
		credentials.TokenURI = "https://oauth2.googleapis.com/token"
	}
	if credentials.ProjectID == "" || credentials.ClientEmail == "" || credentials.PrivateKey == "" {
		return Credentials{}, fmt.Errorf("invalid service account json")
	}
	return credentials, nil
}

func (s *Sender) Send(ctx context.Context, request SendRequest) (*SendResponse, error) {
	if err := validateMessage(request.Message); err != nil {
		return nil, err
	}

	accessToken, err := s.accessToken(ctx)
	if err != nil {
		return nil, err
	}

	var response SendResponse
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(s.messageSendURL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("fcm send failed: status=%d body=%s", resp.StatusCode(), resp.String())
	}

	return &response, nil
}

func (s *Sender) SendToToken(ctx context.Context, token string, notification *Notification, data map[string]string) (*SendResponse, error) {
	return s.Send(ctx, SendRequest{
		Message: Message{
			Token:        strings.TrimSpace(token),
			Notification: notification,
			Data:         data,
		},
	})
}

func (s *Sender) messageSendURL() string {
	return s.baseURL + path.Join("/v1/projects", s.credentials.ProjectID, "messages:send")
}

func (s *Sender) accessToken(ctx context.Context) (string, error) {
	assertion, err := s.signedAssertion()
	if err != nil {
		return "", err
	}

	var tokenResponse accessTokenResponse
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
			"assertion":  assertion,
		}).
		SetResult(&tokenResponse).
		Post(s.credentials.TokenURI)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", fmt.Errorf("fcm auth failed: status=%d body=%s", resp.StatusCode(), resp.String())
	}
	if tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("empty fcm access token")
	}

	return tokenResponse.AccessToken, nil
}

func (s *Sender) signedAssertion() (string, error) {
	privateKey, err := parsePrivateKey([]byte(s.credentials.PrivateKey))
	if err != nil {
		return "", err
	}

	now := s.now()
	claims := jwt.MapClaims{
		"iss":   s.credentials.ClientEmail,
		"scope": s.scope,
		"aud":   s.credentials.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func validateMessage(message Message) error {
	targets := 0
	if strings.TrimSpace(message.Token) != "" {
		targets++
	}
	if strings.TrimSpace(message.Topic) != "" {
		targets++
	}
	if strings.TrimSpace(message.Condition) != "" {
		targets++
	}
	if targets != 1 {
		return fmt.Errorf("exactly one of token, topic, or condition must be set")
	}
	return nil
}

func normalizePrivateKey(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), `\n`, "\n")
}

func parsePrivateKey(raw []byte) (*rsa.PrivateKey, error) {
	if key, err := jwt.ParseRSAPrivateKeyFromPEM(raw); err == nil {
		return key, nil
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("invalid private key pem")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not rsa")
	}
	return key, nil
}
