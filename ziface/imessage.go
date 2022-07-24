package ziface

// IMessage 将请求的消息封装到Message中，定义抽象接口
type IMessage interface {
	// GetMsgId 获取消息的ID
	GetMsgId() uint32

	// GetMsgLen 获取消息的长度
	GetMsgLen() uint32

	// GetData 获取消息的内容
	GetData() []byte

	// SetMsgId 设置消息id
	SetMsgId(uint32)

	// SetData 设置消息内容
	SetData([]byte)

	// SetDataLen 设置消息长度
	SetDataLen(uint32)
}
