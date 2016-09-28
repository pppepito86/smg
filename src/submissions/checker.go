package submissions

import (
	"bytes"
	"db"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Checker() {
	go func() {
		for {
			s := Pop()
			testSubmission(s)
		}
	}()
}

func testSubmission(s db.Submission) {
	fmt.Println("Testing: " + s.SourceFile)
	db.DeleteSubmissionDetails(s.Id)
	compiledFile, err := compile(s)
	if err != nil {
		fmt.Println("compilation error: " + err.Error())
		db.UpdateVerdict(s.Id, "Compilation Failed", err.Error())
		return
	}

	fmt.Println("compilation successful")
	db.UpdateVerdict(s.Id, "Compiled", "")

	ap, _ := db.GetAssignmentProblem(s.ApId)
	testsDir := filepath.Join("problems", strconv.FormatInt(ap.ProblemId, 10))
	files, err := ioutil.ReadDir(testsDir)
	tests := len(files) / 2
	correct := 0
	for i := 1; i <= tests; i++ {
		db.UpdateVerdict(s.Id, "Running test #"+strconv.Itoa(i), "")
		status, reason, time, _ := test(s, compiledFile, testsDir, i)
		if status == "ok" {
			correct++
		}
		db.AddSubmissionDetails(s.Id, "Test #"+strconv.Itoa(i), status, reason, time)
	}
	if correct == tests && tests > 0 {
		db.UpdateVerdict(s.Id, "Accepted", "")
	} else {
		db.UpdateVerdict(s.Id, fmt.Sprintf("%d/%d", correct, tests), "")
	}
}

func test(s db.Submission, compiledFile, testsDir string, testCase int) (string, string, int64, error) {
	cmdArg := "./test"
	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, filepath.Dir(compiledFile))
	if s.Language == "java" {
		cmdArg = "java " + filepath.Base(compiledFile)
	}
	testStr := strconv.Itoa(testCase)
	cmdArg = "cat input" + testStr + "|" + cmdArg + ">output" + testStr + " 2>error" + testStr
	cmd := exec.Command("docker", "run", "-v", dir+":/foo", "-w", "/foo", "-i", "--read-only", "-m", "64M", "--network", "none", "pppepito86/judgebox", "/bin/bash", "-c", cmdArg)
	cmd.Dir = filepath.Dir(compiledFile)
	err := exec.Command("cp",
		filepath.Join(testsDir, fmt.Sprintf("input%d", testCase)),
		filepath.Join(filepath.Dir(compiledFile), fmt.Sprintf("input%d", testCase))).
		Run()
	if err != nil {
		return "system error", "cannot read input", 0, err
	}

	cmd.Start()
	startTime := time.Now()
	chError := make(chan error, 1)
	chTime := make(chan int64, 1)

	go func() {
		err = cmd.Wait()
		chTime <- time.Since(startTime).Nanoseconds() / 1e6
		if err != nil {
			fmt.Println("*****Error", err.Error())
		}
		chError <- err
	}()

	select {
	case err = <-chError:
		{
			durationTime := <-chTime
			if err != nil {
				realErr, _ := ioutil.ReadFile(filepath.Join(filepath.Dir(compiledFile), fmt.Sprintf("error%d", testCase)))
				return "runtime error", err.Error() + " - " + string(realErr), durationTime, nil
			}
			res, _ := ioutil.ReadFile(filepath.Join(filepath.Dir(compiledFile), fmt.Sprintf("output%d", testCase)))
			realOut, _ := ioutil.ReadFile(filepath.Join(testsDir, fmt.Sprintf("output%d", testCase)))
			if len(res) < 1000 {
				//	fmt.Println("***" + string(realOut) + "***")
			}
			if bytes.Equal(res, realOut) {
				return "ok", "", durationTime, nil
			} else {
				if strings.TrimSpace(string(res)) == strings.TrimSpace(string(realOut)) {
					return "presentation error", "", durationTime, nil
				}
				return "wrong answer", "", durationTime, nil
			}
		}
	case <-time.After(time.Second * 5):
		{
			cmd.Process.Kill()
			durationTime := time.Since(startTime).Nanoseconds() / 1e6
			return "time limit exceeded", "", durationTime, nil
		}
	}
}

func compile(s db.Submission) (string, error) {
	destFile := filepath.Join(filepath.Dir(s.SourceFile), "test")
	var cmd *exec.Cmd
	if s.Language == "java" {
		cmd = exec.Command("javac", s.SourceFile)
		destFile = strings.Replace(s.SourceFile, ".java", "", 1)
	} else if s.Language == "c++" {
		cmd = exec.Command("g++", "-o", destFile, s.SourceFile)
	} else {
		return "", errors.New("Language is not supported")
	}
	errPipe, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	errStr, _ := ioutil.ReadAll(errPipe)
	err = cmd.Wait()
	if err != nil {
		return "", errors.New(string(errStr))
	}
	return destFile, nil
}
