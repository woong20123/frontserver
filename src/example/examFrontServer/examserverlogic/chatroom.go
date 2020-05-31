package examserverlogic

import "sync/atomic"

// ChatRoomMgr is
type ChatRoomMgr struct {
	roomSnKey         uint32
	ChatRoomContainer map[uint32]*ChatRoom
}

// Intialize is
func (rm *ChatRoomMgr) Intialize() {
	rm.roomSnKey = 0
	rm.ChatRoomContainer = make(map[uint32]*ChatRoom)
}

// GetUserSn is return Unique User Sn
func (rm *ChatRoomMgr) makeRoomSn() uint32 {
	return atomic.AddUint32(&rm.roomSnKey, 1)
}

func (rm *ChatRoomMgr) makeRoom() *ChatRoom {
	room := new(ChatRoom)
	room.Intialize(rm.makeRoomSn())
	return room
}

func (rm *ChatRoomMgr) addRoom(room *ChatRoom) {
	rm.ChatRoomContainer[room.idx] = room
}

func (rm *ChatRoomMgr) delRoom(room *ChatRoom) {
	_, exist := rm.ChatRoomContainer[room.idx]
	if exist {
		delete(rm.ChatRoomContainer, room.idx)
	}
}

// FindRoom is find room by room idx
func (rm *ChatRoomMgr) FindRoom(roomidx uint32) *ChatRoom {
	value, _ := rm.ChatRoomContainer[roomidx]
	return value
}

// FindRoomByName is find room by room name
func (rm *ChatRoomMgr) FindRoomByName(roomName string) *ChatRoom {
	for _, value := range rm.ChatRoomContainer {
		if value != nil {
			if value.name == roomName {
				return value
			}
		}
	}
	return nil
}

// CreteRoom is
func (rm *ChatRoomMgr) CreateRoom(roomName string) (result bool, room *ChatRoom) {
	room = nil
	result = false
	findresult := rm.FindRoomByName(roomName)

	// 이미 방이 있습니다.
	if nil != findresult {
		return
	}
	result = true
	room = rm.makeRoom()
	room.name = roomName
	rm.addRoom(room)

	return
}

// EnterRoom is
func (rm *ChatRoomMgr) EnterRoom(roomidx uint32, eu *ExamUser) bool {
	room := rm.FindRoom(roomidx)
	if nil != room {
		if room.EnterUser(eu) {
			eu.roomIdx = room.idx
			eu.SetState(UserStateEnum.RoomSTATE)
			return true
		}
	}
	return false
}

// LeaveRoom is
func (rm *ChatRoomMgr) LeaveRoom(roomidx uint32, eu *ExamUser) bool {
	room := rm.FindRoom(roomidx)
	if nil != room {
		if room.LeaveUser(eu) {
			eu.roomIdx = 0
			eu.SetState(UserStateEnum.LobbySTATE)

			// 유저가 없다면 방을 제거 합니다.
			if false == room.HaveUser() {
				rm.delRoom(room)
			}

			return true
		}
	}
	return false
}

// ForEachFunc is
func (rm *ChatRoomMgr) ForEachFunc(roomidx uint32, f func(eu *ExamUser)) {
	room := rm.FindRoom(roomidx)
	if nil != room {
		room.ForEachFunc(f)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// ChatRoom is
type ChatRoom struct {
	idx           uint32
	name          string
	userContainer map[uint64]*ExamUser
}

// Intialize is
func (r *ChatRoom) Intialize(idx uint32) {
	r.idx = idx
	r.name = ""
	r.userContainer = make(map[uint64]*ExamUser)
}

// EnterUser is
func (r *ChatRoom) EnterUser(eu *ExamUser) bool {
	if eu == nil {
		return false
	}
	_, exist := r.userContainer[eu.UserSn()]
	if exist {
		return false
	}
	r.userContainer[eu.UserSn()] = eu
	return true
}

// LeaveUser is
func (r *ChatRoom) LeaveUser(eu *ExamUser) bool {
	if eu == nil {
		return false
	}
	_, exist := r.userContainer[eu.UserSn()]
	if exist {
		delete(r.userContainer, eu.UserSn())
		return true
	}

	return false
}

// HaveUser is
func (r *ChatRoom) HaveUser() bool {
	if 0 == len(r.userContainer) {
		return false
	}
	return true
}

// ForEachFunc is
func (r *ChatRoom) ForEachFunc(f func(eu *ExamUser)) {
	for _, user := range r.userContainer {
		f(user)
	}
}
