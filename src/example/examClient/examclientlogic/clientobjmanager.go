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
	chanUserState      chan ChanUserStateRequest
	chanSceneClose     chan int
	chanRequestFromGui chan ChanRequestFromGui
	chanRequestToGui   chan ChanRequestToGui
}

// ChanUserStateRequest is
type ChanUserStateRequest struct {
	State int
	Msg   string
}

// ChanRequestFromGui
type ChanRequestFromGui struct {
	Type int
	Msg  string
}

type ToGuiType struct {
	TYPEMsgPrint    int
	TYPEWindowClear int
}

// UserStateEnum for public use user state
var ToGUIEnum = &ToGuiType{
	TYPEMsgPrint:    0x10,
	TYPEWindowClear: 0x11,
}

// ChanRequestToGui is
type ChanRequestToGui struct {
	Type int
	Msg  string
}

// Intialize is
func (chanmgr *ClientChanmgr) Intialize() {
	chanmgr.chanUserState = make(chan ChanUserStateRequest)
	chanmgr.chanSceneClose = make(chan int)
	chanmgr.chanRequestFromGui = make(chan ChanRequestFromGui, 512)
	chanmgr.chanRequestToGui = make(chan ChanRequestToGui, 512)
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

// GetchanRequestFromGui is
func (chanmgr *ClientChanmgr) GetchanRequestFromGui() chan ChanRequestFromGui {
	return chanmgr.chanRequestFromGui
}

// SendchanRequestFromGui is
func (chanmgr *ClientChanmgr) SendchanRequestFromGui(msg string) {
	chanmgr.chanRequestFromGui <- ChanRequestFromGui{1, msg}
}

// GetchanRequestToGui is
func (chanmgr *ClientChanmgr) GetchanRequestToGui() chan ChanRequestToGui {
	return chanmgr.chanRequestToGui
}

// SendchanRequestToGui is
func (chanmgr *ClientChanmgr) SendchanRequestToGui(t int, msg string) {
	chanmgr.chanRequestToGui <- ChanRequestToGui{t, msg}
}
