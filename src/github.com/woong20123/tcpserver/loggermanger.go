package tcpserver

import (
	"log"
	"os"
)

// LoggerManager is
type LoggerManager struct {
	loggerObj *log.Logger
}

// Intialize is
func (l *LoggerManager) Intialize() {
	fpLog, err := os.OpenFile("examsvr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l.loggerObj = log.New(fpLog, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// GetLogger is
func (l *LoggerManager) GetLogger() *log.Logger {
	return l.loggerObj
}
