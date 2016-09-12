package submissions

import (
	"db"
	"sync"
)

var mutex = &sync.Mutex{}
var submissions = make(chan db.Submission, 100)

func Push(s db.Submission) {
	mutex.Lock()
	defer mutex.Unlock()
	submissions <- s
}

func Pop() db.Submission {
	return <-submissions
}
