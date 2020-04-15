package examserverlogic

import (
	"example/share"
	"log"
	"net"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/packet"
)

// RegistCommandLogic is
func RegistCommandLogic(lm *logicmanager.LogicManager) {
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {

		Userid := p.ReadString()
		log.Println("C2SPacketCommandLoginUserReq userid = ", Userid)

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandLoginUserRes)
		sendp.Write(uint32(1))
		_, err := conn.Write(sendp.GetByte())
		if err != nil {
			log.Println("Write S2CPacketCommandLoginUserRes", err)
			return
		}

		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalSendMsgReq, func(conn *net.TCPConn, p *packet.Packet) {

		msg := p.ReadString()
		log.Println("C2SPacketCommandGolobalSendMsgReq msg = ", msg)

		// Send 응답 패킷
		sendp := packet.NewPacket()
		sendp.SetHeader(share.ExamplePacketSerialkey, 0, share.S2CPacketCommandGolobalSendMsgRes)
		sendp.Write(uint32(1))
		_, err := conn.Write(sendp.GetByte())
		if err != nil {
			log.Println("Write S2CPacketCommandGolobalSendMsgRes ", err)
			return
		}

		return
	})
}
