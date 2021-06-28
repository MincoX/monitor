package common

import "time"

type Record struct {
	TimeLocal                    time.Time
	BytesSend                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// SystemInfo 系统状态监控
type SystemInfo struct {
	RunTime      string  `json:"runTime"`      // 运行总时间
	Tps          float64 `json:"tps"`          // 系统吞出量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel 长度
	HandleLine   int     `json:"handleLine"`   // 总处理日志行数
	HandleFail   int     `json:"errNum"`       // 错误数
}

// Response HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}
