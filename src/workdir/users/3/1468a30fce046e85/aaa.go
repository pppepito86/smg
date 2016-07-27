package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	test := exec.Command("./test")
	outPipeTest, _ := test.StdoutPipe()
	test.Start()
	out, _ := ioutil.ReadAll(outPipeTest)
	test.Wait()
	fmt.Println(string(out))
}
