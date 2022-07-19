package znet

import (
	"io"
	"net"
	"testing"
)

//只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/

	//1.创建socket TCP
	listenner, err := net.Listen("tcp4", "127.0.0.1:7777")
	if err != nil {
		t.Error("server listen err ", err)
		return
	}
	//创建一个go承载 负责从客户端处理业务
	go func() {
		//2.从客户端读取数据
		for {
			conn, err := listenner.Accept()
			if err != nil {
				t.Error("server accept error", err)
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//----->拆包过程
				dp := NewDataPack()
				for {
					//第一次读head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						t.Error("read head err")
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						t.Error("server unpack err", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//Msg是有数据的，需要第二次读取
						//第二次读data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Error("Server unpack data err", err)
							return
						}

						t.Log("----Recv MsgId:", msg.Id, "dataLen=", msg.DataLen, "data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/**
	模拟客户端
	*/

	conn, err := net.Dial("tcp4", "127.0.0.1:7777")
	if err != nil {
		t.Error("client dial err", err)
	}
	//创建一个封包对象dp
	dp := NewDataPack()

	//第一个包
	msg1 := Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	senData1, err := dp.Pack(&msg1)
	if err != nil {
		t.Error("client pack msg1 err ", err)
	}

	//第二个包
	msg2 := Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte("nihao!!"),
	}
	senData2, err := dp.Pack(&msg2)
	if err != nil {
		t.Error("client pack msg2 err ", err)
	}

	//将两个包粘在一起
	senData1 = append(senData1, senData2...)

	//一次性发送给服务器
	conn.Write(senData1)

	//客户端阻塞
	select {}
}
