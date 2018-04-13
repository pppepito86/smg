package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Submission struct {
	Id                int64
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
	ProblemPoints     int
	Limit             Limit
	TestInfo          string
}

func AddSubmission(s Submission) (Submission, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO submissions(assignmentid, problemid, userid, language, sourcefile, verdict, reason) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return s, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.AssignmentId, s.ProblemId, s.UserId, s.Language, s.SourceFile, s.Verdict, s.Reason)
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
	rows, err := db.Query("select id, assignmentid, problemid, userid, language, sourcefile, time, verdict from submissions order by id desc")
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.AssignmentId, &s.ProblemId, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict)
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
	rows, err := db.Query("select submissions.id, submissions.userid, language, sourcefile, time, verdict, reason, problems.name, assignments.testinfo from submissions"+
		"	inner join problems on problems.id=submissions.problemid and submissions.id=?"+
		"	inner join assignments on assignments.id=submissions.assignmentid", submissionId)
	s := Submission{}
	if err != nil {
		log.Print(err)
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&s.Id, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Reason, &s.ProblemName, &s.TestInfo)
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
		"	inner join problems on problems.id=submissions.problemid and submissions.assignmentid=? and submissions.userid=?", assignmentId, userId)
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
	rows, err := db.Query("select submissions.id as id, time, assignmentid, problemid, submissions.points, problems.name, verdict, language from submissions " +
		"inner join problems on problems.id=submissions.problemid " +
			"where userid=?", userId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.Time, &s.AssignmentId, &s.ProblemId, &s.Points, &s.ProblemName, &s.Verdict, &s.Language)
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

type PointsPerWeek struct {
	Week   string
	Points int
}

func GetPointsPerWeek(userId int64) []PointsPerWeek {

	Response := make([]PointsPerWeek, 0)
	subs, _ := ListMyAllSubmissions(userId)
	if len(subs) == 0 {
		return Response
	}

	monday := func(t time.Time) time.Time {
		tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		tt = tt.AddDate(0, 0, -(int(tt.Weekday())+6)%7)
		return tt
	}

	currWeek := monday(subs[0].Time)
	nextWeek := currWeek.AddDate(0, 0, 7)

	problemPoints := make(map[int64]int, 0)
	totalPoints := 0
	subIdx := 0
	totalPointsLastWeek := 0
	for subIdx < len(subs) {
		currWeekResponse := PointsPerWeek{
			currWeek.String()[:10], 0,
		}

		for subIdx < len(subs) && subs[subIdx].Time.After(currWeek) && subs[subIdx].Time.Before(nextWeek) {
			lastPts, _ := problemPoints[subs[subIdx].ProblemId]
			currPts := subs[subIdx].Points
			if currPts > lastPts {
				totalPoints += currPts - lastPts
				problemPoints[subs[subIdx].ProblemId] = currPts
			}

			subIdx++
		}

		currWeekResponse.Points = totalPoints - totalPointsLastWeek
		totalPointsLastWeek = totalPoints

		Response = append(Response, currWeekResponse)
		// add totalPoints for current week
		currWeek = nextWeek
		nextWeek = nextWeek.AddDate(0, 0, 7)
	}

	return Response
}

func ListSubmissionsForAssignment(assignmentId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, submissions.problemid, language, sourcefile, time, verdict, submissions.points, problems.name, users.id, users.username, users.firstname, users.lastname from submissions"+
		"	inner join problems on problems.id=submissions.problemid and submissions.assignmentid=?"+
		"	inner join users on users.id=submissions.userid", assignmentId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ProblemId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Points, &s.ProblemName, &s.UserId, &s.User, &s.FirstName, &s.LastName)
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

func ListAcceptedSubmissionsForAssignment(assignmentId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, submissions.problemid, language, sourcefile, time, verdict, submissions.points, problems.name, users.id, users.username, users.firstname, users.lastname from submissions"+
		"	inner join problems on problems.id=submissions.problemid and submissions.verdict='Accepted' and submissions.assignmentid=?"+
		"	inner join users on users.id=submissions.userid", assignmentId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	defer rows.Close()
	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ProblemId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Points, &s.ProblemName, &s.UserId, &s.User, &s.FirstName, &s.LastName)
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

func ListMySubmissionsForProblem(userId, aId, pId int64) ([]Submission, error) {
	db := getConnection()
	rows, err := db.Query("select submissions.id, language, sourcefile, verdict from submissions"+
		" where userid=? and assignmentid=? and problemid=?", userId, aId, pId)
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
	rows, err := db.Query("select submissions.id, submissions.problemid, submissions.userid, submissions.language, submissions.sourcefile, submissions.time, submissions.verdict, submissions.reason, assignmentproblems.points from submissions"+
		" inner join assignmentproblems"+
		" where assignmentproblems.problemid=submissions.problemid and assignmentproblems.assignmentid=submissions.assignmentid"+
		" and submissions.problemid=?", pId)
	if err != nil {
		log.Print(err)
		return []Submission{}, err
	}
	submissions := make([]Submission, 0)
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.Id, &s.ProblemId, &s.UserId, &s.Language, &s.SourceFile, &s.Time, &s.Verdict, &s.Reason, &s.ProblemPoints)
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

	if len(reason) > 1000 {
		reason = reason[0:1000]
	}

	stmt, err := db.Prepare("update submissions set verdict=?, reason=?, correct=?, total=?, points=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

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
	defer stmt.Close()

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
	defer stmt.Close()

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
