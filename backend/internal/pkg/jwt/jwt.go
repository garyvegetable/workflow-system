package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtInstance *JWT
var secret []byte

func Init(secretKey string) {
	secret = []byte(secretKey)
	jwtInstance = &JWT{}
}

func NewJWT() *JWT {
	return &JWT{}
}

type JWT struct{}

type Claims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	CompanyID int64  `json:"company_id"`
	jwt.RegisteredClaims
}

func (j *JWT) GenerateToken(userID int64, username string, companyID int64) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Username:  username,
		CompanyID: companyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
