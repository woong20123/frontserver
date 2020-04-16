package tcpserver

// SingletonObj is
type SingletonObj struct {
	packetSerialkey uint32
	lm              *LogicManager
	sm              *SendManager
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

	so.lm = new(LogicManager)
	so.lm.Initialize()

	so.sm = new(SendManager)
	so.sm.Initialize()

	return so
}

// GetLogicManager is return SingletonObj
func (s *SingletonObj) GetLogicManager() *LogicManager {
	return s.lm
}

// GetSendManager is return SingletonObj
func (s *SingletonObj) GetSendManager() *SendManager {
	return s.sm
}

// SetSerialkey is regist server serialkey
func (s *SingletonObj) SetSerialkey(key uint32) {
	s.packetSerialkey = key
}

// GetSerialkey is get server serialkey
func (s *SingletonObj) GetSerialkey() uint32 {
	return s.packetSerialkey
}
