package execx

import (
	"bufio"
	"io"
	"os/exec"
	"runtime"

	"github.com/go-xuan/utilx/errorx"
)

// NewCommand 创建命令
func NewCommand(command string) *Command {
	if runtime.GOOS == "windows" {
		return &Command{exec.Command("cmd", "/C", command)}
	}
	return &Command{exec.Command("/bin/bash", "-c", command)}
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
func (c *Command) Stdin(reader io.Reader) *Command {
	c.cmd.Stdin = io.NopCloser(reader)
	return c
}

// Stdout 设置命令输出
func (c *Command) Stdout(writer io.Writer) *Command {
	c.cmd.Stdout = writer
	return c
}

// Stderr 设置命令错误输出
func (c *Command) Stderr(writer io.Writer) *Command {
	c.cmd.Stderr = writer
	return c
}

// Run 执行命令
func (c *Command) Run() error {
	if c.cmd == nil {
		return errorx.New("command instance is nil")
	}
	if err := c.cmd.Run(); err != nil {
		return errorx.Wrap(err, "command run error")
	}
	return nil
}

// RunRealTime 执行命令并实时输出结果（不要设置Stdout和Stderr）
func (c *Command) RunRealTime(realtime func(b []byte)) error {
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
	go pipeOutput(stdoutPipe, realtime)
	go pipeOutput(stderrPipe, realtime)

	// 等待命令执行完成
	if err = c.cmd.Wait(); err != nil {
		return errorx.Wrap(err, "command wait error")
	}
	return nil
}

// pipeOutput 读取管道输出
func pipeOutput(pipe io.ReadCloser, realtime func(b []byte)) {
	defer errorx.Close(pipe)
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		realtime(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		realtime([]byte("pipe output scan error: \n"))
		realtime([]byte(err.Error()))
	}
}
