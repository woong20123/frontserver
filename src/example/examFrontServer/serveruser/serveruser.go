package serveruser

import "net"

// UStatelogicFunc is
type UStatelogicFunc func()

type statelist struct {
	NoneSTATE      uint32
	ConnectedSTATE uint32
	LoginSTATE     uint32
	RoomEnterSTATE uint32
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
	sn           uint32
	conn         *net.TCPConn
	id           string
	state        uint32
	roomIdx      int
	onSteteLogic map[int]UStatelogicFunc
}

// NewExamUser is make ExamUser
func NewExamUser() *ExamUser {
	eu := ExamUser{}
	eu.state = UserStateEnum.NoneSTATE
	eu.roomIdx = -1
	eu.onSteteLogic = make(map[int]UStatelogicFunc)
	eu.conn = nil
	return &eu
}

// SetConn is set connection obj
func (eu *ExamUser) SetConn(conn *net.TCPConn) {
	eu.conn = conn
}

// GetConn is set connection obj
func (eu *ExamUser) GetConn() *net.TCPConn {
	return eu.conn
}

// SetState is set user's state
func (eu *ExamUser) SetState(state uint32) {
	eu.state = state
}

// GetState is return user's state
func (eu *ExamUser) GetState() uint32 {
	return eu.state
}

// RegistOnStateLogic is return user's state
func (eu *ExamUser) RegistOnStateLogic(state int, logicfunc UStatelogicFunc) {
	eu.onSteteLogic[state] = logicfunc
}
