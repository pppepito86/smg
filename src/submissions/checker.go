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
	testsDir := filepath.Join("workdir", "problems", strconv.FormatInt(s.ProblemId, 10))

	files, _ := ioutil.ReadDir(testsDir)
	tests := len(files) / 2

	if s.Limit.Language == "" {
		db.UpdateVerdict(s.Id, "Language Not Allowed", "", 0, 0, 0)
		return
	}

	if tests == 0 {
		db.UpdateVerdict(s.Id, "System error", "Missing test cases", 0, 0, 0)
		return
	}

	compiledFile, err := compile(s)
	correct := 0

	if err != nil {
		fmt.Println("compilation error: " + err.Error())
		db.UpdateVerdict(s.Id, "Compilation Failed", err.Error(), 0, tests, 0)
		return
	}

	fmt.Println("compilation successful")
	db.UpdateVerdict(s.Id, "Compiled", "", 0, tests, 0)

	for i := 1; i <= tests; i++ {
		db.UpdateVerdict(s.Id, "Running test #"+strconv.Itoa(i), "", correct, tests, correct*s.ProblemPoints/tests)
		status, reason, time, _ := test(s, compiledFile, testsDir, i, s.Limit)
		if status == "ok" {
			correct++
		}
		db.AddSubmissionDetails(s.Id, "Test #"+strconv.Itoa(i), status, reason, time)
	}
	if correct == tests && tests > 0 {
		db.UpdateVerdict(s.Id, "Accepted", "", correct, tests, s.ProblemPoints)
	} else {
		db.UpdateVerdict(s.Id, fmt.Sprintf("%d/%d", correct, tests), "", correct, tests, correct*s.ProblemPoints/tests)
	}
}

func test(s db.Submission, compiledFile, testsDir string, testCase int, limit db.Limit) (string, string, int64, error) {
	cmdArg := "./test"
	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, filepath.Dir(compiledFile))
	if s.Language == "java" {
		cmdArg = "java " + filepath.Base(compiledFile)
	} else if s.Language == "nodejs" {
		cmdArg = "node " + filepath.Base(compiledFile)
	}
	testStr := strconv.Itoa(testCase)
	cmdArg = "cat input" + testStr + "|" + cmdArg + ">output" + testStr + " 2>error" + testStr
	mLimit := strconv.Itoa(limit.MemoryLimit) + "M"
	cmd := exec.Command("docker", "run", "--cidfile", "cid", "-v", dir+":/foo", "-w", "/foo", "-i", "--read-only", "-m", mLimit, "--network", "none", "pppepito86/judgebox", "/bin/bash", "-c", cmdArg)
	cmd.Dir = filepath.Dir(compiledFile)
	exec.Command("rm", filepath.Join(cmd.Dir, "cid")).Run()
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
	case <-time.After(time.Millisecond * time.Duration(limit.TimeLimit)):
		{
			cmd.Process.Kill()
			cid, _ := ioutil.ReadFile(filepath.Join(cmd.Dir, "cid"))
			exec.Command("docker", "kill", string(cid)).Run()
			durationTime := time.Since(startTime).Nanoseconds() / 1e6
			return "time limit exceeded", "", durationTime, nil
		}
	}
}

func compile(s db.Submission) (string, error) {
	destFile := filepath.Join(filepath.Dir(s.SourceFile), "test")
	cmdArg := ""
	if s.Language == "java" {
		cmdArg = "javac " + s.SourceFile
		destFile = strings.Replace(s.SourceFile, ".java", "", 1)
	} else if s.Language == "c++" {
		cmdArg = "g++  -O2 -std=c++11 -o " + destFile + " " + s.SourceFile
	} else if s.Language == "nodejs" {
		return s.SourceFile, nil
	} else {
		return "", errors.New("Language is not supported")
	}

	cmd := exec.Command("docker", "run", "-v", filepath.Dir(destFile)+":/foo", "-w", "/foo", "--network", "none", "pppepito86/judgebox", "/bin/bash", "-c", cmdArg)
	cmd.Dir = filepath.Dir(destFile)
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
	exec.Command("rm", filepath.Join(cmd.Dir, "cid")).Run()
	return destFile, nil
}
