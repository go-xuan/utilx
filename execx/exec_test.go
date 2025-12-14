package execx

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	cmd := "curl -o heap-localhost-20251201155636.pprof http://localhost:6060/debug/pprof/heap"
	var outWriter, errWriter = new(bytes.Buffer), new(bytes.Buffer)
	if err := NewCommand(cmd).Stdout(outWriter).Stderr(errWriter).Run(); err != nil {
		t.Error(err)
	}
	fmt.Println("stdout:", outWriter.String())
	fmt.Println("stderr:", errWriter.String())
}
