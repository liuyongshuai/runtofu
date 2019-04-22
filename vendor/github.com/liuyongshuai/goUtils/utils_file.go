/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-03-27 15:20
 */
package goUtils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**
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

//设置文件
func (fi *FileIterator) SetFile(f string) *FileIterator {
	fi.file = f
	return fi
}

//初始化操作
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

//返回一个可遍历的通道
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

//初始化操作
func (ft *FileTool) Init() {
	ft.flist = make([]string, 0)
}

//判断文件或目录是否存在
func (ft *FileTool) IsExists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//读取目录下的所有文件
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

//文件是否存在
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//打开文件，得到一个句柄。如果不存在创建一个，如果存在再判断
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

//文件最后修改时间
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
