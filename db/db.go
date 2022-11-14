package db

import (
	"fmt"
	"gorm.io/gorm"
	"user-management-api/model"
)

import (
	"gorm.io/driver/mysql"
	"log"
)
var Instance *gorm.DB
var dbError error
func Connect(connectionString string) () {
	Instance, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		fmt.Println("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}
func Migrate() {
	Instance.AutoMigrate(&model.User{})
	log.Println("Database Migration Completed!")
}



