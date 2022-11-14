package model

import "github.com/dgrijalva/jwt-go"

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleLoginRequest struct {
	Credential    string `json:"credential"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}


type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Name      string `json:"name"`
	Sub      	  string `json:"sub"`
}
