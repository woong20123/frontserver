package examclientlogic

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Objmanager is client obj data manage
type Objmanager struct {
	channelmgr *ClientChanmgr
	eu         *ExamUser
}

// Intialize is
func (objmgr *Objmanager) Intialize() {
	objmgr.channelmgr = new(ClientChanmgr)
	objmgr.channelmgr.Intialize()
	objmgr.eu = NewExamUser()
}

// GetChanManager is
func (objmgr *Objmanager) GetChanManager() *ClientChanmgr {
	return objmgr.channelmgr
}

// GetUser is return user info
func (objmgr *Objmanager) GetUser() *ExamUser {
	return objmgr.eu
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
