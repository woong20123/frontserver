package examserverlogic

import (
	"example/examshare"

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
func ChatServerModeRegistCommandLogic(lm *tcpserver.ClientLogicManager) {
	chatServerModeRegistUserCommandLogic(lm)
	// ChatRoom 관련 패킷 로직 등록 함수
	chatServerModeRegistChatRoomCommandLogic(lm)
}

// RegistCommandLogic is regist Packet process logic
func chatServerModeRegistUserCommandLogic(lm *tcpserver.ClientLogicManager) {

	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	lm.RegistLogicfun(int32(examshare.Cmd_C2SLoginUserReq), func(s tcpserver.Session, p *packet.Packet) {
		req := examshare.C2CS_LoginUserReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		switch tcpcs := s.(type) {
		case *tcpserver.TCPClientSession:
			req.IssuedSessionSn = tcpcs.SessionSn()
		default:
			Logger().Println("Session is not TCPClientSession")
		}

		res := examshare.CS2C_LoginUserRes{}

		res.Result = examshare.ErrCode_ResultSuccess
		res.IssuedSessionSn = req.IssuedSessionSn

		// user id가 빈문자열입니다.
		if req.UserID == "" {
			res.Result = examshare.ErrCode_ResultFail
		}

		// user id가 이미 등록되어 있습니다.
		if res.Result == examshare.ErrCode_ResultSuccess && true == Instance().ObjMgr().FindUserString(&req.UserID) {
			res.Result = examshare.ErrCode_ResultExistUserID
		}

		// res에 유저의 정보 셋팅
		res.UserSn = Instance().ObjMgr().MakeUserSn()
		res.UserID = req.UserID

		// front 모드와 main 모드에서 동일한 로직을 사용하는 것을 묶습니다.
		CommonLogicS2CLoginUserRes(&res, s, p)
		return
	})

	// C2SPacketCommandLobbyMsgReq Packet Logic
	// 로비에 전달하는 메시지 패킷 처리 작업 등록
	lm.RegistLogicfun(int32(examshare.Cmd_C2SLobbyMsgReq), func(s tcpserver.Session, p *packet.Packet) {
		req := examshare.C2CS_LobbySendMsgReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		eu := Instance().ObjMgr().FindUserByConn(s.Conn())
		if eu == nil {
			return
		}

		res := examshare.CS2C_LobbySendMsgRes{}
		res.Result = examshare.ErrCode_ResultSuccess
		res.UserID = eu.UserID()
		res.Msg = req.Msg

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		Instance().ObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
			if loop_eu != nil && loop_eu.State() == UserStateEnum.LobbySTATE {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CLobbyMsgRes))
				err := sendp.MarshalFromProto(&res)
				if err == nil {
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Session(), sendp)
				} else {
					Logger().Println(err)
					packet.Pool().ReleasePacket(sendp)
				}
			}
		})
		return
	})
}

func chatServerModeRegistChatRoomCommandLogic(lm *tcpserver.ClientLogicManager) {

	// C2SPacketCommandRoomEnterReq Packet Logic =======================================================================
	// 유저의 방입장 패킷 처리 작업 등록 - 방이 없으면 생성합니다.
	lm.RegistLogicfun(int32(examshare.Cmd_C2SRoomEnterReq), func(s tcpserver.Session, p *packet.Packet) {
		req := examshare.C2CS_RoomEnterReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		res := examshare.CS2C_RoomEnterRes{}
		res.Result = examshare.ErrCode_ResultSuccess

		eu := Instance().ObjMgr().FindUserByConn(s.Conn())
		// 유저의 상태가 정상적인지 확인합니다.
		if false == VerifyUserObj(eu, UserStateEnum.LobbySTATE) {
			res.Result = examshare.ErrCode_ResultUserStateErr
		}

		// 방이 존재하는지 확인합니다.
		if res.Result == examshare.ErrCode_ResultSuccess {
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
				res.Result = examshare.ErrCode_ResultRoomCreateFail
			}
		}

		// 성공시 방에 모든 유저들에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == examshare.ErrCode_ResultSuccess {
			Instance().ChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomEnterRes))
				err = sendp.MarshalFromProto(&res)
				if err == nil {
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Session(), sendp)
				} else {
					packet.Pool().ReleasePacket(sendp)
				}
			})
		} else {
			// Send 응답 패킷
			sendp := packet.Pool().AcquirePacket()
			sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomEnterRes))
			err = sendp.MarshalFromProto(&res)
			if err == nil {
				tcpserver.Instance().SendManager().SendToClientConn(eu.Session(), sendp)
			} else {
				packet.Pool().ReleasePacket(sendp)
			}
		}
	})

	// C2SPacketCommandRoomCreateReq Packet Logic =======================================================================
	// 유저의 방 생성 패킷 처리 작업
	lm.RegistLogicfun(int32(examshare.Cmd_C2SRoomCreateReq), func(s tcpserver.Session, p *packet.Packet) {
		req := examshare.C2CS_RoomEnterReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		res := examshare.CS2C_RoomCreateRes{}
		res.Result = examshare.ErrCode_ResultSuccess

		eu := Instance().ObjMgr().FindUserByConn(s.Conn())
		// 유저의 상태가 정상적인지 확인합니다.
		if false == VerifyUserObj(eu, UserStateEnum.LobbySTATE) {
			res.Result = examshare.ErrCode_ResultUserStateErr
		}

		// 방이 존재하는지 확인합니다.
		if res.Result == examshare.ErrCode_ResultSuccess {
			onceLoop := true
			for onceLoop {
				onceLoop = false
				room := Instance().ChatRoomMgr().FindRoomByName(req.RoomName)
				// 방이 없으면 방을 생성합니다.
				if nil == room {
					_, room = Instance().ChatRoomMgr().CreateRoom(req.RoomName)
				} else {
					res.Result = examshare.ErrCode_ResultRoomAlreadyExist
					break
				}

				if room == nil {
					res.Result = examshare.ErrCode_ResultRoomCreateFail
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
		sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomCreateRes))
		err = sendp.MarshalFromProto(&res)
		if err == nil {
			tcpserver.Instance().SendManager().SendToClientConn(eu.Session(), sendp)
		} else {
			packet.Pool().ReleasePacket(sendp)
		}

	})

	// C2SPacketCommandRoomLeaveReq Packet Logic =======================================================================
	// 유저의 방 퇴장 패킷 처리 작업
	lm.RegistLogicfun(int32(examshare.Cmd_C2SRoomLeaveReq), func(s tcpserver.Session, p *packet.Packet) {

		res := examshare.CS2C_RoomLeaveRes{}
		res.Result = examshare.ErrCode_ResultSuccess

		// 유저의 상태가 정상적인지 확인합니다.
		eu := Instance().ObjMgr().FindUserByConn(s.Conn())
		if false == VerifyUserObj(eu, UserStateEnum.RoomSTATE) {
			res.Result = examshare.ErrCode_ResultUserStateErr
		}

		// 유저를 방에서 퇴장 시킵니다.
		roomidx := eu.UserRoomIdx()
		if true == Instance().ChatRoomMgr().LeaveRoom(roomidx, eu) {
			res.LeaveUserSn = eu.UserSn()
			res.LeaveUserID = eu.UserID()
		} else {
			res.Result = examshare.ErrCode_ResultFail
		}

		// 성공시 방에 모든 유저 & 자기자신에게 알림, 실패시 자기 자신에게만 알림
		if res.Result == examshare.ErrCode_ResultSuccess {
			Instance().ChatRoomMgr().ForEachFunc(roomidx, func(loop_eu *ExamUser) {
				// Send 응답 패킷
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomLeaveRes))
				err := sendp.MarshalFromProto(&res)
				if err == nil {
					Logger().Println(err)
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Session(), sendp)
				}
			})
		}

		// Send 응답 패킷
		sendp := packet.Pool().AcquirePacket()
		sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomLeaveRes))
		err := sendp.MarshalFromProto(&res)
		if err == nil {
			Logger().Println(err)
			tcpserver.Instance().SendManager().SendToClientConn(eu.Session(), sendp)
		}
	})

	// C2SPacketCommandRoomMsgReq Packet Logic
	// 유저의 방에서 패킷을 전송합니다
	lm.RegistLogicfun(int32(examshare.Cmd_C2SRoomMsgReq), func(s tcpserver.Session, p *packet.Packet) {
		req := examshare.C2CS_RoomSendMsgReq{}
		err := p.UnMarshalFromProto(&req)
		if err != nil {
			Logger().Println(err)
			return
		}

		eu := Instance().ObjMgr().FindUserByConn(s.Conn())

		// 비정상적인 유저라면 리턴합니다.

		if false == VerifyUserObj(eu, UserStateEnum.RoomSTATE) {
			return
		}

		if req.Msg == "" {
			return
		}

		res := examshare.CS2C_RoomSendMsgRes{}
		res.Result = examshare.ErrCode_ResultSuccess
		res.Userid = eu.UserID()
		res.Msg = req.Msg

		// 방안에 있는 유저들에게 메시지를 보냅니다.
		Instance().ChatRoomMgr().ForEachFunc(eu.roomIdx, func(loop_eu *ExamUser) {

			// 응답 패킷 전송
			sendp := packet.Pool().AcquirePacket()
			sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CRoomMsgRes))
			err := sendp.MarshalFromProto(&res)
			if err == nil {
				tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Session(), sendp)
			} else {
				packet.Pool().ReleasePacket(sendp)
			}

			//Logger().Println("[Send Room Msg] send user ", res.Userid, " recv user ", loop_eu.UserID(), " :  ", req.Msg)
		})

		return
	})
}
