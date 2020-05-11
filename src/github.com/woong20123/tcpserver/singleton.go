package tcpserver

// SingletonObj is
type SingletonObj struct {
	packetSerialkey           uint32
	logicm                    *LogicManager
	sendm                     *SendManager
	clientSessionHanlder      *SessionHandler
	serverProxySessionHanlder *SessionHandler
	loggerm                   *LoggerManager
	tcpclientm                *TCPClientSessionMgr
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

	so.clientSessionHanlder = new(SessionHandler)
	so.clientSessionHanlder.Initialize()

	so.serverProxySessionHanlder = new(SessionHandler)
	so.serverProxySessionHanlder.Initialize()

	so.loggerm = new(LoggerManager)
	so.loggerm.Intialize()

	so.tcpclientm = new(TCPClientSessionMgr)
	so.tcpclientm.Intialize()

	return so
}

// LogicManager is return SingletonObj
func (s *SingletonObj) LogicManager() *LogicManager {
	return s.logicm
}

// SendManager is return SendManager
func (s *SingletonObj) SendManager() *SendManager {
	return s.sendm
}

// ClientSessionHandler is return clientSessionHanlder
func (s *SingletonObj) ClientSessionHandler() *SessionHandler {
	return s.clientSessionHanlder
}

// ServerProxySessionHandler is return serverProxySessionHanlder
func (s *SingletonObj) ServerProxySessionHandler() *SessionHandler {
	return s.serverProxySessionHanlder
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
func (s *SingletonObj) LoggerMgr() *LoggerManager {
	return s.loggerm
}

func (s *SingletonObj) TCPClientMgr() *TCPClientSessionMgr {
	return s.tcpclientm
}
