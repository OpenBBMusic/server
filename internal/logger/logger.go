package logger

import "log"

// 自用应用log服务
type SvcLogger struct{}

func NewSvcLogger() *SvcLogger {
	return &SvcLogger{}
}

func (l *SvcLogger) Info(message ...string) {
	log.Println("BiliSvc Info | ", message)
}
func (l *SvcLogger) Warn(message ...string) {
	log.Println("BiliSvc Warn | ", message)
}
func (l *SvcLogger) Error(message ...string) {
	log.Println("BiliSvc Err | ", message)
}
