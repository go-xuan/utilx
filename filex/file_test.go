package filex

import (
	"fmt"
	"os"
	"testing"
)

func TestFileSplit(t *testing.T) {
	filePath := "./nohup.log"
	fmt.Println(Analyse(filePath))
	files, err := FileSplit(filePath, 5000)
	if err != nil {
		t.Log(err)
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

func TestFileScan(t *testing.T) {
	files, err := FileScan("../", "go")
	if err != nil {
		t.Log(err)
	}
	for _, file := range files {
		fmt.Printf("path:%s size:%dB mode:%s fname:%s \n", file.Path, file.Info.Size(), file.Info.Mode(), file.Info.Name())
	}
}

func TestFileScanMatch(t *testing.T) {
	files, err := FileScanMatch("../", func(path string, info os.FileInfo) bool {
		return info.IsDir()
	})
	if err != nil {
		t.Log(err)
	}

	for _, file := range files {
		fmt.Printf("path:%s size:%dB mode:%s fname:%s \n", file.Path, file.Info.Size(), file.Info.Mode(), file.Info.Name())
	}
}
