package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"cams.dev/video_upload_backend/internal/domain/entity"
)

// Claims defines the JWT claims
type Claims struct {
	UserID   string          `json:"userId"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTService provides JWT authentication functionality
type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a JWT token for a user
func (s *JWTService) GenerateToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(s.tokenDuration)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cams.dev/video_upload_backend",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword generates a bcrypt hash from a password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if a password matches a bcrypt hash
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
