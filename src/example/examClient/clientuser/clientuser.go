package clientuser

// UStatelogicFunc is
type UStatelogicFunc func()

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
	id           string
	state        int
	roomIdx      int
	onSteteLogic map[int]UStatelogicFunc
}

// SetState is set user's state
func (u *ExamUser) SetState(state int) {
	u.state = state
}

// GetState is return user's state
func (u *ExamUser) GetState() int {
	return u.state
}

// RegistOnStateLogic is return user's state
func (u *ExamUser) RegistOnStateLogic(state int, logicfunc UStatelogicFunc) {
	u.onSteteLogic[state] = logicfunc
}
