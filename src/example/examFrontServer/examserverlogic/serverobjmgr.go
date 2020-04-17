package examserverlogic

import (
	"example/examFrontServer/serveruser"
	"log"
	"net"
)

// SvrObjMgr is
type SvrObjMgr struct {
	userContainer map[*net.TCPConn]*serveruser.ExamUser
}

// Initialize is
func (somgr *SvrObjMgr) Initialize() {
	somgr.userContainer = make(map[*net.TCPConn]*serveruser.ExamUser)
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
