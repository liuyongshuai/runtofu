/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-03-27 15:20
 */
package negoutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

/*
*
***********************************************************************
因为ES经常有如备份数据到文件、从文件中恢复数据的操作，此处写了个方便遍历文件的小东东
***********************************************************************
*/
type FileIterator struct {
	file string        //文件的路径
	fp   *os.File      //文件的打开句柄
	rd   *bufio.Reader //按行读文件
}

type IteratorLineFunc func(string)

func NewFileIterator() *FileIterator {
	return &FileIterator{}
}

// 设置文件
func (fi *FileIterator) SetFile(f string) *FileIterator {
	fi.file = f
	return fi
}

// 初始化操作
func (fi *FileIterator) Init() (*FileIterator, error) {
	if fi.fp != nil {
		fi.fp.Close()
	}
	if !FileExists(fi.file) {
		return nil, fmt.Errorf("file not exist,%s", fi.file)
	}
	fp, err := os.Open(fi.file)
	if err != nil {
		return nil, err
	}
	fi.rd = bufio.NewReader(fp)
	fi.fp = fp
	return fi, nil
}

// 返回一个可遍历的通道
func (fi *FileIterator) IterLine(fn IteratorLineFunc) {
	if fi.fp == nil {
		fi.Init()
	}
	var buf bytes.Buffer
	for {
		//开始按行读取文件
		line, isPrefix, err := fi.rd.ReadLine()
		if len(line) > 0 {
			buf.Write(line)
			if !isPrefix {
				fn(buf.String())
				buf.Reset()
			}
		}
		//读取文件结束
		if err == io.EOF {
			break
		}
		//出错了，退出
		if err != nil && err != io.EOF {
			fmt.Printf("read file occur error\t%v\n", err)
			break
		}
	}
	fi.fp.Close()
}

func NewFileTool() *FileTool {
	ret := &FileTool{}
	return ret
}

type FileTool struct {
	flist []string //读取目录下的文件列表结果
}

// 初始化操作
func (ft *FileTool) Init() {
	ft.flist = make([]string, 0)
}

// 判断文件或目录是否存在
func (ft *FileTool) IsExists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 读取目录下的所有文件
func (ft *FileTool) ReadDirFiles(dir string) ([]string, error) {
	if !ft.IsExists(dir) {
		errmsg := fmt.Sprintf("dir %v not exists", dir)
		return ft.flist, fmt.Errorf(errmsg)
	}
	dir = strings.TrimRight(dir, "/")
	//开始读取目录
	dp, err := os.Open(dir)
	if err != nil {
		return ft.flist, err
	}
	defer dp.Close()
	dlist, err := dp.Readdir(-1)
	if err != nil {
		errmsg := fmt.Sprintf("read dir error %s %v", dir, err)
		return ft.flist, fmt.Errorf(errmsg)
	}
	//开始递归读取目录，忽略掉隐藏目录及文件
	for _, v := range dlist {
		if strings.HasPrefix(v.Name(), ".") {
			continue
		}
		f := dir + "/" + v.Name()
		if v.IsDir() {
			ft.ReadDirFiles(f)
		} else {
			ft.flist = append(ft.flist, f)
		}
	}
	ft.flist = UniqueStrSlice(ft.flist)
	return ft.flist, nil
}

// 文件是否存在
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 打开文件，得到一个句柄。如果不存在创建一个，如果存在再判断
func OpenNewFile(fileName, bakExt string, isBak bool) (fp *os.File, err error) {
	if isBak && len(bakExt) == 0 {
		bakExt = fmt.Sprintf("%s.bak", time.Now().Local().Format("2006-01-02T15:04:05.000"))
	}
	if FileExists(fileName) {
		if isBak {
			err = os.Rename(fileName, fmt.Sprintf("%s.%s", fileName, bakExt))
		} else {
			err = os.Remove(fileName)
		}
		if err != nil {
			return
		}
	} else {
		//判断所在的目录是否存在
		dirName := filepath.Dir(fileName)
		if !FileExists(dirName) {
			err = os.MkdirAll(dirName, 0755)
			if err != nil {
				return
			}
		}
	}
	fp, err = os.Create(fileName)
	return
}

// 文件最后修改时间
func FileModTime(filename string) (int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	fileinfo, err := fd.Stat()
	if err != nil {
		return 0, err
	}
	return fileinfo.ModTime().Unix(), nil
}

// 按行遍历文件的每一行
func FileLineIterator(f string, fn func(line string) error) error {
	if !FileExists(f) {
		return fmt.Errorf("file %s not exists", f)
	}
	fp, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)
	var buf bytes.Buffer
	total := 0
	for {
		//开始按行读取文件
		line, isPrefix, err := rd.ReadLine()
		buf.Write(line)
		if !isPrefix && err == nil {
			//如果回调函数报错则直接返回
			if ferr := fn(buf.String()); ferr != nil {
				return ferr
			}
			buf.Reset()
			//校验数量
			total++
		}
		//读取文件结束
		if err == io.EOF {
			break
		}
		//出错了，退出
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

// 读取目录下的所有文件
func ReadDirFiles(dir string) (ret []string, err error) {
	if !FileExists(dir) {
		err = fmt.Errorf("dir %v not exists", dir)
		return ret, err
	}
	dir = strings.TrimRight(dir, "/")

	//开始读取目录
	var dp *os.File
	dp, err = os.Open(dir)
	if err != nil {
		return ret, err
	}
	if dp == nil {
		err = fmt.Errorf("open %v failed", dir)
		return ret, err
	}

	defer dp.Close()

	var dlist []os.FileInfo
	dlist, err = dp.Readdir(-1)
	if err != nil {
		return ret, err
	}

	//开始递归读取目录，忽略掉隐藏目录及文件
	for _, v := range dlist {
		if strings.HasPrefix(v.Name(), ".") {
			continue
		}
		f := dir + "/" + v.Name()
		if v.IsDir() {
			tret, terr := ReadDirFiles(f)
			if terr == nil && tret != nil && len(tret) > 0 {
				ret = append(ret, tret...)
			}
		} else {
			ret = append(ret, f)
		}
	}
	ret = UniqueStrSlice(ret)
	sort.Strings(ret)
	return ret, nil
}

// 读取文件内容，按行返回
func GetFileContentLines(filePath string) (lines []string, err error) {
	err = FileLineIterator(filePath, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	return
}

// 读取文件内容，按行返回
func GetFileContentStr(filePath string) (str string, err error) {
	var lines []string
	err = FileLineIterator(filePath, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	str = strings.Join(lines, "\n")
	str = strings.TrimSpace(str)
	return
}

// 当前路径
func CurPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// 当前目录
func CurDir() string {
	return filepath.Dir(CurPath())
}
