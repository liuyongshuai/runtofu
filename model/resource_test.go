/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-02-25 20:23
 */
package model

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"github.com/liuyongshuai/runtofu/confutils"
	"github.com/liuyongshuai/runtofu/negoutils"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "../conf/tofu.conf", "server config.")
	flag.Parse()

	// 解析配置。
	if err := confutils.GetConfiger().Init(configPath); err != nil {
		fmt.Printf("fail to read config.||err=%v||config=%v", err, configPath)
		os.Exit(1)
		return
	}
	conf := confutils.GetConfiger()
	fmt.Printf("%# v\n", pretty.Formatter(conf))

	//初始化model层
	err := Init(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 测试
func TestAliyunOSS(t *testing.T) {
	img := "/root/timg.jpeg"
	data, err := ioutil.ReadFile(img)
	if err != nil {
		panic(err)
	}
	content := string(data)
	md5 := negoutils.MD5(content)
	rder := new(bytes.Buffer)
	rder.Write(data)
	err = AliyunOSSBucket.PutObject(md5+".jpg", rder)
	if err != nil {
		panic(err)
	}
}
