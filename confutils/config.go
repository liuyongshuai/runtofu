// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        config.go
// @date        2025-03-25 15:48

package confutils

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kr/pretty"
)

// http服务相关
type HTTPConf struct {
	Port     string `toml:"port"`     //http服务监听地址
	TplDir   string `toml:"tpldir"`   //模板的路径
	TplExt   string `toml:"tplext"`   //模板扩展名
	SiteName string `toml:"sitename"` //站点的名称
}

// 管理后台服务相关
type AdminConf struct {
	Port   string `toml:"port"`   //http服务监听地址
	TplDir string `toml:"tpldir"` //模板的路径
	TplExt string `toml:"tplext"` //模板扩展名
}

// 通用的配置项
type CommonConf struct {
	StaticPrefix string `toml:"static_prefix"` //静态资源的地址前缀
	ImagePrefix  string `toml:"image_prefix"`  //图片地址前缀
	CookieKey    string `toml:"cookie_key"`    //登录时所用的cookie的key
}

// aliyun-oss
type AliyunOSSConf struct {
	EndPoint     string `toml:"end_point"`
	AccessId     string `toml:"access_id"`
	AccessSecret string `toml:"access_secret"`
}

// aliyun总体配置
type AliyunConf struct {
	OSS AliyunOSSConf `toml:"oss"` //oss相关配置
}

// oauth weibo
type OauthWeiboConf struct {
	AppKey      string `toml:"app_key"`
	AppSecret   string `toml:"app_secret"`
	ApiUrl      string `toml:"api_url"`
	CallBackUrl string `toml:"callback_url"`
}

type OauthGithubConf struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	ApiUrl       string `toml:"api_url"`
	CallBackUrl  string `toml:"callback_url"`
}

type OauthConf struct {
	Weibo  OauthWeiboConf  `toml:"weibo"`
	Github OauthGithubConf `toml:"github"`
}

// MySQL的相关配置信息
type MySQLConf struct {
	DbName  string `toml:"db"`      //DB名称
	Charset string `toml:"charset"` //字符编码
	Passwd  string `toml:"passwd"`  //密码
	Host    string `toml:"host"`    //地址
	User    string `toml:"user"`    //用户名
	Port    uint16 `toml:"port"`    //端口
}

// Redis相关配置信息
type RedisConf struct {
	Host   string `toml:"host"`   //地址
	Port   uint16 `toml:"port"`   //端口
	Passwd string `toml:"passwd"` //密码
	Db     int    `toml:"db"`     //哪个库
}

type RuntofuConfig struct {
	MySQL  MySQLConf  `toml:"mysql"`  //mysql配置
	Redis  RedisConf  `toml:"redis"`  //mysql配置
	Common CommonConf `toml:"common"` //通用的配置项
	Http   HTTPConf   `toml:"http"`   //http服务
	Admin  AdminConf  `toml:"admin"`
	Aliyun AliyunConf `toml:"aliyun"` //阿里云相关服务的配置
	Oauth  OauthConf  `toml:"oauth"`
}

var gConfig RuntofuConfig

// 提取配置对象
func GetConfiger() *RuntofuConfig {
	return &gConfig
}

// 初始化配置信息
func (p *RuntofuConfig) Init(configPath string) error {
	fmt.Println("start init config.....")
	gConfig = RuntofuConfig{}
	if _, err := toml.DecodeFile(configPath, &gConfig); err != nil {
		fmt.Printf("fail to read config.||err=%v||config=%v", err, configPath)
		return errors.New("fail to init config")
	}
	fmt.Printf("init config succeed. config=%# v\n", pretty.Formatter(gConfig))
	return nil
}
