package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Group struct {
	Id          int64
	GroupName   string
	Description string
	CreatorId   int64
	Creator     string
}

func CreateGroup(group Group) (Group, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO groups(groupname, description, creatorid) VALUES(?, ?, ?)")
	if err != nil {
		log.Print(err)
		return group, err
	}

	res, err := stmt.Exec(group.GroupName, group.Description, group.CreatorId)
	if err != nil {
		log.Print(err)
		return group, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return group, err
	}

	group.Id = lastId

	//userGroup, _ := CreateUserGroup(group.CreatorId, group.Id)
	//CreateUserGroupRole(userGroup.Id, 1)
	return group, nil
}

func FindGroupByName(name string) (Group, error) {
	db := getConnection()
	rows, err := db.Query("select id, groupname, description, creatorid from groups where groupname=?", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId)
		if err != nil {
			log.Fatal(err)
		}
		return group, nil
	}
	return Group{}, nil
}

func ListGroups() []Group {
	db := getConnection()
	rows, err := db.Query("select groups.id, groups.groupname, groups.description, groups.creatorid, users.username from groups" +
		" inner join users on groups.creatorid = users.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	groups := make([]Group, 0)
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId, &group.Creator)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return groups
}

/*
type GroupRole struct {
	Id          int64
	RoleName    string // admin, student, parent
	Description string
}
*/

type UserGroup struct {
	Id      int64
	UserId  int64
	GroupId int64
}

func CreateUserGroup(userId, groupId int64) (UserGroup, error) {
	db := getConnection()

	userGroup := UserGroup{-1, userId, groupId}

	stmt, err := db.Prepare("INSERT INTO usergroups(userid, groupid) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
		return userGroup, err
	}

	res, err := stmt.Exec(userId, groupId)
	if err != nil {
		log.Fatal(err)
		return userGroup, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return userGroup, err
	}

	userGroup.Id = lastId
	return userGroup, nil
}

func ListUserGroups() []UserGroup {
	db := getConnection()
	rows, err := db.Query("select id, userid, groupid from usergroups")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	userGroups := make([]UserGroup, 0)
	for rows.Next() {
		var userGroup UserGroup
		err := rows.Scan(&userGroup.Id, &userGroup.UserId, &userGroup.GroupId)
		if err != nil {
			log.Fatal(err)
		}
		userGroups = append(userGroups, userGroup)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return userGroups
}

func ListGroupsForUser(userId int64) []Group {
	db := getConnection()
	rows, err := db.Query("select groups.id, groupname, description, creatorid from groups inner join usergroups on groups.id=usergroups.groupid and usergroups.userid=?", userId)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	groups := make([]Group, 0)
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return groups
}

/*
type UserGroupRole struct {
	Id          int64
	UserGroupId int64
	GroupRoleId int64
}

func CreateUserGroupRole(userGroupId, groupRoleId int64) (UserGroupRole, error) {
	db := getConnection()

	userGroupRole := UserGroupRole{-1, userGroupId, groupRoleId}

	stmt, err := db.Prepare("INSERT INTO usergrouproles(usergroupid, grouproleid) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
		return userGroupRole, err
	}

	res, err := stmt.Exec(userGroupId, groupRoleId)
	if err != nil {
		log.Fatal(err)
		return userGroupRole, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return userGroupRole, err
	}

	userGroupRole.Id = lastId
	return userGroupRole, nil
}

func main() {
	err := OpenConnection()
	if err != nil {
		panic("Could not open db connection: " + err.Error())
	}
	defer db.Close()
	CreateUser(User{-1, "user6", "first6", "last6", "email6", "pass6"})
	fmt.Println(ListUsers())
}
*/
