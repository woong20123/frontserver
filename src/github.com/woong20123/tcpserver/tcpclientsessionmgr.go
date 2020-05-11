package tcpserver

import "errors"

// TCPClientSessionMgr is
type TCPClientSessionMgr struct {
	TCPClientConatiner map[uint32]TCPClientSession
}

// Intialize is
func (mgr *TCPClientSessionMgr) Intialize() {
	mgr.TCPClientConatiner = make(map[uint32]TCPClientSession)
}

// AddTCPClient is
func (mgr *TCPClientSessionMgr) AddTCPClient(index uint32, ip string, port int) (err error) {
	Client := NewTCPClientSession()
	err = Client.Connect(ip, port)
	if err == nil {
		mgr.TCPClientConatiner[index] = *Client
	}
	return
}

// TCPClient is
func (mgr *TCPClientSessionMgr) TCPClient(index uint32) (tcpclient *TCPClientSession, err error) {
	val, ok := mgr.TCPClientConatiner[index]
	if ok {
		tcpclient = &val
	} else {
		err = errors.New("Not Find TCPClientConatiner")
	}
	return
}
