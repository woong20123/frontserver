package tcpserver

import (
	"net"
)

// TSessionConnectFunc is
type TSessionConnectFunc func(s Session)

// TSessionRecvFunc is
type TSessionRecvFunc func(s Session, buffers []byte, pos uint32) uint32

type sessionStatelist struct {
	OnConnected uint32
	OnClosed    uint32
	OnRecv      uint32
}

// SessionStateEnum for public use user state
var SessionStateEnum = &sessionStatelist{
	OnConnected: 0x10,
	OnClosed:    0x11,
	OnRecv:      0x12,
}

// SessionHanlder is
type SessionHanlder interface {
	RegistConnectFunc(state uint32, f TSessionConnectFunc)
	RunConnectFunc(state uint32, conn *net.TCPConn)
	RegistRecvFunc(f TSessionRecvFunc)
}

// SessionHandler is
type SessionHandler struct {
	SessionConnectFunCon map[uint32]TSessionConnectFunc
	SessionRecvFun       TSessionRecvFunc
}

// Initialize is intialize to SessionMgr obj
func (sh *SessionHandler) Initialize() {
	sh.SessionConnectFunCon = make(map[uint32]TSessionConnectFunc)
}

// RegistConnectFunc is 접속한 세션을 처리할 로직을 등록합니다.
func (sh *SessionHandler) RegistConnectFunc(state uint32, f TSessionConnectFunc) {
	sh.SessionConnectFunCon[state] = f
}

// RunConnectFunc is
func (sh *SessionHandler) RunConnectFunc(state uint32, s Session) {
	fun, exist := sh.SessionConnectFunCon[state]
	if exist {
		fun(s)
	} else {
		Instance().LoggerMgr().Logger().Println("RunConnectFunc fail ", state)
	}
}

// RegistRecvFunc is
func (sh *SessionHandler) RegistRecvFunc(f TSessionRecvFunc) {
	sh.SessionRecvFun = f
}

// RunRecvFunc is
func (sh *SessionHandler) RunRecvFunc(s Session, buffers []byte, pos uint32) uint32 {
	return sh.SessionRecvFun(s, buffers, pos)
}
