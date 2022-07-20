package utils

import (
	"encoding/json"
	"io/ioutil"

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
	Version        string //当前zinx版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值
}

//定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

//初始化当前的GlobalObject对象
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "v0.7",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4094,
	}

	//应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
