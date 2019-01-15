package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"fmt"
	"log"
)

var db *gorm.DB
var dbErr error

func connectToDatabase() {
	db, dbErr = gorm.Open("mysql", "root:mario123@/home_automation?charset=utf8&parseTime=True&loc=Local")

	if dbErr != nil {
		fmt.Printf("DB connection failed")
		log.Fatal(dbErr)
	}
}
