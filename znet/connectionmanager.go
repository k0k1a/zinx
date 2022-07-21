package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/k0k1a/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

//创建连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加连接
func (c *ConnManager) Add(conn ziface.IConnection) {
	//包含共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn加入connManager中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connId=", conn.GetConnID(), " add to ConnManager successfully:conn num", c.Len())
}

//删除连接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	//包含共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())
	fmt.Println("connId=", conn.GetConnID(), " remove from ConnManager successfully:conn num", c.Len())
}

//根据connId获取连接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//包含共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		//找到了
		return conn, nil
	}
	return nil, errors.New("connection not FOUND")
}

//得到当前连接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

//清除并终止所有的连接
func (c *ConnManager) ClearConn() {
	//包含共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}
	fmt.Println("Clear All connections succ! conn num =", c.Len())
}
