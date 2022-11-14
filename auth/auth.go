package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type JWTClaim struct {
	Telephone string `json:"telephone"`
	Email    string `json:"email"`
	jwt.StandardClaims
}
func GenerateJWT(email string, telephone string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims:= &JWTClaim{
		Email: email,
		Telephone: telephone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtKey = []byte(viper.GetString("auth.jwt_key"))
	tokenString, err = token.SignedString(jwtKey)
	return
}
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth.jwt_key")), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
