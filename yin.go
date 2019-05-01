package yin

import (
	"log"
	"net/http"
)

var defaultEngine Engine  // 默认框架
var defaultAddress string = "0.0.0.0:9000" // 默认框架的监听地址

// 架构引擎接口
type Engine interface {
	GET(string, func(Context))  // 绑定GET请求
	POST(string, func(Context)) // 绑定POST请求
	Run()                       // 启动服务
	Close() error               // 关闭服务
}

// 设置默认地址
func SetDefault(address string){
	defaultAddress = address
}

// 取得默认引擎
func Default() Engine {
	if defaultEngine == nil{
		defaultEngine = New(defaultAddress)
	}
	return defaultEngine
}

// 新建
func New(bind string) Engine {
	core := &coreEngine{
		RouterManager: &coreRouterManager{},
		localAddress:  bind,
	}
	core.start()
	return core
}

type coreEngine struct {
	RouterManager              // 路由管理器
	localAddress  string       // 绑定地址
	server        *http.Server // http服务
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
	router := NewRouter("GET", path, handler)
	core.RouterManager.Add(router)
}

func (core *coreEngine) POST(path string, handler func(Context)) {
	router := NewRouter("POST", path, handler)
	core.RouterManager.Add(router)
}

func (core *coreEngine) Run() {
	log.Printf("http server Running on http://%s\n",core.localAddress)
	go core.server.ListenAndServe()
}

func (core *coreEngine) Close() error {
	return core.server.Close()
}

func (core *coreEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routerHash := getRouterHash(r.Method, r.URL.Path)
	router := core.RouterManager.Get(routerHash)
	if router == nil {
		// TODO 返回not find
		return
	}
	log.Printf("%s %s",router.Method(),router.Path())
	context := NewContext(r, w)
	router.Handler()(context)
}
