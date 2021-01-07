package jwt

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/caffeines/filepile/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/lib/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserID string `json:"id"`
	Phone  string `json:"phone"`
	jwt.StandardClaims
}

func BuildJWTToken(phone, scope, id string) (string, error) {
	claims := Claims{
		UserID: id,
		Phone:  phone,
		StandardClaims: jwt.StandardClaims{
			Audience:  scope,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.GetJWT().TTL)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWT().Secret))
}

func NewRefresToken() string {
	token := fmt.Sprintf("%d_%s", time.Now().Unix(), primitive.NewObjectID())
	return base64.StdEncoding.EncodeToString([]byte(token))
}
func extractTokenFromHeader(ctx echo.Context) string {
	tokenWithBearer := ctx.Request().Header.Get("Authorization")
	token := strings.Replace(tokenWithBearer, "Bearer", "", -1)
	return strings.TrimSpace(token)
}

func ExtractAndValidateToken(ctx echo.Context) (*Claims, *jwt.Token, error) {
	token := extractTokenFromHeader(ctx)
	if token == "" {
		return nil, nil, errors.NewError("Authorization token not found")
	}
	claims := Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.GetJWT().Secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !jwtToken.Valid {
		return nil, nil, errors.NewError("Token is invalid")
	}
	return &claims, jwtToken, nil
}

func ParseRefreshToken(ctx echo.Context) (string, error) {
	refresh := ctx.Request().Header.Get("RefreshToken")
	refreshWithToken := strings.Split(refresh, " ")
	if len(refreshWithToken) != 2 {
		return "", errors.NewError("Refresh token not found")
	}
	return refreshWithToken[1], nil
}
