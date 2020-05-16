package examserverlogic

import (
	"log"
	"net"
	"sync/atomic"
)

// SvrObjMgr is
type SvrObjMgr struct {
	userContainer map[uint32]*ExamUser
	userIDChecker map[string]bool
	userSnKey     uint32
}

// Initialize is
func (somgr *SvrObjMgr) Initialize() {
	somgr.userContainer = make(map[uint32]*ExamUser)
	somgr.userIDChecker = make(map[string]bool)
	somgr.userSnKey = 0
}

// AddUser is add the user in userContainer.
func (somgr *SvrObjMgr) AddUser(usersn uint32, eu *ExamUser) bool {
	_, exist := somgr.userContainer[usersn]
	if true == exist {
		log.Println("already exist user")
		return false
	}
	somgr.userContainer[usersn] = eu
	return true
}

// DelUser is delete the user in userContainer.
func (somgr *SvrObjMgr) DelUser(usersn uint32) {
	delete(somgr.userContainer, usersn)
}

// DelUserByConn is
func (somgr *SvrObjMgr) DelUserByConn(conn *net.TCPConn) {
	eu := somgr.FindUserByConn(conn)
	if eu != nil {
		delete(somgr.userContainer, eu.UserSn())
	}

}

// FindUser is Find the user in userContainer.
func (somgr *SvrObjMgr) FindUser(usersn uint32) *ExamUser {
	eu, exist := somgr.userContainer[usersn]
	if exist {
		return eu
	}
	return nil
}

// FindUserByConn is
func (somgr *SvrObjMgr) FindUserByConn(conn *net.TCPConn) *ExamUser {

	for _, eu := range somgr.userContainer {
		if eu.Conn() == conn {
			return eu
		}
	}
	return nil
}

// ForEachFunc is Run function to All User
func (somgr *SvrObjMgr) ForEachFunc(f func(eu *ExamUser)) {
	for _, eu := range somgr.userContainer {
		f(eu)
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

// MakeUserSn is return Unique User Sn
func (somgr *SvrObjMgr) MakeUserSn() uint32 {
	return atomic.AddUint32(&somgr.userSnKey, 1)
}
