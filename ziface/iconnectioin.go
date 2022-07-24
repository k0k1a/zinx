package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接
	Start()

	// Stop 停止连接
	Stop()

	// GetTCPConnection 获取当前连接的绑定socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接id
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr

	// SendMsg 发送数据
	SendMsg(msgId uint32, data []byte) error

	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})

	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)

	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}
