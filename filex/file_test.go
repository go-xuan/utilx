package filex

import (
	"fmt"
	"testing"
)

func TestFileSplit(t *testing.T) {
	filePath := "./nohup.log"
	fmt.Println(Analyse(filePath))
	files, err := FileSplit(filePath, 5000)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}

func TestFileWrite(t *testing.T) {
	_ = WriteFileString("nohup.log", "1111\n", Append)
	_ = WriteFileString("nohup.log", "2222\n", Append)
	_ = WriteFileString("nohup.log", "3333\n", Append)
	_ = WriteFileLine("nohup.log", []string{
		"4444",
		"5555",
		"6666",
	}, Append)
	_ = WriteFileLine("nohup.log", []string{
		"7777",
		"8888",
		"9999",
	}, Append)
}
