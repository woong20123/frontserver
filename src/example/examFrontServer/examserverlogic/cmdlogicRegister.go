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
		log.Println("C2SPacketCommandLoginUserReq userid = ", req.UserID)

		// set user info
		objmgr := GetInstance().GetObjMgr()
		eu := objmgr.FindUser(conn)
		if eu != nil {
			// SN 키 등록 및 USER ID
			eu.SetUserID(&req.UserID)
			eu.SetUserSn(objmgr.GetUserSn())
			eu.SetState(serveruser.UserStateEnum.LoginSTATE)
		}

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLoginUserRes)
		sendp.Write(uint32(share.ResultSuccess), eu.GetUserSn())
		sendp.WriteString(req.UserID)
		tcpserver.GetObjInstance().GetSendManager().SendToConn(conn, sendp)
		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		req := share.C2SPCLobbySendMsgReq{}
		p.Read(&req.Msg)
		eu := GetInstance().GetObjMgr().FindUser(conn)

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLobbyMsgRes)
		sendp.Write(uint32(1))
		sendp.WriteString(eu.GetUserID())
		sendp.WriteString(req.Msg)

		// 로비에 있는 유저들에게 메시지를 보냅니다.
		GetInstance().GetObjMgr().ForEachFunc(func(eu *serveruser.ExamUser) {
			if eu.GetState() == serveruser.UserStateEnum.LoginSTATE {
				tcpserver.GetObjInstance().GetSendManager().SendToConn(eu.GetConn(), sendp)
			}
		})
		return
	})
}
