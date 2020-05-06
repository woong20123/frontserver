package examserverlogic

import (
	"example/examchatserverPacket"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// VerifyUserObj is
func VerifyUserObj(eu *ExamUser, state uint32) bool {
	if eu == nil || eu.State() != state {
		return false
	}
	return true
}

// ChatServerModeRegistCommandLogic is regist Packet process logic from ChatServerMode
func ChatServerModeRegistCommandLogic(lm *tcpserver.LogicManager) {
	chatServerModeRegistUserCommandLogic(lm)
	// ChatRoom 관련 패킷 로직 등록 함수
	chatServerModeRegistChatRoomCommandLogic(lm)
}

// RegistCommandLogic is regist Packet process logic
func chatServerModeRegistUserCommandLogic(lm *tcpserver.LogicManager) {

	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SLoginUserReq), func(conn *net.TCPConn, p *packet.Packet) {
		req := examchatserverPacket.C2SPCLoginUserReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		res := examchatserverPacket.S2CPCLoginUserRes{}

		res.Result = examchatserverPacket.ResultSuccess
		var userSn uint32 = 0

		// user id가 빈문자열입니다.
		if req.UserID == "" {
			res.Result = examchatserverPacket.ResultFail
		}

		// user id가 이미 등록되어 있습니다.
		if res.Result == examchatserverPacket.ResultSuccess && true == Instance().ObjMgr().FindUserString(&req.UserID) {
			res.Result = examchatserverPacket.ResultExistUserID
		}

		if res.Result == examchatserverPacket.ResultSuccess {
			// user info를 셋팅힙니다.
			eu := Instance().ObjMgr().FindUser(conn)
			if eu != nil {
				// SN 키 등록 및 USER ID
				eu.SetUserID(&req.UserID)
				userSn = Instance().ObjMgr().MakeUserSn()
				eu.SetUserSn(userSn)
				eu.SetState(UserStateEnum.LobbySTATE)

				// 접속한 유저의 ID 등록
				Instance().ObjMgr().AddUserString(&req.UserID)

				res.UserID = req.UserID
				res.UserSn = userSn

				Logger().Println("[", req.UserID, "] 유저가 접속하였습니다.")
			}
		}

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		Instance().ObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
			if loop_eu != nil && loop_eu.State() == UserStateEnum.LobbySTATE {
				// Send 응답 패킷
				if res.UserSn == loop_eu.UserSn() {
					sendp := packet.Pool().AcquirePacket()
					sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CLoginUserRes))
					err := sendp.MarshalFromProto(&res)
					if err == nil {
						tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
					} else {
						Logger().Println(err)
						packet.Pool().ReleasePacket(sendp)
					}
				} else {
					sendp := packet.Pool().AcquirePacket()
					sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CSystemMsgSend))
					sendReq := examchatserverPacket.S2CPCSystemMsgSend{}
					sendReq.Msg = "[" + res.UserID + "] 유저가 로비에 접속하였습니다."
					err := sendp.MarshalFromProto(&sendReq)
					if err == nil {
						tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
					} else {
						Logger().Println(err)
						packet.Pool().ReleasePacket(sendp)
					}

				}

			}
		})
		return
	})

	// C2SPacketCommandLobbyMsgReq Packet Logic
	// 로비에 전달하는 메시지 패킷 처리 작업 등록
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SLobbyMsgReq), func(conn *net.TCPConn, p *packet.Packet) {
		req := examchatserverPacket.C2SPCLobbySendMsgReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		eu := Instance().ObjMgr().FindUser(conn)
		if eu == nil {
			return
		}

		res := examchatserverPacket.S2CPCLobbySendMsgRes{}
		res.Result = examchatserverPacket.ResultSuccess
		res.UserID = eu.UserID()
		res.Msg = req.Msg

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		Instance().ObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
			if loop_eu != nil && loop_eu.State() == UserStateEnum.LobbySTATE {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CLobbyMsgRes))
				err := sendp.MarshalFromProto(&res)
				if err == nil {
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
				} else {
					Logger().Println(err)
					packet.Pool().ReleasePacket(sendp)
				}
			}
		})
		return
	})
}

func chatServerModeRegistChatRoomCommandLogic(lm *tcpserver.LogicManager) {

	// C2SPacketCommandRoomEnterReq Packet Logic =======================================================================
	// 유저의 방입장 패킷 처리 작업 등록 - 방이 없으면 생성합니다.
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SRoomEnterReq), func(conn *net.TCPConn, p *packet.Packet) {
		req := examchatserverPacket.C2SPCRoomEnterReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		res := examchatserverPacket.S2CPCRoomEnterRes{}
		res.Result = examchatserverPacket.ResultSuccess

		eu := Instance().ObjMgr().FindUser(conn)
		// 유저의 상태가 정상적인지 확인합니다.
		if false == VerifyUserObj(eu, UserStateEnum.LobbySTATE) {
			res.Result = examchatserverPacket.ResultUserStateErr
		}

		// 방이 존재하는지 확인합니다.
		if res.Result == examchatserverPacket.ResultSuccess {
			room := Instance().ChatRoomMgr().FindRoomByName(req.RoomName)

			if room != nil {
				// 유저를 방에 입장 시킵니다.
				Instance().ChatRoomMgr().EnterRoom(room.idx, eu)
				Logger().Println("[", eu.UserID(), "] 유저가 [", eu.UserRoomIdx(), "] 방에 접속하였습니다.")

				res.RoomIdx = room.idx
				res.RoomName = req.RoomName
				res.EnterUserSn = eu.UserSn()
				res.EnterUserID = eu.UserID()
			} else {
				res.Result = examchatserverPacket.ResultRoomCreateFail
			}
		}

		// 성공시 방에 모든 유저들에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == examchatserverPacket.ResultSuccess {
			Instance().ChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomEnterRes))
				err = sendp.MarshalFromProto(&res)
				if err == nil {
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
				} else {
					packet.Pool().ReleasePacket(sendp)
				}
			})
		} else {
			// Send 응답 패킷
			sendp := packet.Pool().AcquirePacket()
			sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomEnterRes))
			err = sendp.MarshalFromProto(&res)
			if err == nil {
				tcpserver.Instance().SendManager().SendToClientConn(eu.Conn(), sendp)
			} else {
				packet.Pool().ReleasePacket(sendp)
			}
		}
	})

	// C2SPacketCommandRoomCreateReq Packet Logic =======================================================================
	// 유저의 방 생성 패킷 처리 작업
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SRoomCreateReq), func(conn *net.TCPConn, p *packet.Packet) {
		req := examchatserverPacket.C2SPCRoomEnterReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		res := examchatserverPacket.S2CPCRoomCreateRes{}
		res.Result = examchatserverPacket.ResultSuccess

		eu := Instance().ObjMgr().FindUser(conn)
		// 유저의 상태가 정상적인지 확인합니다.
		if false == VerifyUserObj(eu, UserStateEnum.LobbySTATE) {
			res.Result = examchatserverPacket.ResultUserStateErr
		}

		// 방이 존재하는지 확인합니다.
		if res.Result == examchatserverPacket.ResultSuccess {
			onceLoop := true
			for onceLoop {
				onceLoop = false
				room := Instance().ChatRoomMgr().FindRoomByName(req.RoomName)
				// 방이 없으면 방을 생성합니다.
				if nil == room {
					_, room = Instance().ChatRoomMgr().CreateRoom(req.RoomName)
				} else {
					res.Result = examchatserverPacket.ResultRoomAlreadyExist
					break
				}

				if room == nil {
					res.Result = examchatserverPacket.ResultRoomCreateFail
					break
				}

				// 유저를 방에 입장 시킵니다.
				Instance().ChatRoomMgr().EnterRoom(room.idx, eu)

				res.RoomIdx = room.idx
				res.RoomName = req.RoomName
				res.EnterUserSn = eu.UserSn()
				res.EnterUserID = eu.UserID()
				Logger().Println("[", eu.UserID(), "] 유저가 [", eu.UserRoomIdx(), "] 방을 생성하였습니다.")
			}
		}

		// Send 응답 패킷
		sendp := packet.Pool().AcquirePacket()
		sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomCreateRes))
		err = sendp.MarshalFromProto(&res)
		if err == nil {
			tcpserver.Instance().SendManager().SendToClientConn(eu.Conn(), sendp)
		} else {
			packet.Pool().ReleasePacket(sendp)
		}

	})

	// C2SPacketCommandRoomLeaveReq Packet Logic =======================================================================
	// 유저의 방 퇴장 패킷 처리 작업
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SRoomLeaveReq), func(conn *net.TCPConn, p *packet.Packet) {

		res := examchatserverPacket.S2CPCRoomLeaveRes{}
		res.Result = examchatserverPacket.ResultSuccess

		// 유저의 상태가 정상적인지 확인합니다.
		eu := Instance().ObjMgr().FindUser(conn)
		if false == VerifyUserObj(eu, UserStateEnum.RoomSTATE) {
			res.Result = examchatserverPacket.ResultUserStateErr
		}

		// 유저를 방에서 퇴장 시킵니다.
		roomidx := eu.UserRoomIdx()
		if true == Instance().ChatRoomMgr().LeaveRoom(roomidx, eu) {
			res.LeaveUserSn = eu.UserSn()
			res.LeaveUserID = eu.UserID()
		} else {
			res.Result = examchatserverPacket.ResultFail
		}

		// 성공시 방에 모든 유저 & 자기자신에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == examchatserverPacket.ResultSuccess {
			Instance().ChatRoomMgr().ForEachFunc(roomidx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomLeaveRes))
				err := sendp.MarshalFromProto(&res)
				if err == nil {
					Logger().Println(err)
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
				}
			})
		}

		// Send 응답 패킷
		sendp := packet.Pool().AcquirePacket()
		sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomLeaveRes))
		err := sendp.MarshalFromProto(&res)
		if err == nil {
			Logger().Println(err)
			tcpserver.Instance().SendManager().SendToClientConn(eu.Conn(), sendp)
		}
	})

	// C2SPacketCommandRoomMsgReq Packet Logic
	// 유저의 방에서 패킷을 전송합니다
	lm.RegistLogicfun(int32(examchatserverPacket.Cmd_C2SRoomMsgReq), func(conn *net.TCPConn, p *packet.Packet) {
		req := examchatserverPacket.C2SPCRoomSendMsgReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		eu := Instance().ObjMgr().FindUser(conn)

		// 비정상적인 유저라면 리턴합니다.

		if false == VerifyUserObj(eu, UserStateEnum.RoomSTATE) {
			return
		}

		if req.Msg == "" {
			return
		}

		res := examchatserverPacket.S2CPCRoomSendMsgRes{}
		res.Result = examchatserverPacket.ResultSuccess
		res.Userid = eu.UserID()
		res.Msg = req.Msg

		// 방안에 있는 유저들에게 메시지를 보냅니다.
		Instance().ChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {

			// 응답 패킷 전송
			sendp := packet.Pool().AcquirePacket()
			sendp.SetHeader(examchatserverPacket.ExamplePacketSerialkey, 0, int32(examchatserverPacket.Cmd_S2CRoomMsgRes))
			err := sendp.MarshalFromProto(&res)
			if err == nil {
				tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
			} else {
				packet.Pool().ReleasePacket(sendp)
			}

			//Logger().Println("[Send Room Msg] send user ", res.Userid, " recv user ", loop_eu.UserID(), " :  ", req.Msg)
		})

		return
	})
}
