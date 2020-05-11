package tcpserver

import "errors"

// TCPClientSessionMgr is
type TCPClientSessionMgr struct {
	TCPClientConatiner map[int]TCPClientSession
}

// Intialize is
func (mgr *TCPClientSessionMgr) Intialize() {
	mgr.TCPClientConatiner = make(map[int]TCPClientSession)
}

// AddTCPClientSession is
func (mgr *TCPClientSessionMgr) AddTCPClientSession(tcpclient *TCPClientSession) (err error) {
	index := tcpclient.Index()
	existSession, _ := mgr.TCPClientSession(index)
	if existSession != nil {
		err = errors.New("Exist TCPClientSession")
		return
	}

	mgr.TCPClientConatiner[index] = *tcpclient
	return
}

// TCPClientSession is
func (mgr *TCPClientSessionMgr) TCPClientSession(index int) (tcpclient *TCPClientSession, err error) {
	tcpclient = nil
	val, ok := mgr.TCPClientConatiner[index]
	if ok {
		tcpclient = &val
	} else {
		err = errors.New("Not Find TCPClientConatiner")
	}
	return
}
