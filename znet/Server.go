package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/k0k1a/zinx/utils"
	"github.com/k0k1a/zinx/ziface"
)

//iserver 接口的实现
type server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块，用来绑定MsgId和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	//该server的连接管理器
	ConnManager ziface.IConnManager
	//该Server创建连接之后自动调用的Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server销毁连接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

//定义当前客户端连接所绑定的Handle API
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	fmt.Println("[Conn Handle] CallbackTOClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write callback err", err)
		return errors.New("CallBackTOClient error")
	}

	return nil
}

func (s *server) Start() {
	fmt.Printf("[Zinx] Server Name:%s,listener at IP:%s,Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn %d, MaxPackageSize %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//0.开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		//1.获取TCP Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		//2.监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx Serer succeed", s.Name, "succ,listenning...")
		var cid uint32 = 0

		//3.阻塞的等待连接过来，处理客户端业务
		for {
			//如果有客户端连接过来，阻塞回返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//设置最大连接个数的判断，如果超过最大连接，那么关闭此新连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn=", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)

			cid++
			go dealConn.Start()
		}
	}()

}

func (s *server) Stop() {
	//清除所有连接
	fmt.Println("[STOP] Zinx Server name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *server) Serve() {
	s.Start()

	//阻塞状态
	select {}
}

func (s *server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Succ!!")
}

func New(name string) ziface.IServer {

	s := &server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

//注册OnConnStart钩子函数的方法
func (s *server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数的方法
func (s *server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数的方法
func (s *server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->Call OnConnStart()...")
		s.OnConnStart(conn)
	}

}

//调用OnConnStop钩子函数的方法
func (s *server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----->Call OnConnStop()...")
		s.OnConnStart(conn)
	}
}
