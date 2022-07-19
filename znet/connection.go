package znet

import (
	"fmt"
	"net"

	"github.com/k0k1a/zinx/utils"
	"github.com/k0k1a/zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经停止的channel
	ExitChan chan bool

	//当前连接处理的方法Router
	Router ziface.IRouter
}

//初始化连接的方法

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is Running ...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		//得到当前连接的request数据
		req := Request{
			conn: c,
			data: buf,
		}

		//执行注册的路由方法
		//从路由中，找到注册绑定的conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID=", c.ConnID)
	//启动当前连接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务
}

//停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop .. ConnID =", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

//获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据
func (c *Connection) Send(data []byte) error {
	return nil
}
