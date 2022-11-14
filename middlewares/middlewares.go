package middlewares

import (
	"net/http"
	"strings"
	"user-management-api/auth"
	"user-management-api/util"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		bearerToken := strings.Split(tokenString, " ")
		if tokenString == "" || len(bearerToken) == 1 {
			util.ResponseWithError(w, http.StatusUnauthorized, "request does not contain an access token")
			return
		}
		err:= auth.ValidateToken(bearerToken[1])
		if err != nil {
			util.ResponseWithError(w, http.StatusUnauthorized,  err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}
