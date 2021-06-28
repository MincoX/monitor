package processor

import (
	"monitor/common"
	"monitor/logger"
	"net"
	"net/http"
	"time"
)

var (
	GMonitor *Monitor
)

// Monitor 监控数据信息结构体
type Monitor struct {
	handleLineChan chan int // 处理行数
	handleFailChan chan int // 解析失败行数
	startTime      time.Time
	data           common.SystemInfo
	tpsSli         []int
}

func (_self *Monitor) monitorLoop() {

	var (
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)

	scheduleAfter = 5 * time.Second
	// 从管道中读取错误信息，根据错误类型向监控消息体中更新数据
	scheduleTimer = time.NewTimer(scheduleAfter)
	for {
		select {
		case <-GMonitor.handleLineChan:
			_self.data.HandleLine += 1
		case <-GMonitor.handleFailChan:
			_self.data.HandleFail += 1
		case <-scheduleTimer.C:
			_self.tpsSli = append(_self.tpsSli, _self.data.HandleLine)
			if len(_self.tpsSli) > 2 {
				_self.tpsSli = _self.tpsSli[1:]
			}
			scheduleTimer.Reset(scheduleAfter)
		}
	}

}

func handleMonitor(resp http.ResponseWriter, req *http.Request) {

	var (
		err   error
		bytes []byte
	)

	GMonitor.data.RunTime = time.Now().Sub(GMonitor.startTime).String()
	GMonitor.data.ReadChanLen = len(GResolver.readChan)
	GMonitor.data.WriteChanLen = len(GWriter.writeChan)

	// QPS/TPS: 每秒钟 request 或 处理事务 的数量
	if len(GMonitor.tpsSli) >= 2 {
		GMonitor.data.Tps = float64(GMonitor.tpsSli[1]-GMonitor.tpsSli[0]) / 5
	}

	// 正常应答
	if bytes, err = common.BuildResponse(1, "success", GMonitor.data); err == nil {
		logger.Info.Printf("响应信息：%s", string(bytes))
		resp.Write(bytes)
		return
	}

	// 失败应答
	if bytes, err = common.BuildResponse(0, err.Error(), nil); err == nil {
		logger.Info.Printf("响应信息：%s", string(bytes))
		resp.Write(bytes)
		return
	}

}

func InitMonitor() (err error) {

	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/monitor", handleMonitor)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+common.GConfig.Http.Port); err != nil {
		logger.Error.Printf("启动 TCP 监听失败：%s", err)
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(common.GConfig.Http.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(common.GConfig.Http.WriteTimeout) * time.Second,
		Handler:      mux,
	}

	// 赋值单例
	GMonitor = &Monitor{
		handleLineChan: make(chan int, 200),
		handleFailChan: make(chan int, 200),
		startTime:      time.Now(),
	}

	go GMonitor.monitorLoop()

	// 启动 api 服务
	go httpServer.Serve(listener)
	return
}
