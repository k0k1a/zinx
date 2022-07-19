package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

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
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

		//创建一个拆包对象
		dp := NewDataPack()

		//读取客户端的msg head的二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err", err)
			break
		}

		//拆包得到msgId和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
		}

		//根据msgDataLen再次读取data，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前连接的request数据
		req := Request{
			conn: c,
			msg:  msg,
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

//提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 msgDataLen|MsgID|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	if _, err := c.GetTCPConnection().Write(binaryMsg); err != nil {
		fmt.Println("write msg id ", msgId, "error:", err)
		return errors.New("conn write error")
	}

	return nil
}
