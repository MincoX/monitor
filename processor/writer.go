package processor

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"monitor/common"
	"monitor/logger"
)

var (
	GWriter *Writer
)

type Writer struct {
	writeChan chan *common.Record
}

func (_self *Writer) pushWriteEvent(record *common.Record) {
	_self.writeChan <- record
}

func (_self *Writer) HandleWriteEvent(rec *common.Record) {

	var (
		err error
		cli client.Client
		bp  client.BatchPoints
		pt  *client.Point
	)

	if cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     common.GConfig.Influx.Addr,
		Username: common.GConfig.Influx.User,
		Password: common.GConfig.Influx.Password,
	}); err != nil {
		logger.Error.Printf("连接 influxdb 失败：%s", err)
		return
	}

	// Create a new point batch
	if bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  common.GConfig.Influx.Database,
		Precision: common.GConfig.Influx.Precision,
	}); err != nil {
		logger.Error.Printf("创建 batchPoints 失败：%s", err)
		return
	}

	// Create a point and add to batch
	// Tags: Path, Method, Scheme, Status
	// Scheme： 请求协议
	tags := map[string]string{"Path": rec.Path, "Method": rec.Method, "Scheme": rec.Scheme, "Status": rec.Status}
	// Fields: UpstreamTime, RequestTime, BytesSent
	fields := map[string]interface{}{
		"UpstreamTime": rec.UpstreamTime,
		"RequestTime":  rec.RequestTime,
		"BytesSent":    rec.BytesSend, // 请求数据大小
	}

	if pt, err = client.NewPoint(common.GConfig.Influx.Measurement, tags, fields, rec.TimeLocal); err != nil {
		logger.Error.Printf("创建 point 失败：%s", err)
		return
	}
	bp.AddPoint(pt)

	// Write the batch
	if err = cli.Write(bp); err != nil {
		logger.Error.Printf("数据插入失败：%s", err)
	}
	logger.Info.Printf("数据插入成功：%s", rec)
}

func (_self *Writer) WriteLoop() {
	var (
		record *common.Record
	)

	for record = range _self.writeChan {
		_self.HandleWriteEvent(record)
	}
}

func InitWriter() (err error) {

	GWriter = &Writer{
		writeChan: make(chan *common.Record, 3000),
	}

	for i := 0; i < common.GConfig.Processor.WriterNum; i++ {
		go GWriter.WriteLoop()
	}
	return
}
