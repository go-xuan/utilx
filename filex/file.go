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

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

const (
	Overwrite = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	Append    = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.Wrap(err, "read file error")
	}
	return content, nil
}

// ReadFileLine 按行读取
func ReadFileLine(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	reader := bufio.NewReader(file)
	var lines []string
	for {
		var line []byte
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

// ReadFileLine2Writer 按行读取文件内容并写入到writer
func ReadFileLine2Writer(path string, writer io.Writer) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	reader := bufio.NewReader(file)
	for {
		var line []byte
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			return errorx.Wrap(err, "read line error")
		}
		_, _ = writer.Write(line)
	}
	return nil
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
	defer errorx.Close(file)

	if _, err = file.Write(data); err != nil {
		return errorx.Wrap(err, "file write error")
	}
	return nil
}

// WriteFileString 写入文件
func WriteFileString(path, data string, flag ...int) error {
	file, err := Open(path, flag...)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	if _, err = file.WriteString(data); err != nil {
		return errorx.Wrap(err, "write string error")
	}
	return nil
}

// WriteFileLine 数组按行写入文件
func WriteFileLine(path string, lines []string, flag ...int) error {
	file, err := Open(path, flag...)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = writer.WriteString(line)
		_, _ = writer.WriteString("\n")
	}
	if err = writer.Flush(); err != nil {
		return errorx.Wrap(err, "writer flush error")
	}
	return nil
}

// WriteCSV 写入csv文件
func WriteCSV(path string, data [][]string) error {
	file, err := Open(path)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.UseCRLF = true
	if err = writer.WriteAll(data); err != nil {
		return errorx.Wrap(err, "write csv to file error")
	}
	writer.Flush()
	return nil
}

// Open 打开文件
func Open(path string, flag ...int) (*os.File, error) {
	CreateIfNotExist(path)
	f := Overwrite
	if len(flag) > 0 {
		f = flag[0]
	}
	file, err := os.OpenFile(path, f, 0644)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	return file, nil
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

// Clear 清空文件内容
func Clear(path string) {
	file, _ := os.OpenFile(path, os.O_TRUNC, 0644)
	defer errorx.Close(file)
}

// SplitFile 拆分文件
func SplitFile(path string, size int) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	dir, filename, suffix := AnalysePath(path)
	dir = filepath.Join(dir, filename)
	reader := bufio.NewReader(file)
	count, index := 1, 1
	sb := strings.Builder{}
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
			} else if err != nil {
				return paths, errorx.Wrap(err, "read line error")
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

	return paths, nil
}

// SplitPath 拆分路径为文件夹和文件名
func SplitPath(path string) (string, string) {
	if path != "" {
		if stringx.ContainsAny(path, "/", "\\") {
			return filepath.Split(path)
		}
		return "", path
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
	}
	pwd, _ := filepath.Abs(path[0])
	return pwd
}

// Depth 获取路径深度
func Depth(path string) int {
	return strings.Count(filepath.Clean(path), string(filepath.Separator))
}

// IsDir 判断是否文件夹
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
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
	errorx.Close(file)
	return nil
}

// Copy 复制文件
func Copy(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return errorx.Wrap(err, "open file error")
	}
	defer errorx.Close(file)

	var cp *os.File
	if cp, err = os.Create(dst); err != nil {
		return errorx.Wrap(err, "create copy file error")
	}
	defer errorx.Close(cp)

	if _, err = io.Copy(cp, file); err != nil {
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
func IsEmptyDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer errorx.Close(file)

	var names []string
	if names, err = file.Readdirnames(0); err != nil {
		return false
	}
	return len(names) == 0
}

// GetSuffix 获取后缀
func GetSuffix(path string) string {
	if path != "" {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '.' {
				return path[i+1:]
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
	Path string      // 文件路径
	Info os.FileInfo // 文件信息
}

// FileScan 文件扫描
func FileScan(dir string, match string) ([]*File, error) {
	return FileScanMatch(dir, func(path string, info os.FileInfo) bool {
		switch {
		case match == "" || match == "*":
			return true
		case match == "dir" && info.IsDir():
			return true
		case match == "file" && !info.IsDir():
			return true
		case stringx.Index(info.Name(), match) >= 0:
			return true
		}
		return false
	})
}

// FileScanMatch 文件扫描匹配（包括下级目录）
func FileScanMatch(root string, match func(string, os.FileInfo) bool) ([]*File, error) {
	var files []*File
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if match(path, info) {
			files = append(files, &File{
				Path: path,
				Info: info,
			})
		}
		return nil
	}); err != nil {
		return nil, errorx.Wrap(err, "file scan error")
	}
	return files, nil
}

// DirScanMatch 文件夹扫描匹配（仅当前目录）
func DirScanMatch(dir string, match func(os.DirEntry) bool) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []os.DirEntry
	for _, entry := range entries {
		if match(entry) {
			result = append(result, entry)
		}
	}
	return result, nil
}
