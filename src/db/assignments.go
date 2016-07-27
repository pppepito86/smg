package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Assignment struct {
	Id             int64
	AssignmentName string
	AuthorId       int64
	GroupId        int64
	Author         string
	Group          string
}

func CreateAssignment(a Assignment) (Assignment, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO assignments(name, author, groupid) VALUES(?, ?, ?)")
	if err != nil {
		log.Print(err)
		return a, err
	}

	res, err := stmt.Exec(a.AssignmentName, a.AuthorId, a.GroupId)
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

func ListAssignments() ([]Assignment, error) {
	db := getConnection()
	rows, err := db.Query("select assignments.id, assignments.name, assignments.author, assignments.groupid, users.username, groups.groupname from assignments" +
		" inner join users on assignments.author = users.id" +
		" inner join groups on assignments.groupid = groups.id")
	if err != nil {
		log.Print(err)
		return []Assignment{}, err
	}
	defer rows.Close()
	assignments := make([]Assignment, 0)
	for rows.Next() {
		var a Assignment
		err := rows.Scan(&a.Id, &a.AssignmentName, &a.AuthorId, &a.GroupId, &a.Author, &a.Group)
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

type AssignmentProblem struct {
	Id           int64
	AssignmentId int64
	ProblemId    int64
	Number       int64
	ProblemName  string
}

func AddProblemToAssignment(aId, pId, number int64) (AssignmentProblem, error) {
	db := getConnection()

	ap := AssignmentProblem{-1, aId, pId, number, ""}
	stmt, err := db.Prepare("INSERT INTO assignmentproblems(assignmentid, problemid) VALUES(?, ?)")
	if err != nil {
		log.Print(err)
		return ap, err
	}

	res, err := stmt.Exec(aId, pId)
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
	rows, err := db.Query("select assignmentproblems.id, assignmentid, problemid, problems.name from assignmentproblems"+
		" inner join problems on assignmentproblems.assignmentid=? and assignmentproblems.problemid=problems.id", assignmentid)
	if err != nil {
		log.Print(err)
		return []AssignmentProblem{}, err
	}
	defer rows.Close()
	assignmentProblems := make([]AssignmentProblem, 0)
	for rows.Next() {
		var ap AssignmentProblem
		err := rows.Scan(&ap.Id, &ap.AssignmentId, &ap.ProblemId, &ap.ProblemName)
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
	rows, err := db.Query("select id, assignmentid, problemid from assignmentproblems"+
		" where id = ?", apId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	ap := AssignmentProblem{}
	for rows.Next() {
		err := rows.Scan(&ap.Id, &ap.AssignmentId, &ap.ProblemId)
		if err != nil {
			log.Print(err)
			return AssignmentProblem{}, nil
		}
	}
	return ap, nil
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
