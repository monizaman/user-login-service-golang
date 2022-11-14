package routes

import (
	"github.com/gorilla/mux"
	"user-management-api/controllers"
	"user-management-api/middlewares"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewPostControllerRoute(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (r *UserRouteController) UserRoute(router *mux.Router) {
	/*router.HandleFunc("/api/user/{id}", r.userController.LoginUserController).Methods("GET")*/
	router.HandleFunc("/", r.userController.HomeController).Methods("GET")
	router.HandleFunc("/api/registration", r.userController.CreateUserController).Methods("POST")
	router.HandleFunc("/api/login", r.userController.GenerateToken).Methods("POST")
	router.HandleFunc("/api/google-login", r.userController.GoogleLoginController).Methods("POST")
	router.HandleFunc("/api/user/update", middlewares.Auth(r.userController.UpdateUserController)).Methods("PATCH")
	router.HandleFunc("/api/profile", middlewares.Auth(r.userController.ProfileController)).Methods("GET")
	router.HandleFunc("/api/refresh-token", middlewares.Auth(r.userController.RefreshTokenController)).Methods("GET")
}
