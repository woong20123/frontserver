package examclientlogic

import "example/examClient/clientuser"

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Objmanager is client obj data manage
type Objmanager struct {
	channelmgr *ClientChanmgr
	eu         *clientuser.ExamUser
}

// Intialize is
func (objmgr *Objmanager) Intialize() {
	objmgr.channelmgr = new(ClientChanmgr)
	objmgr.channelmgr.Intialize()
	objmgr.eu = clientuser.NewExamUser()
}

// GetChanManager is
func (objmgr *Objmanager) GetChanManager() *ClientChanmgr {
	return objmgr.channelmgr
}

// GetUser is return user info
func (objmgr *Objmanager) GetUser() *clientuser.ExamUser {
	return objmgr.eu
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ClientChanmgr is client channel data manage
type ClientChanmgr struct {
	chanUserState  chan ChanUserStateRequest
	chanSceneClose chan int
	chanGuiRequest chan ChanGuiRequest
}

// ChanUserStateRequest is
type ChanUserStateRequest struct {
	State int
	Msg   string
}

type ChanGuiRequest struct {
	Msg string
}

// Intialize is
func (chanmgr *ClientChanmgr) Intialize() {
	chanmgr.chanUserState = make(chan ChanUserStateRequest)
	chanmgr.chanSceneClose = make(chan int)
	chanmgr.chanGuiRequest = make(chan ChanGuiRequest, 512)
}

// SendChanUserState is
func (chanmgr *ClientChanmgr) SendChanUserState(state int, msg string) {
	chanmgr.chanUserState <- ChanUserStateRequest{state, msg}
}

// GetChanUserState is
func (chanmgr *ClientChanmgr) GetChanUserState() chan ChanUserStateRequest {
	return chanmgr.chanUserState
}

// SendChanSceneClose is
func (chanmgr *ClientChanmgr) SendChanSceneClose() {
	chanmgr.chanSceneClose <- 1
}

// GetchanGuiRequest is
func (chanmgr *ClientChanmgr) GetchanGuiRequest() chan ChanGuiRequest {
	return chanmgr.chanGuiRequest
}

// SendchanGuiRequest is
func (chanmgr *ClientChanmgr) SendchanGuiRequest(msg string) {
	chanmgr.chanGuiRequest <- ChanGuiRequest{msg}
}
