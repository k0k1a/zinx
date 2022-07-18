package znet

import (
	"fmt"
	"github.com/k0k1a/zinx/ziface"
	"net"
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
}

func (s *server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s,Port :%d, is starting", s.IP, s.Port)

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

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read buf err", err)
						continue
					}
					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
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

func New(name string) ziface.IServer {

	s := &server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
