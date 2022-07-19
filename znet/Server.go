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
	//当前server的router，server连接对应的处理业务
	Router ziface.IRouter
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

		//3.阻塞的等待连接过来，处理客户端业务
		for {
			//如果有客户端连接过来，阻塞回返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			var cid uint32 = 0

			dealConn := NewConnection(conn, cid, s.Router)

			cid++
			go dealConn.Start()
		}
	}()

}

func (s *server) Stop() {
	//TODO
}

func (s *server) Serve() {
	s.Start()

	//阻塞状态
	select {}
}

func (s *server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!!")
}

func New(name string) ziface.IServer {

	s := &server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
