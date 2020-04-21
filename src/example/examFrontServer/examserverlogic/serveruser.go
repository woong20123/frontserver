package examserverlogic

import "net"

// UStatelogicFunc is
type UStatelogicFunc func()

type statelist struct {
	NoneSTATE      uint32
	ConnectedSTATE uint32
	LobbySTATE     uint32
	RoomSTATE      uint32
}

// UserStateEnum for public use user state
var UserStateEnum = &statelist{
	NoneSTATE:      0x10,
	ConnectedSTATE: 0x11,
	LobbySTATE:     0x12,
	RoomSTATE:      0x13,
}

// ExamUser = User object connected to the server
type ExamUser struct {
	sn           uint32
	conn         *net.TCPConn
	id           string
	state        uint32
	roomIdx      uint32
	onSteteLogic map[int]UStatelogicFunc
}

// NewExamUser is make ExamUser
func NewExamUser() *ExamUser {
	eu := ExamUser{}
	eu.state = UserStateEnum.NoneSTATE
	eu.roomIdx = 0
	eu.onSteteLogic = make(map[int]UStatelogicFunc)
	eu.conn = nil
	eu.sn = 0
	return &eu
}

// SetConn is set connection obj
func (eu *ExamUser) SetConn(conn *net.TCPConn) {
	eu.conn = conn
}

// GetConn is set connection obj
func (eu *ExamUser) Conn() *net.TCPConn {
	return eu.conn
}

// SetState is set user's state
func (eu *ExamUser) SetState(state uint32) {
	eu.state = state
}

// GetState is return user's state
func (eu *ExamUser) State() uint32 {
	return eu.state
}

// SetUserID is set user's id
func (eu *ExamUser) SetUserID(id *string) {
	eu.id = *id
}

// GetUserId is return user's id
func (eu *ExamUser) UserID() string {
	return eu.id
}

// SetUserSn is set user's id
func (eu *ExamUser) SetUserSn(sn uint32) {
	eu.sn = sn
}

// GetUserSn is return user's id
func (eu *ExamUser) UserSn() uint32 {
	return eu.sn
}

// SetUserSn is set user's id
func (eu *ExamUser) SetUserRoomIdx(roomIdx uint32) {
	eu.roomIdx = roomIdx
}

// GetUserSn is return user's id
func (eu *ExamUser) UserRoomIdx() uint32 {
	return eu.roomIdx
}

// RegistOnStateLogic is return user's state
func (eu *ExamUser) RegistOnStateLogic(state int, logicfunc UStatelogicFunc) {
	eu.onSteteLogic[state] = logicfunc
}
