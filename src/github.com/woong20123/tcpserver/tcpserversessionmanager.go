package tcpserver

import "errors"

// TCPServerSessionMgr is
type TCPServerSessionMgr struct {
	TCPServerConatiner map[int]TCPServerSession
}

// Intialize is
func (mgr *TCPServerSessionMgr) Intialize() {
	mgr.TCPServerConatiner = make(map[int]TCPServerSession)
}

// AddTCPServerSession is
func (mgr *TCPServerSessionMgr) AddTCPServerSession(tcpsession *TCPServerSession) (err error) {
	index := tcpsession.Index()
	existSession, _ := mgr.TCPServerSession(index)
	if existSession != nil {
		err = errors.New("Exist TCPServerSession")
		return
	}

	mgr.TCPServerConatiner[index] = *tcpsession
	return
}

// TCPServerSession is
func (mgr *TCPServerSessionMgr) TCPServerSession(index int) (tcpsession *TCPServerSession, err error) {
	tcpsession = nil
	val, ok := mgr.TCPServerConatiner[index]
	if ok {
		tcpsession = &val
	} else {
		err = errors.New("Not Find TCPServerConatiner")
	}
	return
}
