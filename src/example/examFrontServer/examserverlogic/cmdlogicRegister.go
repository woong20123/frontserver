package examserverlogic

import (
	"example/examFrontServer/serveruser"
	"example/share"
	"log"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistCommandLogic is
func RegistCommandLogic(lm *tcpserver.LogicManager) {
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {

		req := share.C2SPCLoginUserReq{}
		p.Read(&req.UserID)

		result := share.ResultSuccess
		var userSn uint32 = 0

		// user id가 이미 등록되어 있습니다.
		if true == GetInstance().GetObjMgr().FindUserString(&req.UserID) {
			result = share.ResultExistUserID
		}

		if result == share.ResultSuccess {
			// user info를 셋팅힙니다.
			objmgr := GetInstance().GetObjMgr()
			eu := objmgr.FindUser(conn)
			if eu != nil {
				// SN 키 등록 및 USER ID
				eu.SetUserID(&req.UserID)
				userSn = objmgr.GetUserSn()
				eu.SetUserSn(userSn)
				eu.SetState(serveruser.UserStateEnum.LoginSTATE)

				// 접속한 유저의 ID 등록
				objmgr.AddUserString(&req.UserID)

				log.Println("[", req.UserID, "] 유저가 접속하였습니다.")
			}
		}

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLoginUserRes)
		sendp.Write(result, userSn, &req.UserID)
		tcpserver.GetInstance().GetSendManager().SendToConn(conn, sendp)
		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		req := share.C2SPCLobbySendMsgReq{}
		p.Read(&req.Msg)
		eu := GetInstance().GetObjMgr().FindUser(conn)

		result := share.ResultSuccess

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLobbyMsgRes)
		sendp.Write(result, eu.GetUserID(), &req.Msg)

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		GetInstance().GetObjMgr().ForEachFunc(func(eu *serveruser.ExamUser) {
			if eu.GetState() == serveruser.UserStateEnum.LoginSTATE {
				tcpserver.GetInstance().GetSendManager().SendToConn(eu.GetConn(), sendp)
			}
		})
		return
	})
}
