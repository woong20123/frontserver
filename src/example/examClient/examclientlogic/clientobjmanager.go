package examclientlogic

import "github.com/woong20123/tcpserver"

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Objmanager is client obj data manage
type Objmanager struct {
	channelmgr *ClientChanmgr
	eu         *ExamUser
	chatClient *tcpserver.TCPClientSession
}

// Intialize is
func (objmgr *Objmanager) Intialize() {
	objmgr.channelmgr = new(ClientChanmgr)
	objmgr.channelmgr.Intialize()
	objmgr.eu = NewExamUser()
	objmgr.chatClient = tcpserver.NewTCPClientSession()
}

// ChanManager is
func (objmgr *Objmanager) ChanManager() *ClientChanmgr {
	return objmgr.channelmgr
}

// User is return user info
func (objmgr *Objmanager) User() *ExamUser {
	return objmgr.eu
}

// ChatClient is return user info
func (objmgr *Objmanager) ChatClient() *tcpserver.TCPClientSession {
	return objmgr.chatClient
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
