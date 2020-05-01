package tcpserver

import "errors"

// TCPClientMgr is
type TCPClientMgr struct {
	TCPClientConatiner map[uint32]TCPClient
}

// Intialize is
func (mgr *TCPClientMgr) Intialize() {
	mgr.TCPClientConatiner = make(map[uint32]TCPClient)
}

// AddTCPClient is
func (mgr *TCPClientMgr) AddTCPClient(index uint32, ip string, port int) (err error) {
	Client := NewTCPClient()
	err = Client.Connect(ip, port)
	if err == nil {
		mgr.TCPClientConatiner[index] = *Client
	}
	return
}

// TCPClient is
func (mgr *TCPClientMgr) TCPClient(index uint32) (tcpclient *TCPClient, err error) {
	val, ok := mgr.TCPClientConatiner[index]
	if ok {
		tcpclient = &val
	} else {
		err = errors.New("Not Find TCPClientConatiner")
	}
	return
}
