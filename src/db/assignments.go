package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Assignment struct {
	Id             int64
	AssignmentName string
	AuthorId       int64
	GroupId        int64
	Author         string
	Group          string
	Problems       string
	StartTime      time.Time
	EndTime        time.Time
	TestInfo       string
	Standings      string
	IsActive       bool
	HasFinished    bool
}

func CreateAssignment(a Assignment) (Assignment, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO assignments(name, author, groupid, starttime, endtime, testinfo, standings) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return a, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.AssignmentName, a.AuthorId, a.GroupId, a.StartTime, a.EndTime, a.TestInfo, a.Standings)
	if err != nil {
		log.Print(err)
		return a, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return a, err
	}

	a.Id = lastId
	return a, nil
}

func ListAssignment(aid int64) (Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname, assignments.starttime, assignments.endtime, assignments.testinfo, assignments.standings from assignments"+
		" inner join users on assignments.id=? and assignments.author = users.id"+
		" inner join groups on assignments.groupid = groups.id", aid)
	a := Assignment{}
	if err != nil {
		log.Print(err)
		return a, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group, &a.StartTime, &a.EndTime, &a.TestInfo, &a.Standings)
		if err != nil {
			log.Print(err)
			return a, err
		}
	}
	location, _ := time.LoadLocation("Europe/Sofia")
	a.StartTime = a.StartTime.In(location)
	a.EndTime = a.EndTime.In(location)

	return a, nil
}

func ListAssignments() ([]Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname, assignments.starttime, assignments.endtime from assignments" +
		" inner join users on assignments.author = users.id" +
		" inner join groups on assignments.groupid = groups.id")
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	location, _ := time.LoadLocation("Europe/Sofia")
	defer rows.Close()
	assignments := make([]Assignment, 0)
	for rows.Next() {
		var a Assignment
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group, &a.StartTime, &a.EndTime)
		a.StartTime = a.StartTime.In(location)
		a.EndTime = a.EndTime.In(location)

		if err != nil {
			log.Print(err)
			return []Assignment{}, err
		}
		assignments = append(assignments, a)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	return assignments, nil
}

func ListAssignmentsForGroup(groupId int64) ([]Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname, assignments.starttime, assignments.endtime from assignments"+
		" inner join users on assignments.author = users.id"+
		" inner join groups on assignments.groupid = groups.id and assignments.groupid = ?", groupId)
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	location, _ := time.LoadLocation("Europe/Sofia")
	defer rows.Close()
	assignments := make([]Assignment, 0)
	for rows.Next() {
		var a Assignment
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group, &a.StartTime, &a.EndTime)
		a.StartTime = a.StartTime.In(location)
		a.EndTime = a.EndTime.In(location)

		if err != nil {
			log.Print(err)
			return []Assignment{}, err
		}
		assignments = append(assignments, a)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	return assignments, nil
}

func ListAssignmentsForAuthor(aId int64) ([]Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname, assignments.starttime, assignments.endtime from assignments"+
		" inner join users on assignments.author = users.id and assignments.author = ?"+
		" inner join groups on assignments.groupid = groups.id", aId)
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	location, _ := time.LoadLocation("Europe/Sofia")
	defer rows.Close()
	assignments := make([]Assignment, 0)
	for rows.Next() {
		var a Assignment
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group, &a.StartTime, &a.EndTime)
		a.StartTime = a.StartTime.In(location)
		a.EndTime = a.EndTime.In(location)

		if err != nil {
			log.Print(err)
			return []Assignment{}, err
		}
		assignments = append(assignments, a)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	return assignments, nil
}
func ListAssignmentsForUser(user User) ([]Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname, assignments.starttime, assignments.endtime from assignments"+
		" inner join users on assignments.author = users.id"+
		" inner join groups on assignments.groupid = groups.id"+
		" inner join usergroups on assignments.groupid = usergroups.groupid and usergroups.userid = ?", user.Id)
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	defer rows.Close()
	assignments := make([]Assignment, 0)
	location, _ := time.LoadLocation("Europe/Sofia")
	for rows.Next() {
		var a Assignment
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group, &a.StartTime, &a.EndTime)
		if err != nil {
			log.Print(err)
			return []Assignment{}, err
		}
		a.StartTime = a.StartTime.In(location)
		a.EndTime = a.EndTime.In(location)

		assignments = append(assignments, a)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	return assignments, nil
}

func UpdateAssignment(a Assignment) error {
	db := getConnection()

	stmt, err := db.Prepare("update assignments set name=?,groupid=?,starttime=?,endtime=?,testinfo=?,standings=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.AssignmentName, a.GroupId, a.StartTime, a.EndTime, a.TestInfo, a.Standings, a.Id)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func ListUsersForAssignment(id int64) ([]User, error) {
	db := getConnection()
	rows, err := db.Query("select users.id, users.roleid, username, firstname, lastname, email from users"+
		" inner join usergroups on usergroups.userid = users.id"+
		" inner join assignments on assignments.groupid = usergroups.groupid and assignments.id=?", id)
	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.RoleId, &u.UserName, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			log.Print(err)
			return []User{}, err
		}
		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []User{}, err
	}
	return users, nil
}

type AssignmentProblem struct {
	Id           int64
	AssignmentId int64
	ProblemId    int64
	Number       int64
	ProblemName  string
	Points       int
	Languages    string
}

func AddProblemToAssignment(aId, pId, number int64, points int) (AssignmentProblem, error) {
	db := getConnection()

	ap := AssignmentProblem{
		Id:           -1,
		AssignmentId: aId,
		ProblemId:    pId,
		Number:       number,
	}
	stmt, err := db.Prepare("INSERT INTO assignmentproblems(assignmentid, problemid, number, points) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return ap, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(aId, pId, number, points)
	if err != nil {
		log.Print(err)
		return ap, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return ap, err
	}

	ap.Id = lastId
	return ap, nil
}

func ListAssignmentProblems(assignmentid int64) ([]AssignmentProblem, error) {
	db := getConnection()
	rows, err := db.Query("select assignmentproblems.id, assignmentid, problemid, assignmentproblems.number, assignmentproblems.points, problems.name from assignmentproblems"+
		" inner join problems on assignmentproblems.assignmentid=? and assignmentproblems.problemid=problems.id order by assignmentproblems.number", assignmentid)
	if err != nil {
		log.Print(err)
		return []AssignmentProblem{}, err
	}
	defer rows.Close()
	assignmentProblems := make([]AssignmentProblem, 0)
	for rows.Next() {
		var ap AssignmentProblem
		err := rows.Scan(&ap.Id, &ap.AssignmentId, &ap.ProblemId, &ap.Number, &ap.Points, &ap.ProblemName)
		if err != nil {
			log.Print(err)
			return []AssignmentProblem{}, err
		}
		assignmentProblems = append(assignmentProblems, ap)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []AssignmentProblem{}, err
	}
	return assignmentProblems, nil
}

func GetAssignmentProblem(apId int64) (AssignmentProblem, error) {
	db := getConnection()
	rows, err := db.Query("select a.id, a.assignmentid, a.problemid, a.points, p.languages from assignmentproblems as a"+
		" inner join problems as p on a.problemid=p.id and a.id = ?", apId)
	if err != nil {
		return AssignmentProblem{}, nil
	}
	defer rows.Close()
	ap := AssignmentProblem{}
	for rows.Next() {
		err := rows.Scan(&ap.Id, &ap.AssignmentId, &ap.ProblemId, &ap.Points, &ap.Languages)
		if err != nil {
			log.Print(err)
			return AssignmentProblem{}, nil
		}
	}
	return ap, nil
}

func UpdateAssignmentProblem(apId, problemId int64, points int) error {
	db := getConnection()

	stmt, err := db.Prepare("update assignmentproblems set problemid=?, points=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(problemId, points, apId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func DeleteAssignmentProblem(apId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("delete from assignmentproblems where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(apId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func DeleteAssignmentProblem2(aId, pId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("delete from assignmentproblems where assignmentid=? and problemid=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(aId, pId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func DeleteAssignmentProblems(aId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("delete from assignmentproblems where assignmentid=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(aId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func IsUserAssignedToCompetition(userId, cId int64) (bool, error) {
	db := getConnection()
	rows, err := db.Query("select * from assignments inner join usergroups"+
		" on assignments.id=? and assignments.groupid=usergroups.groupid"+
		" and usergroups.userid=?", cId, userId)
	if err != nil {
		log.Print(err)
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
