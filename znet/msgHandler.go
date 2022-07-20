package znet

import (
	"fmt"
	"strconv"

	"github.com/k0k1a/zinx/ziface"
)

//消息处理模块的实现
type MsgHandle struct {

	//存放每个MsgId所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgId
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is NOT FOUND! NEED Register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

	//2 根据msgId 调度对应router业务即可
}

//为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		//id已经注册
		panic("repeat api, msgId =" + strconv.Itoa(int(msgId)))
	}
	//2.添加msg与APi的绑定
	m.Apis[msgId] = router
	fmt.Println("Add api MsgId=", msgId, "succ!")
}
