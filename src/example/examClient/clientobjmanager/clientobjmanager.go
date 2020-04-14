package clientobjmanager

// objmanager is client obj data manage
type objmanager struct {
	count         int
	chanUserState chan int
}

func (objmgr *objmanager) Intialize() {
	objmgr.chanUserState = make(chan int)
}

func (objmgr *objmanager) GetChanUserState() chan int {
	return objmgr.chanUserState
}

var instance *objmanager

// GetInstance is return singleton objmanager
func GetInstance() *objmanager {
	if instance == nil {
		instance = new(objmanager)
		instance.Intialize()
	}
	return instance
}
