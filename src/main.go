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
	request.Route(w, r)
}
