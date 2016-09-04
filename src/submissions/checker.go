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
		status, reason, _ := test(s, compiledFile, testsDir, i)
		if status == "ok" {
			correct++
		}
		db.AddSubmissionDetails(s.Id, "Test #"+strconv.Itoa(i), status, reason)
	}
	if correct == tests && tests > 0 {
		db.UpdateVerdict(s.Id, "Accepted", "")
	} else {
		db.UpdateVerdict(s.Id, fmt.Sprintf("%d/%d", correct, tests), "")
	}
}

func test(s db.Submission, compiledFile, testsDir string, testCase int) (string, string, error) {
	cmd := exec.Command("./test")
	//step := "Test #" + strconv.Itoa(testCase)
	if s.Language == "java" {
		cmd = exec.Command("java", filepath.Base(compiledFile))
		//pwd, _ := os.Getwd()
		//dir := filepath.Join(pwd, filepath.Dir(compiledFile))
		//cmd = exec.Command("docker", "run", "-v", dir+":/foo", "-w", "/foo", "-i", "--read-only", "-m", "128M", "pppepito86/judgebox", "java", filepath.Base(compiledFile))
	}
	cmd.Dir = filepath.Dir(compiledFile)
	inPipeTest, _ := cmd.StdinPipe()
	outPipeTest, _ := cmd.StdoutPipe()
	errPipeTest, _ := cmd.StderrPipe()
	fmt.Println("testsDir", testsDir)
	in, err := os.Open(filepath.Join(testsDir, fmt.Sprintf("input%d", testCase)))
	if err != nil {
		return "system error", "cannot read input", err
	}
	go func() {
		defer in.Close()
		io.Copy(inPipeTest, in)
		inPipeTest.Close()
	}()

	cmd.Start()
	chError := make(chan error, 2)
	chOutput := make(chan []byte, 2)

	go func() {
		out, _ := ioutil.ReadAll(outPipeTest)
		errOut, _ := ioutil.ReadAll(errPipeTest)
		err = cmd.Wait()
		if err != nil {
			fmt.Println("*****Error", err.Error())
			fmt.Println("e", string(errOut))
		}
		chError <- err
		chOutput <- out
		chOutput <- errOut
	}()
	select {
	case err = <-chError:
		{
			res := <-chOutput
			errOut := <-chOutput
			fmt.Println("***" + string(res) + "***")
			if err != nil {
				return "runtime error", err.Error() + " - " + string(errOut), nil
			}
			realOut, _ := ioutil.ReadFile(filepath.Join(testsDir, fmt.Sprintf("output%d", testCase)))
			fmt.Println("***" + string(realOut) + "***")
			if bytes.Equal(res, realOut) {
				return "ok", "", nil
			} else {
				return "wrong answer", "", nil
			}
		}
	case <-time.After(time.Second * 5):
		{
			cmd.Process.Kill()
			return "time limit exceeded", "", nil
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
