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

		Userid := p.ReadString()
		log.Println("C2SPacketCommandLoginUserReq userid = ", Userid)

		// regist user
		eu := serveruser.NewExamUser(1)
		eu.SetConn(conn)
		instance := GetInstance()
		instance.GetObjMgr().AddUser(1, eu)

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLoginUserRes)
		sendp.Write(uint32(share.ResultSuccess), uint32(100))
		sendp.WriteString(Userid)
		tcpserver.GetObjInstance().GetSendManager().SendToConn(conn, sendp)
		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalMsgReq, func(conn *net.TCPConn, p *packet.Packet) {

		msg := p.ReadString()
		log.Println("C2SPacketCommandGolobalSendMsgReq msg = ", msg)

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandGolobalMsgRes)
		sendp.Write(uint32(1))
		_, err := conn.Write(sendp.GetByte())
		if err != nil {
			log.Println("Write S2CPacketCommandGolobalSendMsgRes ", err)
			return
		}

		return
	})
}
