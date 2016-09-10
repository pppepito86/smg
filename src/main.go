package main

import (
	"db"
	"fmt"
	"log"
	"net/http"
	"request"
	"submissions"
)

func main() {
	db.OpenConnection()
	defer db.Close()

	go submissions.Checker()
	fmt.Println("server started")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := request.GetUser(*r)
	if user.RoleName == "admin" {
		request.HandleAdmin(w, r, user)
	} else if user.RoleName == "teacher" {
		request.HandleTeacher(w, r, user)
	} else if user.RoleName == "user" {
		request.HandleUser(w, r, user)
	} else {
		request.HandleGuest(w, r)
	}
	/*
		if path == "/login.html" {
			http.ServeFile(w, r, "../login.html")
		} else if path == "/login" {
			login(w, r)
		} else if path == "/register" {
			register(w, r)
		} else if userId == "" {
			http.Redirect(w, r, "/login.html", http.StatusFound)
		} else if r.Method == "POST" {
			if path == "/joingroup" {
				joinGroup(w, r)
			}
		} else {
			indexHtml(w, r)
		}
	*/
}

/*
func addGroupHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve addgroup.html")
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/addgroup.html")
	t.Execute(w, nil)
}

func joinGroupHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve joingroup.html")
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/joingroup.html")
	t.Execute(w, nil)
}

func groupsHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve groups.html")
	userId := getUserId(*r)
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/groups.html")
	id, _ := strconv.ParseInt(userId, 10, 64)
	t.Execute(w, db.ListGroupsForUser(id))
}

func indexHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve index.html")
	http.ServeFile(w, r, "../admin/admin.html")
}

func adminHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve admin.html")
	http.ServeFile(w, r, "../admin/admin.html")
}

func joinGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server joingroup")

	userId := getUserId(*r)
	id, _ := strconv.ParseInt(userId, 10, 64)
	r.ParseForm()
	group, _ := db.FindGroupByName(r.Form["groupname"][0])
	userGroup, _ := db.CreateUserGroup(id, group.Id)

	fmt.Println(userGroup)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func addGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server addgroup")
	if r.Method != "POST" {
		return
	}

	userId := getUserId(*r)
	id, _ := strconv.ParseInt(userId, 10, 64)
	r.ParseForm()
	group := db.Group{
		-1,
		r.Form["groupname"][0],
		r.Form["description"][0],
		id,
	}
	group, _ = db.CreateGroup(group)

	fmt.Println(group)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}
*/

/*
func handler2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		source := r.Form["source"][0]
		file := "test.cpp"
		ioutil.WriteFile(file, []byte(source), 0777)
		cmd := exec.Command("g++", "-o", "test", file)
		outPipe, _ := cmd.StdoutPipe()
		errPipe, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		_, _ = ioutil.ReadAll(outPipe)
		errStr, _ := ioutil.ReadAll(errPipe)
		err = cmd.Wait()
		if err != nil {
			io.WriteString(w, "error: "+string(errStr))
		} else {
			//io.WriteString(w, "ok: "+string(outStr))
			test := exec.Command("./test")
			inPipeTest, _ := test.StdinPipe()
			outPipeTest, _ := test.StdoutPipe()
			test.Start()
			inPipeTest.Write([]byte("5 7"))
			inPipeTest.Close()
			out, _ := ioutil.ReadAll(outPipeTest)
			test.Wait()
			io.WriteString(w, "out: "+string(out))
		}
	}
}
*/
