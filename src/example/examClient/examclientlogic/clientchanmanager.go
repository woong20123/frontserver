package examclientlogic

import "github.com/nsf/termbox-go"

// ClientChanmgr is client channel data manage
type ClientChanmgr struct {
	chanUserState      chan ChanUserStateRequest
	chanConnectSrvInfo chan ChanSvrInfoRequest
	chanSceneClose     chan int
	chanRequestFromGui chan ChanRequestFromGui
	chanRequestToGui   chan ChanRequestToGui
}

// SvrInfo is
type ChanSvrInfoRequest struct {
	Ip   string
	Port int
}

// ChanUserStateRequest is
type ChanUserStateRequest struct {
	State int
	Msg   string
}

// ChanRequestFromGui is
type ChanRequestFromGui struct {
	Type int
	Msg  string
}

type toGuiType struct {
	TYPEMsgPrint    int
	TYPEWindowClear int
}

// ToGUIEnum for public use user state
var ToGUIEnum = &toGuiType{
	TYPEMsgPrint:    0x10,
	TYPEWindowClear: 0x11,
}

// ChanRequestToGui is
type ChanRequestToGui struct {
	Type    int
	Msg     string
	wordcol termbox.Attribute
}

// Intialize is
func (chanmgr *ClientChanmgr) Intialize() {
	chanmgr.chanUserState = make(chan ChanUserStateRequest)
	chanmgr.chanConnectSrvInfo = make(chan ChanSvrInfoRequest)
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

// SendChanSrvInfo is
func (chanmgr *ClientChanmgr) SendChanSrvInfo(ip string, port int) {
	chanmgr.chanConnectSrvInfo <- ChanSvrInfoRequest{ip, port}
}

// GetChanSrvInfo is
func (chanmgr *ClientChanmgr) GetChanSrvInfo() chan ChanSvrInfoRequest {
	return chanmgr.chanConnectSrvInfo
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
func (chanmgr *ClientChanmgr) SendchanRequestToGui(t int, msg string, wordcol termbox.Attribute) {
	chanmgr.chanRequestToGui <- ChanRequestToGui{t, msg, wordcol}
}
