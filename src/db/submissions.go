package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Submission struct {
	Id                int64
	ApId              int64
	AssignmentId      int64
	ProblemId         int64
	UserId            int64
	Language          string
	SourceFile        string
	Time              time.Time
	Verdict           string
	Reason            string
	ProblemName       string
	Source            string
	User              string
	FirstName         string
	LastName          string
	SubmissionDetails []SubmissionDetail
	Points            int
}

func AddSubmission(s Submission) (Submission, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO submissions(assignmentproblemid, assignmentid, problemid, userid, language, sourcefile, verdict, reason) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return s, err
	}

	res, err := stmt.Exec(s.ApId, s.AssignmentId, s.ProblemId, s.UserId, s.Language, s.SourceFile, s.Verdict, s.Reason)
	if err != nil {
		log.Print(err)
		return s, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return s, err
	}

	s.Id = lastId
	return s, nil
}

func ListSubmissions() ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select id, assignmentproblemid, userid, language, sourcefile, time, verdict from submissions order by id desc")
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ApId, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func ListSubmission(submissionId int64) (Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, submissions.userid, language, sourcefile, time, verdict, reason, problems.name from submissions"+
		" inner join assignmentproblems on submissions.id=? and assignmentproblems.id=submissions.assignmentproblemid"+
		"	inner join problems on problems.id=assignmentproblems.problemid", submissionId)
	s := Submission{}
	if err != nil {
		log.Print(err)
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&s.Id, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Reason, &s.ProblemName)
		if err != nil {
			log.Print(err)
			return s, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return s, err
	}
	return s, nil
}

func ListMySubmissions(userId, assignmentId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, language, sourcefile, time, verdict, problems.name from submissions"+
		" inner join assignmentproblems on assignmentproblems.id=submissions.assignmentproblemid and assignmentproblems.assignmentid = ? and submissions.userid = ?"+
		"	inner join problems on problems.id=assignmentproblems.problemid", assignmentId, userId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.ProblemName)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func ListMyAllSubmissions(userId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, time, assignmentproblemid, points from submissions where userid=?", userId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.Time, &s.ApId, &s.Points)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func ListSubmissionsForAssignment(assignmentId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, assignmentproblemid, language, sourcefile, time, verdict, submissions.points, problems.name, users.id, users.username, users.firstname, users.lastname from submissions"+
		" inner join assignmentproblems on assignmentproblems.id=submissions.assignmentproblemid and assignmentproblems.assignmentid = ?"+
		"	inner join problems on problems.id=assignmentproblems.problemid"+
		"	inner join users on users.id=submissions.userid", assignmentId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ApId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Points, &s.ProblemName, &s.UserId, &s.User, &s.FirstName, &s.LastName)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func ListMySubmissionsForProblem(userId, apId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, language, sourcefile, verdict from submissions"+
		" where userid=? and assignmentproblemid=?", userId, apId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.Language, &s.SourceFile, &s.Verdict)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func ListProblemSubmissions(pId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, submissions.assignmentproblemid, submissions.userid, language, sourcefile, time, verdict, reason from submissions"+
		" inner join assignmentproblems on assignmentproblems.id=submissions.assignmentproblemid"+
		"	and assignmentproblems.problemid=?", pId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	submissions := make([]Submission, 0)
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ApId, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Reason)
		if err != nil {
			log.Print(err)
			return []Submission{}, err
		}
		submissions = append(submissions, s)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	return submissions, nil
}

func UpdateVerdict(id int64, verdict, reason string, correct, total, points int) error {
	db := getConnection()

	stmt, err := db.Prepare("update submissions set verdict=?, reason=?, correct=?, total=?, points=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(verdict, reason, correct, total, points, id)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

type SubmissionDetail struct {
	Id           int64
	SubmissionId int64
	Step         string
	Verdict      string
	Reason       string
	Time         int64
	Input        string
	Output       string
}

func AddSubmissionDetails(sid int64, step, verdict, reason string, time int64) (SubmissionDetail, error) {
	db := getConnection()
	sd := SubmissionDetail{
		SubmissionId: sid,
		Step:         step,
		Verdict:      verdict,
		Reason:       reason,
		Time:         time,
	}

	stmt, err := db.Prepare("INSERT INTO submissiondetails(submissionid, step, status, reason, time) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return sd, err
	}

	res, err := stmt.Exec(sd.SubmissionId, sd.Step, sd.Verdict, sd.Reason, sd.Time)
	if err != nil {
		log.Print(err)
		return sd, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return sd, err
	}

	sd.Id = lastId
	return sd, nil
}

func DeleteSubmissionDetails(sId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("delete from submissiondetails where submissionid=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(sId)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func ListSubmissionDetails(submissionid int64) ([]SubmissionDetail, error) {
	db := getConnection()
	rows, err := db.Query("select id, submissionid, step, status, reason, time from submissiondetails where submissionid = ? order by id asc", submissionid)
	sds := make([]SubmissionDetail, 0)
	if err != nil {
		log.Print(err)
		return sds, err
	}
	defer rows.Close()
	for rows.Next() {
		var sd SubmissionDetail
		err := rows.Scan(&sd.Id, &sd.SubmissionId, &sd.Step, &sd.Verdict, &sd.Reason, &sd.Time)
		if err != nil {
			log.Print(err)
			return sds, err
		}
		sds = append(sds, sd)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return sds, err
	}
	return sds, nil
}
