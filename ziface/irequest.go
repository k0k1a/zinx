package ziface

// IRequest 请求对象
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgID 得到当前请求的id
	GetMsgID() uint32
}
