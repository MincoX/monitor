package processor

import (
	"bufio"
	"io"
	"monitor/common"
	"monitor/logger"
	"os"
	"time"
)

var (
	GReader *Reader
)

type Reader struct {
}

func (_self *Reader) ReaderLoop() {

	var (
		err    error
		file   *os.File
		line   []byte
		reader *bufio.Reader
	)

	if file, err = os.Open(common.GConfig.Processor.LogPath); err != nil {
		logger.Error.Printf("打开日志文件失败：%s", err)
	}
	_, _ = file.Seek(0, 2)

	reader = bufio.NewReader(file)

	for {
		if line, err = reader.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				logger.Error.Printf("读取日志文件错误：%s", err)
			}
		}
		GResolver.PushReadEvent(line[:len(line)-1])
		GMonitor.handleLineChan <- common.TypeHandleLine
	}

}

func InitReader() (err error) {
	GReader = &Reader{}
	for i := 0; i < common.GConfig.Processor.ReaderNum; i++ {
		go GReader.ReaderLoop()
	}
	return
}
