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
	userIDChecker map[string]bool
	userSnKey     uint32
}

// Initialize is
func (somgr *SvrObjMgr) Initialize() {
	somgr.userContainer = make(map[*net.TCPConn]*serveruser.ExamUser)
	somgr.userIDChecker = make(map[string]bool)
	somgr.userSnKey = 0
}

// AddUser is add the user in userContainer.
func (somgr *SvrObjMgr) AddUser(conn *net.TCPConn, eu *serveruser.ExamUser) bool {
	_, exist := somgr.userContainer[conn]
	if true == exist {
		log.Println("already exist user")
		return false
	}
	somgr.userContainer[conn] = eu
	return true
}

// DelUser is delete the user in userContainer.
func (somgr *SvrObjMgr) DelUser(conn *net.TCPConn) {
	delete(somgr.userContainer, conn)
}

// FindUser is Find the user in userContainer.
func (somgr *SvrObjMgr) FindUser(conn *net.TCPConn) *serveruser.ExamUser {
	eu, exist := somgr.userContainer[conn]
	if exist {
		return eu
	}
	return nil
}

// ForEachFunc is Run function to All User
func (somgr *SvrObjMgr) ForEachFunc(f func(eu *serveruser.ExamUser)) {
	for _, user := range somgr.userContainer {
		f(user)
	}
}

// AddUserString is
func (somgr *SvrObjMgr) AddUserString(id *string) {
	somgr.userIDChecker[*id] = true
}

// DelUserString is
func (somgr *SvrObjMgr) DelUserString(id *string) {
	delete(somgr.userIDChecker, *id)
}

// FindUserString is
func (somgr *SvrObjMgr) FindUserString(id *string) bool {
	_, exist := somgr.userIDChecker[*id]
	return exist
}

// GetUserSn is return Unique User Sn
func (somgr *SvrObjMgr) GetUserSn() uint32 {
	return atomic.AddUint32(&somgr.userSnKey, 1)
}
