package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"user-management-api/config"
	"user-management-api/controllers"
	"user-management-api/model"
	"user-management-api/routes"
	"user-management-api/service"
)

var (
	ctx                 context.Context
	dbConnection        *gorm.DB
	DbError             error
	userService         service.UserService
	UserController      controllers.UserController
	UserRouteController = routes.UserRouteController{}
)

func init() {
	config.LoadAppConfig()
	ctx = context.TODO()
	dsn := viper.GetString("database.db_user") + ":" + viper.GetString("database.db_password") +
		"@tcp(" + viper.GetString("database.db_tcp") + ")/" + viper.GetString("database.db_name") + "?charset=utf8mb4&parseTime=true&loc=Local"
	dbConnection, DbError = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if DbError != nil {
		fmt.Println("error", DbError)
	}
	dbConnection.AutoMigrate(&model.User{})
	// ? Instantiate the Constructors
	userService = service.NewUserService(dbConnection, ctx)
	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewPostControllerRoute(UserController)

}

func main() {
	// Initialize the router
	router := mux.NewRouter() //.StrictSlash(true)
	// Register Routes
	UserRouteController.UserRoute(router)
	http.Handle("/", router)

	err := http.ListenAndServe(":"+viper.GetString("server.port"), nil)
	if err != nil {
		log.Println("An error occurred starting HTTP listener at port " + viper.GetString("server.port"))
		log.Println("Error: " + err.Error())
	}
}
