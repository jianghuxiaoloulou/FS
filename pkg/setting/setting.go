package setting

import (
	"WowjoyProject/FileServer/pkg/loggin"
	"os"
	"path/filepath"
	"time"

	"github.com/Unknwon/goconfig"
)

var (
	Cfg *goconfig.ConfigFile
	// 数据库连接串
	DBConn string
	// 数据库最大连接数
	MaxConn int
	// 源地址根路径
	SrcRoot string
	// 目标地址跟路径
	DestRoot string
	// 源地址code
	SrcCode int
	// 目标地址code
	DestCode int
	//设置最大线程数
	MaxThreads int
	//最大任务数
	MaxTasks int
	// 空闲磁盘大小
	DiskSize int
	// 检查磁盘
	CheckDisk string
	// 程序开始执行时间
	StartHour int
	// 程序结束执行时间
	EndHour int
	// 对象桶Id
	OBJECT_BucketId string
	// 传送数据同步机制
	OBJECT_Sync string
	// 查询对象版本
	OBJECT_GET_Version string
	// 对象上传
	OBJECT_POST_Upload string
	// 对象下载
	OBJECT_GET_Download string
	// 对象删除
	OBJECT_DEL_Delete string
	OBJECT_PATH       string
	// 服务端口
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	// 运行模式
	RUN_MODE string
	// Token
	TOKEN_USERNAME string
	TOKEN_PASSWORD string
	TOKEN_URL      string
)

//初始化全局变量
func init() {
	// 获取可执行文件的路径
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	Cfg, err := goconfig.LoadConfigFile(dir + "/conf/config.ini")
	if err != nil {
		loggin.Error("读取配置文件错误", err)
	}

	LoadBase(Cfg)
	LoadMysql(Cfg)
	LoadObject(Cfg)
	LoadToken(Cfg)
}

func LoadBase(Cfg *goconfig.ConfigFile) {
	RUN_MODE, _ = Cfg.GetValue("General", "RUN_MODE")
	SrcRoot, _ = Cfg.GetValue("General", "SrcRoot")
	DestRoot, _ = Cfg.GetValue("General", "DestRoot")
	SrcCode, _ = Cfg.Int("General", "SrcCode")
	DestCode, _ = Cfg.Int("General", "DestCode")
	MaxThreads, _ = Cfg.Int("General", "MaxThreads")
	MaxTasks, _ = Cfg.Int("General", "MaxTasks")
	DiskSize, _ = Cfg.Int("General", "DiskSize")
	CheckDisk, _ = Cfg.GetValue("General", "CheckDisk")
	StartHour, _ = Cfg.Int("General", "StartHour")
	EndHour, _ = Cfg.Int("General", "EndHour")
}

func LoadMysql(Cfg *goconfig.ConfigFile) {
	DBConn, _ = Cfg.GetValue("Mysql", "DBConn")
	MaxConn, _ = Cfg.Int("Mysql", "MaxConn")
}

func LoadObject(Cfg *goconfig.ConfigFile) {
	OBJECT_BucketId, _ = Cfg.GetValue("Object", "OBJECT_BucketId")
	OBJECT_Sync, _ = Cfg.GetValue("Object", "OBJECT_Sync")
	OBJECT_GET_Version, _ = Cfg.GetValue("Object", "OBJECT_GET_Version")
	OBJECT_POST_Upload, _ = Cfg.GetValue("Object", "OBJECT_POST_Upload")
	OBJECT_GET_Download, _ = Cfg.GetValue("Object", "OBJECT_GET_Download")
	OBJECT_DEL_Delete, _ = Cfg.GetValue("Object", "OBJECT_DEL_Delete")
	OBJECT_PATH, _ = Cfg.GetValue("Object", "OBJECT_PATH")

}

func LoadServer(Cfg *goconfig.ConfigFile) {
	HTTPPort, _ = Cfg.Int("Server", "HTTP_PORT")
	ReadTimeout = time.Duration(Cfg.MustInt("Server", "ReadTimeout", 60)) * time.Second
	WriteTimeout = time.Duration(Cfg.MustInt("Server", "WriteTimeout", 60)) * time.Second
}

func LoadToken(Cfg *goconfig.ConfigFile) {
	TOKEN_USERNAME, _ = Cfg.GetValue("Token", "TOKEN_USERNAME")
	TOKEN_PASSWORD, _ = Cfg.GetValue("Token", "TOKEN_PASSWORD")
	TOKEN_URL, _ = Cfg.GetValue("Token", "TOKEN_URL")
}
