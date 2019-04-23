package yin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
	Method() string // 取得提交方式
	Path() string   // 取得path
	String() string // 取得字符串
	Hash() string   // 取得路由hash值
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
	method string // 提交方式
	path   string // url path
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
	hashStr := fmt.Sprintf("%s-%s-yin", core.method, core.path)
	h := md5.New()
	h.Write([]byte(hashStr))
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
