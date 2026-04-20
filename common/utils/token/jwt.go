package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"cth.release/common"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	SecretKey []byte
	Issuer    string
	Expiry    time.Duration
}

type Service struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
	now       func() time.Time
}

type Claims struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

func NewService(cfg Config) (*Service, error) {
	secretKey := cfg.SecretKey
	if len(secretKey) == 0 {
		var err error
		secretKey, err = generateSecretKey()
		if err != nil {
			return nil, err
		}
	}

	expiry := cfg.Expiry
	if expiry <= 0 {
		expiry = 8 * time.Hour
	}

	return &Service{
		secretKey: secretKey,
		issuer:    cfg.Issuer,
		expiry:    expiry,
		now:       time.Now,
	}, nil
}

func NewServiceFromAppConfig(cfg *common.Config) (*Service, error) {
	return NewService(Config{
		SecretKey: []byte(cfg.JWTSecret),
		Issuer:    cfg.JWTIssuer,
		Expiry:    cfg.JWTExpiry,
	})
}

func (s *Service) CreateToken(uuid string) (string, error) {
	if uuid == "" {
		return "", fmt.Errorf("uuid is required")
	}

	issuedAt := s.now()
	expiresAt := issuedAt.Add(s.expiry)
	claims := Claims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *Service) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *Service) Expiry() time.Duration {
	return s.expiry
}

func generateSecretKey() ([]byte, error) {
	secretKey := make([]byte, 64)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %w", err)
	}
	return secretKey, nil
}

func GenerateEncodedSecret() (string, error) {
	secret, err := generateSecretKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}
