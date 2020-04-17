package tcpserver

// SingletonObj is
type SingletonObj struct {
	packetSerialkey uint32
	logicm          *LogicManager
	sendm           *SendManager
	sessionm        *SessionMgr
}

var instance *SingletonObj = nil

// GetObjInstance is return SingletonObj
func GetObjInstance() *SingletonObj {
	if instance == nil {
		instance = newSingletonObj()
	}
	return instance
}

func newSingletonObj() *SingletonObj {
	so := new(SingletonObj)

	so.logicm = new(LogicManager)
	so.logicm.Initialize()

	so.sendm = new(SendManager)
	so.sendm.Initialize()

	so.sessionm = new(SessionMgr)
	so.sessionm.Initialize()

	return so
}

// GetLogicManager is return SingletonObj
func (s *SingletonObj) GetLogicManager() *LogicManager {
	return s.logicm
}

// GetSendManager is return SendManager
func (s *SingletonObj) GetSendManager() *SendManager {
	return s.sendm
}

// GetSessionMgr is return GetSessionMgr
func (s *SingletonObj) GetSessionMgr() *SessionMgr {
	return s.sessionm
}

// SetSerialkey is regist server serialkey
func (s *SingletonObj) SetSerialkey(key uint32) {
	s.packetSerialkey = key
}

// GetSerialkey is get server serialkey
func (s *SingletonObj) GetSerialkey() uint32 {
	return s.packetSerialkey
}
