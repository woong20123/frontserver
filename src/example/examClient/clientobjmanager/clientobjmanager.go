package clientobjmanager

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

var instance *Objmanager

// GetInstance is return singleton objmanager
func GetInstance() *Objmanager {
	if instance == nil {
		instance = new(Objmanager)
		instance.Intialize()
	}
	return instance
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ClientChanmgr is client channel data manage
type ClientChanmgr struct {
	chanUserState  chan int
	chanSceneClose chan int
}

// NewChannelmanager is
// func NewChannelmanager() *ClientChanmgr {
// 	c := make(ClientChanmgr)
// 	c.Intialize()
// 	return &c
// }

// Intialize is
func (chanmgr *ClientChanmgr) Intialize() {
	chanmgr.chanUserState = make(chan int)
	chanmgr.chanSceneClose = make(chan int)
}

// GetChanUserState is
func (chanmgr *ClientChanmgr) GetChanUserState() chan int {
	return chanmgr.chanUserState
}

// GetChanSceneClose is
func (chanmgr *ClientChanmgr) GetChanSceneClose() chan int {
	return chanmgr.chanSceneClose
}
