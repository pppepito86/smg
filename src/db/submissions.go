package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Submission struct {
	Id                int64
	ApId              int64
	UserId            int64
	Language          string
	SourceFile        string
	Time              string
	Verdict           string
	Reason            string
	ProblemName       string
	Source            string
	User              string
	SubmissionDetails []SubmissionDetail
}

func AddSubmission(s Submission) (Submission, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO submissions(assignmentproblemid, userid, language, sourcefile, verdict) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return s, err
	}

	res, err := stmt.Exec(s.ApId, s.UserId, s.Language, s.SourceFile, s.Verdict)
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

func ListSubmissionsForAssignment(assignmentId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, language, sourcefile, time, verdict, problems.name, users.id, users.username from submissions"+
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
		err := rows.Scan(&s.Id, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.ProblemName, &s.UserId, &s.User)
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

func ListMySubmissionsForProblem(apId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, language, sourcefile, verdict from submissions"+
		" where assignmentproblemid=?", apId)
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

func UpdateVerdict(id int64, verdict, reason string) error {
	db := getConnection()

	stmt, err := db.Prepare("update submissions set verdict=?, reason=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(verdict, reason, id)
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
}

func AddSubmissionDetails(sid int64, step, verdict, reason string) (SubmissionDetail, error) {
	db := getConnection()
	sd := SubmissionDetail{
		SubmissionId: sid,
		Step:         step,
		Verdict:      verdict,
		Reason:       reason,
	}

	stmt, err := db.Prepare("INSERT INTO submissiondetails(submissionid, step, status, reason) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return sd, err
	}

	res, err := stmt.Exec(sd.SubmissionId, sd.Step, sd.Verdict, sd.Reason)
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

func ListSubmissionDetails(submissionid int64) ([]SubmissionDetail, error) {
	db := getConnection()
	rows, err := db.Query("select id, submissionid, step, status, reason from submissiondetails where submissionid = ? order by id asc", submissionid)
	sds := make([]SubmissionDetail, 0)
	if err != nil {
		log.Print(err)
		return sds, err
	}
	defer rows.Close()
	for rows.Next() {
		var sd SubmissionDetail
		err := rows.Scan(&sd.Id, &sd.SubmissionId, &sd.Step, &sd.Verdict, &sd.Reason)
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
