package clientuser

import "log"

// UStatelogicFunc is
type UStatelogicFunc func()

// UStateScene is
type UStateScene func(chan int)

type statelist struct {
	NoneSTATE      int
	ConnectedSTATE int
	LoginSTATE     int
	RoomEnterSTATE int
}

// UserStateEnum for public use user state
var UserStateEnum = &statelist{
	NoneSTATE:      0x10,
	ConnectedSTATE: 0x11,
	LoginSTATE:     0x12,
	RoomEnterSTATE: 0x13,
}

// ExamUser = User object connected to the server
type ExamUser struct {
	id             string
	state          int
	roomIdx        int
	onSteteLogic   map[int]UStatelogicFunc
	SteteScene     map[int]UStateScene
	closeSceneChan chan int
}

// NewExamUser is make ExamUser
func NewExamUser() *ExamUser {
	eu := ExamUser{}
	eu.state = UserStateEnum.NoneSTATE
	eu.roomIdx = -1
	eu.onSteteLogic = make(map[int]UStatelogicFunc)
	eu.SteteScene = make(map[int]UStateScene)
	eu.closeSceneChan = make(chan int)
	return &eu
}

// CloseScene is
func (eu *ExamUser) CloseScene() {
	eu.closeSceneChan <- 1
}

// SetState is set user's state
func (eu *ExamUser) SetState(state int) {
	eu.state = state
}

// GetState is return user's state
func (eu *ExamUser) GetState() int {
	return eu.state
}

// RegistOnStateLogic is return user's state
func (eu *ExamUser) RegistOnStateLogic(state int, logicfunc UStatelogicFunc) {
	eu.onSteteLogic[state] = logicfunc
}

// RegistScene is
func (eu *ExamUser) RegistScene(state int, sceneFunc UStateScene) {
	eu.SteteScene[state] = sceneFunc
}

// RunScene is
func (eu *ExamUser) RunScene(state int) {
	scene, exist := eu.SteteScene[state]
	if exist {
		go scene(eu.closeSceneChan)
	} else {
		log.Println("[ExamUser] RunScene fail ", state)
	}

}
