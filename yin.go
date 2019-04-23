package yin

import (
	"net/http"
)

// 架构引擎接口
type Engine interface {
	GET(string, func(Context))  // 绑定GET请求
	POST(string, func(Context)) // 绑定POST请求
	Run()                       // 启动服务
	Close() error               // 关闭服务
}

// 取得默认引擎
func Default() Engine {
	return New("0.0.0.0:9000")
}

// 新建
func New(bind string) Engine {
	core := &coreEngine{
		localAddress: bind,
	}
	core.start()
	return core
}

type coreEngine struct {
	localAddress string       // 绑定地址
	server       *http.Server // http服务
}

// http服务初始化操作
func (core *coreEngine) start() {
	mux := http.NewServeMux()
	mux.Handle("/", core)
	core.server = &http.Server{
		Addr:    core.localAddress,
		Handler: mux,
	}
}

func (core *coreEngine) GET(path string, handler func(Context)) {

}

func (core *coreEngine) POST(path string, handler func(Context)) {

}

func (core *coreEngine) Run() {
	go core.server.ListenAndServe()
}

func (core *coreEngine) Close() error {
	return core.server.Close()
}

func (core *coreEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
