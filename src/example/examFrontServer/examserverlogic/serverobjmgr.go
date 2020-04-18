package examserverlogic

import (
	"example/examFrontServer/serveruser"
	"log"
	"net"
	"sync/atomic"
)

// SvrObjMgr is
type SvrObjMgr struct {
	userContainer map[*net.TCPConn]*serveruser.ExamUser
	userSnKey     uint32
}

// Initialize is
func (somgr *SvrObjMgr) Initialize() {
	somgr.userContainer = make(map[*net.TCPConn]*serveruser.ExamUser)
	somgr.userSnKey = 1
}

// AddUser is
func (somgr *SvrObjMgr) AddUser(conn *net.TCPConn, eu *serveruser.ExamUser) bool {
	_, exist := somgr.userContainer[conn]
	if true == exist {
		log.Println("already exist user")
		return false
	}
	somgr.userContainer[conn] = eu
	return true
}

// DelUser is
func (somgr *SvrObjMgr) DelUser(conn *net.TCPConn) {
	delete(somgr.userContainer, conn)
}

// FindUser is
func (somgr *SvrObjMgr) FindUser(conn *net.TCPConn) *serveruser.ExamUser {
	eu, _ := somgr.userContainer[conn]
	return eu
}

// ForEachFunc is
func (somgr *SvrObjMgr) ForEachFunc(f func(eu *serveruser.ExamUser)) {
	for _, user := range somgr.userContainer {
		f(user)
	}
}

// GetUserSn is return Unique User Sn
func (somgr *SvrObjMgr) GetUserSn() uint32 {
	return atomic.AddUint32(&somgr.userSnKey, 1)
}
