package filex

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

const (
	Overwrite = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	Append    = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	if content, err := os.ReadFile(path); err != nil {
		return nil, errorx.Wrap(err, "read file error")
	} else {
		return content, nil
	}
}

// ReadFileLine 按行读取
func ReadFileLine(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	// 按行处理txt
	reader := bufio.NewReader(file)
	var lines []string
	for {
		var line []byte
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	if err = file.Close(); err != nil {
		return nil, errorx.Wrap(err, "close file error")
	}
	return lines, nil
}

// Replace 内容替换
func Replace(path string, replaces map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errorx.Wrap(err, "read file error")
	}
	for k, v := range replaces {
		content = bytes.ReplaceAll(content, []byte(k), []byte(v))
	}
	if err = WriteFile(path, content); err != nil {
		return errorx.Wrap(err, "write to file error")
	}
	return nil
}

// WriteFile 写入文件
func WriteFile(path string, data []byte, flag ...int) error {
	file, err := Open(path, flag...)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	if _, err = file.Write(data); err != nil {
		return errorx.Wrap(err, "file write error")
	}
	if err = file.Close(); err != nil {
		return errorx.Wrap(err, "close file error")
	}
	return nil
}

// WriteFileString 写入文件
func WriteFileString(path, data string, flag ...int) error {
	file, err := Open(path, flag...)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	if _, err = file.WriteString(data); err != nil {
		return errorx.Wrap(err, "write string error")
	}
	if err = file.Close(); err != nil {
		return errorx.Wrap(err, "close file error")
	}
	return nil
}

// WriteFileLine 数组按行写入文件
func WriteFileLine(path string, content []string, flag ...int) error {
	file, err := Open(path, flag...)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	writer := bufio.NewWriter(file)
	for _, line := range content {
		_, _ = writer.WriteString(line)
		_, _ = writer.WriteString("\n")
	}
	if err = writer.Flush(); err != nil {
		return errorx.Wrap(err, "writer flush error")
	}
	if err = file.Close(); err != nil {
		return errorx.Wrap(err, "close file error")
	}
	return nil
}

// WriteCSV 写入csv文件
func WriteCSV(path string, data [][]string) error {
	file, err := Open(path)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.UseCRLF = true
	if err = writer.WriteAll(data); err != nil {
		return errorx.Wrap(err, "write csv to file error")
	}
	writer.Flush()
	if err = file.Close(); err != nil {
		return errorx.Wrap(err, "close file error")
	}
	return nil
}

// Open 打开文件
func Open(path string, flag ...int) (*os.File, error) {
	CreateIfNotExist(path)
	file, err := os.OpenFile(path, anyx.Default(Overwrite, flag...), 0644)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	return file, nil
}

// Clear 清空文件内容
func Clear(path string) {
	file, _ := os.OpenFile(path, os.O_TRUNC, 0644)
	_ = file.Close()
}

// MustOpen 强制打开文件
func MustOpen(dir string, name string) (*os.File, error) {
	path, err := filepath.Abs(filepath.Join(dir, name))
	if err != nil {
		return nil, errorx.Wrap(err, "abs path error")
	}
	if _, err = os.Stat(path); os.IsPermission(err) {
		return nil, errorx.Wrap(err, "file permission denied")
	}
	return Open(path, Append)
}

// SplitFile 拆分文件
func SplitFile(path string, size int) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	dir, filename, suffix := AnalysePath(path)
	dir = filepath.Join(dir, filename)
	reader := bufio.NewReader(file)
	count, index := 1, 1
	var sb = strings.Builder{}
	var paths []string
	for {
		if index < size {
			var line []byte
			if line, _, err = reader.ReadLine(); err == io.EOF {
				subpath := filepath.Join(dir, "split_"+strconv.Itoa(count)+suffix)
				if err = WriteFileString(subpath, sb.String()); err != nil {
					return nil, errorx.Wrap(err, "write file error")
				}
				paths = append(paths, subpath)
				break
			}
			sb.WriteString("\n")
			sb.Write(line)
		} else {
			index = 1
			subpath := filepath.Join(dir, "split_"+strconv.Itoa(count)+suffix)
			if err = WriteFileString(subpath, sb.String()); err != nil {
				return nil, errorx.Wrap(err, "write file error")
			}
			paths = append(paths, subpath)
			sb.Reset()
			count++
		}
		index++
	}
	if err = file.Close(); err != nil {
		return nil, errorx.Wrap(err, "close file error")
	}
	return paths, nil
}

// SplitPath 拆分路径为文件夹和文件名
func SplitPath(path string) (string, string) {
	if path != "" {
		if stringx.ContainsAny(path, "/", "\\") {
			return filepath.Split(path)
		} else {
			return "", path
		}
	}
	return "", ""
}

// AnalysePath 拆分路径为文件夹、文件名和文件后缀
func AnalysePath(path string) (dir, name, suffix string) {
	if dir, name = filepath.Split(path); name != "" {
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == '.' {
				name, suffix = name[:i], name[i:]
				return
			}
		}
	}
	return
}

// Pwd 获取绝对路径
func Pwd(path ...string) string {
	if len(path) == 0 {
		_, file, _, _ := runtime.Caller(1)
		pwd, _ := filepath.Split(file)
		return pwd
	} else {
		pwd, _ := filepath.Abs(path[0])
		return pwd
	}
}

// Depth 获取路径深度
func Depth(path string) int {
	return strings.Count(filepath.Clean(path), string(filepath.Separator))
}

// IsDir 判断是否文件夹
func IsDir(path string) bool {
	if fileInfo, err := os.Stat(path); err != nil {
		return false
	} else {
		return fileInfo.IsDir()
	}
}

// Exists 判断所给路径文件或文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// Create 创建文件
func Create(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errorx.Wrap(err, "create error")
	}
	if err = file.Close(); err != nil {
		return errorx.Wrap(err, "close file error")
	}
	return nil
}

// Copy 复制文件
func Copy(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return errorx.Wrap(err, "open source file error")
	}
	defer srcFile.Close()
	// 创建目标文件
	var dstFile *os.File
	if dstFile, err = os.Create(dst); err != nil {
		return errorx.Wrap(err, "create destination file error")
	}
	defer dstFile.Close()
	// 复制文件内容
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return errorx.Wrap(err, "copy file error")
	}

	return nil
}

// CreateIfNotExist 创建文件
func CreateIfNotExist(path string) {
	if !Exists(path) {
		CreateFile(path)
	}
}

// CreateFile 创建文件，但是不打开
func CreateFile(path string) {
	if dir, _ := filepath.Split(path); !Exists(dir) {
		// 先创建文件夹
		_ = os.MkdirAll(dir, os.ModePerm)
		// 再修改权限
		_ = os.Chmod(dir, os.ModePerm)
	}
	_ = Create(path)
}

// CreateDir 创建文件夹
func CreateDir(path string) {
	if !Exists(path) {
		dir, file := filepath.Split(path)
		if stringx.Index(file, ".") < 0 {
			dir = filepath.Join(dir, file)
		}
		// 先创建文件夹
		_ = os.MkdirAll(dir, os.ModePerm)
		// 再修改权限
		_ = os.Chmod(dir, os.ModePerm)
	}
}

// IsEmptyDir 检查给定的目录是否为空
func IsEmptyDir(dir string) bool {
	file, err := os.Open(dir)
	if err != nil {
		return false
	}
	// 读取目录内容
	var names []string
	if names, err = file.Readdirnames(0); err != nil {
		return false
	}
	_ = file.Close()
	// 如果目录内容为空，则目录为空
	return len(names) == 0
}

// SetSuffix 设置后缀
func SetSuffix(path string, suffix string) string {
	if path != "" {
		var name = path
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '.' {
				name = path[:i]
			}
			if path[i] == os.PathSeparator {
				break
			}
		}
		return name + stringx.AddPrefix(suffix, ".")
	}
	return ""
}

// GetSuffix 获取后缀
func GetSuffix(path string, withPoint ...bool) string {
	if path != "" {
		var p = 1
		if len(withPoint) > 0 && withPoint[0] {
			p = 0
		}
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '.' {
				return path[i+p:]
			}
		}
	}
	return ""
}

// FileName 获取文件名(不带后缀)
func FileName(path string) string {
	var fullName = filepath.Base(path)
	return strings.TrimSuffix(fullName, filepath.Ext(fullName))
}

type File struct {
	Path string
	Info os.FileInfo
}

// FileScan 文件扫描
func FileScan(dir string, keyword string) ([]*File, error) {
	var files []*File
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch {
		case keyword == "" || keyword == "*":
			files = append(files, &File{Path: path, Info: info})
		case keyword == "dir" && info.IsDir():
			files = append(files, &File{Path: path, Info: info})
		case keyword == "file" && !info.IsDir():
			files = append(files, &File{Path: path, Info: info})
		case stringx.Index(info.Name(), keyword) >= 0:
			files = append(files, &File{Path: path, Info: info})
		}
		return nil
	}); err != nil {
		return nil, errorx.Wrap(err, "file scan error")
	}
	return files, nil
}

// FileScanMatch 文件扫描匹配
func FileScanMatch(root string, match func(path string, info os.FileInfo) bool) ([]*File, error) {
	var files []*File
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if match(path, info) {
			files = append(files, &File{Path: path, Info: info})
		}
		return nil
	}); err != nil {
		return nil, errorx.Wrap(err, "file scan match error")
	}
	return files, nil
}
