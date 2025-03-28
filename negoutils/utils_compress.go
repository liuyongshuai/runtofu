// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2019-03-11 17:53

package negoutils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// gzip压缩
func GzipEncode(in []byte) (ret []byte, err error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	defer writer.Close()
	_, err = writer.Write(in)
	if err != nil {
		return
	}
	return buffer.Bytes(), nil
}

// gzip解压缩
func GzipDecode(data []byte) (ret []byte, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
