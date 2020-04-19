package examclientlogic

import (
	"log"
	"net"
)

const (
	// NoneSN is
	NoneSN = 0
)

// UStatelogicFunc is
type UStatelogicFunc func()

// UStateScene is
type UStateScene func(chan int)

type statelist struct {
	NoneSTATE      int
	ConnectedSTATE int
	LobbySTATE     int
	RoomSTATE      int
	CloseSTATE     int
}

// UserStateEnum for public use user state
var UserStateEnum = &statelist{
	NoneSTATE:      0x10,
	ConnectedSTATE: 0x11,
	LobbySTATE:     0x12,
	RoomSTATE:      0x13,
	CloseSTATE:     0x14,
}

// ExamUser = User object connected to the server
type ExamUser struct {
	conn           *net.TCPConn
	id             string
	sn             uint32
	state          int
	roomIdx        uint32
	roomName       string
	onSteteLogic   map[int]UStatelogicFunc
	SteteScene     map[int]UStateScene
	closeSceneChan chan int
}

// NewExamUser is make ExamUser
func NewExamUser() *ExamUser {
	eu := ExamUser{}
	eu.state = UserStateEnum.NoneSTATE
	eu.sn = NoneSN
	eu.roomIdx = NoneSN
	eu.onSteteLogic = make(map[int]UStatelogicFunc)
	eu.SteteScene = make(map[int]UStateScene)
	eu.closeSceneChan = make(chan int)
	return &eu
}

// SetConn is
func (eu *ExamUser) SetConn(conn *net.TCPConn) {
	eu.conn = conn
}

// GetConn is
func (eu *ExamUser) GetConn() *net.TCPConn {
	return eu.conn
}

// SetID is set user ID
func (eu *ExamUser) SetID(id string) {
	eu.id = id
}

// GetID is return user ID
func (eu *ExamUser) GetID() string {
	return eu.id
}

// SetSn is set user serial number
func (eu *ExamUser) SetSn(sn uint32) {
	eu.sn = sn
}

// GetSn is return user serial number
func (eu *ExamUser) GetSn() uint32 {
	return eu.sn
}

// SetRoomIdx is set user's room index
func (eu *ExamUser) SetRoomIdx(idx uint32) {
	eu.roomIdx = idx
}

// GetRoomIdx is return user's room index
func (eu *ExamUser) GetRoomIdx() uint32 {
	return eu.roomIdx
}

// SetRoomName is set user ID
func (eu *ExamUser) SetRoomName(name string) {
	eu.roomName = name
}

// GetRoomName is return user ID
func (eu *ExamUser) GetRoomName() string {
	return eu.roomName
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
