package examserverlogic

import (
	"example/share"
	"fmt"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistCommandLogic is regist Packet process logic
func RegistCommandLogic(lm *tcpserver.LogicManager) {

	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {

		req := share.C2SPCLoginUserReq{}
		p.Read(&req.UserID)

		res := share.S2CPCLoginUserRes{}

		res.Result = share.ResultSuccess
		var userSn uint32 = 0

		// user id가 빈문자열입니다.
		if req.UserID == "" {
			res.Result = share.ResultFail
		}

		// user id가 이미 등록되어 있습니다.
		if res.Result == share.ResultSuccess && true == GetInstance().GetObjMgr().FindUserString(&req.UserID) {
			res.Result = share.ResultExistUserID
		}

		if res.Result == share.ResultSuccess {
			// user info를 셋팅힙니다.
			eu := GetInstance().GetObjMgr().FindUser(conn)
			if eu != nil {
				// SN 키 등록 및 USER ID
				eu.SetUserID(&req.UserID)
				userSn = GetInstance().GetObjMgr().MakeUserSn()
				eu.SetUserSn(userSn)
				eu.SetState(UserStateEnum.LobbySTATE)

				// 접속한 유저의 ID 등록
				GetInstance().GetObjMgr().AddUserString(&req.UserID)

				res.UserID = req.UserID
				res.UserSn = userSn

				GetLogger().Println("[", req.UserID, "] 유저가 접속하였습니다.")
			}
		}

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		GetInstance().GetObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
			if loop_eu != nil && loop_eu.GetState() == UserStateEnum.LobbySTATE {
				// Send 응답 패킷
				if res.UserSn == loop_eu.GetUserSn() {
					sendp := packet.GetPool().AcquirePacket()
					sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLoginUserRes)
					sendp.Write(res.Result, res.UserSn, &res.UserID)
					tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
				} else {
					sendp := packet.GetPool().AcquirePacket()
					sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandSystemMsgSend)
					sendp.Write(fmt.Sprint("[", res.UserID, "] 유저가 로비에 접속하였습니다."))
					tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
				}

			}
		})
		return
	})

	// C2SPacketCommandLobbyMsgReq Packet Logic
	// 로비에 전달하는 메시지 패킷 처리 작업 등록
	lm.RegistLogicfun(share.C2SPacketCommandLobbyMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		req := share.C2SPCLobbySendMsgReq{}
		p.Read(&req.Msg)
		eu := GetInstance().GetObjMgr().FindUser(conn)
		if eu == nil {
			return
		}

		res := share.S2CPCLobbySendMsgRes{}
		res.Result = share.ResultSuccess
		res.Userid = eu.GetUserID()
		res.Msg = req.Msg

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		GetInstance().GetObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
			if loop_eu != nil && loop_eu.GetState() == UserStateEnum.LobbySTATE {
				// Send 응답 패킷
				sendp := packet.GetPool().AcquirePacket()
				sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLobbyMsgRes)
				sendp.Write(res.Result, &res.Userid, &res.Msg)
				tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
				GetLogger().Println("[Send Room Msg] send user ", &res.Userid, " recv user ", loop_eu.GetUserID(), " :  ", req.Msg)
			}
		})
		return
	})

	// ChatRoom 관련 패킷 로직 등록 함수
	registChatRoomCommandLogic(lm)
}

func registChatRoomCommandLogic(lm *tcpserver.LogicManager) {
	// C2SPacketCommandRoomEnterReq Packet Logic
	// 유저의 방입장 패킷 처리 작업 등록 - 방이 없으면 생성합니다.
	lm.RegistLogicfun(share.C2SPacketCommandRoomEnterReq, func(conn *net.TCPConn, p *packet.Packet) {
		req := share.C2SPCRoomEnterReq{}
		p.Read(&req.RoomName)

		res := share.S2CPCRoomEnterRes{}
		res.Result = share.ResultSuccess

		// 유저의 상태가 정상적인지 확인합니다.
		eu := GetInstance().GetObjMgr().FindUser(conn)
		if eu == nil && eu.GetState() != UserStateEnum.LobbySTATE {
			res.Result = share.ResultUserStateErr
		}

		// 방이 존재하는지 확인합니다.
		if res.Result == share.ResultSuccess {
			room := GetInstance().GetChatRoomMgr().FindRoomByName(req.RoomName)
			// 방이 없으면 방을 생성합니다.
			if nil == room {
				_, room = GetInstance().GetChatRoomMgr().CreateRoom(req.RoomName)
			}

			if room != nil {
				// 유저를 방에 입장 시킵니다.
				GetInstance().GetChatRoomMgr().EnterRoom(room.idx, eu)
				GetLogger().Println("[", eu.GetUserID(), "] 유저가 [", eu.GetUserRoomIdx(), "] 방에 접속하였습니다.")

				res.RoomIdx = room.idx
				res.RoomName = req.RoomName
				res.EnterUserSn = eu.GetUserSn()
				res.EnterUserid = eu.GetUserID()
			} else {
				res.Result = share.ResultRoomCreateFail
			}
		}

		// 성공시 방에 모든 유저들에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == share.ResultSuccess {
			GetInstance().GetChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.GetPool().AcquirePacket()
				sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandRoomEnterRes)
				sendp.Write(res.Result, res.RoomIdx, &res.RoomName, res.EnterUserSn, &res.EnterUserid)
				tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
			})
		} else {
			// Send 응답 패킷
			sendp := packet.GetPool().AcquirePacket()
			sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandRoomEnterRes)
			sendp.Write(res.Result, res.RoomIdx, &res.RoomName, res.EnterUserSn, &res.EnterUserid)
			tcpserver.Instance().GetSendManager().SendToConn(eu.GetConn(), sendp)
		}
	})

	// C2SPacketCommandRoomLeaveReq Packet Logic
	// 유저의 방 퇴장 패킷 처리 작업
	lm.RegistLogicfun(share.C2SPacketCommandRoomLeaveReq, func(conn *net.TCPConn, p *packet.Packet) {

		res := share.S2CPCRoomLeaveRes{}
		res.Result = share.ResultSuccess

		// 유저의 상태가 정상적인지 확인합니다.
		eu := GetInstance().GetObjMgr().FindUser(conn)
		if eu == nil && eu.GetState() != UserStateEnum.LobbySTATE {
			res.Result = share.ResultUserStateErr
		}

		// 유저를 방에서 퇴장 시킵니다.
		roomidx := eu.GetUserRoomIdx()
		if true == GetInstance().GetChatRoomMgr().LeaveRoom(roomidx, eu) {
			res.LeaveUserSn = eu.GetUserSn()
			res.LeaveUserid = eu.GetUserID()
		} else {
			res.Result = share.ResultFail
		}

		// 성공시 방에 모든 유저 & 자기자신에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == share.ResultSuccess {
			GetInstance().GetChatRoomMgr().ForEachFunc(roomidx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.GetPool().AcquirePacket()
				sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandRoomLeaveRes)
				sendp.Write(res.Result, res.LeaveUserSn, &res.LeaveUserid)
				tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
			})
		}

		// Send 응답 패킷
		sendp := packet.GetPool().AcquirePacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandRoomLeaveRes)
		sendp.Write(res.Result, res.LeaveUserSn, &res.LeaveUserid)
		tcpserver.Instance().GetSendManager().SendToConn(eu.GetConn(), sendp)
	})

	// C2SPacketCommandRoomMsgReq Packet Logic
	// 유저의 방에서 패킷을 전송합니다
	lm.RegistLogicfun(share.C2SPacketCommandRoomMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		req := share.C2SPCRoomSendMsgReq{}
		p.Read(&req.RoomIdx, &req.Msg)
		eu := GetInstance().GetObjMgr().FindUser(conn)

		// 비정상적인 유저라면 리턴합니다.
		if eu == nil || eu.GetState() != UserStateEnum.RoomSTATE {
			return
		}

		if req.Msg == "" {
			return
		}

		res := share.S2CPCRoomSendMsgRes{}
		res.Result = share.ResultSuccess
		res.Userid = eu.GetUserID()

		// 방안에 있는 유저들에게 메시지를 보냅니다.
		GetInstance().GetChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {

			// 응답 패킷 전송
			sendp := packet.GetPool().AcquirePacket()
			sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandRoomMsgRes)
			sendp.Write(res.Result, res.Userid, &req.Msg)
			tcpserver.Instance().GetSendManager().SendToConn(loop_eu.GetConn(), sendp)
			GetLogger().Println("[Send Room Msg] send user ", res.Userid, " recv user ", loop_eu.GetUserID(), " :  ", req.Msg)
		})

		return
	})
}
