package ziface

/*
	路由抽象接口
	路由的数据都是IRquest
*/
type IRouter interface {
	// PreHandle 在处理conn之前的钩子方法hook
	PreHandle(request IRequest)

	// Handle 在处理conn之前的主方法
	Handle(request IRequest)

	// PostHandle 在处理conn之后的钩子方法hook
	PostHandle(request IRequest)
}
