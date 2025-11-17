package execx

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"runtime"

	"github.com/go-xuan/utilx/errorx"
)

// NewCommand 创建命令
func NewCommand(command string) *Command {
	if runtime.GOOS == `windows` {
		return &Command{exec.Command("cmd", `/C`, command)}
	} else {
		return &Command{exec.Command("/bin/bash", `-c`, command)}
	}
}

// Command 命令
type Command struct {
	cmd *exec.Cmd
}

// Dir 设置命令执行目录
func (c *Command) Dir(dir string) *Command {
	c.cmd.Dir = dir
	return c
}

// Stdin 设置命令输入
func (c *Command) Stdin(in io.Reader) *Command {
	c.cmd.Stdin = io.NopCloser(in)
	return c
}

// Run 执行命令
func (c *Command) Run() (string, string, error) {
	if c.cmd == nil {
		return "", "", errorx.New("command instance is nil")
	}

	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	c.cmd.Stdout, c.cmd.Stderr = stdout, stderr
	err := c.cmd.Run()
	return stdout.String(), stderr.String(), err
}

// OutputRun 执行命令并输出结果
func (c *Command) OutputRun(output func(line string)) error {
	if c.cmd == nil {
		return errorx.New("command instance is nil")
	}

	var err error
	var stdoutPipe, stderrPipe io.ReadCloser
	if stdoutPipe, err = c.cmd.StdoutPipe(); err != nil {
		return errorx.Wrap(err, "stdout pipe error")
	}
	if stderrPipe, err = c.cmd.StderrPipe(); err != nil {
		return errorx.Wrap(err, "stderr pipe error")
	}

	// 启动命令
	if err = c.cmd.Start(); err != nil {
		return errorx.Wrap(err, "command start error")
	}

	// 启动 goroutine 读取输出
	go pipeOutput(stdoutPipe, output)
	go pipeOutput(stderrPipe, output)

	// 等待命令执行完成
	if err = c.cmd.Wait(); err != nil {
		return errorx.Wrap(err, "command wait error")
	}
	return nil
}

// pipeOutput 读取管道输出
func pipeOutput(pipe io.ReadCloser, output func(string)) {
	defer pipe.Close()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		output(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		output("pipe output scan error: " + err.Error())
	}
}
