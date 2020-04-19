package examserverlogic

import (
	"log"

	"github.com/woong20123/tcpserver"
)

// GetLogger is
func GetLogger() *log.Logger {
	return tcpserver.GetInstance().GetLoggerMgr().GetLogger()
}

//ExamSvrSingleton is
type ExamSvrSingleton struct {
	objmanager *SvrObjMgr
}

// GetObjMgr is
func (s *ExamSvrSingleton) GetObjMgr() *SvrObjMgr {
	return s.objmanager
}

var instance *ExamSvrSingleton = nil

// GetInstance is
func GetInstance() *ExamSvrSingleton {
	if instance == nil {
		instance = newExamSvrSingleton()
	}
	return instance
}

// newExamSvrSingleton is
func newExamSvrSingleton() *ExamSvrSingleton {
	ess := new(ExamSvrSingleton)
	ess.objmanager = new(SvrObjMgr)
	ess.objmanager.Initialize()

	return ess
}
