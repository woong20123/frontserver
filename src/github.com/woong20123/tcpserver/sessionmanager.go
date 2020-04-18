package tcpserver

import (
	"log"
	"net"
)

// TSessionConnectFunc is
type TSessionConnectFunc func(conn *net.TCPConn)

// TSessionRecvFunc is
type TSessionRecvFunc func(conn *net.TCPConn, buffers []byte, pos uint32) uint32

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

// SessionMgr is
type SessionMgr struct {
	SessionConnectFunCon map[uint32]TSessionConnectFunc
	SessionRecvFun       TSessionRecvFunc
}

// Initialize is intialize to SessionMgr obj
func (s *SessionMgr) Initialize() {
	s.SessionConnectFunCon = make(map[uint32]TSessionConnectFunc)
}

// RegistConnectFunc is
func (s *SessionMgr) RegistConnectFunc(state uint32, f TSessionConnectFunc) {
	s.SessionConnectFunCon[state] = f
}

// RunConnectFunc is
func (s *SessionMgr) RunConnectFunc(state uint32, conn *net.TCPConn) {
	fun, exist := s.SessionConnectFunCon[state]
	if exist {
		fun(conn)
	} else {
		log.Println("RunConnectFunc fail ", state)
	}
}

// RegistRecvFunc is
func (s *SessionMgr) RegistRecvFunc(f TSessionRecvFunc) {
	s.SessionRecvFun = f
}

// RunRecvFunc is
func (s *SessionMgr) RunRecvFunc(conn *net.TCPConn, buffers []byte, pos uint32) uint32 {
	return s.SessionRecvFun(conn, buffers, pos)
}
