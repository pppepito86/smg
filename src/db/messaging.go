package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewMessage(from string, to string, msg string) (error) {
	db := getConnection()
	
	stmt, err := db.Prepare(`
INSERT INTO msg (fromUsr, toUsr, msg) 
SELECT fromUsr.id, toUsr.id, ?
FROM 
	(SELECT id from users where username=?) AS fromUsr
CROSS JOIN
	(SELECT id from users where username=?) AS toUsr;`)
	
	if err != nil {
		log.Print(err)
		return err
	}
	
	
	_ , err = stmt.Exec(msg, from, to)
	if err != nil {
		log.Print(err)
		return err
	}
	
	return nil
}


func GetNumUnread(user string) (int64, error) {
	db := getConnection()
	
	 stmt, err := db.Prepare(`
SELECT COUNT(*) 
FROM msg 
INNER JOIN users
ON msg.toUsr=users.id
WHERE msg.seen=false
AND users.username=?;`)

        if err != nil {
                log.Print(err)
                return -1, err
        }


        rows , err := stmt.Query(user)
        if err != nil {
                log.Print(err)
                return -1, err
        }


	defer rows.Close()
        var numUnread int64

        rows.Next() 
        err = rows.Scan(&numUnread)
        
	if err != nil {
        	log.Print(err)
                return -1, err
        }
        
        return numUnread, nil
}

/*
type Role struct {
	Id          int64
	RoleName    string //admin, teacher
	Description string
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

func CreateUser(user User) (User, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO users(roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return user, err
	}

	res, err := stmt.Exec(user.RoleId, user.UserName, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.PasswordSalt, user.IsDisabled)
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

func GetUser(username string) (User, error) {
	db := getConnection()
	rows, err := db.Query("select users.id, roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, rolename from users inner join roles on username= ? and isdisabled=? and users.roleid=roles.id", username, false)
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

	_, err = stmt.Exec(roleId, userId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
*/

