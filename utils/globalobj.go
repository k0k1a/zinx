package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/k0k1a/zinx/ziface"
)

//存储一切有关zinx框架的全局参数，供其他模块使用
//一些参数是可以通过zinx.json由用户进行配置
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		zinx
	*/
	Version          string //当前zinx版本号
	MaxConn          int    //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //每个worker对应的消息队列的任务的数量最大值
}

//配置文件位置，默认为当前目录下的conf/zinx.json
const CONFIG_FILE_PATH = "conf/zinx.json"

//定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

//初始化当前的GlobalObject对象
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v1.0",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4094,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {

	if exists, _ := PathExists(CONFIG_FILE_PATH); !exists {
		return
	}

	data, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// PathExists 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
