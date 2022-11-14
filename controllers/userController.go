package controllers

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"google.golang.org/api/idtoken"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"user-management-api/auth"
	"user-management-api/model"
	"user-management-api/service"
	"user-management-api/util"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(propertyDetailsGinService service.UserService) UserController {
	return UserController{propertyDetailsGinService}
}

func (u *UserController) LoginUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uId, _ := strconv.ParseUint(id, 10, 32)
	user, _ := u.userService.GetUser(uint(uId))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) UpdateUserController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("authorization")
	bearerToken := strings.Split(tokenString, " ")
	tknStr := bearerToken[1]
	claims := &model.Claims{}
	var jwtKey = []byte(viper.GetString("auth.jwt_key"))
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	userUpdateCondition := model.User{Email: claims.Email}
	user := model.User{}
	json.Unmarshal([]byte(body), &user)
	usr, err := u.userService.UpdateUser(&user, userUpdateCondition)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	userRes := model.UserObject{
		Email: usr.Email,
		Telephone: usr.Telephone,
		Fullname: usr.Fullname,
		GoogleId: usr.GoogleId,
	}
	util.ResponseWithJSON(w, http.StatusOK, userRes)
}

func (u *UserController) CreateUserController(w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.Unmarshal([]byte(body), &user)
	if err := user.HashPassword(user.Password); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _ ,err := u.userService.CreateUser(&user); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := model.UserObject{Email: user.Email, Fullname: user.Fullname, Telephone: user.Telephone}
	util.ResponseWithJSON(w, http.StatusOK, res)
}


func (u *UserController) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var request model.TokenRequest
	var user model.User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.Unmarshal([]byte(body), &request)
	// check if email exists and password is correct
	user, err = u.userService.GetUserByEmail(request.Email)
	if  err != nil {
		util.ResponseWithJSON(w, http.StatusOK, map[string]string{"email": "username entered does not exist"})
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		util.ResponseWithJSON(w, http.StatusOK, map[string]string{"password": "password is incorrect"})
		return
	}
	tokenString, err:= auth.GenerateJWT(user.Email, user.Telephone)
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (u *UserController) GoogleLoginController(w http.ResponseWriter, r *http.Request) {
	var request model.GoogleLoginRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.Unmarshal([]byte(body), &request)
	payload, err := idtoken.Validate(context.Background(), request.Credential, viper.GetString("auth:google_client_id"))
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonBody, err := json.Marshal(payload.Claims)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	claimsStruct := model.GoogleClaims{}
	if err := json.Unmarshal(jsonBody, &claimsStruct); err != nil {
		return
	}
	user, err := u.userService.GetUserByGoogleId(claimsStruct.Sub)
	if err != nil {
		usr := model.User{Email: claimsStruct.Email, Fullname: claimsStruct.Name, GoogleId: claimsStruct.Sub}
		if _ ,err := u.userService.CreateUser(&usr); err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		}
	}else {
		tokenString, err:= auth.GenerateJWT(user.Email, user.Telephone)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.ResponseWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	}
}

func (u *UserController) ProfileController(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("authorization")
	bearerToken := strings.Split(tokenString, " ")
	tknStr := bearerToken[1]
	claims := &model.Claims{}
	var jwtKey = []byte(viper.GetString("auth.jwt_key"))
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := u.userService.GetUserByEmail(claims.Email)
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, err.Error())
		return
	}
	userRes := model.UserObject{
		Email: user.Email,
		Fullname: user.Fullname,
		Telephone: user.Telephone,
		GoogleId: user.GoogleId,
	}
	util.ResponseWithJSON(w, http.StatusOK, userRes)
}

func (u *UserController) RefreshTokenController(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("authorization")
	bearerToken := strings.Split(tokenString, " ")
	tknStr := bearerToken[1]
	claims := &model.Claims{}
	var jwtKey = []byte(viper.GetString("auth.jwt_key"))
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := u.userService.GetUserByEmail(claims.Email)
	newTokenString, err:= auth.GenerateJWT(user.Email, user.Telephone)
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseWithJSON(w, http.StatusOK, map[string]string{"token": newTokenString})
}
