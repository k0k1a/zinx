package znet

import (
	"fmt"
	"strconv"

	"github.com/k0k1a/zinx/utils"
	"github.com/k0k1a/zinx/ziface"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {

	//存放每个MsgId所对应的处理方法 MsgID->router
	Apis map[uint32]ziface.IRouter

	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest

	//业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// AddRouter 为消息添加具体的处理逻辑
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

// StartWorkerPool 启动一个worker工作池（开启工作池动作只能发生一次，一个zinx框架只能有一个worker工作池）
func (m *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启worker，每个worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker启动
		//1 当前的worker对应channel消息队列 开辟空间 第0个worker 对应第0个channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// StartOneWorker 启动一个worker工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, "is Started...")

	//不断阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 根据msgId找到对应的handler
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is NOT FOUND! NEED Register")
		return
	}

	//2 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker进行处理
func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 平均分配给不同的worker
	//根据客户端建立的ConnID进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(),
		"request MsgID=", request.GetMsgID(),
		"to workerID=", workerID)

	//2 将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
