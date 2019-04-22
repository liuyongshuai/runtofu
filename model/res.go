/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-02-13 21:37
 */
package model

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/garyburd/redigo/redis"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/configer"
	"github.com/liuyongshuai/runtofu/model/openapi"
)

const (
	//管理后台存session的前缀
	ADMIN_SESSION_PREFIX = "admin_session_%d"

	//博客前台的session前缀：uid
	BLOG_SESSION_PREFIX = "blog_session_%d"
)

//设置的cookie的key默认值
var CookieKey = "rtf"
var BlogCookieKey = "f" + CookieKey

//模型层的实例化
var MAdminUser = NewAdminUserModel()
var MAdminMenu = NewAdminMenuModel()
var MArticle = NewArticleModel()
var MTag = NewTagModel()
var MRuntofuUser = NewRuntofuUserModel()
var MArticleTag = NewArticleTagModel()

var MWeiboApi = new(openapi.WeiboOpenApi)
var MWeiboUser = NewWeiboUserModel()

var MGithubApi = new(openapi.GithubOpenApi)
var MGithubUser = NewGithubUserModel()

var mConf = goUtils.MakeMySQLConf()
var mSystemMenuList AdminMenuList

//mysql连接
var mDB *goUtils.DBase

//redis连接池
var RedisPool *redis.Pool

//aliyun
var AliyunOSSBucket oss.Bucket

//ID生成器
var MSnowFlake, _ = goUtils.NewIDGenerator().
	SetTimeBitSize(48).
	SetWorkerIdBitSize(3).
	SetSequenceBitSize(12).
	SetWorkerId(1).
	Init()

func Init(conf *configer.RuntofuConfig) error {
	if len(conf.Common.CookieKey) > 0 {
		CookieKey = conf.Common.CookieKey
	}

	var err error

	//初始化MySQL
	mConf.DbName = conf.MySQL.DbName
	mConf.Charset = conf.MySQL.Charset
	mConf.Passwd = conf.MySQL.Passwd
	mConf.Host = conf.MySQL.Host
	mConf.User = conf.MySQL.User
	mConf.Port = conf.MySQL.Port
	mConf.AutoCommit = true
	mDB = goUtils.NewDBase(mConf)
	mDB.Conn()
	_, err = mDB.Conn()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	//初始化Redis连接池
	RedisPool = &redis.Pool{
		IdleTimeout: 600,
		MaxIdle:     20,
		MaxActive:   0,
		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port)
			c, err := redis.Dial(
				"tcp", address,
				redis.DialPassword(conf.Redis.Passwd), //密码
				redis.DialDatabase(conf.Redis.Db),     //DB
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, nil
		},
	}

	//aliyun oss sdk
	aliConf := conf.Aliyun
	OSSClient, err := oss.New(
		aliConf.OSS.EndPoint,
		aliConf.OSS.AccessId,
		aliConf.OSS.AccessSecret,
	)
	if err != nil {
		panic(err)
	}
	AliyunOSSBucket = oss.Bucket{
		Client:     *OSSClient,
		BucketName: "runtofu",
	}

	//weibo openapi
	MWeiboApi.InitConf(conf.Oauth.Weibo)
	MGithubApi.InitConf(conf.Oauth.Github)

	//系统菜单自带的，不准修改
	mSystemMenuList.MenuInfo = AdminMenuInfo{
		MenuId:       -1000000,
		MenuName:     "系统管理",
		MenuPath:     "",
		IconName:     "glyphicon-cog",
		IconColor:    "red",
		ParentMenuId: 0,
	}
	mSystemMenuList.SubMenuList = append(mSystemMenuList.SubMenuList, AdminMenuInfo{
		MenuId:       -1000001,
		MenuName:     "菜单管理",
		MenuPath:     "/system/menu",
		IconName:     "glyphicon-align-justify",
		IconColor:    "red",
		ParentMenuId: -1000000,
	})
	mSystemMenuList.MenuInfo.ChildMenuNum = len(mSystemMenuList.SubMenuList)
	return err
}
