package examserverlogic

import (
	"log"

	"github.com/woong20123/tcpserver"
)

// GetLogger is
func GetLogger() *log.Logger {
	return tcpserver.Instance().GetLoggerMgr().GetLogger()
}

//ExamSvrSingleton is
type ExamSvrSingleton struct {
	objmanager  *SvrObjMgr
	chatroomMgr *ChatRoomMgr
}

// GetObjMgr is
func (s *ExamSvrSingleton) GetObjMgr() *SvrObjMgr {
	return s.objmanager
}

// GetChatRoomMgr is
func (s *ExamSvrSingleton) GetChatRoomMgr() *ChatRoomMgr {
	return s.chatroomMgr
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
	ess.chatroomMgr = new(ChatRoomMgr)
	ess.chatroomMgr.Intialize()

	return ess
}
