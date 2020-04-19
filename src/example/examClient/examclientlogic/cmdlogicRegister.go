package examclientlogic

import (
	"example/share"
	"fmt"
	"net"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// ContructLogicManager is register a logic for the command
func ContructLogicManager(lm *tcpserver.LogicManager) {

	// S2CPacketCommandLoginUserRes에 대한 처리 작업을 등록합니다.
	// 유저의 로그인 패킷 응답 처리 작업 등록
	lm.RegistLogicfun(share.S2CPacketCommandLoginUserRes, func(conn *net.TCPConn, p *packet.Packet) {

		res := share.S2CPCLoginUserRes{}
		p.Read(&res.Result, &res.UserSn, &res.UserID)

		switch res.Result {
		case share.ResultSuccess:
			eu := GetInstance().GetObjMgr().GetUser()
			if eu != nil {
				eu.SetID(res.UserID)
				eu.SetSn(res.UserSn)
				GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "==========================")})
			}
		case share.ResultExistUserID:
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.ConnectedSTATE, []string{fmt.Sprint("==========================", "[ 접 속 화 면 ]", "=========================="), fmt.Sprint("[Login 실패] ", res.UserID, "  유저가 이미 존재합니다.")})
		default:
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.ConnectedSTATE, []string{fmt.Sprint("==========================", "[ 접 속 화 면 ]", "=========================="), "[Login 실패]"})
		}

		return
	})

	// S2CPacketCommandLobbyMsgRes 대한 처리 작업을 등록합니다.
	// 로비에 전달하는 메시지 응답 패킷 처리 작업 등록
	lm.RegistLogicfun(share.S2CPacketCommandLobbyMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCLobbySendMsgRes{}
		p.Read(&res.Result, &res.Userid, &res.Msg)

		var sb strings.Builder
		sb.WriteString(res.Userid)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorDefault)
		return
	})

	lm.RegistLogicfun(share.S2CPacketCommandSystemMsgSend, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCSystemMsgSend{}
		p.Read(&res.Msg)
		GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, res.Msg, termbox.ColorGreen)
		return
	})

	// ChatRoom 관련 패킷 로직 등록 함수
	registChatRoomCommandLogic(lm)

	lm.RunLogicHandle(1)
}

func registChatRoomCommandLogic(lm *tcpserver.LogicManager) {
	// S2CPacketCommandRoomEnterRes Packet Logic
	// 유저의 방입장 패킷 응답 처리 로직
	lm.RegistLogicfun(share.S2CPacketCommandRoomEnterRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCRoomEnterRes{}
		p.Read(&res.Result, &res.RoomIdx, &res.RoomName, &res.EnterUserSn, &res.EnterUserid)
		eu := GetInstance().GetObjMgr().GetUser()
		if eu == nil {
			return
		}

		switch res.Result {
		case share.ResultSuccess:
			// 자기 자신이면 씬전환 작업을 진행합니다. 다른 유저면 입장 메시지 출력합니다.
			if res.EnterUserSn == eu.GetSn() {
				eu.roomIdx = res.RoomIdx
				eu.roomName = res.RoomName
				GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.RoomSTATE, []string{fmt.Sprint("==========================", "[ 채 팅 방 화 면 ]", "==========================")})
			} else {
				var sb strings.Builder
				sb.WriteString("[")
				sb.WriteString(res.EnterUserid)
				sb.WriteString("] 유저가 방에 입장하였습니다.")
				GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorGreen)
			}
		default:

			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "=========================="), fmt.Sprint("[방생성 실패] ", res.RoomName)})
		}
	})

	// S2CPacketCommandRoomLeaveRes Packet Logic
	// 유저의 방 퇴장 패킷 응답 처리 작업
	lm.RegistLogicfun(share.S2CPacketCommandRoomLeaveRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCRoomLeaveRes{}
		p.Read(&res.Result, &res.LeaveUserSn, &res.LeaveUserid)

		eu := GetInstance().GetObjMgr().GetUser()
		if eu == nil {
			return
		}

		switch res.Result {
		case share.ResultSuccess:
			// 자기 자신이면 씬전환 작업을 진행합니다. 다른 유저면 입장 메시지 출력합니다.
			if res.LeaveUserSn == eu.GetSn() {
				eu.roomIdx = 0
				eu.roomName = ""
				GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "==========================")})
			} else {
				var sb strings.Builder
				sb.WriteString("[")
				sb.WriteString(res.LeaveUserid)
				sb.WriteString("] 유저가 퇴장하였습니다.")
				GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorGreen)
			}
		default:
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(UserStateEnum.RoomSTATE, []string{fmt.Sprint("==========================", "[ 채 팅 방 화 면 ]", "==========================")})
		}
	})

	// S2CPacketCommandRoomMsgRes Packet Logic
	// 유저의 방입장 패킷 응답 처리 로직
	lm.RegistLogicfun(share.S2CPacketCommandRoomMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCRoomSendMsgRes{}
		p.Read(&res.Result, &res.Userid, &res.Msg)

		// 전달 받은 메시지를 출력합니다.
		var sb strings.Builder
		sb.WriteString(res.Userid)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorDefault)
		return
	})
}
