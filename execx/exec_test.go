package execx

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	stdout, stderr, err := NewCommand("echo $GOPATH").Run()
	fmt.Println("stdout:", stdout)
	fmt.Println("stderr:", stderr)
	fmt.Println(err)
}

func TestOutputRun(t *testing.T) {
	cmd := "go_build.sh"
	err := NewCommand(cmd).Dir("./").OutputRun(func(line string) {
		fmt.Println(line)
	})
	fmt.Println(err)
}
