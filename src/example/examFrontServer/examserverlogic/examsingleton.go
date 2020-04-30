package examserverlogic

import (
	"log"

	"github.com/woong20123/tcpserver"
)

// GetLogger is
func Logger() *log.Logger {
	return tcpserver.Instance().LoggerMgr().Logger()
}

//ExamSvrSingleton is
type ExamSvrSingleton struct {
	objmanager  *SvrObjMgr
	chatroomMgr *ChatRoomMgr
	configMgr   *ConfigManager
}

// GetObjMgr is
func (s *ExamSvrSingleton) ObjMgr() *SvrObjMgr {
	return s.objmanager
}

// GetChatRoomMgr is
func (s *ExamSvrSingleton) ChatRoomMgr() *ChatRoomMgr {
	return s.chatroomMgr
}

// ConfigMgr is
func (s *ExamSvrSingleton) ConfigMgr() *ConfigManager {
	return s.configMgr
}

var instance *ExamSvrSingleton = nil

// GetInstance is
func Instance() *ExamSvrSingleton {
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
	ess.configMgr = newConfigMgr()

	return ess
}
