package tcpserver

// SingletonObj is
type SingletonObj struct {
	packetSerialkey uint32
	logicm          *LogicManager
	sendm           *SendManager
	sessionm        *SessionMgr
	loggerm         *LoggerManager
}

var instance *SingletonObj = nil

// GetInstance is return SingletonObj
func Instance() *SingletonObj {
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

	so.loggerm = new(LoggerManager)
	so.loggerm.Intialize()

	return so
}

// GetLogicManager is return SingletonObj
func (s *SingletonObj) LogicManager() *LogicManager {
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

// GetLoggerMgr is
func (s *SingletonObj) GetLoggerMgr() *LoggerManager {
	return s.loggerm
}
