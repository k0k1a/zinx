package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/k0k1a/zinx/utils"
	"github.com/k0k1a/zinx/ziface"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer

	//当前连接的socket TCP
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经停止的channel(Reader告知Writer)
	ExitChan chan bool

	//无缓存通道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理msgId和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

//初始化连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
	//将conn加入ConnManager中
	c.TcpServer.GetConnManager().Add(c)

	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is Running ...]")
	defer fmt.Println("[Reader is exit!],connID=", c.ConnID, "remote addr is ", c.RemoteAddr().String())
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的conn对应的router调用
			//根据绑定好的MsgId找到对应处理api业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

//写消息的Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is Running...]")
	defer fmt.Println("[conn Writer exit]", c.RemoteAddr().String())

	//不断的阻塞的等待channel的消息，进行给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Read已经退出，此时Writer也要退出
			return
		}
	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID=", c.ConnID)
	//启动当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建连接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop .. ConnID =", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	//调用开发者注册的 销毁连接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从connManager中删除
	c.TcpServer.GetConnManager().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
	c.msgChan <- binaryMsg

	return nil
}
