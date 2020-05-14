package examclientlogic

import (
	"example/examshare"
	"fmt"
	"net"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// ContructLogicManager is register a logic for the command
func ContructLogicManager(lm *tcpserver.ClientLogicManager) {

	// S2CPacketCommandLoginUserRes에 대한 처리 작업을 등록합니다.
	// 유저의 로그인 패킷 응답 처리 작업 등록
	lm.RegistLogicfun(int32(examshare.Cmd_S2CLoginUserRes), func(conn *net.TCPConn, p *packet.Packet) {

		res := examshare.CS2C_LoginUserRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch examshare.ErrCode(res.Result) {
		case examshare.ErrCode_ResultSuccess:
			eu := Instance().ObjMgr().User()
			if eu != nil {
				eu.SetID(res.UserID)
				eu.SetSn(res.UserSn)
				Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "==========================")})
			}
		case examshare.ErrCode_ResultExistUserID:
			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.ConnectedSTATE, []string{fmt.Sprint("==========================", "[ 접 속 화 면 ]", "=========================="), fmt.Sprint("[Login 실패] ", res.UserID, "  유저가 이미 존재합니다.")})
		default:
			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.ConnectedSTATE, []string{fmt.Sprint("==========================", "[ 접 속 화 면 ]", "=========================="), "[Login 실패]"})
		}

		return
	})

	// S2CPacketCommandLobbyMsgRes 대한 처리 작업을 등록합니다.
	// 로비에 전달하는 메시지 응답 패킷 처리 작업 등록
	lm.RegistLogicfun(int32(examshare.Cmd_S2CLobbyMsgRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_LobbySendMsgRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			fmt.Println(err)
			return
		}

		var sb strings.Builder
		sb.WriteString(res.UserID)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		Instance().ObjMgr().ChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorDefault)
		return
	})

	lm.RegistLogicfun(int32(examshare.Cmd_CS2CSystemMsgSend), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_SystemMsgSend{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			fmt.Println(err)
			return
		}
		Instance().ObjMgr().ChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, res.Msg, termbox.ColorGreen)
		return
	})

	// ChatRoom 관련 패킷 로직 등록 함수
	registChatRoomCommandLogic(lm)

	lm.RunLogicHandler(1)
}

func registChatRoomCommandLogic(lm *tcpserver.ClientLogicManager) {
	// S2CPacketCommandRoomEnterRes Packet Logic
	// 유저의 방입장 패킷 응답 처리 로직
	lm.RegistLogicfun(int32(examshare.Cmd_S2CRoomEnterRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_RoomEnterRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			return
		}

		eu := Instance().ObjMgr().User()
		if eu == nil {
			return
		}

		switch res.Result {
		case examshare.ErrCode_ResultSuccess:
			// 자기 자신이면 씬전환 작업을 진행합니다. 다른 유저면 입장 메시지 출력합니다.
			if res.EnterUserSn == eu.Sn() {
				eu.roomIdx = res.RoomIdx
				eu.roomName = res.RoomName
				Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.RoomSTATE, []string{fmt.Sprint("==========================", "[ 채 팅 방 화 면 ]", "==========================")})
			} else {
				var sb strings.Builder
				sb.WriteString("[")
				sb.WriteString(res.EnterUserID)
				sb.WriteString("] 유저가 방에 입장하였습니다.")
				Instance().ObjMgr().ChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorGreen)
			}
		default:

			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "=========================="), fmt.Sprint("[방생성 실패] ", res.RoomName)})
		}
	})

	// S2CPacketCommandRoomLeaveRes Packet Logic
	// 유저의 방 퇴장 패킷 응답 처리 작업
	lm.RegistLogicfun(int32(examshare.Cmd_S2CRoomCreateRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_RoomCreateRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			return
		}

		eu := Instance().ObjMgr().User()
		if eu == nil {
			return
		}

		switch res.Result {
		case examshare.ErrCode_ResultSuccess:
			eu.roomIdx = res.RoomIdx
			eu.roomName = res.RoomName
			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.RoomSTATE, []string{fmt.Sprint("==========================", "[ 채 팅 방 화 면 ]", "==========================")})
		default:
			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "=========================="), fmt.Sprint("[방생성 실패] ", res.RoomName)})
		}
	})

	// S2CPacketCommandRoomLeaveRes Packet Logic
	// 유저의 방 퇴장 패킷 응답 처리 작업
	lm.RegistLogicfun(int32(examshare.Cmd_S2CRoomLeaveRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_RoomLeaveRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			return
		}

		eu := Instance().ObjMgr().User()
		if eu == nil {
			return
		}

		switch res.Result {
		case examshare.ErrCode_ResultSuccess:
			// 자기 자신이면 씬전환 작업을 진행합니다. 다른 유저면 입장 메시지 출력합니다.
			if res.LeaveUserSn == eu.Sn() {
				eu.roomIdx = 0
				eu.roomName = ""
				Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.LobbySTATE, []string{fmt.Sprint("==========================", "[ 로 비 화 면 ]", "==========================")})
			} else {
				var sb strings.Builder
				sb.WriteString("[")
				sb.WriteString(res.LeaveUserID)
				sb.WriteString("] 유저가 퇴장하였습니다.")
				Instance().ObjMgr().ChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorGreen)
			}
		default:
			Instance().ObjMgr().ChanManager().SendChanUserState(UserStateEnum.RoomSTATE, []string{fmt.Sprint("==========================", "[ 채 팅 방 화 면 ]", "==========================")})
		}
	})

	// S2CPacketCommandRoomMsgRes Packet Logic
	// 유저의 방입장 패킷 응답 처리 로직
	lm.RegistLogicfun(int32(examshare.Cmd_S2CRoomMsgRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_RoomSendMsgRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			return
		}
		// 전달 받은 메시지를 출력합니다.
		var sb strings.Builder
		sb.WriteString(res.Userid)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		Instance().ObjMgr().ChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorDefault)
		return
	})
}
