package examclientlogic

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Objmanager is client obj data manage
type Objmanager struct {
	channelmgr *ClientChanmgr
}

// Intialize is
func (objmgr *Objmanager) Intialize() {
	objmgr.channelmgr = new(ClientChanmgr)
	objmgr.channelmgr.Intialize()
}

// GetChanManager is
func (objmgr *Objmanager) GetChanManager() *ClientChanmgr {
	return objmgr.channelmgr
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ClientChanmgr is client channel data manage
type ClientChanmgr struct {
	chanUserState  chan ChanUserStateRequest
	chanSceneClose chan int
}

type ChanUserStateRequest struct {
	State            int
	DelayMilliSecond int
}

// Intialize is
func (chanmgr *ClientChanmgr) Intialize() {
	chanmgr.chanUserState = make(chan ChanUserStateRequest)
	chanmgr.chanSceneClose = make(chan int)
}

// SendChanUserState is
func (chanmgr *ClientChanmgr) SendChanUserState(state int, delay int) {
	chanmgr.chanUserState <- ChanUserStateRequest{state, delay}
}

// GetChanUserState is
func (chanmgr *ClientChanmgr) GetChanUserState() chan ChanUserStateRequest {
	return chanmgr.chanUserState
}

// SendChanSceneClose is
func (chanmgr *ClientChanmgr) SendChanSceneClose() {
	chanmgr.chanSceneClose <- 1
}
