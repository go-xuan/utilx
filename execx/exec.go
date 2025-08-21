package execx

import (
	"bytes"
	"io"
	"os/exec"
	"runtime"

	"github.com/go-xuan/utilx/errorx"
)

// Command 创建命令
func Command(command string) *Cmd {
	if runtime.GOOS == `windows` {
		return &Cmd{exec.Command("cmd", `/C`, command)}
	} else {
		return &Cmd{exec.Command("/bin/bash", `-c`, command)}
	}
}

// Cmd 命令
type Cmd struct {
	cmd *exec.Cmd
}

// Dir 设置命令执行目录
func (c *Cmd) Dir(dir string) *Cmd {
	c.cmd.Dir = dir
	return c
}

// Stdin 设置命令输入
func (c *Cmd) Stdin(in io.Reader) *Cmd {
	c.cmd.Stdin = io.NopCloser(in)
	return c
}

// Run 执行命令
func (c *Cmd) Run() (string, string, error) {
	if c.cmd == nil {
		return "", "", errorx.New("command instance is nil") // 更合适的错误信息
	}

	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	c.cmd.Stdout, c.cmd.Stderr = stdout, stderr
	err := c.cmd.Run()
	return stdout.String(), stderr.String(), err
}
