package main

import (
	"db"
	"fmt"
)

func main() {
	err := db.OpenConnection()
	if err != nil {
		panic("Could not open db connection: " + err.Error())
	}
	defer db.Close()
	db.CreateUser(db.User{-1, "user7777", "first7777", "last7777", "email7777", "pass7777"})
	fmt.Println(db.ListUsers())
}
