package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	db, _ = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/smg?parseTime=true")
	defer db.Close()
	users, _ := ListUsers()

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error for user ", user.Id, err.Error())
		}
		UpdateUserHash(user.Id, string(hashedPassword))
	}
	fmt.Println("Update finished successfully")
}

type User struct {
	Id           int64
	RoleId       int64
	UserName     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	PasswordSalt string
	IsDisabled   bool
	RoleName     string
}

func ListUsers() ([]User, error) {
	rows, err := db.Query("select users.id, roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, roles.rolename from users inner join roles on users.roleid = roles.id")
	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.RoleName)
		if err != nil {
			log.Print(err)
			return []User{}, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	return users, nil
}

func UpdateUserHash(userId int64, salt string) error {
	stmt, err := db.Prepare("update users set passwordsalt=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(salt, userId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
