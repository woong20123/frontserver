package examserverlogic

import (
	"example/examFrontServer/serveruser"
	"log"
)

// SvrObjMgr is
type SvrObjMgr struct {
	userContainer map[uint32]*serveruser.ExamUser
}

// Initialize is
func (somgr *SvrObjMgr) Initialize() {
	somgr.userContainer = make(map[uint32]*serveruser.ExamUser)
}

// AddUser is
func (somgr *SvrObjMgr) AddUser(sn uint32, eu *serveruser.ExamUser) bool {
	_, exist := somgr.userContainer[sn]
	if true == exist {
		log.Println("already exist user")
		return false
	}
	somgr.userContainer[sn] = eu
	return true
}

// FindUser is
func (somgr *SvrObjMgr) FindUser(sn uint32) *serveruser.ExamUser {
	eu, _ := somgr.userContainer[sn]
	return eu
}
