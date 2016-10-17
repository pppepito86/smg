package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Problem struct {
	Id          int64
	ProblemName string
	Version     string
	Description string
    Tags        []string
	Visibility  string
	Languages   string
	AuthorId    int64
	Author      string
	Tests       string
	Points      int
	LangLimits  Limits
}

type Limit struct {
	Language    string
	TimeLimit   int
	MemoryLimit int
}

type Limits map[string]Limit

func CreateProblem(p Problem) (Problem, error) {
	db := getConnection()

	stmt, err := db.Prepare("INSERT INTO problems(name, version, description, languages, visibility, author, points) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return p, err
	}

	res, err := stmt.Exec(p.ProblemName, p.Version, p.Description, p.Languages, p.Visibility, p.AuthorId, p.Points)
	if err != nil {
		log.Print(err)
		return p, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return p, err
	}

	p.Id = lastId
    
    for _, tag := range p.Tags {
        TagProblem(p.Id, tag)
    }
    
    
	return p, nil
}

func ListProblems() ([]Problem, error) {
	db := getConnection()
	rows, err := db.Query("select problems.id, problems.name, problems.version, problems.points, problems.description, problems.languages, problems.visibility, problems.author, users.username from problems" +
		" inner join users on problems.author = users.id")
	if err != nil {
		log.Print(err)
		return []Problem{}, err
	}
	defer rows.Close()
	problems := make([]Problem, 0)
	for rows.Next() {
		var p Problem
		err := rows.Scan(&p.Id, &p.ProblemName, &p.Version, &p.Points, &p.Description, &p.Languages, &p.Visibility, &p.AuthorId, &p.Author)
		if err != nil {
			log.Print(err)
			return []Problem{}, err
		}
        
        p.Tags, err = GetListOfTags(p.Id)
        if err != nil {
			log.Print(err)
			return []Problem{}, err
		}
		
        problems = append(problems, p)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return []Problem{}, err
	}
    
    
    
	return problems, nil
}

func GetProblem(id int64) (Problem, error) {
	db := getConnection()
	rows, err := db.Query("select id, name, version, points, description, languages, visibility, author from problems where id=?", id)
	if err != nil {
		log.Print(err)
		return Problem{}, nil
	}
	defer rows.Close()
	p := Problem{}
	for rows.Next() {
		err := rows.Scan(&p.Id, &p.ProblemName, &p.Version, &p.Points, &p.Description, &p.Languages, &p.Visibility, &p.AuthorId)
		if err != nil {
			log.Print(err)
			return Problem{}, nil
		}
	}
	if err != nil {
		log.Print(err)
		return Problem{}, nil
	}
    
    p.Tags, err = GetListOfTags(p.Id)
    if err != nil {
        log.Print(err)
        return Problem{}, err
    }
    
	return p, nil
}

func TagProblem(id int64, tag string) error {
    db := getConnection()
    
    stmt, err := db.Prepare("INSERT INTO tags(problemid, tag) VALUES(?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(id, tag)
	if err != nil {
		log.Print(err)
		return err
	}
    
    return nil
}

func GetListOfTags(problemid int64) ([]string, error) {
    db := getConnection()
	rows, err := db.Query("select tag from tags where problemid=?", problemid)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
    
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
        return []string{}, err
	}
	return tags, nil
}

func RemoveTagsForProblem(problemid int64) error {
    db := getConnection()
    
    stmt, err := db.Prepare("DELETE FROM tags WHERE problemid=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(problemid)
	if err != nil {
		log.Print(err)
		return err
	}
    
    return nil
}

func UpdateProblem(p Problem) error {
	db := getConnection()

	stmt, err := db.Prepare("update problems set name=?, version=?, description=?, languages=?, points=? where id=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(p.ProblemName, p.Version, p.Description, p.Languages, p.Points, p.Id)
	if err != nil {
		log.Print(err)
		return err
	}
    
    err = RemoveTagsForProblem(p.Id)
    if err != nil {
		log.Print(err)
		return err
	}
    
    for _, tag := range p.Tags {
        TagProblem(p.Id, tag)
    }

	return nil
}

/*
func GetUser(username string) User {
	fmt.Println("find user:", username)
	db := getConnection()
	rows, err := db.Query("select users.id, roleid, username, firstname, lastname, email, passwordhash, passwordsalt, isdisabled, rolename from users inner join roles on username= ? and isdisabled=? and users.roleid=roles.id", username, false)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.RoleId, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.RoleName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("found user:", user)
		return user
	}
	return User{}
}

func UpdateUserRole(userId, roleId int64) error {
	db := getConnection()

	stmt, err := db.Prepare("update users set roleid=? where id=?")
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(roleId, userId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
*/
