package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/k0k1a/zinx/utils"
	"github.com/k0k1a/zinx/ziface"
)

// DataPack 封包、拆包的具体模块
type DataPack struct{}

// NewDataPack 初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头的长度
func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32 （4字节）
	return 8
}

// Pack 封包方法
// 消息格式dataLen|msgId|data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	//将datalen写进缓冲
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将msgID写进缓冲
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data写进缓冲
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

// Unpack 拆包方法 (将包的Head信息读出来) 之后再根据head信息的data的长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuf := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和msgID
	msg := &Message{}

	//读DataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgId
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
