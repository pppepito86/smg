package submissions

import (
	"bytes"
	"db"
	"errors"
	"fmt"
	"io"
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
		err := test(s, compiledFile, testsDir, i)
		if err == nil {
			correct++
		}
	}
	if correct == tests && tests > 0 {
		db.UpdateVerdict(s.Id, "Accepted", "")
	} else {
		db.UpdateVerdict(s.Id, fmt.Sprintf("%d/%d", correct, tests), "")
	}
}

func test(s db.Submission, compiledFile, testsDir string, testCase int) error {
	cmd := exec.Command("./test")
	if s.Language == "java" {
		cmd = exec.Command("java", filepath.Base(compiledFile))
	}
	cmd.Dir = filepath.Dir(compiledFile)
	inPipeTest, _ := cmd.StdinPipe()
	outPipeTest, _ := cmd.StdoutPipe()
	fmt.Println("testsDir", testsDir)
	in, err := os.Open(filepath.Join(testsDir, fmt.Sprintf("input%d", testCase)))
	if err != nil {
		return err
	}
	defer in.Close()
	io.Copy(inPipeTest, in)
	inPipeTest.Close()
	cmd.Start()

	ch := make(chan []byte, 1)

	go func() {
		out, _ := ioutil.ReadAll(outPipeTest)
		cmd.Wait()
		ch <- out
	}()
	select {
	case res := <-ch:
		{
			fmt.Println("***" + string(res) + "***")
			realOut, _ := ioutil.ReadFile(filepath.Join(testsDir, fmt.Sprintf("output%d", testCase)))
			fmt.Println("***" + string(realOut) + "***")
			if bytes.Equal(res, realOut) {
				return nil
			} else {
				return errors.New("WA")
			}
		}
	case <-time.After(time.Second * 2):
		{
			cmd.Process.Kill()
			return errors.New("TLE")
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
