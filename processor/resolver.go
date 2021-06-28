package processor

import (
	"monitor/common"
	"monitor/logger"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	GResolver *Resolver
)

type Resolver struct {
	readChan chan []byte
}

func (_self *Resolver) PushReadEvent(bt []byte) {
	_self.readChan <- bt
}

func (_self *Resolver) handleReadEvent(line []byte) {
	var (
		err       error
		regexRes  []string
		parseTime time.Time
		parseUrl  *url.URL
		reqSli    []string
		record    *common.Record
	)

	re := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	regexRes = re.FindStringSubmatch(string(line))

	loc, _ := time.LoadLocation("Asia/Shanghai")
	if len(regexRes) != 14 {
		GMonitor.handleFailChan <- common.TypeErrNum
		logger.Error.Printf("日志长度 %s 错误：%s", len(regexRes), regexRes)
		return
	}

	if parseTime, err = time.ParseInLocation("02/Jan/2006:15:04:05 +0000", regexRes[4], loc); err != nil {
		GMonitor.handleFailChan <- common.TypeErrNum
		logger.Error.Printf("日志时间解析错误：%s", err)
		return
	}

	byteSend, _ := strconv.Atoi(regexRes[8])

	// GET /foo?query=t HTTP/1.0
	reqSli = strings.Split(regexRes[6], " ")
	if len(reqSli) != 3 {
		GMonitor.handleFailChan <- common.TypeErrNum
		logger.Error.Printf("解析请求 %s 错误", regexRes[6])
		return
	}

	if parseUrl, err = url.Parse(reqSli[1]); err != nil {
		GMonitor.handleFailChan <- common.TypeErrNum
		logger.Error.Printf("解析 url 失败：%s", err)
		return
	}

	upstreamTime, _ := strconv.ParseFloat(regexRes[12], 64)
	requestTime, _ := strconv.ParseFloat(regexRes[13], 64)

	record = &common.Record{
		TimeLocal:    parseTime,
		BytesSend:    byteSend, // 请求大小
		Method:       reqSli[0],
		Path:         parseUrl.Path,
		Scheme:       regexRes[5], // 请求协议
		Status:       regexRes[7],
		UpstreamTime: upstreamTime, // 响应时间
		RequestTime:  requestTime,  // 请求总时间
	}
	GWriter.pushWriteEvent(record)

}

func (_self *Resolver) ResolverLoop() {

	var (
		readLine []byte
	)

	for readLine = range _self.readChan {
		_self.handleReadEvent(readLine)
	}

}

func InitResolver() (err error) {
	GResolver = &Resolver{
		readChan: make(chan []byte, 300),
	}

	for i := 0; i < common.GConfig.Processor.ResolverNum; i++ {
		go GResolver.ResolverLoop()
	}
	return
}
