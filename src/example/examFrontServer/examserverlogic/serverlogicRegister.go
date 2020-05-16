package examserverlogic

import (
	"example/examshare"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistServerLogic is regist Packet process logic from ChatServerMode
func RegistServerLogic(slm *tcpserver.ServerLogicManager) {
	RegistSystemServerLogic(slm)
}

// RegistSystemServerLogic is regist Packet process logic
func RegistSystemServerLogic(slm *tcpserver.ServerLogicManager) {
	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	slm.RegistLogicfun(int32(examshare.Cmd_CS2FServerRegistRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2F_ServerRegistRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			Logger().Println(err)
			return
		}

		// 서버에 등록 성공
		if res.Result == examshare.ErrCode_ResultSuccess {
			Logger().Println("Registed Chat server!!")
			println("Registed Chat server!!")
		} else {
			Logger().Fatal("Registed Chat server fall")
			println("Registed Chat server fall")
		}
	})

	slm.RegistLogicfun(int32(examshare.Cmd_S2CLoginUserRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2C_LoginUserRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			Logger().Println(err)
			return
		}
		CommonLogicS2CLoginUserRes(&res, conn, p)

		return
	})
}

// CommonLogicS2CLoginUserRes is
func CommonLogicS2CLoginUserRes(res *examshare.CS2C_LoginUserRes, conn *net.TCPConn, p *packet.Packet) {
	if res.Result == examshare.ErrCode_ResultSuccess {
		// user를 셋팅힙니다.
		eu := NewExamUser()
		eu.SetUserSn(res.UserSn)
		eu.SetUserID(&res.UserID)
		eu.SetConn(conn)
		eu.SetState(UserStateEnum.LobbySTATE)
		Instance().ObjMgr().AddUser(res.UserSn, eu)
		// 접속한 유저의 ID 등록
		Instance().ObjMgr().AddUserString(&res.UserID)
		Logger().Println("[", res.UserID, "] 유저가 접속하였습니다.")
	}

	// 로비에 있는 유저들에게 메시지를 보냅니다.
	Instance().ObjMgr().ForEachFunc(func(loop_eu *ExamUser) {
		if loop_eu != nil && loop_eu.State() == UserStateEnum.LobbySTATE {
			// Send 응답 패킷
			if res.UserSn == loop_eu.UserSn() {
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_S2CLoginUserRes))
				err := sendp.MarshalFromProto(res)
				if err == nil {
					tcpserver.Instance().SendManager().SendToClientConn(loop_eu.Conn(), sendp)
				} else {
					Logger().Println(err)
					packet.Pool().ReleasePacket(sendp)
				}
			} else {
				sendp := packet.Pool().AcquirePacket()
				sendp.SetHeaderByDefaultKey(0, int32(examshare.Cmd_CS2CSystemMsgSend))
				sendReq := examshare.CS2C_SystemMsgSend{}
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
}
