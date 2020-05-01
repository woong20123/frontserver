package tcpserver

import (
	"log"
	"os"
	"time"
)

// LoggerManager is
type LoggerManager struct {
	loggerObj   *log.Logger
	logfilename string
}

// Intialize is
func (l *LoggerManager) Intialize() {
	l.loggerObj = nil
	l.logfilename = "examsvr"
}

func (l *LoggerManager) makeLoggerObj() {
	filename := l.logfilename
	filename += "_"
	filename += time.Now().Format("01_02")
	filename += ".log"

	fpLog, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l.loggerObj = log.New(fpLog, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetLogFileName is
func (l *LoggerManager) SetLogFileName(name string) {
	l.logfilename = name
}

// Logger is
func (l *LoggerManager) Logger() *log.Logger {
	if l.loggerObj == nil {
		l.makeLoggerObj()
	}

	return l.loggerObj
}
