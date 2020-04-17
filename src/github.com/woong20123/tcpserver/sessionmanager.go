package tcpserver

import (
	"log"
	"net"
)

// TlogicFunc is trans request struct
type TSessionStateFunc func(conn *net.TCPConn)

type sessionStatelist struct {
	OnConnected uint32
	OnClosed    uint32
}

// SessionStateEnum for public use user state
var SessionStateEnum = &sessionStatelist{
	OnConnected: 0x10,
	OnClosed:    0x11,
}

// SessionMgr is
type SessionMgr struct {
	SessionStateFunCon map[uint32]TSessionStateFunc
}

// Initialize is intialize to SessionMgr obj
func (s *SessionMgr) Initialize() {
	s.SessionStateFunCon = make(map[uint32]TSessionStateFunc)
}

// RegistStateFunc is
func (s *SessionMgr) RegistStateFunc(state uint32, f TSessionStateFunc) {
	s.SessionStateFunCon[state] = f
}

// RunStateFunc is
func (s *SessionMgr) RunStateFunc(state uint32, conn *net.TCPConn) {
	fun, exist := s.SessionStateFunCon[state]
	if exist {
		fun(conn)
	} else {
		log.Println("RunStateFunc fail ", state)
	}
}
