// 执行模板的辅助类
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-29 11:59

package negoutils

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	//路径分隔符
	pathSep = string(os.PathSeparator)
)

// 模板构造器
func NewTplBuilder() *TplBuilder {
	return &TplBuilder{
		TplFuncMap: make(template.FuncMap),
		TplNameMap: make(map[string]string),
		TplCache:   make(map[string]*template.Template),
		TplExt:     "tpl",
		isHaveInit: false,
		lock:       new(sync.RWMutex),
	}
}

// 模板结构体
type TplBuilder struct {
	TplRootPathDir string                        //模板的路径目录信息
	TplFuncMap     template.FuncMap              //注册的模板函数
	TplNameMap     map[string]string             //模板名称对应模板的绝对路径
	TplCache       map[string]*template.Template //编译好的名称对应对象的关系
	TplExt         string                        //模板的扩展类型，如"html/tpl..."
	isHaveInit     bool                          //是否已经初始化
	lock           *sync.RWMutex                 //同步用的
}

// 初始化相关系统
func (tb *TplBuilder) initTplBuilder() error {
	if tb.isHaveInit {
		return nil
	}
	tb.lock.Lock()
	defer tb.lock.Unlock()
	if len(tb.TplNameMap) <= 0 {
		err := tb.getAllTplFiles()
		if err != nil {
			return err
		}
	}
	if len(tb.TplNameMap) <= 0 {
		return fmt.Errorf("connt find tpl files")
	}
	for k, fn := range CommonTplFuncs {
		tb.TplFuncMap[k] = fn
	}
	tb.isHaveInit = true
	return nil
}

// 设置根目录
func (tb *TplBuilder) SetRootPathDir(dir string) *TplBuilder {
	tb.isHaveInit = false
	tb.TplRootPathDir = dir
	return tb
}

// 设置扩展
func (tb *TplBuilder) SetTplExt(ext string) *TplBuilder {
	tb.isHaveInit = false
	tb.TplExt = ext
	return tb
}

// 添加用于模板上的函数
func (tb *TplBuilder) AddTplFunc(name string, fn interface{}) *TplBuilder {
	tb.TplFuncMap[name] = fn
	return tb
}

// 添加模板
func (tb *TplBuilder) AddTplFile(tplFile string) error {
	return tb.parseTpl(tplFile)
}

// 批量添加添加用于模板上的函数
func (tb *TplBuilder) AddTplFuncs(fns map[string]interface{}) *TplBuilder {
	for n, f := range fns {
		tb.TplFuncMap[n] = f
	}
	return tb
}

// 执行指定的模板
func (tb *TplBuilder) ExecuteTpl(wr io.Writer, name string, data interface{}) error {
	if !tb.isHaveInit {
		err := tb.initTplBuilder()
		if err != nil {
			fmt.Println("initTplBuilder:", err)
			return err
		}
	}
	//第一步，将名称规范化
	name = tb.formatTplName(name)

	//查缓存里有没有相应的模板信息
	if t, ok := tb.TplCache[name]; ok {
		err := t.ExecuteTemplate(wr, name, data)
		if err != nil {
			fmt.Println("Execute template err:", err)
			return err
		}
		return nil
	}
	tb.lock.Lock()
	defer tb.lock.Unlock()

	//解析出所有的相关模板
	relatedFiles := make(map[string]string)
	err := tb.getTplRelated(name, relatedFiles)
	if err != nil {
		fmt.Println("getTplRelated:", err)
		return err
	}
	t := template.New(name).Funcs(tb.TplFuncMap)
	for tplName, tplFile := range relatedFiles {
		data, err := ioutil.ReadFile(tplFile)
		if err != nil {
			fmt.Println("ReadFile:", err)
			return err
		}
		_, err = t.New(tplName).Parse(string(data))
		if err != nil {
			fmt.Println("Parse:", err, "tplName:", tplName, "tplFile", tplFile)
			return err
		}
	}
	err = t.ExecuteTemplate(wr, name, data)
	if err != nil {
		fmt.Println("ExecuteTemplate:", err, "name:", name, "data:", data)
		return err
	}
	tb.TplCache[name] = t
	return nil
}

// 提取所有的模板文件列表信息
// 就是要遍历指定目录下的所有的文件
func (tb *TplBuilder) getAllTplFiles(tplDir ...string) error {
	rootDir := tb.TplRootPathDir
	if len(tplDir) > 0 {
		rootDir = tplDir[0]
	}
	//要读取的目录
	dir := strings.TrimRight(rootDir, pathSep)
	ftool := NewFileTool()
	//读取所有的文件
	tplFileList, err := ftool.ReadDirFiles(dir)
	if err != nil {
		return err
	}
	//逐个处理模板
	for _, tplFile := range tplFileList {
		//非模板文件
		if !strings.HasSuffix(tplFile, tb.TplExt) {
			continue
		}
		fmt.Println("fint TPL:\t", tplFile)
		err := tb.parseTpl(tplFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// 解析模板
func (tb *TplBuilder) parseTpl(tplFile string) error {
	fmt.Println("parseTpl", tplFile)
	tn := tb.formatTplName(tplFile)
	//将模板的内容完全读出来
	data, err := ioutil.ReadFile(tplFile)
	if err != nil {
		return err
	}
	//先找模板开头的define指令
	reg := regexp.MustCompile(`{{(?:.*?)define(?:\s+)\"(.*)\"(?:.*?)}}`)
	allSub := reg.FindAllStringSubmatch(string(data), 1)
	//找着了就以指定的命名
	if len(allSub) > 0 {
		tn = allSub[0][1]
	}
	if tmpTpl, ok := tb.TplNameMap[tn]; ok {
		return fmt.Errorf("duplicate template name %s for %s and %s", tn, tmpTpl, tplFile)
	}
	fmt.Println("find template [" + tn + "] " + tplFile)
	tb.TplNameMap[tn] = tplFile
	return nil
}

// 提取所有相关的模板系统
func (tb *TplBuilder) getTplRelated(tplName string, ret map[string]string) error {
	tf, ok := tb.TplNameMap[tplName]
	if !ok {
		tplName = tb.formatTplName(tplName)
		tf, ok = tb.TplNameMap[tplName]
		if !ok {
			return fmt.Errorf("invalid tplName %s", tplName)
		}
	}
	ret[tplName] = tf
	data, err := ioutil.ReadFile(tf)
	if err != nil {
		return err
	}
	//解析依赖问题
	reg := regexp.MustCompile(`{{(?:.*?)template(?:\s+)\"(.*)\"(?:.*?)}}`)
	allSub := reg.FindAllStringSubmatch(string(data), -1)
	for _, sub := range allSub {
		if len(sub) < 2 {
			continue
		}
		dn := sub[1]
		if len(dn) > 0 {
			err := tb.getTplRelated(dn, ret)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 规范化模板名称，就是去掉根目录即可
func (tb *TplBuilder) formatTplName(tname string) string {
	if strings.HasPrefix(tname, tb.TplRootPathDir) {
		tname = tname[len(tb.TplRootPathDir):]
		tname = strings.TrimLeft(tname, "/")
	}
	tname = strings.TrimLeft(tname, pathSep)
	fmt.Println("formatTplName=", tname)
	return tname
}
