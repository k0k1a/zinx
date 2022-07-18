package znet

import "github.com/k0k1a/zinx/ziface"

//实现router时，先嵌入baserouter基类，然后根据需要对这个基类进行重写
type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
