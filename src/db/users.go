package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Role struct {
	Id          int64
	RoleName    string //admin, teacher
	Description string
}

type User struct {
	Id             int64
	RoleId         int64
	UserName       string
	FirstName      string
	LastName       string
	Email          string
	PasswordHash   string
	PasswordSalt   string
	IsDisabled     bool
	RoleName       string
	ValidationCode string
}

func CreateUser(user User) (User, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO users(roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, validationcode) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return user, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.RoleId, user.UserName, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.PasswordSalt, user.IsDisabled, user.ValidationCode)
	if err != nil {
		log.Print(err)
		return user, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return user, err
	}

	user.Id = lastId
	return user, nil
}

func ListUsers() ([]User, error) {
	db := getConnection()
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

func ListMyUsers(id int64) ([]User, error) {
	db := getConnection()
	rows, err := db.Query("select distinct u.id, u.roleid, u.username, u.firstname, u.lastname, u.email, u.isdisabled, r.rolename from usergroups as ug"+
		" inner join users as u on u.id = ug.userid"+
		" inner join groups as g on ug.groupid=g.id and g.creatorid = ?"+
		" inner join roles as r on u.roleid = r.id", id)

	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.IsDisabled, &user.RoleName)
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

func ListUsersInGroup(groupId int64) ([]User, error) {
	db := getConnection()
	rows, err := db.Query("select distinct u.id, u.roleid, u.username, u.firstname, u.lastname, u.email, u.isdisabled, r.rolename from usergroups as ug"+
		" inner join users as u on u.id = ug.userid and ug.groupid = ? "+
		" inner join roles as r on u.roleid = r.id", groupId)

	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.IsDisabled, &user.RoleName)
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

func GetUser(username string) (User, error) {
	db := getConnection()
	rows, err := db.Query("select users.id, roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, validationcode, rolename from users inner join roles on username= ? and isdisabled=? and users.roleid=roles.id", username, false)
	if err != nil {
		log.Print(err)
		return User{}, nil
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.ValidationCode, &user.RoleName)
		if err != nil {
			log.Print(err)
			return User{}, err
		}
		return user, nil
	}
	return user, nil
}

func GetUserById(userId int64) (User, error) {
	db := getConnection()
	rows, err := db.Query("select users.id, roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, rolename from users inner join roles on users.id= ? and isdisabled=? and users.roleid=roles.id", userId, false)
	if err != nil {
		log.Print(err)
		return User{}, nil
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.RoleName)
		if err != nil {
			log.Print(err)
			return User{}, err
		}
		return user, nil
	}
	return user, nil
}

func UpdateUserRole(userId, roleId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("update users set roleid=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleId, userId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func ValidateEmail(code string) error {
	db := getConnection()
	stmt, err := db.Prepare("update users set validationcode=? where validationcode=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec("", code)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
