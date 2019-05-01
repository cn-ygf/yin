package yin

import (
	"sync"
	"sync/atomic"
)

// 路由管理器
type RouterManager interface {
	Add(Router)        // 添加路由
	Remove(Router)     // 移除路由
	Get(string) Router // 取得路由
	Count() int64      // 获取路由总数
}

// 路由接口
type Router interface {
	Method() string         // 取得提交方式
	Path() string           // 取得path
	String() string         // 取得字符串
	Hash() string           // 取得路由hash值
	Handler() func(Context) // 取得处理程序
}

func NewRouter(method string, path string, handler func(Context)) Router {
	r := &coreRouter{
		method:  method,
		path:    path,
		handler: handler,
	}
	return r
}

type coreRouterManager struct {
	routers sync.Map //保存路由map
	count   int64    // 路由计数
}

func (core *coreRouterManager) Add(router Router) {
	atomic.AddInt64(&core.count, 1)
	core.routers.Store(router.Hash(), router)
}

func (core *coreRouterManager) Remove(router Router) {
	core.routers.Delete(router.Hash())
	atomic.AddInt64(&core.count, -1)
}

func (core *coreRouterManager) Get(hash string) Router {
	if v, ok := core.routers.Load(hash); ok {
		return v.(Router)
	}
	return nil
}

func (core *coreRouterManager) Count() int64 {
	return atomic.LoadInt64(&core.count)
}

type coreRouter struct {
	method  string        // 提交方式
	path    string        // url path
	handler func(Context) // 处理函数
}

func (core *coreRouter) Method() string {
	return core.method
}
func (core *coreRouter) Path() string {
	return core.path
}

func (core *coreRouter) String() string {
	// TODO
	return ""
}

func (core *coreRouter) Hash() string {
	return getRouterHash(core.method, core.path)
}

func (core *coreRouter) Handler() func(Context) {
	return core.handler
}
