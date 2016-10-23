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
	defer stmt.Close()

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

	return group, nil
}

func FindGroupByName(name string) (Group, error) {
	db := getConnection()
	rows, err := db.Query("select id, groupname, description, creatorid from groups where groupname=?", name)
	if err != nil {
		return Group{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId)
		if err != nil {
			return Group{}, err
		}
		return group, nil
	}
	return Group{}, nil
}

func ListGroups() []Group {
	db := getConnection()
	rows, err := db.Query("select groups.id, groups.groupname, groups.description, groups.creatorid, users.username from groups" +
		" inner join users on groups.creatorid = users.id")
	groups := make([]Group, 0)
	if err != nil {
		return groups
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId, &group.Creator)
		if err != nil {
			return groups
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		return groups
	}
	return groups
}

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
		log.Print(err)
		return userGroup, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(userId, groupId)
	if err != nil {
		log.Print(err)
		return userGroup, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return userGroup, err
	}

	userGroup.Id = lastId
	return userGroup, nil
}

func ListUserGroups() ([]UserGroup, error) {
	db := getConnection()
	rows, err := db.Query("select id, userid, groupid from usergroups")
	if err != nil {
		log.Print(err)
		return []UserGroup{}, err
	}

	defer rows.Close()
	userGroups := make([]UserGroup, 0)
	for rows.Next() {
		var userGroup UserGroup
		err := rows.Scan(&userGroup.Id, &userGroup.UserId, &userGroup.GroupId)
		if err != nil {
			log.Print(err)
			return []UserGroup{}, err
		}
		userGroups = append(userGroups, userGroup)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []UserGroup{}, err
	}
	return userGroups, nil
}

func ListGroupsForUser(userId int64) ([]Group, error) {
	db := getConnection()
	rows, err := db.Query("select groups.id, groupname, description, creatorid, users.username from groups"+
		" inner join usergroups on groups.id=usergroups.groupid and usergroups.userid=?"+
		" inner join users on groups.creatorid = users.id", userId)

	if err != nil {
		log.Print(err)
		return []Group{}, err
	}
	defer rows.Close()
	groups := make([]Group, 0)
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.GroupName, &group.Description, &group.CreatorId, &group.Creator)
		if err != nil {
			log.Print(err)
			return []Group{}, err
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Group{}, err
	}
	return groups, nil
}
