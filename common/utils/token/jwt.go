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

type Claims struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

type Service struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
	now       func() time.Time
}

func NewService(cfg Config) (*Service, error) {
	key := cfg.SecretKey
	if len(key) == 0 {
		var err error
		key, err = generateKey()
		if err != nil {
			return nil, err
		}
	}

	expiry := cfg.Expiry
	if expiry <= 0 {
		expiry = 8 * time.Hour
	}

	return &Service{
		secretKey: key,
		issuer:    cfg.Issuer,
		expiry:    expiry,
		now:       time.Now,
	}, nil
}

func NewServiceFromAppConfig(cfg *common.Config) (*Service, error) {
	return NewService(Config{
		SecretKey: []byte(cfg.JWT.Secret),
		Issuer:    cfg.JWT.Issuer,
		Expiry:    cfg.JWT.Expiry,
	})
}

func (s *Service) CreateToken(uuid string) (string, error) {
	if uuid == "" {
		return "", fmt.Errorf("uuid is required")
	}
	now := s.now()
	claims := Claims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secretKey)
}

func (s *Service) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
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

// GenerateEncodedSecret creates a random 64-byte secret and returns it as a base64 string.
func GenerateEncodedSecret() (string, error) {
	key, err := generateKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %w", err)
	}
	return key, nil
}
