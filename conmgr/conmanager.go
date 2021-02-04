package conmgr

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// HomeDataCache 首页缓存信息
	HomeDataCache = "首页缓存信息"
)

// Conmgr 程序唯一的一个连接管理器
// 处理定时任务
var Conmgr *ConnManager

// ConnManager 连接管理器
type ConnManager struct {
	start        int32
	stop         int32
	quit         chan struct{}
	cacheMap     map[string]interface{}
	cacheMapLock sync.RWMutex
}

// Start 启动连接管理器
func (cm *ConnManager) Start() {
	// 是否已经启动
	if atomic.AddInt32(&cm.start, 1) != 1 {
		return
	}
	log.Println("启动连接管理器")
	// 定时任务
	go cronTaskStart(cm)

}

// Stop 停止连接管理器
func (cm *ConnManager) Stop() {
	if atomic.AddInt32(&cm.stop, 1) != 1 {
		log.Println("连接管理器已经关闭")
		return
	}
	close(cm.quit)
	log.Println("关闭连接管理器")
}

// New 新建一个连接管理器
func New() {
	cm := ConnManager{
		quit:     make(chan struct{}),
		cacheMap: make(map[string]interface{}),
	}
	Conmgr = &cm
	Conmgr.Start()
}

// cronTaskStart 启动定时任务
func cronTaskStart(cm *ConnManager) {
	log.Println("启动定时任务")
	err := ReExeProcessByErrLog()
	if err != nil {
		println(err)
	}
out:
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 4, 0, 0, next.Location())
		// next := now.Add(time.Second * 10)
		t := time.NewTimer(next.Sub(now))
		select {
		// 连接管理器终止时退出
		case <-cm.quit:
			break out
		case <-t.C:
			// 刷新缓存表
			go RefreshCacheMap()
			// 处理分布式错误
			go ReExeProcessByErrLog()
			// 月度考核自动减分程序
			go AutoDeductWithMonthProcess()
		}
	}
}

// RefreshCacheMap 刷新cacheMap中的内容
func RefreshCacheMap() {
	clearCacheMap()
}

// clearCacheMap 清空ClearCacheMap
func clearCacheMap() {
	Conmgr.cacheMapLock.Lock()
	defer Conmgr.cacheMapLock.Unlock()

	//清空 map 的唯一办法就是重新 make 一个新的 map，不用担心垃圾回收的效率，Go语言中的并行垃圾回收效率比写一个清空函数要高效的多。
	Conmgr.cacheMap = map[string]interface{}{}

}
