package service

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
    secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
    return &JWTMaker{secretKey: secretKey}
}

type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func (maker *JWTMaker) CreateToken(userID, username string, duration time.Duration) (string, error) {
    claims := &JWTClaims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &JWTClaims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(maker.secretKey), nil
        },
    )
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok {
        return nil, jwt.ErrInvalidKey
    }

    return claims, nil
}