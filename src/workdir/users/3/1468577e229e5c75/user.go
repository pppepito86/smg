package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	FirstName  string
	SecondName string
	LastName   string
	UserName   string
	Email      string
	Id         string
}

func CreateUser(user User) error {
	if user.UserName == "" {
		return errors.New("Username is empty")
	}
	if ContainsUser(user.UserName) {
		return errors.New("User already exists")
	}
	users := GetAllUsers()
	users = append(users, user)
	return writeEntities("users", users)
}

func GetAllUsers() []User {
	return getAllEntities("users").([]User)
}

func ContainsUser(username string) bool {
	user := FindUser(username)
	return user.UserName != ""
}

func FindUser(username string) User {
	users := GetAllUsers()
	for _, user := range users {
		if user.UserName == username {
			return user
		}
	}
	return User{}
}

func writeEntities(name string, entities interface{}) error {
	file := name + ".json"
	json, _ := json.Marshal(entities)

	os.Create(file)
	return ioutil.WriteFile(file, json, 0700)
}

func getAllEntities(name string) interface{} {
	users := make([]User, 0)
	file := ReadFile(name + ".json")
	json.Unmarshal(file, &users)
	return users
}

func ReadFile(filename string) []byte {
	fileAsByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	fmt.Println(string(fileAsByte))
	return fileAsByte
}
